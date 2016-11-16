// Copyright 2016 The wasm Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

import (
	"encoding/binary"
	"errors"
	"io"
)

var order = binary.LittleEndian

type (
	varuint32 uint32
	varuint7  uint32
	varuint1  uint32

	varint64 int64
	varint32 int32
	varint7  int32
)

func (v *varuint32) read(r io.Reader) (int, error) {
	vv, n, err := uvarint(r)
	if err != nil {
		return n, err
	}
	*v = varuint32(vv)
	return n, nil
}

func (v *varuint7) read(r io.Reader) (int, error) {
	vv, n, err := uvarint(r)
	if err != nil {
		return n, err
	}
	*v = varuint7(vv)
	return n, nil
}

func uvarint(r io.Reader) (uint32, int, error) {
	var x uint32
	var s uint
	var buf = make([]byte, 1)
	for i := 0; ; i++ {
		_, err := r.Read(buf)
		if err != nil {
			return 0, i, err
		}
		b := buf[0]
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, i, errors.New("wasm: overflow")
			}
			return x | uint32(b)<<s, i, nil
		}
		x |= uint32(b&0x7f) << s
		s += 7
	}
	panic("unreachable")
}

type ValueType varint7

type BlockType varint7
type ElemType varint7
type FuncType struct {
	Form        varint7   // value for the 'func' type constructor
	ParamCount  varuint32 // number of parameters to the function
	ParamTypes  []ValueType
	ReturnCount varuint1 // number of results from the function
	ReturnType  []ValueType
}

// GlobalType describes a global variable
type GlobalType struct {
	ContentType ValueType
	Mutability  varuint1 // 0:immutable, 1:mutable
}

// TableType describes a table
type TableType struct {
	ElementType ElemType // the type of elements
	Limits      ResizableLimits
}

// MemoryType describes a memory
type MemoryType struct {
	Limits ResizableLimits
}

// ExternalKind indicates the kind of definition being imported or defined:
// 0: indicates a Function import or definition
// 1: indicates a Table import or definition
// 2: indicates a Memory import or definition
// 3: indicates a Global import or definition
type ExternalKind byte

// 0: indicates a Function import or definition
// 1: indicates a Table import or definition
// 2: indicates a Memory import or definition
// 3: indicates a Global import or definition
const (
	FunctionKind ExternalKind = 0
	TableKind                 = 1
	MemoryKind                = 2
	GlobalKind                = 3
)

// ResizableLimits describes the limits of a table or memory
type ResizableLimits struct {
	Flags   varuint32 // bit 0x1 is set if the maximum field is present
	Initial varuint32 // initial length (in units of table elements or wasm pages)
	Maximum varuint32 // only present if specified by Flags
}

// InitExpr encodes an initializer expression.
// FIXME(sbinet)
type InitExpr struct {
	Expr []byte
	End  byte
}
