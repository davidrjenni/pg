// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser_test

import (
	"testing"

	"github.com/davidrjenni/pg/ast"
	"github.com/davidrjenni/pg/parser"
)

func TestParseErrors(t *testing.T) {
	errors := []struct {
		src string
		err string
	}{
		{"E E E", `test:1:3: expected →, got E (and 1 more error)`},
		{"E ->", `test:1:5: expected an expression (and 1 more error)`},
		{"E -> .", `test:1:6: expected an expression`},
		{"E -> T ", `test:1:8: production not terminated with .`},
		{"E -> T | | D | | .", `test:1:10: expected an expression (and 2 more errors)`},
		{"E -> T | | D.", `test:1:10: expected an expression`},
		{"E -> T F -> D.", `test:1:10: unexpected ->`},
		{`"foo"`, `test:1:1: expected a production, got "foo"`},
		{"?", `test:1:1: syntax error: illegal character U+003F '?'`},
	}

	for i, e := range errors {
		_, err := parser.Parse([]byte(e.src), "test")
		if err == nil {
			t.Errorf("%d: got no error, want %q", i, e.err)
		} else if err.Error() != e.err {
			t.Errorf("%d: got error %q, want %q", i, err.Error(), e.err)
		}
	}
}

func TestParse(t *testing.T) {
	const src = `Start -> Expr.
Expr → Term "+" Expr | Term "-" Expr | Term | ε.
Term → Factor "*" Term | Factor "/" Term | Factor .
Factor → "(" Expr ")" | Number .
Number → Digit | Digit Number .
Digit → "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" .`

	expected := ast.Grammar([]*ast.Production{
		{
			Name: &ast.Name{Name: "Start"},
			Expr: &ast.Name{Name: "Expr"},
		},
		{
			Name: &ast.Name{Name: "Expr"},
			Expr: ast.Alternative([]ast.Expression{
				ast.Sequence([]ast.Expression{
					&ast.Name{Name: "Term"},
					&ast.Terminal{Terminal: "+"},
					&ast.Name{Name: "Expr"},
				}),
				ast.Sequence([]ast.Expression{
					&ast.Name{Name: "Term"},
					&ast.Terminal{Terminal: "-"},
					&ast.Name{Name: "Expr"},
				}),
				&ast.Name{Name: "Term"},
				&ast.Epsilon{Epsilon: "ε"},
			}),
		},
		{
			Name: &ast.Name{Name: "Term"},
			Expr: ast.Alternative([]ast.Expression{
				ast.Sequence([]ast.Expression{
					&ast.Name{Name: "Factor"},
					&ast.Terminal{Terminal: "*"},
					&ast.Name{Name: "Term"},
				}),
				ast.Sequence([]ast.Expression{
					&ast.Name{Name: "Factor"},
					&ast.Terminal{Terminal: "/"},
					&ast.Name{Name: "Term"},
				}),
				&ast.Name{Name: "Factor"},
			}),
		},
		{
			Name: &ast.Name{Name: "Factor"},
			Expr: ast.Alternative([]ast.Expression{
				ast.Sequence([]ast.Expression{
					&ast.Terminal{Terminal: "("},
					&ast.Name{Name: "Expr"},
					&ast.Terminal{Terminal: ")"},
				}),
				&ast.Name{Name: "Number"},
			}),
		},
		{
			Name: &ast.Name{Name: "Number"},
			Expr: ast.Alternative([]ast.Expression{
				&ast.Name{Name: "Digit"},
				ast.Sequence([]ast.Expression{
					&ast.Name{Name: "Digit"},
					&ast.Name{Name: "Number"},
				}),
			}),
		},
		{
			Name: &ast.Name{Name: "Digit"},
			Expr: ast.Alternative([]ast.Expression{
				&ast.Terminal{Terminal: "0"}, &ast.Terminal{Terminal: "1"},
				&ast.Terminal{Terminal: "2"}, &ast.Terminal{Terminal: "3"},
				&ast.Terminal{Terminal: "4"}, &ast.Terminal{Terminal: "5"},
				&ast.Terminal{Terminal: "6"}, &ast.Terminal{Terminal: "7"},
				&ast.Terminal{Terminal: "8"}, &ast.Terminal{Terminal: "9"},
			}),
		},
	})

	g, err := parser.Parse([]byte(src), "test")
	if err != nil {
		t.Errorf("error: %v", err)
	}
	check(t, g, expected)
}

// check checks whether two grammars are the same.
// The position does not matter.
func check(t *testing.T, actual, expected ast.Grammar) {
	if len(actual) != len(expected) {
		t.Errorf("got %d productions, want %d", len(actual), len(expected))
	}
	for i, p := range actual {
		ep := expected[i]
		if p.Name.Name != ep.Name.Name {
			t.Errorf("%d: got production name %q, want %q", i, p.Name.Name, expected[i].Name.Name)
		}
		checkExpr(t, p.Expr, ep.Expr)
	}
}

// checkExpr checks whether two expressions are the same.
// The position does not matter.
func checkExpr(t *testing.T, actual, expected ast.Expression) {
	switch expr := expected.(type) {
	case ast.Alternative:
		a, ok := actual.(ast.Alternative)
		if !ok {
			t.Fatalf("got %T, want %T", actual, expr)
		}
		if len(a) != len(expr) {
			t.Fatalf("got alternative length %d, want %d", len(a), len(expr))
		}
		for i, e := range expr {
			checkExpr(t, a[i], e)
		}
	case ast.Sequence:
		s, ok := actual.(ast.Sequence)
		if !ok {
			t.Fatalf("got %T, want %T", actual, expr)
		}
		if len(s) != len(expr) {
			t.Fatalf("got sequence length %d, want %d", len(s), len(expr))
		}
		for i, e := range expr {
			checkExpr(t, s[i], e)
		}
	case *ast.Name:
		n, ok := actual.(*ast.Name)
		if !ok {
			t.Fatalf("got %T, want %T", actual, expr)
		}
		if n.Name != expr.Name {
			t.Errorf("got %q, want %q", n.Name, expr.Name)
		}
	case *ast.Terminal:
		term, ok := actual.(*ast.Terminal)
		if !ok {
			t.Fatalf("got %T, want %T", actual, expr)
		}
		if term.Terminal != expr.Terminal {
			t.Errorf("got %q, want %q", term.Terminal, expr.Terminal)
		}
	case *ast.Epsilon:
		epsilon, ok := actual.(*ast.Epsilon)
		if !ok {
			t.Fatalf("got %T, want %T", actual, expr)
		}
		if epsilon.Epsilon != expr.Epsilon {
			t.Errorf("got %q, want %q", epsilon.Epsilon, expr.Epsilon)
		}
	default:
		t.Errorf("unknown expression of type %T", expr)
	}
}
