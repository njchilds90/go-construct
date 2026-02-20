# go-construct

A simple Golang port of Python's [construct library](https://construct.readthedocs.io/) for declarative binary data parsing and building.

## Features
- Declarative structures using composable fields
- Parse binary data into Go values
- Build Go values into binary data
- No external dependencies — uses only the standard library
- Easy to extend with new field types via the `Field` interface
- Designed to be readable by both humans and AI agents (doc comments + examples)

## Installation

```bash
go get github.com/njchilds90/go-construct
Quick Example
Gopackage main

import (
	"bytes"
	"fmt"
	"github.com/njchilds90/go-construct"
)

func main() {
	// Define a simple binary structure (like a tiny header + message)
	header := construct.Struct{
		construct.Byte{},           // version (1 byte)
		construct.Int32be{},        // length (4 bytes, big-endian)
		construct.String{Length: 8}, // fixed 8-byte name field (null-padded)
	}

	// Example binary data
	data := []byte{
		0x01,                               // version = 1
		0x00, 0x00, 0x00, 0x0A,             // length = 10
		'H', 'e', 'l', 'l', 'o', 0x00, 0x00, 0x00, // "Hello" padded to 8 bytes
	}

	// Parse
	values, err := header.Parse(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	version := values[0].(byte)
	length := values[1].(int32)
	name := values[2].(string)

	fmt.Printf("Version: %d\nLength: %d\nName: %q\n", version, length, name)
	// Output:
	// Version: 1
	// Length: 10
	// Name: "Hello"

	// Build (reverse)
	var buf bytes.Buffer
	err = header.Build(&buf, []any{
		byte(2),           // new version
		int32(7),          // new length
		"World!!",         // will be padded/truncated to 8 bytes
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Built binary: % x\n", buf.Bytes())
	// Example output: 02 00 00 00 07 57 6f 72 6c 64 21 21 00 00
}
Supported Field Types (v0.1.0)

































TypeDescriptionGo TypeEndiannessNotesByteSingle bytebyte—Int32be32-bit signed integerint32Big-endianAdd more (Int16le, Uint64be, etc) laterStringFixed-length string (null-padded)string—Trims nulls on parse, pads on build
More field types (arrays, structs nesting, bitfields, enums, etc.) can be added in future versions.
How to Extend
Implement the Field interface:
Gotype Field interface {
	Parse(r io.Reader) (any, error)
	Build(w io.Writer, v any) error
}
Then add your new type to any Struct.
