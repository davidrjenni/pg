// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package generator

import (
	"testing"

	"github.com/davidrjenni/pg/ast"
)

var testGrammar = ast.Grammar([]*ast.Production{
	{
		Name: &ast.Name{Name: "E"},
		Expr: ast.Alternative([]ast.Expression{
			&ast.Name{Name: "T"},
			ast.Sequence([]ast.Expression{
				&ast.Name{Name: "E"},
				&ast.Terminal{Terminal: "+"},
				&ast.Name{Name: "T"},
			}),
		}),
	},
	{
		Name: &ast.Name{Name: "T"},
		Expr: ast.Alternative([]ast.Expression{
			&ast.Name{Name: "F"},
			ast.Sequence([]ast.Expression{
				&ast.Name{Name: "T"},
				&ast.Terminal{Terminal: "*"},
				&ast.Name{Name: "F"},
			}),
		}),
	},
	{
		Name: &ast.Name{Name: "F"},
		Expr: ast.Alternative([]ast.Expression{
			ast.Sequence([]ast.Expression{
				&ast.Terminal{Terminal: "("},
				&ast.Name{Name: "E"},
				&ast.Terminal{Terminal: ")"},
			}),
			&ast.Terminal{Terminal: "id"},
		}),
	},
})

var testGrammar2 = ast.Grammar([]*ast.Production{
	{
		Name: &ast.Name{Name: "E"},
		Expr: ast.Sequence([]ast.Expression{
			&ast.Name{Name: "T"},
			&ast.Name{Name: "X"},
		}),
	},
	{
		Name: &ast.Name{Name: "X"},
		Expr: ast.Sequence([]ast.Expression{
			&ast.Terminal{Terminal: "+"},
			&ast.Name{Name: "T"},
			&ast.Name{Name: "X"},
		}),
	},
	{
		Name: &ast.Name{Name: "X"},
		Expr: &ast.Epsilon{},
	},
	{
		Name: &ast.Name{Name: "T"},
		Expr: ast.Sequence([]ast.Expression{
			&ast.Name{Name: "F"},
			&ast.Name{Name: "Y"},
		}),
	},
	{
		Name: &ast.Name{Name: "Y"},
		Expr: ast.Sequence([]ast.Expression{
			&ast.Terminal{Terminal: "*"},
			&ast.Name{Name: "F"},
			&ast.Name{Name: "Y"},
		}),
	},
	{
		Name: &ast.Name{Name: "Y"},
		Expr: &ast.Epsilon{},
	},
	{
		Name: &ast.Name{Name: "F"},
		Expr: ast.Sequence([]ast.Expression{
			&ast.Terminal{Terminal: "("},
			&ast.Name{Name: "E"},
			&ast.Terminal{Terminal: ")"},
		}),
	},
	{
		Name: &ast.Name{Name: "F"},
		Expr: &ast.Terminal{Terminal: "id"},
	},
})

func TestTransform(t *testing.T) {
	g, err := transform(testGrammar)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(g.prods) != 7 {
		t.Errorf("got %d, want %d", len(g.prods), 7)
	}
	if len(g.symbols) != 9 {
		t.Errorf("got %d, want %d", len(g.symbols), 9)
	}
}

func TestClosure(t *testing.T) {
	grammar, err := transform(testGrammar)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	var closures = []struct {
		items   itemSet
		closure itemSet
	}{
		{
			items:   itemSet{newItem(1, 4)},
			closure: itemSet{newItem(1, 4)},
		},
		{
			items: itemSet{
				newItem(1, 4),
				newItem(2, 4),
			},
			closure: itemSet{
				newItem(1, 4),
				newItem(2, 4),
				newItem(0, 5),
				newItem(0, 6),
			},
		},
		{
			items: itemSet{newItem(1, 5)},
			closure: itemSet{
				newItem(1, 5),
				newItem(0, 2),
				newItem(0, 1),
				newItem(0, 4),
				newItem(0, 3),
				newItem(0, 5),
				newItem(0, 6),
			},
		},
	}
	g := generator{grammar: grammar}

	for i, c := range closures {
		actual := g.closure(c.items)
		if len(actual) != len(c.closure) {
			t.Fatalf("%d: got %d, want %d", i, len(actual), len(c.closure))
		}
		for _, item := range c.closure {
			if !actual.contains(item) {
				t.Errorf("%d: want %v", i, item)
			}
		}
	}
}

func TestGoto(t *testing.T) {
	grammar, err := transform(testGrammar)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	var gotos = []struct {
		items  itemSet
		symbol symbol
		goTo   itemSet
	}{
		{
			items: []item{
				newItem(1, 0),
				newItem(1, 2),
			},
			symbol: symbol{str: "+", term: true},
			goTo: []item{
				newItem(2, 2),
				newItem(0, 4),
				newItem(0, 3),
				newItem(0, 5),
				newItem(0, 6),
			},
		},
		{
			items: []item{
				newItem(1, 0),
				newItem(2, 2),
			},
			symbol: symbol{str: "+", term: true},
			goTo:   itemSet{},
		},
	}
	gen := generator{grammar: grammar}

	for i, g := range gotos {
		actual := gen.goTo(g.items, g.symbol)
		if len(actual) != len(g.goTo) {
			t.Fatalf("%d: got %d, want %d", i, len(actual), len(g.goTo))
		}
		for _, item := range g.goTo {
			if !actual.contains(item) {
				t.Errorf("%d: want %v", i, item)
			}
		}
	}
}

func TestGenerateItems(t *testing.T) {
	grammar, err := transform(testGrammar)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	itemSets := []itemSet{
		{
			newItem(0, 0),
			newItem(0, 1),
			newItem(0, 2),
			newItem(0, 3),
			newItem(0, 4),
			newItem(0, 5),
			newItem(0, 6),
		},
		{
			newItem(1, 0),
			newItem(1, 2),
		},
		{
			newItem(1, 1),
			newItem(1, 4),
		},
		{newItem(1, 3)},
		{
			newItem(2, 2),
			newItem(0, 3),
			newItem(0, 4),
			newItem(0, 5),
			newItem(0, 6),
		},
		{newItem(1, 6)},
		{
			newItem(1, 5),
			newItem(0, 1),
			newItem(0, 2),
			newItem(0, 3),
			newItem(0, 4),
			newItem(0, 5),
			newItem(0, 6),
		},
		{
			newItem(2, 4),
			newItem(0, 5),
			newItem(0, 6),
		},
		{
			newItem(2, 5),
			newItem(1, 2),
		},
		{
			newItem(3, 2),
			newItem(1, 4),
		},
		{newItem(3, 4)},
		{newItem(3, 5)},
	}
	g := generator{grammar: grammar}
	g.generateItems()

	if len(g.items) != len(itemSets) {
		t.Errorf("got %d, want %d", len(g.items), len(itemSets))
	}
	for _, s := range itemSets {
		if !containsSet(g.items, s) {
			t.Fatalf("want %v", s)
		}
	}
}

func TestFollow(t *testing.T) {
	grammar, err := transform(testGrammar2)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	g := generator{grammar: grammar}
	g.computeFirstSets()
	g.computeFollowSets()

	expect := func(x string, symbols ...string) {
		follow := g.followSets[symbol{str: x}]
		if len(follow) != len(symbols) {
			t.Errorf("got %d, want %d", len(follow), len(symbols))
		}
		for _, s := range symbols {
			if _, ok := follow[symbol{str: s, term: true}]; !ok {
				t.Errorf("want %s in FOLLOW(%s)", s, x)
			}
		}
	}

	expect("E", "$", ")")
	expect("F", "*", "$", ")", "+")
	expect("T", "+", "$", ")")
	expect("X", "$", ")")
	expect("Y", "+", "$", ")")
}

func TestFirst(t *testing.T) {
	grammar, err := transform(testGrammar2)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	g := generator{grammar: grammar}
	g.computeFirstSets()

	expect := func(x string, symbols ...string) {
		firstSet := g.firstSets[symbol{str: x}]
		if len(firstSet) != len(symbols) {
			t.Errorf("got %d, want %d", len(firstSet), len(symbols))
		}
		first := make(map[symbol]bool)
		for s := range firstSet {
			first[s] = true
		}
		for _, s := range symbols {
			if _, ok := first[symbol{str: s, term: true}]; !ok {
				t.Errorf("want %s in FIRST(%s)", s, x)
			}
		}
	}

	expect("E", "(", "id")
	expect("T", "(", "id")
	expect("F", "(", "id")
	expect("X", "+", "ε")
	expect("Y", "*", "ε")
}
