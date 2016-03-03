// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token_test

import (
	"testing"

	"github.com/davidrjenni/pg/token"
)

func TestTypeString(t *testing.T) {
	tokens := []struct {
		tt  token.Type
		str string
	}{
		{token.ILLEGAL, "ILLEGAL"},
		{token.EOF, "EOF"},
		{token.IDENT, "IDENT"},
		{token.STRING, "STRING"},
		{token.ARROW, "ARROW"},
		{token.PERIOD, "PERIOD"},
		{token.PIPE, "PIPE"},
		{token.LBRACK, "LBRACK"},
		{token.RBRACK, "RBRACK"},
		{token.EPSILON, "EPSILON"},
	}

	for i, token := range tokens {
		if token.tt.String() != token.str {
			t.Errorf("%d: got %q want %q", i, token.tt, token.str)
		}
	}
}

func TestTypeIsLiteral(t *testing.T) {
	tokens := []struct {
		tt        token.Type
		isLiteral bool
	}{
		{token.ILLEGAL, false},
		{token.EOF, false},
		{token.IDENT, true},
		{token.STRING, true},
		{token.ARROW, false},
		{token.PERIOD, false},
		{token.PIPE, false},
		{token.LBRACK, false},
		{token.RBRACK, false},
		{token.EPSILON, false},
	}

	for i, token := range tokens {
		if token.tt.IsLiteral() != token.isLiteral {
			t.Errorf("%d: got %v want %v", i, token.tt.IsLiteral(), token.isLiteral)
		}
	}
}

func TestTypeIsOperator(t *testing.T) {
	tokens := []struct {
		tt   token.Type
		isOp bool
	}{
		{token.ILLEGAL, false},
		{token.EOF, false},
		{token.IDENT, false},
		{token.STRING, false},
		{token.ARROW, true},
		{token.PERIOD, true},
		{token.PIPE, true},
		{token.LBRACK, true},
		{token.RBRACK, true},
		{token.EPSILON, false},
	}

	for i, token := range tokens {
		if token.tt.IsOperator() != token.isOp {
			t.Errorf("%d: got %v want %v", i, token.tt.IsOperator(), token.isOp)
		}
	}
}
