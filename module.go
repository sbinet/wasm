// Copyright 2016 The wasm Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

import (
	"fmt"
	"io"
	"os"
)

// Module is a WebAssembly module.
type Module struct {
	Header   ModuleHeader
	Sections []Section
}

func Open(name string) (Module, error) {
	f, err := os.Open(name)
	if err != nil {
		return Module{}, err
	}
	defer f.Close()

	dec := decoder{r: f}
	return dec.readModule()
}

type ModuleHeader struct {
	Magic   [4]byte // wasm magic number (0x6d736100, ie: "\0asm")
	Version uint32  // version number
}

func (hdr ModuleHeader) String() string {
	return fmt.Sprintf("ModuleHeader{Magic=%q Version=0x%x}", hdr.Magic, hdr.Version)
}

// Section represents a section in a wasm module.
type Section interface {
	ID() SectionID
}

// SectionID represents the specific kind of section that a Section represents.
type SectionID byte

const (
	UnknownID  SectionID = 0  // User section ID
	TypeID               = 1  // Function signature declarations
	ImportID             = 2  // Import declarations
	FunctionID           = 3  // Function declarations
	TableID              = 4  // Indirect function table and other tables
	MemoryID             = 5  // Memory attributes
	GlobalID             = 6  // Global declarations
	ExportID             = 7  // Exports
	StartID              = 8  // Start function declaration
	ElementID            = 9  // Elements section
	CodeID               = 10 // Function bodies (code)
	DataID               = 11 // Data segments
)

func (TypeSection) ID() SectionID     { return TypeID }
func (ImportSection) ID() SectionID   { return ImportID }
func (FunctionSection) ID() SectionID { return FunctionID }
func (TableSection) ID() SectionID    { return TableID }
func (MemorySection) ID() SectionID   { return MemoryID }
func (GlobalSection) ID() SectionID   { return GlobalID }
func (ExportSection) ID() SectionID   { return ExportID }
func (StartSection) ID() SectionID    { return StartID }
func (ElementSection) ID() SectionID  { return ElementID }
func (CodeSection) ID() SectionID     { return CodeID }
func (DataSection) ID() SectionID     { return DataID }
func (NameSection) ID() SectionID     { return UnknownID }

type TypeSection struct {
	types []FuncType // type entries
}

func (s *TypeSection) readWasm(r io.Reader) error {
	var err error
	return err
}

type ImportSection struct {
	imports []ImportEntry
}

type ImportEntry struct {
	module string
	field  string
	kind   ExternalKind // the kind of definition being imported

	typ interface{} // imported value

	/*
		FunctionType varuint32  // type index of the function signature (if Kind is Function)
		TableType    TableType  // type of imported table (if Kind==Table)
		MemoryType   MemoryType // type of imported memory (if Kind==Memory)
		GlobalType   GlobalType // type of importer global
	*/
}

// FunctionSection declares the signature of all functions in the module
type FunctionSection struct {
	types []uint32 // indices into the type sections
}

// TableSection encodes a table
type TableSection struct {
	tables []TableType
}

// MemorySection encodes a memory
type MemorySection struct {
	memories []MemoryType
}

// GlobalSection encodes the global section
type GlobalSection struct {
	globals []GlobalVariable
}

// GlobalVariable represents a single global variable of a given type,
// mutability and with the given initializer.
type GlobalVariable struct {
	Type GlobalType // type of the variables
	Init InitExpr   // initial value of the global
}

// ExportSection encodes the export section
type ExportSection struct {
	exports []ExportEntry
}

// ExportEntry represents an exported entity.
type ExportEntry struct {
	field string
	kind  ExternalKind // kind of definition being exported
	index uint32       // index into the corresponding index space
}

// StartSection declares the start function
type StartSection struct {
	Index uint32 // start function index
}

// ElementSection encodes the elements section
type ElementSection struct {
	elements []ElemSegment
}

type ElemSegment struct {
	Index  uint32   // the table index
	Offset InitExpr // an i32 initializer expression that computes the offset at which to place the elements
	Elems  []uint32 // sequence of function indices
}

// CodeSection contains a body for every function in the module.
// The count of function declared in the function section and function bodies
// defined in this section must be the same and the i-th declaration corresponds
// to the i-th function body.
type CodeSection struct {
	Bodies []FunctionBody
}

// DataSection declares the initialized data that is loaded into linear memory
type DataSection struct {
	segments []DataSegment
}

type DataSegment struct {
	Index  uint32   // the linear memory index
	Offset InitExpr // an i32 initializer expression that computes the offset at which to place the data
	Data   []byte
}

// NameSection describes user-defined sections
type NameSection struct {
	name  string
	funcs []FunctionNames
}

type FunctionNames struct {
	name   string
	locals []LocalName
}

type LocalName struct {
	name string
}

type FunctionBody struct {
	BodySize   uint32       // size of function body to follow, in bytes
	LocalCount varuint32    // number of local entries
	Locals     []LocalEntry // local variables
	Code       Code         // bytecode of the function
}

type Code struct {
	Code []byte // bytecode of the function
	End  byte   // 0x0b, indicating the end of the body
}

type LocalEntry struct {
	Count uint32    // number of local variables of the following type
	Type  ValueType // type of the variables
}
