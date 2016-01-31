// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Command calc implements a calculator using pg.
To re-generate the parser, run go generate.
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

//go:generate pg -f grammar -o parser.go

var l = &lexer{}

func pgLex() (typ, val string) { return l.lex() }

func pgError(err error) { fmt.Printf("error: %v\n", err) }

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := r.ReadString('\n')
		if err != nil {
			fmt.Println("cannot read line:", err)
			continue
		}
		l.init(input)
		expr := pgParse()
		fmt.Println(calc(expr))
	}
}

func calc(n pgNode) float64 {
	switch n.typ {
	case "Expr":
		return calcExpr(n)
	case "Term":
		return calcTerm(n)
	case "Factor":
		return calcFactor(n)
	case "NUMBER":
		i, err := strconv.Atoi(n.val)
		if err != nil {
			panic(err)
		}
		return float64(i)
	case "error":
		return 0.0
	default:
		panic("unknown node type")
	}
}

func calcExpr(n pgNode) float64 {
	a := calc(n.children[0])
	if len(n.children) == 3 {
		switch n.children[1].val {
		case "+":
			return a + calc(n.children[2])
		case "-":
			return a - calc(n.children[2])
		}
	}
	return a
}

func calcTerm(n pgNode) float64 {
	a := calc(n.children[0])
	if len(n.children) == 3 {
		switch n.children[1].val {
		case "*":
			return a * calc(n.children[2])
		case "/":
			return a / calc(n.children[2])
		}
	}
	return a
}

func calcFactor(n pgNode) float64 {
	// Expression in parens.
	if len(n.children) == 3 {
		return calc(n.children[1])
	}
	return calc(n.children[0])
}
