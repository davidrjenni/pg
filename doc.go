// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
pg - toy parser generator

pg provides packages to lex, parse and pretty-print context-free
grammars. Furthermore it provides a package for generating SLR(1)
parsers. The command pg implements a parser generator using these
packages. Package example contains example programs which use the
command pg.

A grammar is specified using BNF, which is a set of derivation
rules (productions). The following grammar specifies BNF
(represented itself in BNF); any input must satisfy this grammar:

	Production -> "PRODUCTION_NAME" "->" Expression "." .
	Expression -> Expression "|" Alternative | Alternative .
	Alternative -> Alternative Terminal | Terminal .
	Terminal -> "PRODUCTION_NAME" | "TOKEN" | "e" .

Production names and tokens are symbols of the grammar. The name
of the first production of the grammar is the start symbol. A production
name is an identifier, a token is a string. The arrow "->" is interchangeable
with the UTF-8 character U+2192 "→".
The symbol "e" indicates the empty symbol (epsilon). "e" is interchangeable
with the UTF-8 character U+03B5 "ε".

The arrow means that the symbol on the left must be replaced with
the expression on the right. An expression consists of one or more
sequences of symbols. More sequences are separated by a vertical
bar, indicating a choice. Multiple lines are allowed. A production
is terminated by a dot. The arrow means that the symbol on the
left must be replaced with the expression on the right.
*/
package pg
