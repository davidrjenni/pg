// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package scanner implements a scanner for pg.
package scanner

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/davidrjenni/pg/token"
)

const (
	bom = 0xFEFF // byte order mark, only permitted as very first character
	eof = -1     // eof is a marker rune for the end of the reader.
)

// Scanner represents a lexical scanner for pg.
type Scanner struct {
	src []byte // source buffer for advancing and scanning

	// scanning state
	ch         rune      // current character
	pos        token.Pos // current position
	rdOffset   int       // reading offset (position after current character)
	lineOffset int       // current line offset

	Err      func(token.Pos, string) // error reporting; or nil
	ErrCount int                     // number of errors encountered
}

// New creates and initializes a new instance of Scanner using src as
// its source content and filename as the current filename.
func New(src []byte, filename string) *Scanner {
	return &Scanner{src: src, ch: ' ', pos: token.Pos{Filename: filename, Line: 1}}
}

func (s *Scanner) error(pos token.Pos, msg string) {
	if s.Err != nil {
		s.Err(pos, msg)
	}
	s.ErrCount++
}

// Read the next Unicode char into s.ch.
// s.ch might be eof (end-of-file).
func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.pos.Offset = s.rdOffset
		s.pos.Column = s.rdOffset - s.lineOffset + 1
		if s.ch == '\n' {
			s.lineOffset = s.rdOffset
			s.pos.Line++
			s.pos.Column = 1
		}
		r, w := rune(s.src[s.rdOffset]), 1
		switch {
		case r == 0:
			s.error(s.pos, "illegal character NUL")
		case r >= 0x80:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				s.error(s.pos, "illegal UTF-8 encoding")
			} else if r == bom && s.pos.Offset > 0 {
				s.error(s.pos, "illegal byte order mark")
			}
		}
		s.rdOffset += w
		s.ch = r
	} else {
		s.pos.Offset = len(s.src)
		s.pos.Column++
		if s.ch == '\n' {
			s.pos.Line++
			s.pos.Column = 1
		}
		s.ch = eof
	}
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func (s *Scanner) scanIdentifier() string {
	offs := s.pos.Offset
	for isLetter(s.ch) || isDigit(s.ch) {
		s.next()
	}
	return string(s.src[offs:s.pos.Offset])
}

func (s *Scanner) scanString(quotePos token.Pos) string {
	for {
		ch := s.ch
		if ch == '\n' || ch < 0 {
			s.error(quotePos, "string literal not terminated")
			break
		}
		s.next()
		if ch == '"' {
			break
		}
		if ch == '\\' {
			s.scanEscape('"')
		}
	}
	return string(s.src[quotePos.Offset:s.pos.Offset])
}

// scanEscape parses an escape sequence where rune is the accepted
// escaped quote. In case of a syntax error, it stops at the offending
// character (without consuming it).
func (s *Scanner) scanEscape(quote rune) {
	pos := s.pos

	var n int
	var base, max uint32
	switch s.ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		s.next()
		return
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		s.next()
		n, base, max = 2, 16, 255
	case 'u':
		s.next()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		s.next()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		msg := "unknown escape sequence"
		if s.ch < 0 {
			msg = "escape sequence not terminated"
		}
		s.error(pos, msg)
		return
	}

	var x uint32
	for n > 0 {
		d := uint32(digitVal(s.ch))
		if d >= base {
			msg := fmt.Sprintf("illegal character %#U in escape sequence", s.ch)
			if s.ch < 0 {
				msg = "escape sequence not terminated"
			}
			s.error(s.pos, msg)
			return
		}
		x = x*base + d
		s.next()
		n--
	}
	if x > max || 0xD800 <= x && x < 0xE000 {
		s.error(pos, "escape sequence is invalid Unicode code point")
	}
}

// Scan scans the next token and returns its position, type and literal.
func (s *Scanner) Scan() (pos token.Pos, typ token.Type, lit string) {
	s.skipWhitespace()
	pos = s.pos

	// determine token value
	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		if lit == "e" || lit == "ε" {
			typ = token.EPSILON
		} else {
			typ = token.IDENT
		}
	default:
		s.next() // always make progress
		switch ch {
		case eof:
			typ = token.EOF
		case '"':
			typ = token.STRING
			lit = s.scanString(pos)
		case '.':
			typ = token.PERIOD
		case '|':
			typ = token.PIPE
		case '→':
			typ = token.ARROW
			lit = "→"
		case '-':
			if s.ch == '>' {
				s.next()
				typ = token.ARROW
				lit = "->"
			}
		default:
			// next reports unexpected BOMs - don't repeat
			if ch != bom {
				s.error(pos, fmt.Sprintf("illegal character %#U", ch))
			}
			typ = token.ILLEGAL
			lit = string(ch)
		}
	}
	return
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || ch >= 0x80 && unicode.IsDigit(ch)
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= 0x80 && unicode.IsLetter(ch)
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16 // larger than any legal digit val
}
