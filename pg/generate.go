// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/davidrjenni/pg/generator"
	"github.com/davidrjenni/pg/parser"
)

func gen(args []string) {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	out := flags.String("o", "out.go", "output file")

	if len(args) == 0 {
		log.SetPrefix("")
		log.Fatal("Usage: pg gen [flags] <file>\nFlags:\n\t-o output file (instead of out.go)")
	}
	in := args[len(args)-1]
	flags.Parse(args[:len(args)-1])

	f, err := os.Open(in)
	if err != nil {
		log.Fatalf("cannot open file: %v", err)
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("cannot read file: %v", err)
	}

	g, err := parser.Parse(src, in)
	if err != nil {
		log.Fatalf(err.Error())
	}

	buf, err := generator.GenerateSLR(g)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if err = ioutil.WriteFile(*out, buf, 0644); err != nil {
		log.Fatalf("cannot write file: %v", err)
	}
}
