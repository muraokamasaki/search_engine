package main

import (
	"sort"
	"testing"
)

func TestTokenize(t *testing.T) {
	pairs := []struct {
		str   string
		token []string
	}{
		{"Test string.", []string{"test", "string"}},
		{"I'm 23 years old.", []string{"i", "m", "23", "years", "old"}},
		{"3d!e-fg.", []string{"3d", "e", "fg"}},
	}
	for _, pair := range pairs {
		tok := Tokenize(pair.str)
		if len(tok) == len(pair.token) {
			for i := range tok {
				if tok[i] != pair.token[i] {
					t.Errorf("Wrong token: Got %s, Wanted %s.", tok[i], pair.token[i])
				}
			}
		} else {
			t.Errorf("Different number of token: Got %d, Wanted %d.", len(tok), len(pair.token))
		}
	}
}

func TestTokenizeWildcard(t *testing.T) {
	pairs := []struct {
		str string
		token []string
	}{
		{"Test string.", []string{"test", "string"}},
		{"W?ld*rd.", []string{"w?ld*rd"}},
		{"*me ?? *.", []string{"*me", "??", "*"}},
	}
	for _, pair := range pairs {
		tok := TokenizeWildcard(pair.str)
		if len(tok) == len(pair.token) {
			for i := range tok {
				if tok[i] != pair.token[i] {
					t.Errorf("Wrong token: Got %s, Wanted %s.", tok[i], pair.token[i])
				}
			}
		} else {
			t.Errorf("Different number of token: Got %d, Wanted %d.", len(tok), len(pair.token))
		}
	}
}

func TestMin(t *testing.T) {
	pairs := []struct{
		nums []int
		answer int
	}{
		{[]int{4,3,4,5,6}, 3},
		{[]int{5,4,2,9,2,3,2,9}, 2},
		{[]int{5,4,2,0,2}, 0},
		{[]int{-1,-4,-3,0,-2}, -4},
	}
	for _, pair := range pairs {
		ans := min(pair.nums...)
		if pair.answer != ans {
			t.Errorf("Wrong answer: Got %d, Wanted %d.", ans, pair.answer)
		}
	}
}

func TestMax(t *testing.T) {
	pairs := []struct{
		nums []int
		answer int
	}{
		{[]int{4,3,4,5,6}, 6},
		{[]int{5,4,2,9,2,3,2,9}, 9},
		{[]int{-1,-4,-3,0,-2}, 0},
		{[]int{-3,-4,-1,-9,-2}, -1},
	}
	for _, pair := range pairs {
		ans := max(pair.nums...)
		if pair.answer != ans {
			t.Errorf("Wrong answer: Got %d, Wanted %d.", ans, pair.answer)
		}
	}
}

func TestEditDistance(t *testing.T) {
	pairs := []struct{
		str []string
		answer int
	}{
		{[]string{"fast", "cats"}, 3},
		{[]string{"gopher", "python"}, 5},
		{[]string{"hello", ""}, 5},
		{[]string{"", "world"}, 5},
	}
	for _, pair := range pairs {
		ans := EditDistance(pair.str[0], pair.str[1])
		if pair.answer != ans {
			t.Errorf("Wrong answer: Got %d, Wanted %d.", ans, pair.answer)
		}
	}
}

func TestWildcardMatch(t *testing.T) {
	pairs := []struct{
		str []string
		answer bool
	}{
		{[]string{"time", "time"}, true},
		{[]string{"tome", "time"}, false},
		{[]string{"t?me", "time"}, true},
		{[]string{"t?e", "time"}, false},
		{[]string{"?ime", "time"}, true},
		{[]string{"t*e", "time"}, true},
		{[]string{"t*", "time"}, true},
		{[]string{"*e", "time"}, true},
		{[]string{"t*er", "time"}, false},
		{[]string{"*m*", "time"}, true},
		{[]string{"*m?", "time"}, true},
	}
	for _, pair := range pairs {
		ans := WildcardMatch(pair.str[0], pair.str[1])
		if pair.answer != ans {
			t.Errorf("Pattern: %s, String: %s. Got %v, Wanted %v.", pair.str[0], pair.str[1], ans, pair.answer)
		}
	}
}

func SetUpInvertedIndex() (ii *InvertedIndex) {
	ii = NewInvertedIndex()
	ii.addIDToPostingsList("hello", 1)
	ii.addIDToPostingsList("hello", 2)
	ii.addIDToPostingsList("world", 1)
	ii.addIDToPostingsList("world", 3)
	return
}

func TestInvertedIndex_Intersect(t *testing.T) {
	ii := SetUpInvertedIndex()
	ans := ii.Intersect([]string{"hello", "world"})
	if len(ans) != 1 || ans[0] != 1 {
		t.Errorf("Wrong slice: Got %v, Wanted []int{1}.", ans)
	}
}

func TestIntersectPair(t *testing.T) {
	pairs := []struct{
		plist1 []int
		plist2 []int
		result []int
	}{
		{[]int{1, 2, 3, 4, 5, 6}, []int{2, 4, 6}, []int{2, 4, 6}},
		{[]int{2, 3, 5}, []int{1, 2, 3, 4, 5, 6}, []int{2, 3, 5}},
		{[]int{1, 2, 3}, []int{}, []int{}},
		{[]int{1, 2, 3}, []int{4, 5, 6}, []int{}},
		{[]int{1}, []int{1}, []int{1}},
		{[]int{}, []int{1, 2, 3}, []int{}},
	}
	for _, pair := range pairs {
		res := IntersectPosting(pair.plist1, pair.plist2)
		if len(res) != len(pair.result) {
			t.Errorf("Wrong number of documents: Got %v, Wanted %v.", res, pair.result)
		}
		for i := range res {
			if res[i] != pair.result[i] {
				t.Errorf("Wrong k-grams: Got %v, Wanted %v", res, pair.result)
			}
		}
	}
}

func TestUnionPair(t *testing.T) {
	pairs := []struct{
		plist1 []int
		plist2 []int
		result []int
	}{
		{[]int{1, 3, 5}, []int{2, 4, 6}, []int{1, 2, 3, 4, 5, 6}},
		{[]int{1, 2, 3, 4, 5, 6}, []int{2, 4, 6}, []int{1, 2, 3, 4, 5, 6}},
		{[]int{2, 3, 5}, []int{1, 2, 3, 4, 5, 6}, []int{1, 2, 3, 4, 5, 6}},
		{[]int{}, []int{}, []int{}},
		{[]int{1}, []int{2}, []int{1, 2}},

	}
	for _, pair := range pairs {
		res := UnionPosting(pair.plist1, pair.plist2)
		if len(res) != len(pair.result) {
			t.Errorf("Wrong number of documents: Got %v, Wanted %v.", res, pair.result)
		}
		for i := range res {
			if res[i] != pair.result[i] {
				t.Errorf("Wrong k-grams: Got %v, Wanted %v", res, pair.result)
			}
		}
	}
}

func TestInvertedIndex_Union(t *testing.T) {
	ii := SetUpInvertedIndex()
	ans := ii.Union([]string{"hello", "world"})
	sort.Ints(ans)
	if len(ans) != 3 || ans[0] != 1 || ans[1] != 2 || ans[2] != 3 {
		t.Errorf("Wrong slice: Got %v, Wanted []int{1,2,3}.", ans)
	}
}

func SetUpKGramIndex(k int) (ki *KGramIndex) {
	ki = NewKGramIndex(k)
	ki.addWordToPostingsList("hello")
	ki.addWordToPostingsList("helicopter")
	ki.addWordToPostingsList("man")
	return
}

func TestBuildKGrams(t *testing.T) {
	pairs := []struct {
		str string
		k int
		grams []string
	}{
		{"hello", 3, []string{"$$h", "$he", "hel", "ell", "llo", "lo$", "o$$"}},
		{"hi", 3, []string{"$$h", "$hi", "hi$", "i$$"}},
		{"i", 3, []string{"i"}},

	}
	for _, pair := range pairs {
		gr := BuildKGrams(pair.str, pair.k)
		if len(gr) == len(pair.grams) {
			sort.Strings(gr)
			sort.Strings(pair.grams)
			for i := range gr {
				if gr[i] != pair.grams[i] {
					t.Errorf("Wrong token: Got %s, Wanted %s.", gr[i], pair.grams[i])
				}
			}
		} else {
			t.Errorf("Different number of token: Got %d, Wanted %d.", len(gr), len(pair.grams))
		}
	}
}

func TestKGramIndex_KGramOverlap(t *testing.T) {
	ki := SetUpKGramIndex(3)
	pairs := []struct{
		query string
		count map[string]int
	}{
		{"hello", map[string]int{"hello": 7, "helicopter": 3}},
		{"help", map[string]int{"hello": 3, "helicopter": 3}},
		{"man", map[string]int{"man": 5}},
		{"an", map[string]int{"man": 2}},
		{"a", map[string]int{}},
	}
	for _, pair := range pairs {
		c := ki.KGramOverlap(pair.query)
		if len(c) != len(pair.count) {
			t.Errorf("Wrong number of documents: Got %v, Wanted %v.", c, pair.count)
		}
		for k, v := range c {
			if pair.count[k] != v {
				t.Errorf("Wrong k-grams: Got %v, Wanted %v", c, pair.count)
			}
		}
	}
}

func TestKGramIndex_KGramMatch(t *testing.T) {
	ki := SetUpKGramIndex(3)
	pairs := []struct{
		q string
		terms []string
	}{
		{"he*", []string{"hello", "helicopter"}},
		{"hell*", []string{"hello"}},
		{"m?n", []string{"man"}},
		{"*n", []string{"man"}},
	}
	for _, pair := range pairs {
		c := ki.KGramMatch(pair.q)
		if len(c) != len(pair.terms) {
			t.Errorf("Wrong number of documents: Got %v, Wanted %v.", c, pair.terms)
		}
		sort.Strings(c)
		sort.Strings(pair.terms)
		for k, v := range c {
			if pair.terms[k] != v {
				t.Errorf("Wrong k-grams: Got %v, Wanted %v", c, pair.terms)
			}
		}
	}
}