package toolkit

var OP_IMMEDIATES = map[string]string{
	"block":          "block_type",
	"loop":           "block_type",
	"if":             "block_type",
	"br":             "varuint32",
	"br_if":          "varuint32",
	"br_table":       "br_table",
	"call":           "varuint32",
	"call_indirect":  "call_indirect",
	"get_local":      "varuint32",
	"set_local":      "varuint32",
	"tee_local":      "varuint32",
	"get_global":     "varuint32",
	"set_global":     "varuint32",
	"load":           "memory_immediate",
	"load8_s":        "memory_immediate",
	"load8_u":        "memory_immediate",
	"load16_s":       "memory_immediate",
	"load16_u":       "memory_immediate",
	"load32_s":       "memory_immediate",
	"load32_u":       "memory_immediate",
	"store":          "memory_immediate",
	"store8":         "memory_immediate",
	"store16":        "memory_immediate",
	"store32":        "memory_immediate",
	"current_memory": "varuint1",
	"grow_memory":    "varuint1",
	"i32":            "varint32",
	"i64":            "varint64",
	"f32":            "uint32",
	"f64":            "uint64",
}

type JSON = map[string]interface{}

type SectionHeader struct {
	Id   byte   `json:"id"`
	Name string `json:"name"`
	Size uint64 `json:"size"`
}

type OP struct {
	Name       string      `json:"name"`
	ReturnType string      `json:"return_type"`
	Type       string      `json:"type"`
	Immediates interface{} `json:"immediates"`
}

type Table struct {
	ElementType string    `json:"element_type"`
	Limits      MemLimits `json:"limits"`
}

type MemLimits struct {
	Flags   uint64      `json:"flags"`
	Intial  uint64      `json:"intial"`
	Maximum interface{} `json:"maximum"` // to distinguish the field is nil or uint64(0)
}

type Global struct {
	ContentType string `json:"content_type"`
	Mutability  byte   `json:"mutability"`
}

// Section data structures.
type CustomSec struct {
	Name        string `json:"name"`
	SectionName string `json:"section_name"`
	Payload     string `json:"payload"`
}

type TypeEntry struct {
	Form       string   `json:"form"`
	Params     []string `json:"params"`
	ReturnType string   `json:"return_type"`
}

type TypeSec struct {
	Name    string      `json:"name"`
	Entries []TypeEntry `json:"entries"`
}

type ImportEntry struct {
	ModuleStr string      `json:"module_str"`
	FieldStr  string      `json:"field_str"`
	Kind      string      `json:"kind"`
	Type      interface{} `json:"type"`
}

type ImportSec struct {
	Name    string        `json:"name"`
	Entries []ImportEntry `json:"entries"`
}

type FuncSec struct {
	Name    string   `json:"name"`
	Entries []uint64 `json:"entries"`
}

type TableSec struct {
	Name    string  `json:"name"`
	Entries []Table `json:"entries"`
}

type MemSec struct {
	Name    string      `json:"name"`
	Entries []MemLimits `json:"entries"`
}

type GlobalEntry struct {
	Type Global `json:"type"`
	Init OP     `json:"init"`
}

type GlobalSec struct {
	Name    string        `json:"name"`
	Entries []GlobalEntry `json:"entries"`
}

type ExportEntry struct {
	FieldStr string `json:"field_str"`
	Kind     string `json:"kind"`
	Index    uint32 `json:"index"`
}

type ExportSec struct {
	Name    string        `json:"name"`
	Entries []ExportEntry `json:"entries"`
}

type StartSec struct {
	Name  string `json:"name"`
	Index uint32 `json:"index"`
}

type ElementEntry struct {
	Index    uint32   `json:"index"`
	Offset   OP       `json:"offset"`
	Elements []uint64 `json:"elements"`
}

type ElementSec struct {
	Name    string         `json:"name"`
	Entries []ElementEntry `json:"entries"`
}

type LocalEntry struct {
	Count uint32 `json:"count"`
	Type  string `json:"type"`
}

type CodeBody struct {
	Locals []LocalEntry `json:"locals"`
	Code   []OP         `json:"code"`
}

type CodeSec struct {
	Name    string     `json:"name"`
	Entries []CodeBody `json:"entries"`
}

type DataSegment struct {
	Index  uint32 `json:"index"`
	Offset OP     `json:"offset"`
	Data   []byte `json:"data"`
}

type DataSec struct {
	Name    string        `json:"name"`
	Entries []DataSegment `json:"entries"`
}
