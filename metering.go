package go_wasm_metering

import (
	"encoding/json"
	"fmt"
	"github.com/yyh1102/go-wasm-metering/toolkit"
	"io/ioutil"
	"reflect"
	"strconv"
)

const (
	defaultCostTablePath = "defaultCostTable.json"
	defaultModuleStr     = "metering"
	defaultFieldStr      = "usegas"
	defaultMeterType     = "i64"
	defaultCost          = uint64(0)
)

var (
	branchOps = map[string]struct{}{
		"grow_memory": struct{}{},
		"end":         struct{}{},
		"br":          struct{}{},
		"br_table":    struct{}{},
		"br_if":       struct{}{},
		"if":          struct{}{},
		"else":        struct{}{},
		"return":      struct{}{},
		"loop":        struct{}{},
	}
)

// MeterWASM injects metering into WebAssembly binary code.
// This func is the real exported function used by outer callers.
func MeterWASM(wasm []byte, opts *Options) ([]byte, error) {
	module := toolkit.Wasm2Json(wasm)
	if opts == nil {
		opts = &Options{}
	}
	metering, err := newMetring(*opts)
	if err != nil {
		return nil, err
	}
	module, err = metering.meterJSON(module)
	if err != nil {
		return nil, err
	}
	return toolkit.Json2Wasm(module), nil
}

type Options struct {
	CostTable string // path of cost table file.
	ModuleStr string // the import string for metering function.
	FieldStr  string // the field string for the metering function.
	MeterType string // the register type that is used to meter. Can be `i64`, `i32`, `f64`, `f32`.
}

type Metering struct {
	costTable toolkit.JSON
	opts      Options
}

func newMetring(opts Options) (*Metering, error) {
	// set defaults.
	if opts.CostTable == "" {
		opts.CostTable = defaultCostTablePath
	}

	if opts.ModuleStr == "" {
		opts.ModuleStr = defaultModuleStr
	}

	if opts.FieldStr == "" {
		opts.FieldStr = defaultFieldStr
	}

	if opts.MeterType == "" {
		opts.MeterType = defaultMeterType
	}

	table, err := ioutil.ReadFile(opts.CostTable)
	if err != nil {
		return nil, err
	}

	costTable := make(toolkit.JSON)
	if err := json.Unmarshal(table, &costTable); err != nil {
		return nil, err
	}

	return &Metering{
		costTable: costTable,
		opts:      opts,
	}, nil
}

// meterJSON injects metering into a JSON output of Wasm2Json.
func (m *Metering) meterJSON(module []toolkit.JSON) ([]toolkit.JSON, error) {
	// find section.
	findSection := func(module []toolkit.JSON, sectionName string) toolkit.JSON {
		for _, section := range module {
			if name, exist := section["name"]; exist {
				if name.(string) == sectionName {
					return section
				}
			}
		}
		return nil
	}

	// create section.
	createSection := func(module []toolkit.JSON, sectionName string) []toolkit.JSON {
		newSectionId := toolkit.J2W_SECTION_IDS[sectionName]
		for i, section := range module {
			name, exist := section["name"]
			if exist {
				secId, exist := toolkit.J2W_SECTION_IDS[name.(string)]
				//fmt.Printf("%v %v %v\n", name, secId, exist)
				if exist && secId > 0 && newSectionId < secId {
					rest := append([]toolkit.JSON{}, module[i:]...)
					// insert the section at pos `i`
					module = append(module[:i], toolkit.JSON{
						"name": sectionName,
					})
					module = append(module, rest...)
					break
				}
			}
		}
		return module
	}
	//fmt.Printf("%#v\n", module)

	// add necessary `type` and `import` sections if and only if they don't exist.
	if findSection(module, "type") == nil {
		module = createSection(module, "type")
	}
	if findSection(module, "import") == nil {
		module = createSection(module, "import")
	}

	importEntry := toolkit.ImportEntry{
		ModuleStr: m.opts.ModuleStr,
		FieldStr:  m.opts.FieldStr,
		Kind:      "function",
	}

	importType := toolkit.TypeEntry{
		Form:   "func",
		Params: []string{m.opts.MeterType},
	}

	var (
		typeModule     toolkit.JSON
		functionModule toolkit.JSON
		funcIndex      int
		newModule      = make([]toolkit.JSON, len(module))
	)

	copy(newModule, module)
	//fmt.Printf("%#v", newModule)

	for _, section := range newModule {
		sectionName, exist := section["name"]
		if !exist {
			continue
		}
		switch sectionName.(string) {
		case "type":
			var entries []toolkit.TypeEntry
			ientries, exist := section["entries"]
			if exist {
				entries = ientries.([]toolkit.TypeEntry)
			}
			//fmt.Printf("Entries %#v\n", entries)
			entries = append(entries, importType)
			section["entries"] = entries
			importEntry.Type = uint64(len(entries) - 1)
			// save for use for the code section.
			typeModule = section
		case "function":
			// save for use for the code section.
			functionModule = section
		case "import":
			var entries []toolkit.ImportEntry
			ientries, exist := section["entries"]
			if exist {
				entries = ientries.([]toolkit.ImportEntry)
			}
			for _, entry := range entries {
				if entry.ModuleStr == m.opts.ModuleStr && entry.FieldStr == m.opts.FieldStr {
					return nil, ErrImportMeterFunc
				}

				if entry.Kind == "function" {
					funcIndex += 1
				}
			}
			// append the metering import.
			section["entries"] = append(entries, importEntry)
		case "export":
			entries := section["entries"].([]toolkit.ExportEntry)
			for i, entry := range entries {
				if entry.Kind == "function" && entry.Index >= uint32(funcIndex) {
					entries[i].Index = entry.Index + 1
				}
			}
		case "element":
			entries := section["entries"].([]toolkit.ElementEntry)
			for _, entry := range entries {
				// remap element indices.
				newElements := make([]uint64, 0, len(entry.Elements))
				for _, el := range entry.Elements {
					if el >= uint64(funcIndex) {
						el += 1
					}
					newElements = append(newElements, el)
				}
				entry.Elements = newElements
			}
		case "start":
			index := section["index"].(uint32)
			if index >= uint32(funcIndex) {
				index += 1
			}
			section["index"] = index
		case "code":
			entries := section["entries"].([]toolkit.CodeBody)
			funcEntries := functionModule["entries"].([]uint64)
			typEntries := typeModule["entries"].([]toolkit.TypeEntry)
			for i, entry := range entries {
				typeIndex := funcEntries[i]
				typ := typEntries[typeIndex]
				cost := getCost(typ, m.costTable["type"].(toolkit.JSON), defaultCost)

				entries[i] = meterCodeEntry(entry, m.costTable["code"].(toolkit.JSON), m.opts.MeterType, funcIndex, cost)
			}
		}
	}
	return newModule, nil
}

// getCost returns the cost of an operation for the entry in a section from the cost table.
func getCost(j interface{}, costTable toolkit.JSON, defaultCost uint64) (cost uint64) {
	if dc, exist := costTable["DEFAULT"]; exist {
		defaultCost = uint64(dc.(float64))
	}
	rval := reflect.ValueOf(j)
	kind := rval.Type().Kind()
	if kind == reflect.Slice {
		for i := 0; i < rval.Len(); i++ {
			cost += getCost(rval.Index(i).Interface(), costTable, 0)
		}
	} else if kind == reflect.Struct {
		rtype := rval.Type()
		for i := 0; i < rval.NumField(); i++ {
			rv := rval.Field(i)
			propCost, exist := costTable[toolkit.Lcfirst(rtype.Field(i).Name)]
			if exist {
				cost += getCost(rv.Interface(), propCost.(toolkit.JSON), defaultCost)
			}
		}
	} else if kind == reflect.String {
		c, exist := costTable[j.(string)]
		if exist {
			cost = uint64(c.(float64))
		} else {
			cost = defaultCost
		}
	} else {
		cost = defaultCost
	}
	//fmt.Printf("json %#v cost %v\n", j, cost)
	return
}

// meterCodeEntry meters a single code entry (see toolkit.CodeBody).
func meterCodeEntry(entry toolkit.CodeBody, costTable toolkit.JSON, meterType string, meterFuncIndex int, cost uint64) toolkit.CodeBody {
	getImmediateFromOP := func(name, opType string) string {
		var immediatesKey string
		if name == "const" {
			immediatesKey = opType
		} else {
			immediatesKey = name
		}
		return toolkit.OP_IMMEDIATES[immediatesKey]
	}

	meteringStatement := func(cost uint64, meteringImportIndex int) (ops []toolkit.OP) {
		opsJson := toolkit.Text2Json(fmt.Sprintf("%s.const %v call %v", meterType, cost, meteringImportIndex))
		for _, op := range opsJson {
			//immediates, _ := strconv.ParseUint(op["immediates"].(string), 10, 64)

			oop := toolkit.OP{
				Name: op["name"].(string),
			}

			// convert immediates.
			imm := getImmediateFromOP(oop.Name, meterType)
			if imm != "" {
				opImm := op["immediates"]
				switch imm {
				case "varuint1":
					imme, _ := strconv.ParseInt(opImm.(string), 10, 8)
					oop.Immediates = int8(imme)
				case "varuint32":
					imme, _ := strconv.ParseUint(opImm.(string), 10, 32)
					oop.Immediates = uint32(imme)
				case "varint32":
					imme, _ := strconv.ParseInt(opImm.(string), 10, 32)
					oop.Immediates = int32(imme)
				case "varint64":
					imme, _ := strconv.ParseInt(opImm.(string), 10, 64)
					oop.Immediates = int64(imme)
				case "uint32":
					oop.Immediates = opImm.([]byte)
				case "uint64":
					oop.Immediates = opImm.([]byte)
				case "block_type":
					oop.Immediates = opImm.(string)
				case "br_table", "call_indirect", "memory_immediate":
					oop.Immediates = opImm.(toolkit.JSON)
				}
				//fmt.Printf("immediate %v %v\n", imm, oop.Immediates)
			}

			if rt, ok := op["return_type"]; ok {
				oop.ReturnType = rt.(string)
			}

			if rt, ok := op["type"]; ok {
				oop.Type = rt.(string)
			}

			ops = append(ops, oop)
		}
		return
	}

	remapOp := func(op toolkit.OP, funcIndex int) {
		if op.Name == "call" {
			immediates := op.Immediates.(string)
			rv, _ := strconv.ParseInt(immediates, 10, 64)
			if rv >= int64(funcIndex) {
				rv += 1
			}
			op.Immediates = strconv.FormatInt(rv, 10)
		}
	}

	meterTheMeteringStatement := func() uint64 {
		code := meteringStatement(0, 0)
		// sum the operations cost
		sum := uint64(0)
		for _, op := range code {
			sum += getCost(op.Name, costTable["code"].(toolkit.JSON), defaultCost)
		}
		return sum
	}

	var (
		meteringCost = meterTheMeteringStatement()
		code         = make([]toolkit.OP, len(entry.Code))
		meteredCode  []toolkit.OP
	)
	//fmt.Printf("meter the meter cost %d\n", meteringCost)

	// create a code copy.
	copy(code, entry.Code)

	cost += getCost(entry.Locals, costTable["locals"].(toolkit.JSON), defaultCost)

	for len(code) > 0 {
		i := 0

		// meter a segment of wasm code.
		for {
			op := code[i]
			i += 1

			remapOp(op, meterFuncIndex)
			cost += getCost(op.Name, costTable["code"].(toolkit.JSON), defaultCost)

			if _, exist := branchOps[op.Name]; exist {
				break
			}
		}

		// add the metering statement.
		if cost != 0 {
			// add the cost of metering
			cost += meteringCost
			meteredCode = append(meteredCode, meteringStatement(cost, meterFuncIndex)...)
		}

		meteredCode = append(meteredCode, code[:i]...)
		code = code[i:]
		cost = 0
	}

	entry.Code = meteredCode
	return entry
}
