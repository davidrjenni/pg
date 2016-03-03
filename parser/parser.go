// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package parser implements a parser for pg.
package parser

import (
	"fmt"

	"github.com/davidrjenni/pg/ast"
	"github.com/davidrjenni/pg/scanner"
	"github.com/davidrjenni/pg/token"
)

type errors []error

func (e errors) err() error {
	if len(e) == 0 {
		return nil
	}
	return e
}

func (e errors) Error() string {
	switch len(e) {
	case 1:
		return e[0].Error()
	case 2:
		return fmt.Sprintf("%s (and %d more error)", e[0], len(e)-1)
	default:
		return fmt.Sprintf("%s (and %d more errors)", e[0], len(e)-1)
	}
}

type parser struct {
	grammar ast.Grammar
	scanner *scanner.Scanner
	errs    errors

	// Last token
	pos token.Pos
	typ token.Type
	lit string

	// Set to true to go one back
	unscan bool
}

func (p *parser) errorf(pos token.Pos, format string, args ...interface{}) {
	p.errs = append(p.errs, fmt.Errorf(fmt.Sprintf("%s: %s", pos, format), args...))
}

func (p *parser) next() {
	if p.unscan {
		p.unscan = false
		return
	}
	p.pos, p.typ, p.lit = p.scanner.Scan()
	for p.typ == token.ILLEGAL {
		p.pos, p.typ, p.lit = p.scanner.Scan()
	}
}

// Parse parses the source code and returns the abstract syntax tree.
func Parse(src []byte, filename string) (ast.Grammar, error) {
	p := &parser{scanner: scanner.New(src, filename)}
	p.scanner.Err = func(pos token.Pos, msg string) {
		p.errorf(pos, "syntax error: %s", msg)
	}
	p.parse()
	return p.grammar, p.errs.err()
}

func (p *parser) parse() {
	for {
		switch p.next(); p.typ {
		case token.EOF:
			return
		case token.IDENT:
			prod := &ast.Production{Name: &ast.Name{Name: p.lit, StartPos: p.pos}}
			if p.next(); p.typ != token.ARROW {
				p.unscan = true
				p.errorf(p.pos, "expected →, got %s", p.lit)
			}
			prod.Expr = p.parseExpression()
			p.grammar = append(p.grammar, prod)
		default:
			p.errorf(p.pos, "expected a production, got %s", p.lit)
		}
	}
}

func (p *parser) parseExpression() ast.Expression {
	var alt ast.Alternative
	for {
		alt = append(alt, p.parseSequence())
		if p.typ != token.PIPE && p.typ != token.RBRACK {
			if p.typ != token.PERIOD {
				p.errorf(p.pos, "production not terminated with .")
			}
			if len(alt) == 1 {
				return alt[0]
			}
			return alt
		}
	}
}

func (p *parser) parseSequence() ast.Expression {
	var seq ast.Sequence
Loop:
	for {
		switch p.next(); p.typ {
		case token.IDENT:
			seq = append(seq, &ast.Name{Name: p.lit, StartPos: p.pos})
		case token.STRING:
			seq = append(seq, &ast.Terminal{Terminal: p.lit[1 : len(p.lit)-1], QuotePos: p.pos})
		case token.EPSILON:
			seq = append(seq, &ast.Epsilon{Epsilon: p.lit, Start: p.pos})
		case token.PIPE, token.RBRACK, token.PERIOD, token.EOF:
			if len(seq) == 0 {
				p.errorf(p.pos, "expected an expression")
			}
			break Loop
		case token.LBRACK:
			lbrack := p.pos
			expr := p.parseOption()
			seq = append(seq, &ast.Option{Expr: expr, Lbrack: lbrack})
		default:
			p.errorf(p.pos, "unexpected %s", p.lit)
		}
	}
	if len(seq) == 1 {
		return seq[0]
	}
	return seq
}

func (p *parser) parseOption() ast.Expression {
	var alt ast.Alternative
	for {
		alt = append(alt, p.parseSequence())
		if p.typ != token.PIPE {
			if p.typ != token.RBRACK {
				p.errorf(p.pos, "option not terminated with ]")
			}
			if len(alt) == 1 {
				return alt[0]
			}
			return alt
		}
	}
}
