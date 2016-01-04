// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package printer_test

import (
	"bytes"
	"testing"

	"github.com/davidrjenni/pg/ast"
	"github.com/davidrjenni/pg/printer"
)

func TestFprint(t *testing.T) {
	const expected = `Expr → Term "+" Expr | Term "-" Expr | Term | ε .
Term → Factor "*" Term | Factor "/" Term | Factor .
Factor → "(" Expr ")" | Number .
Number → Digit | Digit Number .
Digit → "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" .`

	g := ast.Grammar([]*ast.Production{
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
				&ast.Epsilon{Epsilon: "e"},
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

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, g); err != nil {
		t.Errorf("error: %v", err)
	}
	if actual := buf.String(); actual != expected {
		t.Errorf("got\n'%s'\nwant\n'%s'", actual, expected)
	}
}
