// Copyright 2016 The wasm Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

import (
	"encoding/binary"
	"io"
)

type Decoder struct {
	r   io.Reader
	hdr ModuleHeader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (dec *Decoder) Decode(ptr interface{}) error {
	return nil
}

// File represents an open WASM file.
type File struct {
	Header   ModuleHeader
	Sections []*Section
}

func NewFile(r io.Reader) (*File, error) {
	var f File
	err := binary.Read(r, order, &f.Header)
	if err != nil {
		return nil, err
	}
	for {
		var section Section
		err = section.read(r)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		f.Sections = append(f.Sections, &section)
	}
	return &f, err
}
