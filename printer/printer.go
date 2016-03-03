// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package printer implements pretty-printing of AST nodes.
package printer

import (
	"bytes"
	"io"

	"github.com/davidrjenni/pg/ast"
)

// Fprint "pretty-prints" an AST node to output.
func Fprint(output io.Writer, node ast.Node) (err error) {
	switch n := node.(type) {
	case ast.Grammar:
		_, err = output.Write(grammar(n))
	case *ast.Production:
		_, err = output.Write(production(n))
	case ast.Expression:
		_, err = output.Write(expression(n))
	}
	return err
}

func grammar(g ast.Grammar) []byte {
	var buf bytes.Buffer
	var sep string
	for _, p := range g {
		buf.WriteString(sep)
		buf.Write(production(p))
		sep = "\n"
	}
	return buf.Bytes()
}

func production(p *ast.Production) []byte {
	var buf bytes.Buffer
	buf.Write(name(p.Name))
	buf.WriteString(" → ")
	buf.Write(expression(p.Expr))
	buf.WriteString(" .")
	return buf.Bytes()
}

func expression(expr ast.Expression) []byte {
	switch e := expr.(type) {
	case ast.Alternative:
		return alternative(e)
	case ast.Sequence:
		return sequence(e)
	case *ast.Option:
		return option(e)
	case *ast.Name:
		return name(e)
	case *ast.Terminal:
		return []byte(`"` + e.Terminal + `"`)
	case *ast.Epsilon:
		return []byte("ε")
	default:
		panic("not an expression type")
	}
}

func alternative(a ast.Alternative) []byte {
	var buf bytes.Buffer
	var sep string
	for _, e := range a {
		buf.WriteString(sep)
		buf.Write(expression(e))
		sep = " | "
	}
	return buf.Bytes()
}

func sequence(s ast.Sequence) []byte {
	var buf bytes.Buffer
	var sep string
	for _, e := range s {
		buf.WriteString(sep)
		buf.Write(expression(e))
		sep = " "
	}
	return buf.Bytes()
}

func option(o *ast.Option) []byte {
	var buf bytes.Buffer
	buf.WriteString("[ ")
	buf.Write(expression(o.Expr))
	buf.WriteString(" ]")
	return buf.Bytes()
}

func name(n *ast.Name) []byte {
	return []byte(n.Name)
}
