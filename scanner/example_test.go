// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner_test

import (
	"fmt"

	"github.com/davidrjenni/pg/scanner"
	"github.com/davidrjenni/pg/token"
)

func ExampleScanner_Scan() {
	src := []byte(`E -> T "+" T | T | ε .`)
	s := scanner.New(src, "example")

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", pos, tok, lit)
	}

	// output:
	// example:1:1	IDENT	"E"
	// example:1:3	ARROW	"->"
	// example:1:6	IDENT	"T"
	// example:1:8	STRING	"\"+\""
	// example:1:12	IDENT	"T"
	// example:1:14	PIPE	""
	// example:1:16	IDENT	"T"
	// example:1:18	PIPE	""
	// example:1:20	EPSILON	"ε"
	// example:1:23	PERIOD	""
}
