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
	Id   byte   `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Size uint64 `json:"size,omitempty"`
}

type OP struct {
	Name       string      `json:"name,omitempty"`
	ReturnType string      `json:"return_type,omitempty"`
	Type       string      `json:"type,omitempty"`
	Immediates interface{} `json:"immediates,omitempty"`
}

type Table struct {
	ElementType string    `json:"element_type,omitempty"`
	Limits      MemLimits `json:"limits,omitempty"`
}

type MemLimits struct {
	Flags   uint64      `json:"flags,omitempty"`
	Intial  uint64      `json:"intial,omitempty"`
	Maximum interface{} `json:"maximum,omitempty"` // to distinguish the field is nil or uint64(0)
}

type Global struct {
	ContentType string `json:"content_type,omitempty"`
	Mutability  byte   `json:"mutability,omitempty"`
}

// Section data structures.
type CustomSec struct {
	Name        string `json:"name,omitempty"`
	SectionName string `json:"section_name,omitempty"`
	Payload     string `json:"payload,omitempty"`
}

type TypeEntry struct {
	Form       string   `json:"form,omitempty"`
	Params     []string `json:"params"`
	ReturnType string   `json:"return_type,omitempty"`
}

type TypeSec struct {
	Name    string      `json:"name,omitempty"`
	Entries []TypeEntry `json:"entries"`
}

type ImportEntry struct {
	ModuleStr string      `json:"module_str,omitempty"`
	FieldStr  string      `json:"field_str,omitempty"`
	Kind      string      `json:"kind,omitempty"`
	Type      interface{} `json:"type,omitempty"`
}

type ImportSec struct {
	Name    string        `json:"name,omitempty"`
	Entries []ImportEntry `json:"entries"`
}

type FuncSec struct {
	Name    string   `json:"name,omitempty"`
	Entries []uint64 `json:"entries"`
}

type TableSec struct {
	Name    string  `json:"name,omitempty"`
	Entries []Table `json:"entries"`
}

type MemSec struct {
	Name    string      `json:"name,omitempty"`
	Entries []MemLimits `json:"entries"`
}

type GlobalEntry struct {
	Type Global `json:"type,omitempty"`
	Init OP     `json:"init,omitempty"`
}

type GlobalSec struct {
	Name    string        `json:"name,omitempty"`
	Entries []GlobalEntry `json:"entries"`
}

type ExportEntry struct {
	FieldStr string `json:"field_str,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Index    uint32 `json:"index,omitempty"`
}

type ExportSec struct {
	Name    string        `json:"name,omitempty"`
	Entries []ExportEntry `json:"entries"`
}

type StartSec struct {
	Name  string `json:"name,omitempty"`
	Index uint32 `json:"index,omitempty"`
}

type ElementEntry struct {
	Index    uint32   `json:"index,omitempty"`
	Offset   OP       `json:"offset,omitempty"`
	Elements []uint64 `json:"elements"`
}

type ElementSec struct {
	Name    string         `json:"name,omitempty"`
	Entries []ElementEntry `json:"entries"`
}

type LocalEntry struct {
	Count uint32 `json:"count,omitempty"`
	Type  string `json:"type,omitempty"`
}

type CodeBody struct {
	Locals []LocalEntry `json:"locals"`
	Code   []OP         `json:"code"`
}

type CodeSec struct {
	Name    string     `json:"name,omitempty"`
	Entries []CodeBody `json:"entries"`
}

type DataSegment struct {
	Index  uint32 `json:"index,omitempty"`
	Offset OP     `json:"offset,omitempty"`
	Data   []byte `json:"data"`
}

type DataSec struct {
	Name    string        `json:"name,omitempty"`
	Entries []DataSegment `json:"entries"`
}
