// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import "fmt"

// Pos describes an arbitrary source position
// including the file, line, and column location.
// A position is valid if the line number is > 0.
type Pos struct {
	Filename string // filename, if any
	Offset   int    // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (character count)
}

// String returns a string in one of several forms:
//
//	file:line:column    valid position with filename
//	line:column         valid position without filename
//	file                invalid position with filename
//	-                   invalid position without filename
func (p Pos) String() string {
	s := p.Filename
	if p.Line > 0 {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d:%d", p.Line, p.Column)
	}
	if s == "" {
		s = "-"
	}
	return s
}
