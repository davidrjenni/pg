// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast_test

import (
	"reflect"
	"testing"

	"github.com/davidrjenni/pg/ast"
)

func TestWalk(t *testing.T) {
	g := ast.Grammar([]*ast.Production{
		{
			Name: &ast.Name{Name: "E"},
			Expr: ast.Alternative([]ast.Expression{
				ast.Sequence([]ast.Expression{
					&ast.Name{Name: "T"},
					&ast.Terminal{Terminal: "+"},
					&ast.Name{Name: "T"},
				}),
				&ast.Name{Name: "T"},
				&ast.Epsilon{Epsilon: "e"},
			}),
		},
		{
			Name: &ast.Name{Name: "T"},
			Expr: ast.Alternative([]ast.Expression{
				ast.Sequence([]ast.Expression{
					&ast.Terminal{Terminal: "("},
					&ast.Name{Name: "T"},
					&ast.Terminal{Terminal: ")"},
				}),
				&ast.Name{Name: "T"},
			}),
		},
	})

	order := []string{
		"ast.Grammar",
		"*ast.Production",
		"*ast.Name",
		"ast.Alternative",
		"ast.Sequence",
		"*ast.Name",
		"*ast.Terminal",
		"*ast.Name",
		"*ast.Name",
		"*ast.Epsilon",
		"*ast.Production",
		"*ast.Name",
		"ast.Alternative",
		"ast.Sequence",
		"*ast.Terminal",
		"*ast.Name",
		"*ast.Terminal",
		"*ast.Name",
	}

	i := 0
	ast.Walk(func(n ast.Node) bool {
		if n == nil {
			return false
		}
		if typ := reflect.TypeOf(n).String(); order[i] != typ {
			t.Errorf("got %q want %q", typ, order[i])
		}
		i++
		return true
	}, g)
}
