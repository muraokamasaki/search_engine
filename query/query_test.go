package query

import (
	"testing"
)

func TestParseInfix(t *testing.T) {
	pairs := []struct{
		expr string
		token []string
	}{
		{"A || B", []string{"A", "||", "B"}},
		{"A&&B", []string{"A", "&&", "B"}},
		{"A&& B &&C||D", []string{"A", "&&", "B", "&&", "C", "||", "D"}},
		{"A&&B&&C&&D", []string{"A", "&&", "B", "&&", "C", "&&", "D"}},
		{"A||B|| C ||D", []string{"A", "||", "B", "||", "C", "||", "D"}},
		{" A &&&&C|| D ", []string{"A", "&&", "", "&&", "C", "||", "D"}},
		{"A||B&&||", []string{"A", "||", "B", "&&", "", "||", ""}},

	}
	for _, pair := range pairs {
		tok := parseInfix(pair.expr)
		if len(tok) == len(pair.token) {
			for i := range tok {
				if tok[i] != pair.token[i] {
					t.Errorf("Wrong token: Got %v, Wanted %v.", tok, pair.token)
				}
			}
		} else {
			t.Errorf("Different number of token: Got %v, Wanted %v.", tok, pair.token)
		}
	}
}

func TestShuntingYard(t *testing.T) {
	pairs := []struct{
		token []string
		output []string
	}{
		{[]string{"A", "||", "B"}, []string{"A", "B", "||"}},
		{[]string{"A", "||", "B", "&&", "C"}, []string{"A", "B", "C", "&&", "||"}},
		{[]string{"A", "||", "B", "&&", "C", "||", "D"}, []string{"A", "B", "C", "&&", "||", "D", "||"}},
		{[]string{"A", "||", "B", "&&", "", "&&", "C"}, []string{"A", "B", "", "&&", "C", "&&", "||"}},

	}
	for _, pair := range pairs {
		out := shuntingYard(pair.token)
		if len(out) == len(pair.output) {
			for i := range out {
				if out[i] != pair.output[i] {
					t.Errorf("Wrong token: Got %v, Wanted %v.", out, pair.output)
				}
			}
		} else {
			t.Errorf("Different number of token: Got %v, Wanted %v.", out, pair.output)
		}
	}
}

func TestSplitAndTrim(t *testing.T) {
	pairs := []struct{
		str string
		output []string
	}{
		{"Hello && World", []string{"Hello", "World"}},
		{"Hello, Goodbye && Tomorrow && Time ", []string{"Hello, Goodbye", "Tomorrow", "Time"}},
	}
	for _, pair := range pairs {
		out := splitAndTrim(pair.str, "&&")
		if len(out) == len(pair.output) {
			for i := range out {
				if out[i] != pair.output[i] {
					t.Errorf("Wrong token: Got %v, Wanted %v.", out, pair.output)
				}
			}
		} else {
			t.Errorf("Different number of token: Got %v, Wanted %v.", out, pair.output)
		}
	}
}

func SetUpSearcher() (s *Searcher) {
	s = NewSearcher(3)
	s.ii.BuildFromTextFile("example.txt")
	s.ki.BuildFromTextFile("example.txt")
	return
}

func TestSearcher_TermQuery(t *testing.T) {
	pairs := []struct{
		query string
		results []int
	}{
		{"is a statistic", []int{1}},
		{"language", []int{2}},
		{"is", []int{1, 2, 3}},

	}
	s := SetUpSearcher()
	for _, pair := range pairs {
		res := s.TermQuery(pair.query)
		if len(res) == len(pair.results) {
			for i := range res {
				if res[i] != pair.results[i] {
					t.Errorf("Wrong id: Got %v, Wanted %v.", res, pair.results)
				}
			}
		} else {
			t.Errorf("Different number of results: Got %v, Wanted %v.", res, pair.results)
		}
	}
}

func TestSearcher_BooleanQuery(t *testing.T) {
	pairs := []struct{
		query string
		results []int
	}{
		{"statistic && coefficient && items", []int{1}},
		{"reliability || technologies", []int{1, 3}},
		{"qualitative || semantics && reliability || technologies", []int{1, 3}},
		{"|| technique && language && processing", []int{2}},


	}
	s := SetUpSearcher()
	for _, pair := range pairs {
		res, err := s.BooleanQuery(pair.query)
		if err != nil {
			t.Error(err)
		}
		if len(res) == len(pair.results) {
			for i := range res {
				if res[i] != pair.results[i] {
					t.Errorf("Wrong id: Got %v, Wanted %v.", res, pair.results)
				}
			}
		} else {
			t.Errorf("Different number of results: Got %v, Wanted %v.", res, pair.results)
		}
	}
}