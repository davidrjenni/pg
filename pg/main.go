// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
pg is tool for managing context-free grammars.

pg offers the following commands:
	gen	generate parser

"pg gen" converts a context-free grammar in Backus-Naur Form (BNF)
into parse tables for an SLR(1) parser. The input must satisfy the
grammar specified in package github.com/davidrjenni/pg.

The option is
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
	"log"
	"os"
)

var commands = map[string]func(args []string){
	"gen": gen,
}

func main() {
	log.SetPrefix("pg: ")
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.SetPrefix("")
		log.Fatal(`Usage: pg <command> [arguments]
Commands:
	gen	generate parser
`)
	}

	cmd, ok := commands[os.Args[1]]
	if !ok {
		log.Fatalf("unknown command %q", os.Args[1])
	}
	cmd(os.Args[2:])
}
