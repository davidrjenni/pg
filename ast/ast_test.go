// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast_test

import (
	"testing"

	"github.com/davidrjenni/pg/ast"
)

func TestNodes(t *testing.T) {
	var e ast.Expression
	var _ ast.Node = e
	var _ ast.Node = ast.Grammar{}
	var _ ast.Node = &ast.Production{}
	var _ ast.Node = ast.Alternative{}
	var _ ast.Node = ast.Sequence{}
	var _ ast.Node = &ast.Name{}
	var _ ast.Node = &ast.Terminal{}
	var _ ast.Node = &ast.Epsilon{}
}

func TestExpressions(t *testing.T) {
	var _ ast.Expression = ast.Alternative{}
	var _ ast.Expression = ast.Sequence{}
	var _ ast.Expression = &ast.Name{}
	var _ ast.Expression = &ast.Terminal{}
	var _ ast.Expression = &ast.Epsilon{}
}
