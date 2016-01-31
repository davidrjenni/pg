package main

import "fmt"

type pgElem struct {
	sym	string
	state	int
}

type pgStack []pgElem

func (s pgStack) top() pgElem		{ return s[len(s)-1] }
func (s *pgStack) pop(n int)		{ *s = (*s)[:len(*s)-n] }
func (s *pgStack) push(e pgElem)	{ *s = append(*s, e) }

type pgNode struct {
	typ		string
	val		string
	children	[]pgNode
}

func pgParse() pgNode {
	var (
		table		= map[string][][2]int{"Expr'": [][2]int{[2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}}, "(": [][2]int{[2]int{1, 5}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{1, 5}, [2]int{1, 5}, [2]int{1, 5}, [2]int{1, 5}, [2]int{1, 5}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}}, ")": [][2]int{[2]int{3, 0}, [2]int{3, 0}, [2]int{2, 3}, [2]int{2, 6}, [2]int{2, 8}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{1, 15}, [2]int{2, 2}, [2]int{2, 1}, [2]int{2, 5}, [2]int{2, 4}, [2]int{2, 7}}, "Term": [][2]int{[2]int{4, 2}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{4, 2}, [2]int{4, 11}, [2]int{4, 12}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}}, "$": [][2]int{[2]int{3, 0}, [2]int{0, 0}, [2]int{2, 3}, [2]int{2, 6}, [2]int{2, 8}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{2, 2}, [2]int{2, 1}, [2]int{2, 5}, [2]int{2, 4}, [2]int{2, 7}}, "+": [][2]int{[2]int{3, 0}, [2]int{1, 7}, [2]int{2, 3}, [2]int{2, 6}, [2]int{2, 8}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{1, 7}, [2]int{2, 2}, [2]int{2, 1}, [2]int{2, 5}, [2]int{2, 4}, [2]int{2, 7}}, "*": [][2]int{[2]int{3, 0}, [2]int{3, 0}, [2]int{1, 9}, [2]int{2, 6}, [2]int{2, 8}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{1, 9}, [2]int{1, 9}, [2]int{2, 5}, [2]int{2, 4}, [2]int{2, 7}}, "Factor": [][2]int{[2]int{4, 3}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{4, 3}, [2]int{4, 3}, [2]int{4, 3}, [2]int{4, 13}, [2]int{4, 14}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}}, "-": [][2]int{[2]int{3, 0}, [2]int{1, 6}, [2]int{2, 3}, [2]int{2, 6}, [2]int{2, 8}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{1, 6}, [2]int{2, 2}, [2]int{2, 1}, [2]int{2, 5}, [2]int{2, 4}, [2]int{2, 7}}, "/": [][2]int{[2]int{3, 0}, [2]int{3, 0}, [2]int{1, 8}, [2]int{2, 6}, [2]int{2, 8}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{1, 8}, [2]int{1, 8}, [2]int{2, 5}, [2]int{2, 4}, [2]int{2, 7}}, "NUMBER": [][2]int{[2]int{1, 4}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{1, 4}, [2]int{1, 4}, [2]int{1, 4}, [2]int{1, 4}, [2]int{1, 4}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}}, "Expr": [][2]int{[2]int{4, 1}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{4, 10}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}, [2]int{3, 0}}}
		count		= []int{1, 3, 3, 1, 3, 3, 1, 3, 1}
		names		= []string{"Expr'", "Expr", "Expr", "Expr", "Term", "Term", "Term", "Factor", "Factor"}
		tree		= make([]pgNode, 0)
		stack		= &pgStack{pgElem{state: 0}}
		typ, tok	= pgLex()
	)

	for {
		s := stack.top()
		var column [][2]int

		if typ != "" {
			column = table[typ]
		} else {
			column = table[tok]
		}
		if column == nil {
			pgError(fmt.Errorf("unexpected token %q (type: %q)", tok, typ))
			typ, tok = pgLex()
			if tok == "$" {
				if len(tree) == 0 {
					return pgNode{typ: "error"}
				}
				return tree[0]
			}
			continue
		}
		entry := column[s.state]
		switch entry[0] {
		case 2:
			c := count[entry[1]]
			name := names[entry[1]]
			stack.pop(2 * c)
			s = stack.top()
			stack.push(pgElem{sym: name})
			stack.push(pgElem{state: table[name][s.state][1]})
			var rest []pgNode
			for _, n := range tree[:len(tree)-c] {
				rest = append(rest, n)
			}
			tree = append(rest, pgNode{typ: name, val: name, children: tree[len(tree)-c:]})
		case 1:
			stack.push(pgElem{sym: tok})
			stack.push(pgElem{state: entry[1]})
			tree = append(tree, pgNode{typ: typ, val: tok})
			typ, tok = pgLex()
		case 0:
			if tok == "$" {
				return tree[0]
			}
		default:
			if tok == "$" {
				pgError(fmt.Errorf("unexpected end of input"))
				if len(tree) == 0 {
					return pgNode{typ: "error"}
				}
				return tree[0]
			}
			pgError(fmt.Errorf("unexpected token %q (type: %q)", tok, typ))
			typ, tok = pgLex()
		}
	}
	return tree[0]
}
