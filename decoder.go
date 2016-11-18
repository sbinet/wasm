// Copyright 2016 The wasm Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type decoder struct {
	r   io.Reader
	err error
}

func (d *decoder) readVarI7(r io.Reader, v *int32) {
	// FIXME(sbinet) ?
	d.readVarI32(r, v)
}

func (d *decoder) readVarI32(r io.Reader, v *int32) {
	if d.err != nil {
		return
	}
	*v, _, d.err = varint(r)
}

func (d *decoder) readVarU1(r io.Reader, v *uint32) {
	// FIXME(sbinet) ?
	d.readVarU32(r, v)
}

func (d *decoder) readVarU7(r io.Reader, v *uint32) {
	// FIXME(sbinet) ?
	d.readVarU32(r, v)
}

func (d *decoder) readVarU32(r io.Reader, v *uint32) {
	if d.err != nil {
		return
	}
	*v, _, d.err = uvarint(r)
}

func (d *decoder) readString(r io.Reader, s *string) {
	if d.err != nil {
		return
	}
	var sz uint32
	d.readVarU32(r, &sz)
	var buf = make([]byte, sz)
	d.read(r, buf)
	*s = string(buf)
}

func (d *decoder) read(r io.Reader, buf []byte) {
	if d.err != nil || len(buf) == 0 {
		return
	}
	_, d.err = r.Read(buf)
}

func (d *decoder) readModule() (Module, error) {
	var (
		m   Module
		err error
	)

	if d.err != nil {
		err = d.err
		return m, err
	}

	d.readHeader(d.r, &m.Header)
	for {
		s := d.readSection()
		if s == nil {
			return m, d.err
		}
		m.Sections = append(m.Sections, s)
	}
	return m, d.err
}

func (d *decoder) readHeader(r io.Reader, hdr *ModuleHeader) {
	if d.err != nil {
		return
	}
	d.err = binary.Read(r, order, hdr)
	if d.err != nil {
		return
	}

	if hdr.Magic != magicWASM {
		d.err = fmt.Errorf("wasm: invalid magic number (%q)", string(hdr.Magic[:]))
		return
	}
}

func (d *decoder) readSection() Section {
	var (
		id  uint32
		sz  uint32
		sec Section
	)

	d.readVarU32(d.r, &id)
	if d.err != nil {
		if d.err == io.EOF {
			d.err = nil
		}
		return nil
	}
	d.readVarU32(d.r, &sz)

	r := &io.LimitedReader{R: d.r, N: int64(sz)}
	switch SectionID(id) {
	case UnknownID:
		var s NameSection
		d.readNameSection(r, &s)
		// fmt.Printf("--- name: %q, funcs: %d\n", s.name, len(s.funcs))
		sec = s

	case TypeID:
		var s TypeSection
		d.readTypeSection(r, &s)
		// fmt.Printf("--- types: %d\n", len(s.types))
		sec = s

	case ImportID:
		var s ImportSection
		d.readImportSection(r, &s)
		// fmt.Printf("--- imports: %d\n", len(s.imports))
		/*
			for ii, imp := range s.imports {
				fmt.Printf("    entry[%d]: %q|%q|%x\n", ii, imp.module, imp.field, imp.kind)
			}
		*/
		sec = s

	case FunctionID:
		var s FunctionSection
		d.readFunctionSection(r, &s)
		// fmt.Printf("--- functions: %d\n", len(s.types))
		sec = s

	case TableID:
		var s TableSection
		d.readTableSection(r, &s)
		// fmt.Printf("--- tables: %d\n", len(s.tables))
		sec = s

	case MemoryID:
		var s MemorySection
		d.readMemorySection(r, &s)
		// fmt.Printf("--- memories: %d\n", len(s.memories))
		sec = s

	case GlobalID:
		var s GlobalSection
		d.readGlobalSection(r, &s)
		// fmt.Printf("--- globals: %d\n", len(s.globals))
		/*
			for ii, ge := range s.globals {
				fmt.Printf("   ge[%d]: type={%x, 0x%x} init=%d\n",
					ii, ge.Type.ContentType, ge.Type.Mutability, len(ge.Init.Expr),
				)
			}
		*/
		sec = s

	case ExportID:
		var s ExportSection
		d.readExportSection(r, &s)
		// fmt.Printf("--- exports: %d\n", len(s.exports))
		sec = s

	case StartID:
		var s StartSection
		d.readStartSection(r, &s)
		// fmt.Printf("--- start: 0x%x\n", s.Index)
		sec = s

	case ElementID:
		var s ElementSection
		d.readElementSection(r, &s)
		// fmt.Printf("--- elements: %d\n", len(s.elements))
		sec = s

	case CodeID:
		var s CodeSection
		d.readCodeSection(r, &s)
		// fmt.Printf("--- func-bodies: %d\n", len(s.Bodies))
		sec = s

	case DataID:
		var s DataSection
		d.readDataSection(r, &s)
		// fmt.Printf("--- data-segments: %d\n", len(s.segments))
		sec = s

	default:
		d.err = fmt.Errorf("wasm: invalid section ID")

	}

	if r.N != 0 {
		log.Printf("wasm: N=%d bytes unread! (section=%d)\n", r.N, sec.ID())
		buf := make([]byte, r.N)
		d.read(r, buf)
	}

	return sec
}

func (d *decoder) readNameSection(r io.Reader, s *NameSection) {
	if d.err != nil {
		return
	}

	d.readString(r, &s.name)
	var n uint32
	d.readVarU32(r, &n)
	s.funcs = make([]FunctionNames, int(n))
	for i := range s.funcs {
		d.readFunctionNames(r, &s.funcs[i])
	}
}

func (d *decoder) readFunctionNames(r io.Reader, f *FunctionNames) {
	if d.err != nil {
		return
	}

	d.readString(r, &f.name)
	var n uint32
	d.readVarU32(r, &n)
	f.locals = make([]LocalName, int(n))
	for i := range f.locals {
		d.readLocalName(r, &f.locals[i])
	}
}

func (d *decoder) readLocalName(r io.Reader, local *LocalName) {
	if d.err != nil {
		return
	}

	d.readString(r, &local.name)
}

func (d *decoder) readTypeSection(r io.Reader, s *TypeSection) {
	if d.err != nil {
		return
	}

	var n uint32
	d.readVarU32(r, &n)
	s.types = make([]FuncType, int(n))
	for i := range s.types {
		d.readFuncType(r, &s.types[i])
	}
}

func (d *decoder) readFuncType(r io.Reader, ft *FuncType) {
	if d.err != nil {
		return
	}

	d.readVarU7(r, &ft.form)

	var params uint32
	d.readVarU32(r, &params)
	ft.params = make([]ValueType, int(params))
	for i := range ft.params {
		d.readValueType(r, &ft.params[i])
	}

	var results uint32
	d.readVarU32(r, &results)
	ft.results = make([]ValueType, int(results))
	for i := range ft.results {
		d.readValueType(r, &ft.results[i])
	}
}

func (d *decoder) readValueType(r io.Reader, vt *ValueType) {
	if d.err != nil {
		return
	}

	var v int32
	d.readVarI7(r, &v)
	*vt = ValueType(v)
}

func (d *decoder) readImportSection(r io.Reader, s *ImportSection) {
	if d.err != nil {
		return
	}

	var sz uint32
	d.readVarU32(r, &sz)
	s.imports = make([]ImportEntry, int(sz))
	for i := range s.imports {
		d.readImportEntry(r, &s.imports[i])
	}
}

func (d *decoder) readImportEntry(r io.Reader, ie *ImportEntry) {
	if d.err != nil {
		return
	}

	d.readString(r, &ie.module)
	d.readString(r, &ie.field)
	d.readExternalKind(r, &ie.kind)

	switch ie.kind {
	case FunctionKind:
		var idx uint32
		d.readVarU32(r, &idx)
		ie.typ = idx

	case TableKind:
		var tt TableType
		d.readTableType(r, &tt)
		ie.typ = tt

	case MemoryKind:
		var mt MemoryType
		d.readMemoryType(r, &mt)
		ie.typ = mt

	case GlobalKind:
		var gt GlobalType
		d.readGlobalType(r, &gt)
		ie.typ = gt

	default:
		fmt.Printf("module=%q field=%q\n", ie.module, ie.field)
		d.err = fmt.Errorf("wasm: invalid ExternalKind (%d)", byte(ie.kind))
	}
}

func (d *decoder) readExternalKind(r io.Reader, ek *ExternalKind) {
	if d.err != nil {
		return
	}

	var v [1]byte
	d.read(r, v[:])
	*ek = ExternalKind(v[0])
}

func (d *decoder) readTableType(r io.Reader, tt *TableType) {
	if d.err != nil {
		return
	}

	d.readElemType(r, &tt.ElemType)
	d.readResizableLimits(r, &tt.Limits)
}

func (d *decoder) readElemType(r io.Reader, et *ElemType) {
	if d.err != nil {
		return
	}

	var v int32
	d.readVarI7(r, &v)
	*et = ElemType(v)
}

func (d *decoder) readResizableLimits(r io.Reader, tl *ResizableLimits) {
	if d.err != nil {
		return
	}

	d.readVarU32(r, &tl.Flags)
	d.readVarU32(r, &tl.Initial)
	if tl.Flags&0x1 != 0 {
		d.readVarU32(r, &tl.Maximum)
	}
}

func (d *decoder) readMemoryType(r io.Reader, mt *MemoryType) {
	if d.err != nil {
		return
	}

	d.readResizableLimits(r, &mt.Limits)
}

func (d *decoder) readGlobalType(r io.Reader, gt *GlobalType) {
	if d.err != nil {
		return
	}

	d.readValueType(r, &gt.ContentType)
	var mut uint32
	d.readVarU1(r, &mut)
	gt.Mutability = varuint1(mut)
}

func (d *decoder) readFunctionSection(r io.Reader, s *FunctionSection) {
	if d.err != nil {
		return
	}

	var sz uint32
	d.readVarU32(r, &sz)
	s.types = make([]uint32, int(sz))
	for i := range s.types {
		d.readVarU32(r, &s.types[i])
	}
}

func (d *decoder) readTableSection(r io.Reader, s *TableSection) {
	if d.err != nil {
		return
	}

	var sz uint32
	d.readVarU32(r, &sz)
	s.tables = make([]TableType, int(sz))
	for i := range s.tables {
		d.readTableType(r, &s.tables[i])
	}
}

func (d *decoder) readMemorySection(r io.Reader, s *MemorySection) {
	if d.err != nil {
		return
	}

	var sz uint32
	d.readVarU32(r, &sz)
	s.memories = make([]MemoryType, int(sz))
	for i := range s.memories {
		d.readMemoryType(r, &s.memories[i])
	}
}

func (d *decoder) readGlobalSection(r io.Reader, s *GlobalSection) {
	if d.err != nil {
		return
	}

	var sz uint32
	d.readVarU32(r, &sz)
	s.globals = make([]GlobalVariable, int(sz))
	for i := range s.globals {
		d.readGlobalVariable(r, &s.globals[i])
	}
}

func (d *decoder) readGlobalVariable(r io.Reader, gv *GlobalVariable) {
	if d.err != nil {
		return
	}

	out := new(bytes.Buffer)
	r = io.TeeReader(r, out)
	d.readGlobalType(r, &gv.Type)
	d.readInitExpr(r, &gv.Init)
}

func (d *decoder) readInitExpr(r io.Reader, ie *InitExpr) {
	if d.err != nil {
		return
	}

	var err error
	var n int
	var buf [1]byte
	for {
		n, err = r.Read(buf[:])
		if err != nil || n <= 0 {
			break
		}
		v := buf[0]
		if v == Op_end {
			ie.End = byte(Op_end)
			break
		}
		ie.Expr = append(ie.Expr, v)
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		d.err = err
	}
}

func (d *decoder) readExportSection(r io.Reader, s *ExportSection) {
	if d.err != nil {
		return
	}

	var sz uint32
	d.readVarU32(r, &sz)
	s.exports = make([]ExportEntry, int(sz))
	for i := range s.exports {
		d.readExportEntry(r, &s.exports[i])
	}
}

func (d *decoder) readExportEntry(r io.Reader, ee *ExportEntry) {
	if d.err != nil {
		return
	}

	d.readString(r, &ee.field)
	d.readExternalKind(r, &ee.kind)
	d.readVarU32(r, &ee.index)
}

func (d *decoder) readStartSection(r io.Reader, s *StartSection) {
	if d.err != nil {
		return
	}

	d.readVarU32(r, &s.Index)
}

func (d *decoder) readElementSection(r io.Reader, s *ElementSection) {
	if d.err != nil {
		return
	}

	var sz uint32
	d.readVarU32(r, &sz)
	s.elements = make([]ElemSegment, int(sz))
	for i := range s.elements {
		d.readElemSegment(r, &s.elements[i])
	}
}

func (d *decoder) readElemSegment(r io.Reader, es *ElemSegment) {
	if d.err != nil {
		return
	}

	d.readVarU32(r, &es.Index)
	d.readInitExpr(r, &es.Offset)

	var sz uint32
	d.readVarU32(r, &sz)
	es.Elems = make([]uint32, int(sz))
	for i := range es.Elems {
		d.readVarU32(r, &es.Elems[i])
	}
}

func (d *decoder) readCodeSection(r io.Reader, s *CodeSection) {
	if d.err != nil {
		return
	}

	var sz uint32
	d.readVarU32(r, &sz)
	s.Bodies = make([]FunctionBody, int(sz))
	for i := range s.Bodies {
		d.readFunctionBody(r, &s.Bodies[i])
	}
}

func (d *decoder) readFunctionBody(r io.Reader, fb *FunctionBody) {
	if d.err != nil {
		return
	}

	d.readVarU32(r, &fb.BodySize)
	r = io.LimitReader(r, int64(fb.BodySize))
	var locals uint32
	d.readVarU32(r, &locals)
	fb.Locals = make([]LocalEntry, int(locals))
	for i := range fb.Locals {
		d.readLocalEntry(r, &fb.Locals[i])
	}

	rcode := new(bytes.Buffer)
	io.Copy(rcode, r)
	d.readCode(rcode, &fb.Code)
}

func (d *decoder) readCode(r io.Reader, code *Code) {
	if d.err != nil {
		return
	}

	var err error
	var n int
	var buf [1]byte
	for {
		n, err = r.Read(buf[:])
		if err != nil || n <= 0 {
			break
		}
		v := buf[0]
		if v == Op_end {
			code.End = byte(Op_end)
			break
		}
		code.Code = append(code.Code, v)
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		d.err = err
	}
}

func (d *decoder) readLocalEntry(r io.Reader, le *LocalEntry) {
	if d.err != nil {
		return
	}

	d.readVarU32(r, &le.Count)
	d.readValueType(r, &le.Type)
}

func (d *decoder) readDataSection(r io.Reader, s *DataSection) {
	if d.err != nil {
		return
	}

	var sz uint32
	d.readVarU32(r, &sz)
	s.segments = make([]DataSegment, int(sz))
	for i := range s.segments {
		d.readDataSegment(r, &s.segments[i])
	}
}

func (d *decoder) readDataSegment(r io.Reader, ds *DataSegment) {
	if d.err != nil {
		return
	}

	d.readVarU32(r, &ds.Index)
	d.readInitExpr(r, &ds.Offset)

	var sz uint32
	d.readVarU32(r, &sz)
	ds.Data = make([]byte, int(sz))
	d.read(r, ds.Data)
}
