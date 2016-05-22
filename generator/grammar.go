// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package generator

import (
	"errors"

	"github.com/davidrjenni/pg/ast"
)

// symbol represents a single grammar symbol.
type symbol struct {
	str   string // name of a production or terminal literal
	term  bool   // is terminal
	start bool   // is start symbol
}

var (
	// end represents the endmarker of the input.
	end = symbol{str: "$", term: true}

	// epsilon represents the empty symbol.
	epsilon = symbol{str: "ε", term: true}
)

// grammar represents a BNF grammar as
// used by the generator.
type grammar struct {
	prods   []prod            // set of productions
	symbols map[string]symbol // all symbols
}

// prod represents a single BNF production as
// used by the generator. A production consists
// of a sequence of one or more symbols.
type prod struct {
	lhs symbol   // name of the production
	rhs []symbol // expression on the right hand side
}

// transform transforms an AST grammar into a set of productions.
// Alternatives are rewritten by adding new productions for each
// alternative. These productions have the same name and one
// choice of the alternative expression as their right hand side.
// transform also adds a start symbol.
func transform(g ast.Grammar) (grammar, error) {
	var prods []prod
	symbols := make(map[string]symbol)

	if len(g) == 0 {
		return grammar{}, errors.New("grammar must not be empty")
	}

	start := prod{
		lhs: symbol{str: g[0].Name.Name + "'", term: false, start: true},
		rhs: []symbol{{str: g[0].Name.Name}},
	}
	symbols[start.lhs.str] = start.lhs
	prods = append(prods, start)

	for _, p := range g {
		lhs := symbol{str: p.Name.Name, term: false}
		symbols[lhs.str] = lhs
		if alt, ok := p.Expr.(ast.Alternative); ok {
			for _, expr := range alt {
				rhs := transformExpr(expr, symbols)
				prods = append(prods, prod{lhs: lhs, rhs: rhs})
			}
		} else {
			rhs := transformExpr(p.Expr, symbols)
			prods = append(prods, prod{lhs: lhs, rhs: rhs})
		}
	}
	return grammar{prods: prods, symbols: symbols}, nil
}

// transformExpr transforms an AST expression into
// a set of grammar symbols.
func transformExpr(expr ast.Expression, symbols map[string]symbol) (rhs []symbol) {
	switch expr := expr.(type) {
	case ast.Sequence:
		for _, e := range expr {
			rhs = append(rhs, transformExpr(e, symbols)...)
		}
	case *ast.Name:
		s := symbol{str: expr.Name, term: false}
		rhs = append(rhs, s)
		symbols[s.str] = s
	case *ast.Terminal:
		s := symbol{str: expr.Terminal, term: true}
		rhs = append(rhs, s)
		symbols[s.str] = s
	case *ast.Epsilon:
		// ignore ε
	}
	return rhs
}
