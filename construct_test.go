package construct

import (
	"bytes"
	"testing"
)

func TestStruct_Parse(t *testing.T) {
	s := Struct{
		Byte{},
		Int32be{},
		String{Length: 5},
	}
	data := []byte{0xFF, 0x00, 0x00, 0x00, 0x64, 'T', 'e', 's', 't', 0x00}
	values, err := s.Parse(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 3 || values[0] != byte(0xFF) || values[1] != int32(100) || values[2] != "Test" {
		t.Errorf("unexpected values: %v", values)
	}
}

func TestStruct_Build(t *testing.T) {
	s := Struct{
		Byte{},
		Int32be{},
		String{Length: 5},
	}
	values := []any{byte(0xFF), int32(100), "Test"}
	var buf bytes.Buffer
	err := s.Build(&buf, values)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte{0xFF, 0x00, 0x00, 0x00, 0x64, 'T', 'e', 's', 't', 0x00}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Errorf("unexpected output: %x", buf.Bytes())
	}
}
