// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
pg is a toy SLR parser generator.

pg converts a context-free grammar in Backus-Naur Form (BNF) into
parse tables for an SLR(1) parser. The input must satisfy the
grammar specified in package github.com/davidrjenni/pg.

The options are
	-f input	Input file containing the grammar
	-o output	Direct output to the specified file instead of out.go

The output file contains the parse tables and the function
"pgParse() (node, error)" which parses input according to the given
grammar rules, using "pgLex() (string, string)" to obtain the input and
"pgError(msg string)" to report errors. The documentation for pgParse,
pgLex and pgError can be found in package github.com/davidrjenni/pg/generator.

The package github.com/davidrjenni/pg/example contains working examples.
*/
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/davidrjenni/pg/generator"
	"github.com/davidrjenni/pg/parser"
)

func main() {
	in := flag.String("f", "", "input file")
	out := flag.String("o", "out.go", "output file")

	log.SetPrefix("pg: ")
	log.SetFlags(0)

	flag.Parse()
	if *in == "" {
		flag.Usage()
		os.Exit(2)
	}

	f, err := os.Open(*in)
	if err != nil {
		log.Fatalf("cannot open file: %v", err)
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("cannot read file: %v", err)
	}

	g, err := parser.Parse(src, *in)
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
