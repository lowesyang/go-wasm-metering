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

var sevenbits = [...]byte{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
	0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f,
	0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f,
	0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f,
	0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f,
}

// EncodeULEB128 appends v to b using unsigned LEB128 encoding.
func EncodeULEB128(v uint64, stream *Stream) (out []byte) {
	for {
		c := uint8(v & 0x7f)
		v >>= 7
		if v != 0 {
			c |= 0x80
		}
		out = append(out, c)
		if c&0x80 == 0 {
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

func ReadWasmFromFile(path string) JSON {
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

func Ucfirst(str string) string {
	if z := rune(str[0]); unicode.IsLower(z) {
		return string(unicode.ToUpper(z)) + str[1:]
	} else {
		return str
	}
}
