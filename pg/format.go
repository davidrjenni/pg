// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/davidrjenni/pg/parser"
	"github.com/davidrjenni/pg/printer"
)

func format(args []string) {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	write := flags.Bool("w", false, "write to file (instead of stdout)")

	if len(args) == 0 {
		log.SetPrefix("")
		log.Fatal("Usage: pg fmt [flags] <file>\nFlags:\n\t-w write to file (instead of stdout)")
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

	if *write {
		f.Close()
		var buf bytes.Buffer
		if err := printer.Fprint(&buf, g); err != nil {
			log.Fatalf("cannot print grammar: %v", err)
		}
		if err = ioutil.WriteFile(in, buf.Bytes(), 0644); err != nil {
			log.Fatalf("cannot write grammar: %v", err)
		}
		return
	}
	if err := printer.Fprint(os.Stdout, g); err != nil {
		log.Fatalf("cannot print grammar: %v", err)
	}
}
