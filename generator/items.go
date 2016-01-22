// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package generator

// item represents an LR(0) item.
type item struct {
	n   int // number of the production
	dot int // position of the dot
}

func newItem(dot, n int) item {
	return item{dot: dot, n: n}
}

type itemSet []item

func (set itemSet) contains(item item) bool {
	for _, i := range set {
		if i.dot == item.dot && i.n == item.n {
			return true
		}
	}
	return false
}

func (set itemSet) equal(items itemSet) bool {
	if len(set) != len(items) {
		return false
	}
	for _, i := range set {
		if !items.contains(i) {
			return false
		}
	}
	return true
}

func containsSet(itemSets []itemSet, items itemSet) bool {
	for _, s := range itemSets {
		if items.equal(s) {
			return true
		}
	}
	return false
}
