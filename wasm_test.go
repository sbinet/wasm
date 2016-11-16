// Copyright 2016 The wasm Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

import (
	"fmt"
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	r, err := os.Open("testdata/hello.wasm")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	f, err := NewFile(r)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("module header: %v\n", f.Header)
	fmt.Printf("#sections: %d\n", len(f.Sections))
	for _, section := range f.Sections {
		fmt.Printf("section: %v\n", *section)
	}
}
