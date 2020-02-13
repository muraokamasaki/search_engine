package indices

import (
	"sort"
	"testing"
)

func TestTokenize(t *testing.T) {
	pairs := []struct {
		str string
		token []string
	}{
		{"Test string.", []string{"Test", "string"}},
		{"I'm 23 years old.", []string{"I", "m", "23", "years", "old"}},
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

func TestPrefixEditDistance(t *testing.T) {
	pairs := []struct{
		str []string
		answer int
	}{
		{[]string{"uni", "university"}, 0},
		{[]string{"tome", "tomorrow"}, 1},
		{[]string{"time", "tomorrow"}, 2},
		{[]string{"", "world"}, 0},
	}
	for _, pair := range pairs {
		ans := PrefixEditDistance(pair.str[0], pair.str[1], 2)
		if pair.answer != ans {
			t.Errorf("Wrong answer: Got %d, Wanted %d.", ans, pair.answer)
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
	ki.addWordToPostingsList("Hello world.")
	ki.addWordToPostingsList("Helicopter.")
	ki.addWordToPostingsList("Man-eating sharks!")
	return
}

func TestBuildKGrams(t *testing.T) {
	pairs := []struct {
		str string
		k int
		grams []string
	}{
		{"Te-st str.", 3, []string{"$$t", "$te", "tes", "est", "st ", "t s", " st", "str", "tr$", "r$$"}},
		{"I'm 90lbs", 3, []string{"$$i", "$im", "im ", "m 9", " 90", "90l", "0lb", "lbs", "bs$", "s$$"}},
		{"Hi!!", 3, []string{"hi"}},
	}
	for _, pair := range pairs {
		gr := buildKGrams(pair.str, pair.k)
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
		{"hello", map[string]int{"Hello world.": 5, "Helicopter.": 3}},
		{"Help", map[string]int{"Hello world.": 3, "Helicopter.": 3}},
		{"World", map[string]int{"Hello world.": 5}},
		{"Maneating", map[string]int{"Man-eating sharks!": 9}},
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