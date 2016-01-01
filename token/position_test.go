// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token_test

import (
	"testing"

	"github.com/davidrjenni/pg/token"
)

func TestPosString(t *testing.T) {
	positions := []struct {
		p   token.Pos
		str string
	}{
		{token.Pos{Filename: "foo", Line: 1, Column: 3}, "foo:1:3"},
		{token.Pos{Line: 1, Column: 3}, "1:3"},
		{token.Pos{Filename: "foo", Column: 3}, "foo"},
		{token.Pos{Column: 3}, "-"},
	}

	for i, pos := range positions {
		if actual := pos.p.String(); actual != pos.str {
			t.Errorf("%d: got %q, want %q", i, pos.p, pos.str)
		}
	}
}
