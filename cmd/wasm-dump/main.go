// Copyright 2016 The wasm Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sbinet/wasm"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("wasm>> ")

	flag.Parse()

	fname := flag.Arg(0)
	mod, err := wasm.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("module header: %v\n", mod.Header)
	fmt.Printf("#sections: %d\n", len(mod.Sections))
	for _, section := range mod.Sections {
		fmt.Printf("section: %2d (%T)\n", section.ID(), section)
	}
}
