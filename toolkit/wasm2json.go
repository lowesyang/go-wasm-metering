package toolkit

import (
	"strings"
)

var (
	// https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#language-types
	// All types are distinguished by a negative varint7 values that is the first
	// byte of their encoding (representing a type constructor)
	W2J_LANGUAGE_TYPES = map[byte]string{
		0x7f: "i32",
		0x7e: "i64",
		0x7d: "f32",
		0x7c: "f64",
		0x70: "anyFunc",
		0x60: "func",
		0x40: "block_type",
	}

	// https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#external_kind
	// A single-byte unsigned integer indicating the kind of definition being imported or defined:
	W2J_EXTERNAL_KIND = map[byte]string{
		0x00: "function",
		0x01: "table",
		0x02: "memory",
		0x03: "global",
	}

	W2J_OPCODES = map[byte]string{
		// flow control
		0x0: "unreachable",
		0x1: "nop",
		0x2: "block",
		0x3: "loop",
		0x4: "if",
		0x5: "else",
		0xb: "end",
		0xc: "br",
		0xd: "br_if",
		0xe: "br_table",
		0xf: "return",

		// calls
		0x10: "call",
		0x11: "call_indirect",

		// Parametric operators
		0x1a: "drop",
		0x1b: "select",

		// Varibale access
		0x20: "get_local",
		0x21: "set_local",
		0x22: "tee_local",
		0x23: "get_global",
		0x24: "set_global",

		// Memory-related operators
		0x28: "i32.load",
		0x29: "i64.load",
		0x2a: "f32.load",
		0x2b: "f64.load",
		0x2c: "i32.load8_s",
		0x2d: "i32.load8_u",
		0x2e: "i32.load16_s",
		0x2f: "i32.load16_u",
		0x30: "i64.load8_s",
		0x31: "i64.load8_u",
		0x32: "i64.load16_s",
		0x33: "i64.load16_u",
		0x34: "i64.load32_s",
		0x35: "i64.load32_u",
		0x36: "i32.store",
		0x37: "i64.store",
		0x38: "f32.store",
		0x39: "f64.store",
		0x3a: "i32.store8",
		0x3b: "i32.store16",
		0x3c: "i64.store8",
		0x3d: "i64.store16",
		0x3e: "i64.store32",
		0x3f: "current_memory",
		0x40: "grow_memory",

		// Constants
		0x41: "i32.const",
		0x42: "i64.const",
		0x43: "f32.const",
		0x44: "f64.const",

		// Comparison operators
		0x45: "i32.eqz",
		0x46: "i32.eq",
		0x47: "i32.ne",
		0x48: "i32.lt_s",
		0x49: "i32.lt_u",
		0x4a: "i32.gt_s",
		0x4b: "i32.gt_u",
		0x4c: "i32.le_s",
		0x4d: "i32.le_u",
		0x4e: "i32.ge_s",
		0x4f: "i32.ge_u",
		0x50: "i64.eqz",
		0x51: "i64.eq",
		0x52: "i64.ne",
		0x53: "i64.lt_s",
		0x54: "i64.lt_u",
		0x55: "i64.gt_s",
		0x56: "i64.gt_u",
		0x57: "i64.le_s",
		0x58: "i64.le_u",
		0x59: "i64.ge_s",
		0x5a: "i64.ge_u",
		0x5b: "f32.eq",
		0x5c: "f32.ne",
		0x5d: "f32.lt",
		0x5e: "f32.gt",
		0x5f: "f32.le",
		0x60: "f32.ge",
		0x61: "f64.eq",
		0x62: "f64.ne",
		0x63: "f64.lt",
		0x64: "f64.gt",
		0x65: "f64.le",
		0x66: "f64.ge",

		// Numeric operators
		0x67: "i32.clz",
		0x68: "i32.ctz",
		0x69: "i32.popcnt",
		0x6a: "i32.add",
		0x6b: "i32.sub",
		0x6c: "i32.mul",
		0x6d: "i32.div_s",
		0x6e: "i32.div_u",
		0x6f: "i32.rem_s",
		0x70: "i32.rem_u",
		0x71: "i32.and",
		0x72: "i32.or",
		0x73: "i32.xor",
		0x74: "i32.shl",
		0x75: "i32.shr_s",
		0x76: "i32.shr_u",
		0x77: "i32.rotl",
		0x78: "i32.rotr",
		0x79: "i64.clz",
		0x7a: "i64.ctz",
		0x7b: "i64.popcnt",
		0x7c: "i64.add",
		0x7d: "i64.sub",
		0x7e: "i64.mul",
		0x7f: "i64.div_s",
		0x80: "i64.div_u",
		0x81: "i64.rem_s",
		0x82: "i64.rem_u",
		0x83: "i64.and",
		0x84: "i64.or",
		0x85: "i64.xor",
		0x86: "i64.shl",
		0x87: "i64.shr_s",
		0x88: "i64.shr_u",
		0x89: "i64.rotl",
		0x8a: "i64.rotr",
		0x8b: "f32.abs",
		0x8c: "f32.neg",
		0x8d: "f32.ceil",
		0x8e: "f32.floor",
		0x8f: "f32.trunc",
		0x90: "f32.nearest",
		0x91: "f32.sqrt",
		0x92: "f32.add",
		0x93: "f32.sub",
		0x94: "f32.mul",
		0x95: "f32.div",
		0x96: "f32.min",
		0x97: "f32.max",
		0x98: "f32.copysign",
		0x99: "f64.abs",
		0x9a: "f64.neg",
		0x9b: "f64.ceil",
		0x9c: "f64.floor",
		0x9d: "f64.trunc",
		0x9e: "f64.nearest",
		0x9f: "f64.sqrt",
		0xa0: "f64.add",
		0xa1: "f64.sub",
		0xa2: "f64.mul",
		0xa3: "f64.div",
		0xa4: "f64.min",
		0xa5: "f64.max",
		0xa6: "f64.copysign",

		// Conversions
		0xa7: "i32.wrap/i64",
		0xa8: "i32.trunc_s/f32",
		0xa9: "i32.trunc_u/f32",
		0xaa: "i32.trunc_s/f64",
		0xab: "i32.trunc_u/f64",
		0xac: "i64.extend_s/i32",
		0xad: "i64.extend_u/i32",
		0xae: "i64.trunc_s/f32",
		0xaf: "i64.trunc_u/f32",
		0xb0: "i64.trunc_s/f64",
		0xb1: "i64.trunc_u/f64",
		0xb2: "f32.convert_s/i32",
		0xb3: "f32.convert_u/i32",
		0xb4: "f32.convert_s/i64",
		0xb5: "f32.convert_u/i64",
		0xb6: "f32.demote/f64",
		0xb7: "f64.convert_s/i32",
		0xb8: "f64.convert_u/i32",
		0xb9: "f64.convert_s/i64",
		0xba: "f64.convert_u/i64",
		0xbb: "f64.promote/f32",

		// Reinterpretations
		0xbc: "i32.reinterpret/f32",
		0xbd: "i64.reinterpret/f64",
		0xbe: "f32.reinterpret/i32",
		0xbf: "f64.reinterpret/i64",
	}

	W2J_SECTION_IDS = map[byte]string{
		0:  "custom",
		1:  "type",
		2:  "import",
		3:  "function",
		4:  "table",
		5:  "memory",
		6:  "global",
		7:  "export",
		8:  "start",
		9:  "element",
		10: "code",
		11: "data",
	}

	immeParsers = immediataryParsers{}
	tParsers    = typeParsers{}
	secParsers  = sectionParsers{}
)

type immediataryParsers struct{}

func (immediataryParsers) Varuint1(stream *Stream) int8 {
	return int8(stream.ReadByte())
}

func (immediataryParsers) Varuint32(stream *Stream) uint32 {
	return uint32(DecodeULEB128(stream))
}

func (immediataryParsers) Varint32(stream *Stream) int32 {
	return int32(DecodeSLEB128(stream))
}

func (immediataryParsers) Varint64(stream *Stream) int64 {
	return int64(DecodeSLEB128(stream))
}

func (immediataryParsers) Uint32(stream *Stream) []byte {
	return stream.Read(4)
}

func (immediataryParsers) Uint64(stream *Stream) []byte {
	return stream.Read(8)
}

func (immediataryParsers) BlockType(stream *Stream) string {
	return W2J_LANGUAGE_TYPES[stream.ReadByte()]
}

func (immediataryParsers) BrTable(stream *Stream) JSON {
	jsonObj := make(JSON)
	targets := []uint64{}

	num := DecodeULEB128(stream)
	for i := uint64(0); i < num; i++ {
		target := DecodeULEB128(stream)
		targets = append(targets, target)
	}

	jsonObj["targets"] = targets
	jsonObj["default_target"] = DecodeULEB128(stream)
	return jsonObj
}

func (immediataryParsers) CallIndirect(stream *Stream) JSON {
	jsonObj := make(JSON)
	jsonObj["index"] = DecodeULEB128(stream)
	jsonObj["reserved"] = stream.ReadByte()
	return jsonObj
}

func (immediataryParsers) MemoryImmediate(stream *Stream) JSON {
	jsonObj := make(JSON)
	jsonObj["flags"] = DecodeULEB128(stream)
	jsonObj["offset"] = DecodeULEB128(stream)
	return jsonObj
}

type typeParsers struct{}

func (typeParsers) Function(stream *Stream) uint64 {
	return DecodeULEB128(stream)
}

func (t typeParsers) Table(stream *Stream) Table {
	typ := stream.ReadByte()
	return Table{
		ElementType: W2J_LANGUAGE_TYPES[typ],
		Limits:      t.Memory(stream),
	}
}

func (typeParsers) Global(stream *Stream) Global {
	typ := stream.ReadByte()
	mutability := stream.ReadByte()
	return Global{
		ContentType: W2J_LANGUAGE_TYPES[typ],
		Mutability:  mutability,
	}
}

func (typeParsers) Memory(stream *Stream) MemLimits {
	flags := DecodeULEB128(stream)
	intial := DecodeULEB128(stream)
	limits := MemLimits{
		Flags:  flags,
		Intial: intial,
	}
	if flags == 1 {
		limits.Maximum = DecodeULEB128(stream)
	}
	return limits
}

func (typeParsers) InitExpr(stream *Stream) OP {
	op := ParseOp(stream)
	stream.ReadByte() // skip the `end`
	return op
}

type sectionParsers struct{}

func (sectionParsers) Custom(stream *Stream, header SectionHeader) CustomSec {
	sec := CustomSec{Name: "custom"}

	// create a new stream to read.
	section := NewStream(stream.Read(int(header.Size)))
	nameLen := DecodeULEB128(section)
	name := section.Read(int(nameLen))

	sec.SectionName = string(name)
	sec.Payload = section.String()
	return sec
}

func (sectionParsers) Type(stream *Stream) TypeSec {
	numberOfEntries := DecodeULEB128(stream)
	typSec := TypeSec{
		Name:    "type",
		Entries: []TypeEntry{},
	}

	for i := uint64(0); i < numberOfEntries; i++ {
		typ := stream.ReadByte()
		entry := TypeEntry{
			Form:   W2J_LANGUAGE_TYPES[typ],
			Params: []string{},
		}

		paramCount := DecodeULEB128(stream)

		// parse the entries.
		for j := uint64(0); j < paramCount; j++ {
			typ := stream.ReadByte()
			entry.Params = append(entry.Params, W2J_LANGUAGE_TYPES[typ])
		}

		numOfReturns := DecodeULEB128(stream)
		if numOfReturns > 0 {
			typ = stream.ReadByte()
			entry.ReturnType = W2J_LANGUAGE_TYPES[typ]
		}

		typSec.Entries = append(typSec.Entries, entry)
	}

	return typSec
}

func (s sectionParsers) Import(stream *Stream) ImportSec {
	numberOfEntries := DecodeULEB128(stream)
	importSec := ImportSec{
		Name:    "import",
		Entries: []ImportEntry{},
	}

	for i := uint64(0); i < numberOfEntries; i++ {
		moduleLen := DecodeULEB128(stream)
		moduleStr := stream.Read(int(moduleLen))

		fieldLen := DecodeULEB128(stream)
		fieldStr := stream.Read(int(fieldLen))

		kind := stream.ReadByte()
		externalKind := W2J_EXTERNAL_KIND[kind]
		var returned interface{}
		switch externalKind {
		case "function":
			returned = tParsers.Function(stream)
		case "table":
			returned = tParsers.Table(stream)
		case "memory":
			returned = tParsers.Memory(stream)
		case "global":
			returned = tParsers.Global(stream)
		}

		entry := ImportEntry{
			ModuleStr: string(moduleStr),
			FieldStr:  string(fieldStr),
			Kind:      externalKind,
			Type:      returned,
		}

		importSec.Entries = append(importSec.Entries, entry)
	}

	return importSec
}

func (sectionParsers) Function(stream *Stream) FuncSec {
	numberOfEntries := DecodeULEB128(stream)
	funcSec := FuncSec{
		Name:    "function",
		Entries: []uint64{},
	}

	for i := uint64(0); i < numberOfEntries; i++ {
		entry := DecodeULEB128(stream)
		funcSec.Entries = append(funcSec.Entries, entry)
	}
	return funcSec
}

func (s sectionParsers) Table(stream *Stream) TableSec {
	numberOfEntries := DecodeULEB128(stream)
	tableSec := TableSec{
		Name:    "table",
		Entries: []Table{},
	}

	// parse table_type.
	tparser := typeParsers{}
	for i := uint64(0); i < numberOfEntries; i++ {
		entry := tparser.Table(stream)
		tableSec.Entries = append(tableSec.Entries, entry)
	}

	return tableSec
}

func (sectionParsers) Memory(stream *Stream) MemSec {
	numberOfEntries := DecodeULEB128(stream)
	memSec := MemSec{
		Name:    "memory",
		Entries: []MemLimits{},
	}

	tparser := typeParsers{}
	for i := uint64(0); i < numberOfEntries; i++ {
		entry := tparser.Memory(stream)
		memSec.Entries = append(memSec.Entries, entry)
	}
	return memSec
}

func (sectionParsers) Global(stream *Stream) GlobalSec {
	numberOfEntries := DecodeULEB128(stream)
	globalSec := GlobalSec{
		Name:    "global",
		Entries: []GlobalEntry{},
	}

	tparser := typeParsers{}
	for i := uint64(0); i < numberOfEntries; i++ {
		entry := GlobalEntry{
			Type: tparser.Global(stream),
			Init: tparser.InitExpr(stream),
		}

		globalSec.Entries = append(globalSec.Entries, entry)
	}

	return globalSec
}

func (sectionParsers) Export(stream *Stream) ExportSec {
	numberOfEntries := DecodeULEB128(stream)
	exportSec := ExportSec{
		Name:    "export",
		Entries: []ExportEntry{},
	}

	for i := uint64(0); i < numberOfEntries; i++ {
		strLength := DecodeULEB128(stream)
		fieldStr := string(stream.Read(int(strLength)))
		kind := stream.ReadByte()
		index := DecodeULEB128(stream)

		entry := ExportEntry{
			FieldStr: fieldStr,
			Kind:     W2J_EXTERNAL_KIND[kind],
			Index:    uint32(index),
		}

		exportSec.Entries = append(exportSec.Entries, entry)
	}

	return exportSec
}

func (sectionParsers) Start(stream *Stream) StartSec {
	startSec := StartSec{
		Name:  "start",
		Index: uint32(DecodeULEB128(stream)),
	}
	return startSec
}

func (sectionParsers) Element(stream *Stream) ElementSec {
	numberOfEntries := DecodeULEB128(stream)
	elSec := ElementSec{
		Name:    "element",
		Entries: []ElementEntry{},
	}

	tparser := typeParsers{}
	for i := uint64(0); i < numberOfEntries; i++ {
		entry := ElementEntry{}
		entry.Index = uint32(DecodeULEB128(stream))
		entry.Offset = tparser.InitExpr(stream)

		numElem := DecodeULEB128(stream)
		for j := uint64(0); j < numElem; j++ {
			elem := DecodeULEB128(stream)
			entry.Elements = append(entry.Elements, elem)
		}

		elSec.Entries = append(elSec.Entries, entry)
	}

	return elSec
}

func (sectionParsers) Code(stream *Stream) CodeSec {
	numberOfEntries := DecodeULEB128(stream)
	codeSec := CodeSec{
		Name:    "code",
		Entries: []CodeBody{},
	}

	for i := uint64(0); i < numberOfEntries; i++ {
		codeBody := CodeBody{
			Locals: []LocalEntry{},
			Code:   []OP{},
		}

		bodySize := DecodeULEB128(stream)
		endBytes := stream.bytesRead + int(bodySize)

		// parse locals
		localCount := DecodeULEB128(stream)
		for j := uint64(0); j < localCount; j++ {
			local := LocalEntry{}
			local.Count = uint32(DecodeULEB128(stream))
			local.Type = W2J_LANGUAGE_TYPES[stream.ReadByte()]
			codeBody.Locals = append(codeBody.Locals, local)
		}

		// parse code
		for stream.bytesRead < endBytes {
			op := ParseOp(stream)
			codeBody.Code = append(codeBody.Code, op)
		}

		codeSec.Entries = append(codeSec.Entries, codeBody)
	}

	return codeSec
}

func (sectionParsers) Data(stream *Stream) DataSec {
	numberOfEntries := DecodeULEB128(stream)
	dataSec := DataSec{
		Name:    "data",
		Entries: []DataSegment{},
	}

	tparser := typeParsers{}
	for i := uint64(0); i < numberOfEntries; i++ {
		entry := DataSegment{}
		entry.Index = uint32(DecodeULEB128(stream))
		entry.Offset = tparser.InitExpr(stream)
		segmentSize := DecodeULEB128(stream)
		entry.Data = append([]byte{}, stream.Read(int(segmentSize))...)

		dataSec.Entries = append(dataSec.Entries, entry)
	}

	return dataSec
}

// Wasm2Json convert the wasm binary to a JSON array output.
func Wasm2Json(buf []byte) []JSON {
	stream := NewStream(buf)
	preramble := ParsePreramble(stream)
	resJson := []JSON{preramble}

	for stream.Len() != 0 {
		header := ParseSectionHeader(stream)
		//fmt.Printf("%#v\n", header)
		jsonObj := make(JSON)
		switch header.Name {
		case "custom":
			rsec := secParsers.Custom(stream, header)
			jsonObj["name"] = rsec.Name
			jsonObj["section_name"] = rsec.SectionName
			jsonObj["payload"] = rsec.Payload
		case "type":
			rsec := secParsers.Type(stream)
			jsonObj["name"] = rsec.Name
			jsonObj["entries"] = rsec.Entries
		case "import":
			rsec := secParsers.Import(stream)
			jsonObj["name"] = rsec.Name
			jsonObj["entries"] = rsec.Entries
		case "function":
			rsec := secParsers.Function(stream)
			jsonObj["name"] = rsec.Name
			jsonObj["entries"] = rsec.Entries
		case "table":
			rsec := secParsers.Table(stream)
			jsonObj["name"] = rsec.Name
			jsonObj["entries"] = rsec.Entries
		case "memory":
			rsec := secParsers.Memory(stream)
			jsonObj["name"] = rsec.Name
			jsonObj["entries"] = rsec.Entries
		case "global":
			rsec := secParsers.Global(stream)
			jsonObj["name"] = rsec.Name
			jsonObj["entries"] = rsec.Entries
		case "export":
			rsec := secParsers.Export(stream)
			jsonObj["name"] = rsec.Name
			jsonObj["entries"] = rsec.Entries
		case "start":
			rsec := secParsers.Start(stream)
			jsonObj["name"] = rsec.Name
			jsonObj["index"] = rsec.Index
		case "element":
			rsec := secParsers.Element(stream)
			jsonObj["name"] = rsec.Name
			jsonObj["entries"] = rsec.Entries
		case "code":
			rsec := secParsers.Code(stream)
			jsonObj["name"] = rsec.Name
			jsonObj["entries"] = rsec.Entries
		case "data":
			rsec := secParsers.Data(stream)
			jsonObj["name"] = rsec.Name
			jsonObj["entries"] = rsec.Entries
		}

		resJson = append(resJson, jsonObj)
	}

	return resJson
}

func ParsePreramble(stream *Stream) JSON {
	magic := stream.Read(4)
	version := stream.Read(4)

	jsonObj := make(JSON)
	jsonObj["name"] = "preramble"
	jsonObj["magic"] = magic
	jsonObj["version"] = version

	return jsonObj
}

func ParseSectionHeader(stream *Stream) SectionHeader {
	id := stream.ReadByte()
	return SectionHeader{
		Id:   id,
		Name: W2J_SECTION_IDS[id],
		Size: DecodeULEB128(stream),
	}
}

func ParseOp(stream *Stream) OP {
	finalOP := OP{}
	op := stream.ReadByte()
	fullName := strings.Split(W2J_OPCODES[op], ".")
	var (
		typ           = fullName[0]
		name          string
		immediatesKey string
	)

	if len(fullName) < 2 {
		name = typ
	} else {
		name = fullName[1]
		finalOP.ReturnType = typ
	}

	finalOP.Name = name

	if name == "const" {
		immediatesKey = typ
	} else {
		immediatesKey = name
	}
	immediates, exist := OP_IMMEDIATES[immediatesKey]
	if exist {
		var returned interface{}
		switch immediates {
		case "block_type":
			returned = immeParsers.BlockType(stream)
		case "call_indirect":
			returned = immeParsers.CallIndirect(stream)
		case "varuint32":
			returned = immeParsers.Varuint32(stream)
		case "varuint1":
			returned = immeParsers.Varuint1(stream)
		case "varint32":
			returned = immeParsers.Varint32(stream)
		case "varint64":
			returned = immeParsers.Varint64(stream)
		case "uint32":
			returned = immeParsers.Uint32(stream)
		case "uint64":
			returned = immeParsers.Uint64(stream)
		case "br_table":
			returned = immeParsers.BrTable(stream)
		case "memory_immediate":
			returned = immeParsers.MemoryImmediate(stream)
		}
		finalOP.Immediates = returned
	}

	return finalOP
}
