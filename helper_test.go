package main

import "testing"

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
		tok := tokenize(pair.str)
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
		tok := tokenizeWildcard(pair.str)
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
		ans := editDistance(pair.str[0], pair.str[1])
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
		ans := wildcardMatch(pair.str[0], pair.str[1])
		if pair.answer != ans {
			t.Errorf("Pattern: %s, String: %s. Got %v, Wanted %v.", pair.str[0], pair.str[1], ans, pair.answer)
		}
	}
}