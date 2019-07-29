package toolkit

import (
	"bytes"
	"encoding/json"
	"os"
	"unicode"
)

type Stream struct {
	length     int
	bytesRead  int
	bytesWrote int
	buffer     *bytes.Buffer
}

func NewStream(buf []byte) *Stream {
	return &Stream{
		buffer: bytes.NewBuffer(buf),
		length: len(buf),
	}
}

// ReadBytes reads and returns the next byte from the buffer.
func (s *Stream) ReadByte() byte {
	b, _ := s.buffer.ReadByte()
	s.bytesRead += 1
	return b
}

// Read returns a slice containing the next n bytes from the buffer.
func (s *Stream) Read(n int) []byte {
	s.bytesRead += n
	return s.buffer.Next(n)
}

// Len returns the number of bytes of the unread portion of the buffer;
func (s *Stream) Len() int {
	return s.buffer.Len()
}

func (s *Stream) Bytes() []byte {
	return s.buffer.Bytes()
}

func (s *Stream) String() string {
	return s.buffer.String()
}

func (s *Stream) Write(buf []byte) (n int, err error) {
	n, err = s.buffer.Write(buf)
	s.bytesWrote += n
	return
}

func (s *Stream) WriteByte(c byte) error {
	err := s.buffer.WriteByte(c)
	if err != nil {
		return err
	}
	s.bytesWrote += 1
	return nil
}

// EncodeULEB128 appends v to b using unsigned LEB128 encoding.
func EncodeULEB128(v uint64, stream *Stream) (out []byte) {
	for {
		c := uint8(v & 0x7f)
		v >>= 7
		if v > 0 {
			c |= 0x80
		}
		out = append(out, c)
		if v == 0 {
			break
		}
	}
	stream.Write(out)
	return
}

// EncodeSLEB128 appends v to b using signed LEB128 encoding.
func EncodeSLEB128(v int64, stream *Stream) (out []byte) {
	for {
		c := uint8(v & 0x7f)
		s := uint8(v & 0x40)
		v >>= 7

		if (v != -1 || s == 0) && (v != 0 || s != 0) {
			c |= 0x80
		}

		out = append(out, c)

		if c&0x80 == 0 {
			break
		}
	}

	stream.Write(out)
	return out
}

// DecodeULEB128 decodes bytes from stream with unsigned LEB128 encoding.
func DecodeULEB128(stream *Stream) (u uint64) {
	var shift uint
	for {
		b := stream.ReadByte()
		u |= uint64(b&0x7f) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}

	return
}

// DecodeSLEB128 decodes bytes from stream with signed LEB128 encoding.
func DecodeSLEB128(stream *Stream) (s int64) {
	var shift uint
	for {
		b := stream.ReadByte()
		s |= int64(b&0x7f) << shift
		shift += 7
		if b&0x80 == 0 {
			// If it's signed
			if b&0x40 != 0 {
				s |= ^0 << shift
			}
			break
		}
	}

	return
}

func ReadFromFile(path string) JSON {
	obj := make(JSON)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	jsonParser := json.NewDecoder(file)
	if err := jsonParser.Decode(&obj); err != nil {
		panic(err)
	}
	return obj
}

// Camel-Case to underline
func Lcfirst(str string) string {
	//if z := rune(str[0]); unicode.IsLower(z) {
	//	return string(unicode.ToUpper(z)) + str[1:]
	//} else {
	//	return str
	//}
	var newStr string
	for i, b := range str {
		if unicode.IsUpper(b) {
			cm := string(unicode.ToLower(b))
			if i != 0 {
				cm = "_" + cm
			}
			newStr += cm
		} else {
			newStr += string(b)
		}
	}
	return newStr
}

// underline to Camel-Case
func Ucfirst(str string) string {
	//if z := rune(str[0]); unicode.IsUpper(z) {
	//	return string(unicode.ToLower(z)) + str[1:]
	//} else {
	//	return str
	//}
	var newStr string
	for i := 0; i < len(str); i++ {
		b := str[i]
		if unicode.IsLower(rune(b)) && i == 0 {
			newStr += string(unicode.ToUpper(rune(b)))
		} else if b == '_' {
			n := str[i+1]
			if unicode.IsLower(rune(n)) {
				newStr += string(unicode.ToUpper(rune(n)))
				i += 1
			}
		} else {
			newStr += string(b)
		}
	}
	return newStr
}

func Interface2Bytes(arr interface{}) (out []byte) {
	switch v := arr.(type) {
	case []interface{}:
		for _, b := range v {
			out = append(out, byte(b.(float64)))
		}
	case []byte:
		out = v
	case string:
		out = []byte(v)
	}
	return
}
