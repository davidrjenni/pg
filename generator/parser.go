// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
The generated parser provides the function pgParse:

	// pgParse returns an AST node.
	func pgParse() pgNode

pgParse returns an AST node. The returned node is of
the type of the start production of the grammar. If
no production could be applied, a node with type
"error" is returned.

An node looks like this:

	// pgNode is an element in the abstract syntax tree.
	type pgNode struct {
		typ      string // type, as defined in the grammar or "error"
		val      string // actual value or empty for non-terminal nodes
		children []pgNode // child nodes, empty for terminal nodes
	}

pgParse uses pgLex to obtain the next lexical token.
The client package must implement a function pgLex:

	// pgLex is called to obtain the next
	// lexical token tok of type typ. pgLex
	// returns "$" to indicate end of input.
	pgLex() (typ, tok string)

The client package must also implement a function
pgError, which is called if an error occurs while
parsing.

	// pgError is called if an error
	// occured while parsing.
	pgError(err error)
*/
package generator

const parserTmpl = `package main

import "fmt"

type pgElem struct {
	sym   string
	state int
}

type pgStack []pgElem

func (s pgStack) top() pgElem    { return s[len(s)-1] }
func (s *pgStack) pop(n int)     { *s = (*s)[:len(*s)-n] }
func (s *pgStack) push(e pgElem) { *s = append(*s, e) }

type pgNode struct {
	typ      string
	val      string
	children []pgNode
}

func pgParse() pgNode {
	var (
		table    = {{ printf "%#v" .Table }}
		count    = {{ printf "%#v" .Count }}
		names    = {{ printf "%#v" .Names }}
		tree     = make([]pgNode, 0)
		stack    = &pgStack{pgElem{state: 0}}
		typ, tok = pgLex()
	)

	for {
		s := stack.top()
		var column [][2]int
		// Use type if available.
		if typ != "" {
			column = table[typ]
		} else {
			column = table[tok]
		}
		if column == nil {
			pgError(fmt.Errorf("unexpected token %q (type: %q)", tok, typ))
			typ, tok = pgLex()
			if tok == "$" {
				if len(tree) == 0 {
					return pgNode{typ: "error"}
				}
				return tree[0]
			}
			continue
		}
		entry := column[s.state]
		switch entry[0] {
		case 2: // Reduce
			c := count[entry[1]]
			name := names[entry[1]]
			stack.pop(2 * c)
			s = stack.top()
			stack.push(pgElem{sym: name})
			stack.push(pgElem{state: table[name][s.state][1]})
			var rest []pgNode
			for _, n := range tree[:len(tree)-c] {
				rest = append(rest, n)
			}
			tree = append(rest, pgNode{typ: name, val: name, children: tree[len(tree)-c:]})
		case 1: // Shift
			stack.push(pgElem{sym: tok})
			stack.push(pgElem{state: entry[1]})
			tree = append(tree, pgNode{typ: typ, val: tok})
			typ, tok = pgLex()
		case 0: // Accept
			if tok == "$" {
				return tree[0]
			}
		default:
			if tok == "$" {
				pgError(fmt.Errorf("unexpected end of input"))
				if len(tree) == 0 {
					return pgNode{typ: "error"}
				}
				return tree[0]
			}
			pgError(fmt.Errorf("unexpected token %q (type: %q)", tok, typ))
			typ, tok = pgLex()
		}
	}
	return tree[0]
}
`
