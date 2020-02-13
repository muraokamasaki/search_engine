package query

import "github.com/muraokamasaki/search_engine/indices"

type Searcher struct {
	ii indices.InvertedIndex
	ki indices.KGramIndex
}

func NewSearcher(k int) *Searcher {
	return &Searcher{ii: *indices.NewInvertedIndex(), ki: *indices.NewKGramIndex(k)}
}
