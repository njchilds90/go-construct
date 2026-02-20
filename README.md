# go-construct

> A simple, idiomatic Go port inspired by Python’s `construct` library for declarative binary parsing and building.

[![Go Reference](https://pkg.go.dev/badge/github.com/njchilds90/go-construct.svg)](https://pkg.go.dev/github.com/njchilds90/go-construct)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

---

## ✨ Overview

`go-construct` provides a lightweight, dependency-free way to describe binary data structures declaratively in Go and then:

- ✅ Parse binary data into Go values
- ✅ Build binary data from Go values
- ✅ Compose reusable binary layouts
- ✅ Extend with custom field types

The goal is to keep things **simple, readable, and idiomatic** — useful for both humans and AI agents.

---

## 🚀 Features

- 🧱 Declarative struct-style field composition
- 🔄 Bidirectional parsing and building
- 📦 Zero external dependencies (standard library only)
- 🛠 Easy extensibility via a small `Field` interface
- 🧼 Minimal API surface

---

## 📦 Installation

```bash
go get github.com/njchilds90/go-construct
```

Or with Go modules:

```bash
go install github.com/njchilds90/go-construct@latest
```

---

## 🧪 Quick Example

Define a binary structure, parse raw bytes, and rebuild them:

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/njchilds90/go-construct"
)

func main() {
	header := construct.Struct{
		construct.Byte{},            // 1 byte version
		construct.Int32be{},         // 4 byte big-endian integer
		construct.String{Length: 8}, // fixed 8-byte string (null-padded)
	}

	data := []byte{
		0x01,                    // version = 1
		0x00, 0x00, 0x00, 0x0A, // length = 10
		'H', 'e', 'l', 'l', 'o', 0x00, 0x00, 0x00,
	}

	values, err := header.Parse(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	version := values[0].(byte)
	length := values[1].(int32)
	name := values[2].(string)

	fmt.Printf("Version: %d\nLength: %d\nName: %q\n", version, length, name)

	var buf bytes.Buffer
	err = header.Build(&buf, []any{
		byte(2),
		int32(7),
		"World!!",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Built binary: % x\n", buf.Bytes())
}
```

---

## 🧩 Supported Field Types

| Type        | Description                      |
|-------------|----------------------------------|
| `Byte`      | Single raw byte                  |
| `Int32be`   | 32-bit signed big-endian integer |
| `String`    | Fixed-length string (null padded)|

Future roadmap ideas:

- `Int16le`, `Int64be`
- Arrays / slices
- Nested structs
- Enums
- Conditional fields
- Length-prefixed fields

---

## 🛠 Extending the Library

To create a custom field type, implement the `Field` interface:

```go
type Field interface {
	Parse(r io.Reader) (any, error)
	Build(w io.Writer, v any) error
}
```

Once implemented, your type can be included in any `Struct`.

---

## 🧪 Testing

Run tests with:

```bash
go test ./...
```

Keep test cases small, readable, and binary-focused.

---

## 🤝 Contributing

Pull requests and issues are welcome.

If contributing:

- Keep the API minimal
- Maintain zero dependencies
- Add tests for new field types
- Prefer clarity over cleverness

---

## 📄 License

MIT License — see the `LICENSE` file for details.

---

## 🎯 Project Philosophy

`go-construct` is designed to be:

- Simple enough to understand in minutes
- Powerful enough for real binary protocols
- Clean enough for AI-assisted code generation
- Stable and predictable

No magic. No reflection-heavy abstractions. Just composable binary primitives.

---

Maintained by: **Nicholas Childs**  
GitHub: https://github.com/njchilds90
