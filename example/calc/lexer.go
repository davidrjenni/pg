// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "unicode"

type lexer struct {
	input string
	index int
}

func (l *lexer) init(input string) {
	l.index = 0
	l.input = input + "$"
}

func (l *lexer) next() rune {
	defer func() { l.index++ }()
	return rune(l.input[l.index])
}

func (l *lexer) lex() (typ, tok string) {
	for unicode.IsSpace(l.next()) {
	}
	l.index--

	switch r := l.next(); {
	case unicode.IsDigit(r):
		for unicode.IsDigit(r) {
			tok += string(r)
			r = l.next()
		}
		l.index--
		return "NUMBER", tok
	default:
		return "", string(r)
	}
}
