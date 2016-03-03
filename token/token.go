// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package token defines constants representing the lexical tokens for pg.
package token

import "strconv"

// Type is the set of lexical tokens of pg.
type Type int

const (
	// Special tokens
	ILLEGAL Type = iota
	EOF

	// Identifiers and literals
	literalBeg
	IDENT  // Foo
	STRING // "abc"
	literalEnd

	// Operators and delimiters
	operatorBeg
	ARROW  // -> or →
	PERIOD // .
	PIPE   // |
	LBRACK // [
	RBRACK // ]
	operatorEnd

	// Keyword
	EPSILON // e or ε
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT:  "IDENT",
	STRING: "STRING",

	ARROW:  "ARROW",
	PERIOD: "PERIOD",
	PIPE:   "PIPE",
	LBRACK: "LBRACK",
	RBRACK: "RBRACK",

	EPSILON: "EPSILON",
}

// String returns the string corresponding to the token.
func (t Type) String() string {
	s := ""
	if 0 <= t && t < Type(len(tokens)) {
		s = tokens[t]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

// IsLiteral returns true for tokens corresponding to literals; it
// returns false otherwise.
func (t Type) IsLiteral() bool { return literalBeg < t && t < literalEnd }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
func (t Type) IsOperator() bool { return operatorBeg < t && t < operatorEnd }
