// Copyright 2016 The wasm Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

import (
	"fmt"
	"io"
)

type ModuleHeader struct {
	Magic   [4]byte // wasm magic number (0x6d736100, ie: "\0asm")
	Version uint32  // version number
}

func (hdr ModuleHeader) String() string {
	return fmt.Sprintf("ModuleHeader{Magic=%q Version=0x%x}", hdr.Magic, hdr.Version)
}

type Section struct {
	ID      varuint7  // section code
	Len     varuint32 // size of this section in bytes
	NameLen varuint32 // length of the section name in bytes, present if ID==0
	Name    []byte    // section name string, present if ID==0
	Data    []byte    // content of this section, of length Len-len(Name)-len(NameLen)
}

func (s Section) String() string {
	return fmt.Sprintf(
		"Section{ID: %d, Len: %d, Name: %q, Data: %d}",
		s.ID, s.Len, string(s.Name), len(s.Data),
	)
}

func (s *Section) read(r io.Reader) error {
	_, err := s.ID.read(r)
	if err != nil {
		return err
	}

	_, err = s.Len.read(r)
	if err != nil {
		return err
	}

	nlen := 0
	if s.ID == 0 {
		nlen, err = s.NameLen.read(r)
		if err != nil {
			return err
		}

		s.Name = make([]byte, int(s.NameLen))
		_, err = r.Read(s.Name)
		if err != nil {
			return err
		}
	}

	s.Data = make([]byte, int(s.Len)-len(s.Name)-nlen)
	_, err = r.Read(s.Data)
	return err
}

const (
	TypeID     varuint7 = 1  // Function signature declarations
	ImportID            = 2  // Import declarations
	FunctionID          = 3  // Function declarations
	TableID             = 4  // Indirect function table and other tables
	MemoryID            = 5  // Memory attributes
	GlobalID            = 6  // Global declarations
	ExportID            = 7  // Exports
	StartID             = 8  // Start function declaration
	ElementID           = 9  // Elements section
	CodeID              = 10 // Function bodies (code)
	DataID              = 11 // Data segments
)

type TypeSection struct {
	Count   varuint32  // count of type entries to follow
	Entries []FuncType // repeated type entries
}

type ImportSection struct {
	Count   varuint32 // count of import entries to follow
	Entries []ImportEntry
}

type ImportEntry struct {
	ModuleLen varuint32    // module string length
	ModuleStr []byte       // module string of ModuleLen bytes
	FieldLen  varuint32    // field name length
	FieldStr  []byte       // field name string of FieldLen bytes
	Kind      ExternalKind // the kind of definition being imported

	FunctionType varuint32  // type index of the function signature (if Kind is Function)
	TableType    TableType  // type of imported table (if Kind==Table)
	MemoryType   MemoryType // type of imported memory (if Kind==Memory)
	GlobalType   GlobalType // type of importer global
}

// FunctionSection declares the signature of all function in the module
type FunctionSection struct {
	Count varuint32   // count of signature indices to follow
	Types []varuint32 // sequence of indices into the type section
}

// TableSection encodes a table
type TableSection struct {
	Count   varuint32 // count of tables defined by the module
	Entries []TableType
}

// MemorySection encodes a memory
type MemorySection struct {
	Count   varuint32 // number of memories defined by the module
	Entries []MemoryType
}

// GlobalSection encodes the global section
type GlobalSection struct {
	Count   varuint32 // number of global variable entries
	Globals []GlobalVariable
}

// GlobalVariable represents a single global variable of a given type,
// mutability and with the given initializer.
type GlobalVariable struct {
	Type GlobalType // type of the variables
	Init InitExpr   // initial value of the global
}

// ExportSection encodes the export section
type ExportSection struct {
	Count   varuint32 // number of export entries
	Entries []ExportEntry
}

// ExportEntry represents an exported entity.
type ExportEntry struct {
	FieldLen varuint32    // field name string length
	FieldStr []byte       // field name string of FieldLen bytes
	Kind     ExternalKind // kind of definition being exported
	Index    varuint32    // index into the corresponding index space
}

// StartSection declares the start function
type StartSection struct {
	Index varuint32 // start function index
}

// ElementSection encodes the elements section
type ElementSection struct {
	Count   varuint32 // number of element segments
	Entries []ElemSegment
}

type ElemSegment struct {
	Index   varuint32   // the table index
	Offset  InitExpr    // an i32 initializer expression that computes the offset at which to place the elements
	NumElem varuint32   // number of elements
	Elems   []varuint32 // sequence of function indices
}

// CodeSection contains a body for every function in the module.
// The count of function declared in the function section and function bodies
// defined in this section must be the same and the i-th declaration corresponds
// to the i-th function body.
type CodeSection struct {
	Count  varuint32 // number of function bodies to follow
	Bodies []FunctionBody
}

// DataSection declares the initialized data that is loaded into linear memory
type DataSection struct {
	Count   varuint32
	Entries []DataSegment
}

type DataSegment struct {
	Index  varuint32 // the linear memory index
	Offset InitExpr  // an i32 initializer expression that computes the offset at which to place the data
	Size   varuint32 // size of Data in bytes
	Data   []byte
}

// NameSection describes user-defined sections
type NameSection struct {
	Count   varuint32
	Entries []FunctionNames
}

type FunctionNames struct {
	FunNameLen varuint32
	FunNameStr []byte
	LocalCount varuint32
	LocalNames []LocalName
}

type LocalName struct {
	LocalNameLen varuint32
	LocalNameStr []byte
}

type FunctionBody struct {
	BodySize   varuint32    // size of function body to follow, in bytes
	LocalCount varuint32    // number of local entries
	Locals     []LocalEntry // local variables
	Code       []byte       // bytecode of the function
	End        byte         // 0x0b, indicating the end of the body
}

type LocalEntry struct {
	Count varuint32 // number of local variables of the following type
	Type  ValueType // type of the variables
}

// FIXME(sbinet): opcodes
// https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#control-flow-operators-described-here
// https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#call-operators-described-here
// https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#parametric-operators-described-here
// https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#variable-access-described-here
// https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#memory-related-operators-described-here
// https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#constants-described-here
// https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#comparison-operators-described-here
// https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#numeric-operators-described-here
// https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#conversions-described-here
// https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#reinterpretations-described-here
