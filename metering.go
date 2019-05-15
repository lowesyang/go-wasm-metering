package go_wasm_metering

import (
	"encoding/json"
	"fmt"
	"github.com/yyh1102/go-wasm-metering/toolkit"
	"io/ioutil"
	"strconv"
)

const (
	defaultCostTablePath = "defaultCostTable.json"
	defaultModuleStr     = "metering"
	defaultFieldStr      = "usegas"
	defaultMeterType     = "i32"
	defaultCost          = 0
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
			if name, exist := section["Name"]; exist {
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
			name, exist := section["Name"]
			if exist {
				secId, exist := toolkit.J2W_SECTION_IDS[name.(string)]
				fmt.Printf("%v %v %v\n", name, secId, exist)
				if exist && secId > 0 && newSectionId < secId {
					fmt.Println("inject section")
					rest := append([]toolkit.JSON{}, module[i+1:]...)
					// insert the section at pos `i`
					module = append(module[:i], toolkit.JSON{
						"Name":    sectionName,
						"Entries": []interface{}{},
					})
					module = append(module, rest...)
					break
				}
			}
		}
		return module
	}

	// add necessary `type` and `import` sections if and only if they don't exist.
	if findSection(module, "type") == nil {
		module = createSection(module, "type")
	}
	if findSection(module, "import") == nil {
		module = createSection(module, "import")
	}

	importEntry := toolkit.JSON{
		"ModuleStr": m.opts.ModuleStr,
		"FieldStr":  m.opts.FieldStr,
		"Kind":      "function",
	}

	importType := toolkit.JSON{
		"Form":   "func",
		"Params": []string{m.opts.MeterType},
	}

	var (
		typeModule     toolkit.JSON
		functionModule toolkit.JSON
		funcIndex      int
		newModule      = make([]toolkit.JSON, len(module))
	)

	copy(newModule, module)

	for _, section := range newModule {
		sectionName, exist := section["Name"]
		if !exist {
			continue
		}
		switch sectionName.(string) {
		case "type":
			entries := section["Entries"].([]interface{})
			entries = append(entries, importType)
			section["Entries"] = entries
			importEntry["Type"] = len(entries) - 1
			// save for use for the code section.
			typeModule = section
		case "function":
			// save for use for the code section.
			functionModule = section
		case "import":
			entries := section["Entries"].([]interface{})
			for _, entry := range entries {
				importEntry := entry.(toolkit.JSON)
				if importEntry["ModuleStr"].(string) == m.opts.ModuleStr && importEntry["FieldStr"].(string) == m.opts.FieldStr {
					return nil, ErrImportMeterFunc
				}

				if importEntry["Kind"] == "function" {
					funcIndex += 1
				}
			}
			entries = append(entries, importEntry)
			// append the metering import.
			section["Entries"] = entries
			fmt.Printf("%#v\n", section)
		case "export":
			entries := section["Entries"].([]interface{})
			for _, entry := range entries {
				exportEntry := entry.(toolkit.JSON)
				entryIndex := exportEntry["Index"].(uint32)
				if exportEntry["Kind"].(string) == "function" && entryIndex >= uint32(funcIndex) {
					entryIndex += 1
					exportEntry["Index"] = entryIndex
				}
			}
		case "element":
			entries := section["Entries"].([]interface{})
			for _, entry := range entries {
				// remap element indices.
				elemEntry := entry.(toolkit.JSON)
				elements := elemEntry["Elements"].([]interface{})
				newElements := make([]uint64, 0, len(elements))
				for _, element := range elements {
					el := element.(uint64)
					if el >= uint64(funcIndex) {
						el += 1
					}
					newElements = append(newElements, el)
				}
				elemEntry["Elements"] = newElements
			}
		case "start":
			index := section["Index"].(uint32)
			if index >= uint32(funcIndex) {
				index += 1
			}
			section["Index"] = index
		case "code":
			entries := section["Entries"].([]interface{})
			for i, entry := range entries {
				typeIndex := functionModule["Entries"].([]interface{})[i].(uint64)
				typ := typeModule["Entries"].([]interface{})[typeIndex].(toolkit.JSON)
				cost := m.getCost(typ, m.costTable["Type"].(toolkit.JSON))

				m.meterCodeEntry(entry.(toolkit.CodeBody), m.costTable["Code"].(toolkit.JSON), funcIndex, cost)
			}
		}
	}
	return newModule, nil
}

// getCost returns the cost of an operation for the entry in a section from the cost table.
func (m *Metering) getCost(j interface{}, costTable toolkit.JSON) (cost uint64) {
	var costDefault uint64
	if dc, exist := costTable["DEFAULT"]; exist {
		costDefault = dc.(uint64)
	} else {
		costDefault = defaultCost
	}

	switch jj := j.(type) {
	case []interface{}:
		for _, el := range jj {
			cost += m.getCost(el, costTable)
		}
	case toolkit.JSON:
		for propName := range jj {
			propCost, exist := costTable[propName]
			if exist {
				cost += m.getCost(jj[propName], propCost.(toolkit.JSON))
			}
		}
	case string:
		c, exist := costTable[jj]
		if exist {
			cost = c.(uint64)
		} else {
			cost = costDefault
		}
	default:
		cost = costDefault
	}
	return
}

// meterCodeEntry meters a single code entry (see toolkit.CodeBody).
func (m *Metering) meterCodeEntry(entry toolkit.CodeBody, costTable toolkit.JSON, meterFuncIndex int, cost uint64) interface{} {
	meterType := m.opts.MeterType

	meteringStatement := func(cost uint64, meteringImportIndex int) (ops []toolkit.OP) {
		opsJson := toolkit.Text2Json(fmt.Sprintf("%s.const %v call %v", meterType, cost, meteringImportIndex))
		for _, op := range opsJson {
			oop := toolkit.OP{
				Name:       op["Name"].(string),
				Type:       op["Type"].(string),
				Immediates: op["Immediates"],
			}

			if rt, ok := op["ReturnType"]; ok {
				oop.ReturnType = rt.(string)
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
			sum += m.getCost(op.Name, m.costTable["Code"].(toolkit.JSON))
		}
		return sum
	}

	var (
		meteringCost = meterTheMeteringStatement()
		code         = make([]toolkit.OP, len(entry.Code))
		meteredCode  []toolkit.OP
	)
	// create a code copy.
	copy(code, entry.Code)

	cost += m.getCost(entry.Locals, costTable["Locals"].(toolkit.JSON))

	for len(code) > 0 {
		i := 0

		// meter a segment of wasm code.
		for {
			op := code[i]
			i += 1

			remapOp(op, meterFuncIndex)
			cost += m.getCost(op.Name, costTable["Code"].(toolkit.JSON))

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
