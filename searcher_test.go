package main

import (
	"testing"
)

func TestParseInfix(t *testing.T) {
	pairs := []struct{
		expr string
		token []string
	}{
		{"A || B", []string{"a", "||", "b"}},
		{"A&&B", []string{"a", "&&", "b"}},
		{"A&& B &&C||D", []string{"a", "&&", "b", "&&", "c", "||", "d"}},
		{"A&&B&&C&&D", []string{"a", "&&", "b", "&&", "c", "&&", "d"}},
		{"A||B|| C ||D", []string{"a", "||", "b", "||", "c", "||", "d"}},
		{" A &&&&C|| D ", []string{"a", "&&", "", "&&", "c", "||", "d"}},
		{"A||B&&||", []string{"a", "||", "b", "&&", "", "||", ""}},

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

func TestSplitTrimToLower(t *testing.T) {
	pairs := []struct{
		str string
		output []string
	}{
		{"Hello && World", []string{"hello", "world"}},
		{"Hello, Goodbye && Tomorrow && Time ", []string{"hello, goodbye", "tomorrow", "time"}},
	}
	for _, pair := range pairs {
		out := splitTrimToLower(pair.str, "&&")
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
	s.BuildFromTextFile("example.txt")
	return
}

func TestSearcher_TermQuery(t *testing.T) {
	pairs := []struct{
		query string
		results []int
	}{
		{"is a statistic", []int{0}},
		{"language", []int{1}},
		{"is", []int{0, 1, 2}},

	}
	s := SetUpSearcher()
	for _, pair := range pairs {
		res := s.TermsQuery(pair.query)
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
		{"", []int{}},
		{"statistic && coefficient", []int{0}},
		{"statistic && coefficient && items", []int{0}},
		{"sTatistic && coeffIcient &&items", []int{0}},
		{"reliability || technologies", []int{0, 2}},
		{"qualitative || semantics && reliability || technologies", []int{0, 2}},
		{"|| technique && language && processing", []int{1}},
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

func TestSearcher_WildcardQuery(t *testing.T) {
	pairs := []struct{
		query string
		results []int
	}{
		{"cohe*", []int{0}},
		{"ch?ce", []int{}},
		{"ch?nc?", []int{0}},
		{"sem*t*c", []int{1}},
		{"sem*ts*c", []int{}},
		{"con*s related", []int{1}},
	}
	s := SetUpSearcher()
	for _, pair := range pairs {
		res := s.WildcardQuery(pair.query)
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

func TestSearcher_FuzzyQuery(t *testing.T) {
	pairs := []struct{
		query string
		results []int
	}{
		{"cohdn", []int{0}},
		{"latent semantic", []int{1}},
		{"by various radi communication techologies", []int{2}},
		{"i", []int{}},
	}
	s := SetUpSearcher()
	for _, pair := range pairs {
		res := s.FuzzyQuery(pair.query)
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

func TestSearcher_VectorSpaceQuery(t *testing.T) {
	pairs := []struct{
		query string
		results []int
	}{
		{"cohen", []int{0}},
		{"latent semantic", []int{1}},
		{"statistic that", []int{0, 1}},
		{"matrix communication channel", []int{2, 1}},


	}
	s := SetUpSearcher()
	for _, pair := range pairs {
		res := s.VectorSpaceQuery(pair.query)
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

func TestSearcher_BM25Query(t *testing.T) {
	pairs := []struct{
		query string
		results []int
	}{
		{"cohen", []int{0}},
		{"latent semantic", []int{1}},
		{"statistic that", []int{0, 1}},
		{"matrix communication channel", []int{2, 1}},


	}
	s := SetUpSearcher()
	for _, pair := range pairs {
		res := s.BM25Query(pair.query, 1.2, 0.4)
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