// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ast declares the types used to represent syntax trees for pg.
package ast

import "github.com/davidrjenni/pg/token"

type (
	// Node is an element in the abstract syntax tree.
	Node interface {
		// Pos returns the position of the first character of the expression.
		Pos() token.Pos
		node()
	}

	// Grammar represents a set of EBNF productions.
	Grammar []*Production

	// Production represents a single EBNF production.
	Production struct {
		Name *Name      // name of the production (lhs)
		Expr Expression // expression of the production (rhs)
	}

	// Expression represents a production expression.
	Expression interface {
		Node
		expr()
	}

	// Alternative represents a list of alternative expressions.
	Alternative []Expression

	// Sequence represents a list of sequential expressions.
	Sequence []Expression

	// Name represents a production name.
	Name struct {
		Name     string    // name of the production
		StartPos token.Pos // position of the first character
	}

	// Terminal represents a terminal.
	Terminal struct {
		Terminal string    // terminal literal
		QuotePos token.Pos // position of "
	}

	// Epsilon represents the epsilon keyword.
	Epsilon struct {
		Epsilon string    // epsilon keyword
		Start   token.Pos // position of e or ε
	}
)

func (g Grammar) Pos() token.Pos     { return g[0].Pos() }
func (p *Production) Pos() token.Pos { return p.Name.Pos() }
func (a Alternative) Pos() token.Pos { return a[0].Pos() }
func (s Sequence) Pos() token.Pos    { return s[0].Pos() }
func (n *Name) Pos() token.Pos       { return n.StartPos }
func (t *Terminal) Pos() token.Pos   { return t.QuotePos }
func (e *Epsilon) Pos() token.Pos    { return e.Start }

func (Grammar) node()     {}
func (Production) node()  {}
func (Alternative) node() {}
func (Sequence) node()    {}
func (Name) node()        {}
func (Terminal) node()    {}
func (Epsilon) node()     {}

func (Alternative) expr() {}
func (Sequence) expr()    {}
func (Name) expr()        {}
func (Terminal) expr()    {}
func (Epsilon) expr()     {}
