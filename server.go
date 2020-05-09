package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const ResultsPerPage = 5

type SERP struct {
	Query string
	Page int
	Results []Document
	Algorithm string
	NextURL string
	PrevURL string
}

// Returns a slice of the results based on the page number.
func paginateResult(results []Document, page int) (resSlice []Document) {
	if len(results) >= (page - 1) * ResultsPerPage {
		resSlice = results[(page - 1) * ResultsPerPage : min(page * ResultsPerPage, len(results))]
	}
	return
}

// Creates a new URL from an existing URL with a different page number.
func changePageURL(u *url.URL, page int) string {
	u, _ = url.Parse(u.String())
	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()
	return u.String()
}

func (s *Searcher) mapNameToFunc(funcName string) (f queryFunc) {
	funcMap := map[string]queryFunc{
		"BM25": s.BM25Query,
		"Classic TF-IDF": s.VectorSpaceQuery,
		"Boolean": s.BooleanQuery,
		"Terms": s.TermsQuery,
		"Fuzzy": s.FuzzyQuery,
		"Wildcard": s.WildcardQuery,
	}
	f, ok := funcMap[funcName]
	if !ok {
		f = s.BM25Query  // Defaults to BM25
	}
	return
}

func (s *Searcher) queryHandler(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query().Get("q")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	searchAlgorithm := r.URL.Query().Get("alg")
	res := s.Query(queryString, s.mapNameToFunc(searchAlgorithm))
	resultSlice := paginateResult(res, page)

	// Create URLs for pagination.
	var nextURL, prevURL string
	if len(res) > page * ResultsPerPage {
		nextURL = changePageURL(r.URL, page + 1)
	} else {
		nextURL = "#"
	}
	if page > 1 && len(resultSlice) > 0 {
		prevURL = changePageURL(r.URL, page - 1)
	} else {
		prevURL = "#"
	}

	resultPage :=  &SERP{
		Query: queryString,
		Page:      page,
		Results:   resultSlice,
		Algorithm: searchAlgorithm,
		NextURL: nextURL,
		PrevURL : prevURL,
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

func RunServer(k int, store DocumentStorage) {
	s := NewSearcher(k, store)
	s.BuildIndices()
	http.HandleFunc("/", s.queryHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}