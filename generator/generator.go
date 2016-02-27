// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package generator implements a parser generator.
package generator

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"text/template"

	"github.com/davidrjenni/pg/ast"
)

// Actions for the parse tables.
const (
	actionAccept = iota
	actionShift
	actionReduce
	actionError
	actionGoto
)

// generator holds the state during
// the generation of the parse.
type generator struct {
	grammar    grammar
	items      []itemSet
	firstSets  map[symbol]map[symbol]bool
	followSets map[symbol]map[symbol]bool
	Table      map[string][][2]int
	Names      []string
	Count      []int
}

// symbolAfterDot returns the symbol after
// the dot of a given item.
func (g *generator) symbolAfterDot(i item) (symbol, bool) {
	if p := g.grammar.prods[i.n]; len(p.rhs) > i.dot {
		return p.rhs[i.dot], true
	}
	return symbol{}, false
}

// closure computes the closure of an item set.
func (g *generator) closure(items itemSet) itemSet {
	for i := 0; i < len(items); i++ {
		if s, ok := g.symbolAfterDot(items[i]); ok && !s.term {
			for j, p := range g.grammar.prods {
				item := newItem(0, j)
				if p.lhs == s && !items.contains(item) {
					items = append(items, item)
				}
			}
		}
	}
	return items
}

// goTo computes the goto function for a given item set and a symbol.
func (g *generator) goTo(items itemSet, sym symbol) (res itemSet) {
	for _, i := range items {
		s, ok := g.symbolAfterDot(i)
		if ok && s == sym {
			res = append(res, newItem(i.dot+1, i.n))
		}
	}
	return g.closure(res)
}

// GenerateSLR generates an SLR(1) parser with suitable
// parse tables for a given grammar. The generated
// parser is gofmt'ed Go code.
func GenerateSLR(grammar ast.Grammar) ([]byte, error) {
	g, err := transform(grammar)
	if err != nil {
		panic(err)
	}
	gen := generator{grammar: g}

	for _, p := range g.prods {
		gen.Names = append(gen.Names, p.lhs.str)
		gen.Count = append(gen.Count, len(p.rhs))
	}

	gen.generateItems()
	gen.computeFirstSets()
	gen.computeFollowSets()

	if err := gen.buildTable(); err != nil {
		return nil, err
	}
	return gen.generateParser()
}

// generateItems generates the canonical collection
// of sets of LR(0) items.
func (g *generator) generateItems() {
	start := itemSet{newItem(0, 0)}
	g.items = []itemSet{g.closure(start)}

	for i := 0; i < len(g.items); i++ {
		for _, sym := range g.grammar.symbols {
			gotoSet := g.goTo(g.items[i], sym)
			if len(gotoSet) > 0 && !containsSet(g.items, gotoSet) {
				g.items = append(g.items, gotoSet)
			}
		}
	}
}

// buildTable builds the parse table.
func (g *generator) buildTable() error {
	g.Table = make(map[string][][2]int, len(g.grammar.symbols))
	g.grammar.symbols[end.str] = end
	for _, s := range g.grammar.symbols {
		g.Table[s.str] = make([][2]int, len(g.items))
		for i := range g.Table[s.str] {
			g.Table[s.str][i] = [2]int{actionError, 0}
		}
	}

	for i, state := range g.items {
		for _, item := range state {
			s, ok := g.symbolAfterDot(item)
			if !ok {
				for s := range g.followSets[g.grammar.prods[item.n].lhs] {
					entry := [2]int{actionReduce, item.n}
					if item.n == 0 {
						entry[0] = actionAccept
					}
					if err := g.assign(s.str, i, entry); err != nil {
						return err
					}
				}
				continue
			}
			n := g.index(g.goTo(state, s))
			entry := [2]int{actionGoto, n}
			if s.term {
				entry[0] = actionShift
			}
			if err := g.assign(s.str, i, entry); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *generator) assign(sym string, i int, entry [2]int) error {
	if x := g.Table[sym][i][0]; x != actionError && x != entry[0] {
		return fmt.Errorf("shift/reduce conflict for symbol %q", sym)
	}
	g.Table[sym][i] = entry
	return nil
}

func (g *generator) index(items itemSet) int {
	for i, set := range g.items {
		if set.equal(items) {
			return i
		}
	}
	return -1
}

// follow computes the FOLLOW function for
// all grammar symbols.
func (g *generator) computeFollowSets() {
	g.followSets = make(map[symbol]map[symbol]bool)
	g.followSets[g.grammar.prods[0].lhs] = map[symbol]bool{end: true}
	modified := true

	for modified {
		modified = false
		for _, p := range g.grammar.prods {
			for i, sym := range p.rhs {
				if sym.term {
					continue
				}
				rest := g.first(p.rhs[i+1:])
				if g.addFollow(sym, rest) {
					modified = true
				}
				if rest[epsilon] && g.addFollow(sym, g.followSets[p.lhs]) {
					modified = true
				}
			}
		}
	}
}

func (g *generator) addFollow(x symbol, symbols map[symbol]bool) bool {
	modified := false
	for s := range symbols {
		if s == epsilon {
			continue
		}
		if g.followSets[x] == nil {
			g.followSets[x] = make(map[symbol]bool)
		}
		if !g.followSets[x][s] {
			g.followSets[x][s] = true
			modified = true
		}
	}
	return modified
}

// computeFirstSets computes the FIRST sets for
// all grammar symbols.
func (g *generator) computeFirstSets() {
	g.firstSets = make(map[symbol]map[symbol]bool)
	modified := true

	for modified {
		modified = false
		for _, p := range g.grammar.prods {
			if _, ok := g.firstSets[p.lhs]; !ok {
				g.firstSets[p.lhs] = make(map[symbol]bool)
			}
			for s := range g.first(p.rhs) {
				if !g.firstSets[p.lhs][s] {
					modified = true
					g.firstSets[p.lhs][s] = true
				}
			}
		}
	}
}

func (g *generator) first(symbols []symbol) map[symbol]bool {
	set := make(map[symbol]bool)
	for _, s := range symbols {
		if s.term {
			set[s] = true
			return set
		}
		f := g.firstSets[s]
		for s := range f {
			set[s] = true
		}
		if !f[epsilon] {
			return set
		}
		delete(set, epsilon)
	}
	set[epsilon] = true
	return set
}

// generateParser returns the generated parser
// with its parse tables.
func (g *generator) generateParser() ([]byte, error) {
	var buf bytes.Buffer
	template.Must(template.New("parser").Parse(parserTmpl)).Execute(&buf, g)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", buf.Bytes(), parser.DeclarationErrors)
	if err != nil {
		return nil, err
	}

	buf.Reset()
	if err = printer.Fprint(&buf, fset, f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
