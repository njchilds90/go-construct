// Package construct provides a simple declarative way to parse and build binary data structures,
// inspired by Python's construct library. This is a basic implementation supporting common field types.
// Usable for users (easy API with examples) and AI agents (structured with interfaces for extension).

package construct

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Field is an interface for individual fields in a binary structure.
// Implement Parse to read from input, and Build to write to output.
type Field interface {
	Parse(r io.Reader) (any, error)
	Build(w io.Writer, v any) error
}

// Struct is a composable binary structure made of multiple Fields.
// It parses/builds fields in sequence.
type Struct []Field

// Parse reads the entire structure from the reader and returns a slice of values.
func (s Struct) Parse(r io.Reader) ([]any, error) {
	values := make([]any, len(s))
	for i, field := range s {
		v, err := field.Parse(r)
		if err != nil {
			return nil, err
		}
		values[i] = v
	}
	return values, nil
}

// Build writes the structure to the writer using the provided slice of values.
func (s Struct) Build(w io.Writer, values []any) error {
	if len(values) != len(s) {
		return errors.New("mismatch between fields and values")
	}
	for i, field := range s {
		if err := field.Build(w, values[i]); error != nil {
			return err
		}
	}
	return nil
}

// Byte is a field for a single byte.
type Byte struct{}

// Parse reads one byte.
func (Byte) Parse(r io.Reader) (any, error) {
	var b byte
	err := binary.Read(r, binary.BigEndian, &b)
	return b, err
}

// Build writes one byte.
func (Byte) Build(w io.Writer, v any) error {
	b, ok := v.(byte)
	if !ok {
		return errors.New("expected byte")
	}
	return binary.Write(w, binary.BigEndian, b)
}

// Int32be is a field for a 32-bit integer (big-endian).
type Int32be struct{}

// Parse reads a big-endian int32.
func (Int32be) Parse(r io.Reader) (any, error) {
	var i int32
	err := binary.Read(r, binary.BigEndian, &i)
	return i, err
}

// Build writes a big-endian int32.
func (Int32be) Build(w io.Writer, v any) error {
	i, ok := v.(int32)
	if !ok {
		return errors.New("expected int32")
	}
	return binary.Write(w, binary.BigEndian, i)
}

// String is a field for a fixed-length string (null-padded if shorter).
type String struct {
	Length int
}

// Parse reads a fixed-length string, trimming nulls.
func (s String) Parse(r io.Reader) (any, error) {
	buf := make([]byte, s.Length)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	return string(bytes.TrimRight(buf, "\x00")), nil
}

// Build writes a fixed-length string, padding with nulls.
func (s String) Build(w io.Writer, v any) error {
	str, ok := v.(string)
	if !ok {
		return errors.New("expected string")
	}
	buf := make([]byte, s.Length)
	copy(buf, str)
	_, err := w.Write(buf)
	return err
}

// Example usage (for users and AI: create a Struct and parse/build binary data):
// myStruct := construct.Struct{
// 	construct.Byte{},
// 	construct.Int32be{},
// 	construct.String{Length: 10},
// }
// data := []byte{0x01, 0x00, 0x00, 0x00, 0x0A, 'H', 'e', 'l', 'l', 'o', 0x00, 0x00, 0x00, 0x00, 0x00}
// values, err := myStruct.Parse(bytes.NewReader(data))
// // values = [1, 10, "Hello"]
