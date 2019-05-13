package go_wasm_metering

type JSON map[string]interface{}

type SectionHeader struct {
	Id   byte   `json:"id"`
	Name string `json:"name"`
	Size uint64 `json:"size"`
}

type OP struct {
	Name       string      `json:"name"`
	ReturnType string      `json:"return_type"`
	Immediates interface{} `json:"immediates"`
}

type Table struct {
	ElementType string    `json:"element_type"`
	Limits      MemLimits `json:"limits"`
}

type MemLimits struct {
	Flags   uint64 `json:"flags"`
	Intial  uint64 `json:"intial"`
	Maximum uint64 `json:"maximum"`
}

type Global struct {
	ContentType string `json:"content_type"`
	Mutability  byte   `json:"mutability"`
}

// Section data structures.
type CustomSec struct {
	Name        string `json:"name"`
	SectionName string `json:"section_name"`
	Payload     []byte `json:"payload"`
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
	ModuleStr string `json:"module_str"`
	FieldStr  string `json:"field_str"`
	Kind      string `json:"kind"`
	Type      interface{}
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
