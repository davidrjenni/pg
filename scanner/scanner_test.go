// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner_test

import (
	"io/ioutil"
	"testing"

	"github.com/davidrjenni/pg/scanner"
	"github.com/davidrjenni/pg/token"
)

func TestScan(t *testing.T) {
	tokens := []struct {
		tok token.Type
		lit string
	}{
		{token.IDENT, "foobar"},
		{token.IDENT, "a۰۱۸"},
		{token.IDENT, "foo६४"},
		{token.IDENT, "bar９８７６"},
		{token.IDENT, "ŝ"},
		{token.IDENT, "ŝfoo"},
		{token.STRING, `"foobar"`},
		{token.STRING, `"\r"`},
		{token.STRING, `"foo\r\nbar"`},
		{token.ARROW, "→"},
		{token.ARROW, "->"},
		{token.PERIOD, "."},
		{token.PIPE, "|"},
		{token.EPSILON, "ε"},
		{token.EPSILON, "e"},
	}

	const (
		filename    = "scan_test"
		whitespaces = "  \t  \n\n\n"
	)

	epos := token.Pos{
		Filename: filename,
		Offset:   0,
		Line:     1,
		Column:   1,
	}

	var src []byte
	for _, t := range tokens {
		src = append(src, t.lit...)
		src = append(src, whitespaces...)
	}

	s := scanner.New(src, filename)
	s.Err = func(_ token.Pos, msg string) {
		t.Errorf("error handler called (msg = %s)", msg)
	}

	for i, tt := range tokens {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		checkPos(t, i, pos, epos)
		if tok != tt.tok {
			t.Errorf("%d: got token %v, want %v", i, tok, tt.tok)
		}
		if tok == token.IDENT || tok == token.STRING {
			if lit != tt.lit {
				t.Errorf("%d: got literal %q, want %q", i, lit, tt.lit)
			}
		}
		epos.Offset += len(tt.lit) + len(whitespaces)
		epos.Line += 3
	}
}

func checkPos(t *testing.T, i int, pos, epos token.Pos) {
	if pos.Filename != epos.Filename {
		t.Errorf("%d: got filename %q, want %q", i, pos.Filename, epos.Filename)
	}
	if pos.Offset != epos.Offset {
		t.Errorf("%d: got offset %v, want %v", i, pos.Offset, epos.Offset)
	}
	if pos.Column != epos.Column {
		t.Errorf("%d: got column %v, want %v", i, pos.Column, epos.Column)
	}
	if pos.Line != epos.Line {
		t.Errorf("%d: got line %v, want %v", i, pos.Line, epos.Line)
	}
}

func TestScanErrors(t *testing.T) {
	errors := []struct {
		src string
		tok token.Type
		col int
		lit string
		err string
	}{
		{"\a", token.ILLEGAL, 1, "", "illegal character U+0007"},
		{"1", token.ILLEGAL, 1, "", "illegal character U+0031 '1'"},
		{`#`, token.ILLEGAL, 1, "", "illegal character U+0023 '#'"},
		{`…`, token.ILLEGAL, 1, "", "illegal character U+2026 '…'"},
		{`"abc`, token.STRING, 1, `"abc`, "string literal not terminated"},
		{"\"abc\n", token.STRING, 1, `"abc`, "string literal not terminated"},
		{"\"abc\n   ", token.STRING, 1, `"abc`, "string literal not terminated"},
		{`"`, token.STRING, 1, `"`, "string literal not terminated"},
		{"\"abc\x00def\"", token.STRING, 5, "\"abc\x00def\"", "illegal character NUL"},
		{"\"abc\x80def\"", token.STRING, 5, "\"abc\x80def\"", "illegal UTF-8 encoding"},
		{"\ufeff\ufeff", token.ILLEGAL, 4, "\ufeff\ufeff", "illegal byte order mark"},        // only first BOM is ignored
		{"\"abc\ufeffdef\"", token.STRING, 5, "\"abc\ufeffdef\"", "illegal byte order mark"}, // only first BOM is ignored
	}

	for i, e := range errors {
		s := scanner.New([]byte(e.src), "error")
		s.Err = func(pos token.Pos, msg string) {
			if pos.Column != e.col {
				t.Errorf("%d: got column %v, want %v", i, pos.Column, e.col)
			}
			if msg != e.err {
				t.Errorf("%d: got error %q, want %q", i, msg, e.err)
			}
		}

		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok != e.tok {
			t.Errorf("%d: got token %v, want %v", i, tok, e.tok)
		}
		if e.tok != token.ILLEGAL && lit != e.lit {
			t.Errorf("%d: got literal %q, want %q", i, lit, e.lit)
		}
		if s.ErrCount != 1 {
			t.Errorf("got error count %v, want 1", len(errors))
		}
	}
}

func BenchmarkScan(b *testing.B) {
	b.StopTimer()
	const filename = "scanner.go"
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	b.SetBytes(int64(len(src)))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := scanner.New(src, filename)
		for {
			_, tok, _ := s.Scan()
			if tok == token.EOF {
				break
			}
		}
	}
}
