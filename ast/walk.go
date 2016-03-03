// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

// Visitor represents a function to be called for each node
// during a Walk. Walking stops if the visitor returns false.
type Visitor func(Node) bool

// Walk traverses an AST in depth-first order: It starts by calling v(node);
// node must not be nil. If v returns true, Walk invokes v recursively for
// each of the non-nil children of node, followed by a call of v(nil).
func Walk(v Visitor, node Node) {
	if !v(node) {
		return
	}
	switch n := node.(type) {
	case Alternative:
		for _, e := range n {
			Walk(v, e)
		}
	case Grammar:
		for _, p := range n {
			Walk(v, p)
		}
	case *Production:
		Walk(v, n.Name)
		Walk(v, n.Expr)
	case Sequence:
		for _, e := range n {
			Walk(v, e)
		}
	case *Option:
		Walk(v, n.Expr)
	}
	v(nil)
}
