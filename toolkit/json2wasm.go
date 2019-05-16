package toolkit

import (
	"reflect"
)

var (
	J2W_LANGUAGE_TYPES = map[string]byte{
		"i32":        0x7f,
		"i64":        0x7e,
		"f32":        0x7d,
		"f64":        0x7c,
		"anyFunc":    0x70,
		"func":       0x60,
		"block_type": 0x40,
	}

	J2W_EXTERNAL_KIND = map[string]byte{
		"function": 0,
		"table":    1,
		"memory":   2,
		"global":   3,
	}

	J2W_SECTION_IDS = map[string]byte{
		"custom":   0,
		"type":     1,
		"import":   2,
		"function": 3,
		"table":    4,
		"memory":   5,
		"global":   6,
		"export":   7,
		"start":    8,
		"element":  9,
		"code":     10,
		"data":     11,
	}

	J2W_OPCODES = map[string]byte{
		"unreachable":         0x0,
		"nop":                 0x1,
		"block":               0x2,
		"loop":                0x3,
		"if":                  0x4,
		"else":                0x5,
		"end":                 0xb,
		"br":                  0xc,
		"br_if":               0xd,
		"br_table":            0xe,
		"return":              0xf,
		"call":                0x10,
		"call_indirect":       0x11,
		"drop":                0x1a,
		"select":              0x1b,
		"get_local":           0x20,
		"set_local":           0x21,
		"tee_local":           0x22,
		"get_global":          0x23,
		"set_global":          0x24,
		"i32.load":            0x28,
		"i64.load":            0x29,
		"f32.load":            0x2a,
		"f64.load":            0x2b,
		"i32.load8_s":         0x2c,
		"i32.load8_u":         0x2d,
		"i32.load16_s":        0x2e,
		"i32.load16_u":        0x2f,
		"i64.load8_s":         0x30,
		"i64.load8_u":         0x31,
		"i64.load16_s":        0x32,
		"i64.load16_u":        0x33,
		"i64.load32_s":        0x34,
		"i64.load32_u":        0x35,
		"i32.store":           0x36,
		"i64.store":           0x37,
		"f32.store":           0x38,
		"f64.store":           0x39,
		"i32.store8":          0x3a,
		"i32.store16":         0x3b,
		"i64.store8":          0x3c,
		"i64.store16":         0x3d,
		"i64.store32":         0x3e,
		"current_memory":      0x3f,
		"grow_memory":         0x40,
		"i32.const":           0x41,
		"i64.const":           0x42,
		"f32.const":           0x43,
		"f64.const":           0x44,
		"i32.eqz":             0x45,
		"i32.eq":              0x46,
		"i32.ne":              0x47,
		"i32.lt_s":            0x48,
		"i32.lt_u":            0x49,
		"i32.gt_s":            0x4a,
		"i32.gt_u":            0x4b,
		"i32.le_s":            0x4c,
		"i32.le_u":            0x4d,
		"i32.ge_s":            0x4e,
		"i32.ge_u":            0x4f,
		"i64.eqz":             0x50,
		"i64.eq":              0x51,
		"i64.ne":              0x52,
		"i64.lt_s":            0x53,
		"i64.lt_u":            0x54,
		"i64.gt_s":            0x55,
		"i64.gt_u":            0x56,
		"i64.le_s":            0x57,
		"i64.le_u":            0x58,
		"i64.ge_s":            0x59,
		"i64.ge_u":            0x5a,
		"f32.eq":              0x5b,
		"f32.ne":              0x5c,
		"f32.lt":              0x5d,
		"f32.gt":              0x5e,
		"f32.le":              0x5f,
		"f32.ge":              0x60,
		"f64.eq":              0x61,
		"f64.ne":              0x62,
		"f64.lt":              0x63,
		"f64.gt":              0x64,
		"f64.le":              0x65,
		"f64.ge":              0x66,
		"i32.clz":             0x67,
		"i32.ctz":             0x68,
		"i32.popcnt":          0x69,
		"i32.add":             0x6a,
		"i32.sub":             0x6b,
		"i32.mul":             0x6c,
		"i32.div_s":           0x6d,
		"i32.div_u":           0x6e,
		"i32.rem_s":           0x6f,
		"i32.rem_u":           0x70,
		"i32.and":             0x71,
		"i32.or":              0x72,
		"i32.xor":             0x73,
		"i32.shl":             0x74,
		"i32.shr_s":           0x75,
		"i32.shr_u":           0x76,
		"i32.rotl":            0x77,
		"i32.rotr":            0x78,
		"i64.clz":             0x79,
		"i64.ctz":             0x7a,
		"i64.popcnt":          0x7b,
		"i64.add":             0x7c,
		"i64.sub":             0x7d,
		"i64.mul":             0x7e,
		"i64.div_s":           0x7f,
		"i64.div_u":           0x80,
		"i64.rem_s":           0x81,
		"i64.rem_u":           0x82,
		"i64.and":             0x83,
		"i64.or":              0x84,
		"i64.xor":             0x85,
		"i64.shl":             0x86,
		"i64.shr_s":           0x87,
		"i64.shr_u":           0x88,
		"i64.rotl":            0x89,
		"i64.rotr":            0x8a,
		"f32.abs":             0x8b,
		"f32.neg":             0x8c,
		"f32.ceil":            0x8d,
		"f32.floor":           0x8e,
		"f32.trunc":           0x8f,
		"f32.nearest":         0x90,
		"f32.sqrt":            0x91,
		"f32.add":             0x92,
		"f32.sub":             0x93,
		"f32.mul":             0x94,
		"f32.div":             0x95,
		"f32.min":             0x96,
		"f32.max":             0x97,
		"f32.copysign":        0x98,
		"f64.abs":             0x99,
		"f64.neg":             0x9a,
		"f64.ceil":            0x9b,
		"f64.floor":           0x9c,
		"f64.trunc":           0x9d,
		"f64.nearest":         0x9e,
		"f64.sqrt":            0x9f,
		"f64.add":             0xa0,
		"f64.sub":             0xa1,
		"f64.mul":             0xa2,
		"f64.div":             0xa3,
		"f64.min":             0xa4,
		"f64.max":             0xa5,
		"f64.copysign":        0xa6,
		"i32.wrap/i64":        0xa7,
		"i32.trunc_s/f32":     0xa8,
		"i32.trunc_u/f32":     0xa9,
		"i32.trunc_s/f64":     0xaa,
		"i32.trunc_u/f64":     0xab,
		"i64.extend_s/i32":    0xac,
		"i64.extend_u/i32":    0xad,
		"i64.trunc_s/f32":     0xae,
		"i64.trunc_u/f32":     0xaf,
		"i64.trunc_s/f64":     0xb0,
		"i64.trunc_u/f64":     0xb1,
		"f32.convert_s/i32":   0xb2,
		"f32.convert_u/i32":   0xb3,
		"f32.convert_s/i64":   0xb4,
		"f32.convert_u/i64":   0xb5,
		"f32.demote/f64":      0xb6,
		"f64.convert_s/i32":   0xb7,
		"f64.convert_u/i32":   0xb8,
		"f64.convert_s/i64":   0xb9,
		"f64.convert_u/i64":   0xba,
		"f64.promote/f32":     0xbb,
		"i32.reinterpret/f32": 0xbc,
		"i64.reinterpret/f64": 0xbd,
		"f32.reinterpret/i32": 0xbe,
		"f64.reinterpret/i64": 0xbf,
	}

	typeGen  = reflect.ValueOf(typeGenerators{})
	immeGen  = reflect.ValueOf(immediataryGenerators{})
	entryGen = reflect.ValueOf(entryGenerators{})
)

type typeGenerators struct{}

func (t typeGenerators) Function(num uint64, stream *Stream) {
	EncodeULEB128(num, stream)
}

func (typeGenerators) Table(table Table, stream *Stream) {
	stream.Write([]byte{J2W_LANGUAGE_TYPES[table.ElementType]})
	typeGenerators{}.Memory(table.Limits, stream)
}

// Generates a [`global_type`](https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#global_type)
func (typeGenerators) Global(global Global, stream *Stream) {
	stream.Write([]byte{J2W_LANGUAGE_TYPES[global.ContentType]})
	stream.Write([]byte{global.Mutability})
}

// Generates a [resizable_limits](https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#resizable_limits)
func (typeGenerators) Memory(mem MemLimits, stream *Stream) {
	if mem.Maximum != nil {
		EncodeULEB128(1, stream)
		EncodeULEB128(mem.Intial, stream)
		EncodeULEB128(mem.Maximum.(uint64), stream)
	} else {
		EncodeULEB128(0, stream)
		EncodeULEB128(mem.Intial, stream)
	}

}

func (typeGenerators) InitExpr(op OP, stream *Stream) {
	GenerateOP(op, stream)
	GenerateOP(OP{
		Name: "end",
		Type: "void",
	}, stream)
}

type immediataryGenerators struct{}

func (immediataryGenerators) Varuint1(j int8, stream *Stream) *Stream {
	stream.Write([]byte{byte(j)})
	return stream
}

func (immediataryGenerators) Varuint32(j uint32, stream *Stream) *Stream {
	EncodeULEB128(uint64(j), stream)
	return stream
}

func (immediataryGenerators) Varint32(j int32, stream *Stream) *Stream {
	EncodeSLEB128(int64(j), stream)
	return stream
}

func (immediataryGenerators) Varint64(j int64, stream *Stream) *Stream {
	EncodeSLEB128(j, stream)
	return stream
}

func (immediataryGenerators) Uint32(j []byte, stream *Stream) *Stream {
	stream.Write(j)
	return stream
}

func (immediataryGenerators) Uint64(j []byte, stream *Stream) *Stream {
	stream.Write(j)
	return stream
}

func (immediataryGenerators) BlockType(j string, stream *Stream) *Stream {
	stream.Write([]byte{J2W_LANGUAGE_TYPES[j]})
	return stream
}

func (immediataryGenerators) BrTable(j JSON, stream *Stream) *Stream {
	targets := j["targets"].([]uint64)
	EncodeULEB128(uint64(len(targets)), stream)

	for _, target := range targets {
		EncodeULEB128(target, stream)
	}
	EncodeULEB128(j["default_target"].(uint64), stream)
	return stream
}

func (immediataryGenerators) CallIndirect(j JSON, stream *Stream) *Stream {
	index := j["index"]
	reserved := j["reserved"].(byte)
	EncodeULEB128(index.(uint64), stream)
	stream.Write([]byte{reserved})
	return stream
}

func (immediataryGenerators) MemoryImmediate(j JSON, stream *Stream) *Stream {
	EncodeULEB128(j["flags"].(uint64), stream)
	EncodeULEB128(j["offset"].(uint64), stream)
	return stream
}

type entryGenerators struct{}

func (entryGenerators) Type(entry TypeEntry, stream *Stream) []byte {
	// a single type entry binary encoded
	stream.WriteByte(J2W_LANGUAGE_TYPES[entry.Form])

	// number of parameters
	paramsLen := len(entry.Params)
	EncodeULEB128(uint64(paramsLen), stream)
	if paramsLen != 0 {
		paramsType := make([]byte, 0, paramsLen)
		for _, typ := range entry.Params {
			paramsType = append(paramsType, J2W_LANGUAGE_TYPES[typ])
		}
		stream.Write(paramsType)
	}

	// number of return types
	if entry.ReturnType != "" {
		stream.Write([]byte{1})
		stream.Write([]byte{J2W_LANGUAGE_TYPES[entry.ReturnType]})
	} else {
		stream.Write([]byte{0})
	}

	return stream.Bytes()
}

func (entryGenerators) Import(entry ImportEntry, stream *Stream) {
	// write the module string
	moduleStr := entry.ModuleStr
	EncodeULEB128(uint64(len(moduleStr)), stream)
	stream.Write([]byte(moduleStr))
	// write the field string
	fieldStr := entry.FieldStr
	EncodeULEB128(uint64(len(fieldStr)), stream)
	stream.Write([]byte(fieldStr))

	stream.Write([]byte{J2W_EXTERNAL_KIND[entry.Kind]})

	typeGen.MethodByName(Ucfirst(entry.Kind)).Call([]reflect.Value{reflect.ValueOf(entry.Type), reflect.ValueOf(stream)})
}

func (entryGenerators) Function(entry uint64, stream *Stream) []byte {
	EncodeULEB128(entry, stream)
	return stream.Bytes()
}

func (entryGenerators) Table(j Table, stream *Stream) {
	typeGenerators{}.Table(j, stream)
}

func (entryGenerators) Global(entry GlobalEntry, stream *Stream) *Stream {
	typeGen := typeGenerators{}
	typeGen.Global(entry.Type, stream)
	typeGen.InitExpr(entry.Init, stream)
	return stream
}

func (entryGenerators) Memory(entry MemLimits, stream *Stream) {
	typeGenerators{}.Memory(entry, stream)
}

func (entryGenerators) Export(entry ExportEntry, stream *Stream) *Stream {
	EncodeULEB128(uint64(len(entry.FieldStr)), stream)
	stream.Write([]byte(entry.FieldStr))
	stream.Write([]byte{J2W_EXTERNAL_KIND[entry.Kind]})
	EncodeULEB128(uint64(entry.Index), stream)
	return stream
}

func (entryGenerators) Element(entry ElementEntry, stream *Stream) *Stream {
	EncodeULEB128(uint64(entry.Index), stream)
	typeGenerators{}.InitExpr(entry.Offset, stream)
	EncodeULEB128(uint64(len(entry.Elements)), stream)
	for _, elem := range entry.Elements {
		EncodeULEB128(elem, stream)
	}
	return stream
}

func (entryGenerators) Code(entry CodeBody, stream *Stream) *Stream {
	codeStream := NewStream(nil)
	// write the locals
	EncodeULEB128(uint64(len(entry.Locals)), codeStream)
	for _, local := range entry.Locals {
		EncodeULEB128(uint64(local.Count), codeStream)
		codeStream.Write([]byte{J2W_LANGUAGE_TYPES[local.Type]})
	}

	// write opcode
	for _, op := range entry.Code {
		GenerateOP(op, codeStream)
	}

	EncodeULEB128(uint64(codeStream.bytesWrote), stream)
	stream.Write(codeStream.Bytes())
	return stream
}

func (entryGenerators) Data(entry DataSegment, stream *Stream) *Stream {
	EncodeULEB128(uint64(entry.Index), stream)
	typeGenerators{}.InitExpr(entry.Offset, stream)
	EncodeULEB128(uint64(len(entry.Data)), stream)
	stream.Write(entry.Data)
	return stream
}

// Json2Wasm converts a JSON array to wasm binary.
func Json2Wasm(j []JSON) []byte {
	stream := NewStream(nil)
	preamble := j[0]
	GeneratePreramble(preamble, stream)
	rest := j[1:]
	for _, item := range rest {
		GenerateSection(item, stream)
	}

	return stream.Bytes()
}

func GeneratePreramble(j JSON, stream *Stream) *Stream {
	if stream == nil {
		stream = NewStream(nil)
	}

	stream.Write(Interface2Bytes(j["magic"]))

	stream.Write(Interface2Bytes(j["version"]))

	return stream
}

func GenerateOP(op OP, stream *Stream) *Stream {
	if stream == nil {
		stream = NewStream(nil)
	}

	name := op.Name
	if op.ReturnType != "" {
		name = op.ReturnType + "." + name
	}

	stream.Write([]byte{J2W_OPCODES[name]})

	immediateKey := op.Name
	if immediateKey == "const" {
		immediateKey = op.ReturnType
	}
	immediates, exist := OP_IMMEDIATES[immediateKey]
	if exist {
		//fmt.Printf("imm %s %v\n", immediates, op.Immediates)
		immeGen.MethodByName(Ucfirst(immediates)).Call([]reflect.Value{
			reflect.ValueOf(op.Immediates),
			reflect.ValueOf(stream),
		})
	}
	return stream
}

func GenerateSection(j JSON, stream *Stream) *Stream {
	if stream == nil {
		stream = NewStream(nil)
	}

	var name string
	nameinterf, exist := j["name"]
	if exist {
		name = nameinterf.(string)
	}
	payload := NewStream(nil)
	stream.Write([]byte{J2W_SECTION_IDS[name]})

	if name == "custom" {
		sectionName := j["section_name"].(string)
		EncodeULEB128(uint64(len(sectionName)), payload)
		payload.Write([]byte(sectionName))
		payload.Write([]byte(j["payload"].(string)))
	} else if name == "start" {
		EncodeULEB128(uint64(j["index"].(uint32)), payload)
	} else {
		entries, exist := j["entries"]
		if exist {
			entries := reflect.ValueOf(entries)
			EncodeULEB128(uint64(entries.Len()), payload)
			//fmt.Printf("Gen %v\n", name)
			method := entryGen.MethodByName(Ucfirst(name))
			for i := 0; i < entries.Len(); i++ {
				entry := entries.Index(i)
				method.Call([]reflect.Value{entry, reflect.ValueOf(payload)})
			}
		}
	}

	// write the size of the payload.
	EncodeULEB128(uint64(payload.bytesWrote), stream)
	stream.Write(payload.Bytes())
	return stream
}
