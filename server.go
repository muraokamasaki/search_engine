package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type SERP struct {
	Query string
	Page int
	Results []Document
	SearchAlgorithms []queryFunc
}

func (s Searcher) mapNameToFunc(funcName string) queryFunc {
	funcMap := map[string]queryFunc{
		"BM25": s.BM25Query,
		"Classic TF-IDF": s.VectorSpaceQuery,
		"Boolean": s.BooleanQuery,
		"Terms": s.TermsQuery,
		"Fuzzy": s.FuzzyQuery,
		"Wildcard": s.WildcardQuery,
	}
	return funcMap[funcName]
}

func (s Searcher) queryHandler(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query().Get("q")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	searchAlgorithm := r.URL.Query().Get("alg")
	if searchAlgorithm == "" {
		searchAlgorithm = "BM25"  // Defaults to BM25
	}
	res := s.Query(queryString, s.mapNameToFunc(searchAlgorithm)) // js injections?
	resultPage := &SERP{
		Query: queryString,
		Page:      page,
		Results:   res,
		SearchAlgorithms: []queryFunc{s.BM25Query},
	}

	t, err := template.ParseFiles("templates/main.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, resultPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
}

func RunServer() {
	s := NewSearcher(3)
	s.BuildFromCSV("example.csv")
	http.HandleFunc("/", s.queryHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}