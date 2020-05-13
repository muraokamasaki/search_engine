package main

import (
	"sort"
	"strings"
)

// The Searcher type is an implementation of a search engine.
type Searcher struct {
	ii InvertedIndex
	ki KGramIndex
	docLen DocumentLengths
	storage DocumentStorage
}

func NewSearcher(k int, storage DocumentStorage) *Searcher {
	return &Searcher{ii: *NewInvertedIndex(), ki: *NewKGramIndex(k), docLen: DocumentLengths{}, storage:storage}
}

// queryFunc defines methods that take in a query string and
// returns a list of document IDs that are relevant to the query.
type queryFunc func(string) []int

// Query returns a list of documents that are relevant to the query,
// where relevance is defined by the given queryFunc.
func (s *Searcher) Query(query string, fn queryFunc) []Document {
	resultIDs := fn(query)
	return s.storage.Get(resultIDs)
}

// Query Methods

// TermsQuery returns documents that contain an exact match of
// all of the words in the query.
func (s *Searcher) TermsQuery(query string) (results []int) {
	results = s.ii.Intersect(tokenize(query))
	return
}

// BooleanQuery returns documents based on the boolean retrieval model.
// Only supports AND (&&) and OR (||).
func (s *Searcher) BooleanQuery(query string) (results []int) {
	// Allows boolean operations between terms. Terms should only consist of a single word.
	var queryTerms []string
	intersectFlag := strings.Contains(query, "&&")
	unionFlag := strings.Contains(query, "||")
	// If the query does not contain both && or ||, we do not need to create postfix tree.
	if intersectFlag && unionFlag {
		queryTerms = shuntingYard(parseInfix(query))
		var stack [][]int
		for i := 0; i < len(queryTerms); i++ {
			if queryTerms[i] == "&&" {
				if len(stack) >= 2 {
					stack[len(stack)-2] = IntersectPosting(stack[len(stack)-1], stack[len(stack)-2])
					stack = stack[:len(stack)-1]
				} else {
					stack = [][]int{}
					break
				}
			} else if queryTerms[i] == "||" {
				if len(stack) >= 2 {
					stack[len(stack)-2] = UnionPosting(stack[len(stack)-1], stack[len(stack)-2])
					stack = stack[:len(stack)-1]
				} else {
					stack = [][]int{}
					break
				}
			}  else {
				stack = append(stack, s.ii.PostingsList(queryTerms[i]))
			}
		}
		if len(stack) == 1 {
			results = stack[0]
		}
	} else if unionFlag {
		results = s.ii.Union(splitTrimToLower(query, "||"))
	} else {
		results = s.ii.Intersect(splitTrimToLower(query, "&&"))
	}
	return
}

// parseInfix parses the expression and splits it into a set of tokens,
// where each token is either a word or an operator.
func parseInfix(expr string) (output []string) {
	for _, i := range splitTrimToLower(expr, "&&") {
		for _, j := range splitTrimToLower(i, "||") {
			output = append(output, j, "||")
		}
		if len(output) > 0 {
			output[len(output) - 1] = "&&"
		}
	}
	output = output[:len(output) - 1]
	return
}

// shuntingYard converts an expression given in infix notation to a
// postfix "Reverse Polish" notation.
// This is a partial implementation of Shunting-yard algorithm,
// which only considers left associative, binary operators.
func shuntingYard(tokens []string) (output []string) {
	// orderOfOperations is a map containing operations used in our BooleanQuery.
	// Larger value implies higher precedence; Values are arbitrary.
	orderOfOperations := map[string]int{"||": 1, "&&": 2}
	var operatorStack []string
	for _, token := range tokens {
		if _, ok := orderOfOperations[token]; !ok {
			output = append(output, token)
		} else {
			for i := len(operatorStack) - 1; i>= 0; i-- {
				// If equal, assumes token is left associative.
				if orderOfOperations[operatorStack[i]] >= orderOfOperations[token] {
					output = append(output, operatorStack[i])
					operatorStack = operatorStack[:i]
				} else {
					break
				}
			}
			operatorStack = append(operatorStack, token)
		}
	}
	for i := len(operatorStack) - 1; i>= 0; i-- {
		output = append(output, operatorStack[i])
	}
	return
}

// splitTrimToLower splits the input string with the provided token,
// trims whitespace and converts the output tokens to lowercase.
func splitTrimToLower(str string, split string) (out []string) {
	out = strings.Split(str, split)
	for i := range out {
		out[i] = strings.TrimSpace(strings.ToLower(out[i]))
	}
	return
}

// FuzzyQuery returns documents that contain all of the provided terms.
// Each term also accepts other terms that are within a certain edit distance.
// For example "Fizzy" will match the query "Fuzzy".
func (s *Searcher) FuzzyQuery(query string) (results []int) {
	for _, queryTerm := range tokenize(query) {
		fuzziness := getFuzziness(queryTerm)
		terms := s.ki.GetCloseTerms(queryTerm, fuzziness)

		if len(results) == 0 {
			results = s.ii.Union(terms)
		} else {
			results = IntersectPosting(results, s.ii.Union(terms))
		}
	}
	return
}

// getFuzziness determines the edit distance for each term
// based on its length.
// Longer words are allowed more spelling mistakes.
func getFuzziness(str string) (fuzziness int) {
	if len(str) <= 2 {
		fuzziness = 0
	} else if len(str) <= 5 {
		fuzziness = 1
	} else {
		fuzziness = 2
	}
	return
}

// WildcardQuery returns documents that contain all of the provided terms.
// Terms can contain the characters '?' which represents a single character,
// and '*' which can be expanded into one or more characters.
func (s Searcher) WildcardQuery(query string) (results []int) {
	for _, queryTerm := range tokenizeWildcard(query) {
		terms := s.ki.KGramMatch(queryTerm)
		var partialResult  []string
		for _, res := range terms {
			if wildcardMatch(queryTerm, res) {
				partialResult = append(partialResult, res)
			}
		}
		if len(results) == 0 {
			results = s.ii.Union(partialResult)
		} else {
			results = IntersectPosting(results, s.ii.Union(partialResult))
		}
	}
	return
}

// ScoringList stores a id, score pair.
// Implements sort.Interface for sorting by descending score.
type ScoringList struct {
	ids []int
	scores []float64
}

func (r ScoringList) Len() int { return len(r.ids) }
func (r ScoringList) Swap(i, j int) {
	r.scores[i], r.scores[j] = r.scores[j], r.scores[i]
	r.ids[i], r.ids[j] = r.ids[j], r.ids[i]
}
func (r ScoringList) Less(i, j int) bool { return r.scores[i] > r.scores[j] }

// VectorSpaceQuery returns a ranked list of results sorted by
// cosine similarity using the vector space model.
// Scores are calculated using tf-idf and document length normalization.
func (s Searcher) VectorSpaceQuery(query string) (results []int) {
	resList := &ScoringList{}
	for _, queryTerm := range tokenize(query) {
		for _, docID := range s.ii.PostingsList(queryTerm) {
			resultsIndex := findIndexInArray(resList.ids, docID)
			// Calculate tf-idf score
			score := float64(s.ii.TermFrequency(queryTerm, docID)) * s.ii.InverseDocumentFrequency(queryTerm)
			if resultsIndex == -1 {
				// Document ID not yet in results.
				resList.ids = append(resList.ids, docID)
				resList.scores = append(resList.scores, score)
			} else {
				resList.scores[resultsIndex] += score
			}
		}
	}
	for i := range resList.ids {
		resList.scores[i] /= float64(s.docLen.docLength(resList.ids[i]))
	}
	sort.Sort(resList)
	results = resList.ids
	return
}

// findIndexInArray returns the index of the given int
// if it exists in the array, else returns -1.
func findIndexInArray(arr []int, value int) int {
	for i, v := range arr {
		if v == value {
			return i
		}
	}
	return -1
}

// BM25Query returns a ranked list of results scored by the Okapi BM25 algorithm.
// Scores are calculated using tf-idf and document length normalization.
func (s Searcher) BM25Query(query string) (results []int) {
	k1 := 0.9
	b := 0.4
	resList := &ScoringList{}
	for _, queryTerm := range tokenize(query) {
		for _, docID := range s.ii.PostingsList(queryTerm) {
			resultsIndex := findIndexInArray(resList.ids, docID)
			// Calculate BM25 score
			tf := float64(s.ii.TermFrequency(queryTerm, docID))
			idf := s.ii.InverseDocumentFrequency(queryTerm)
			score :=  idf * (k1 + 1) * tf / (k1 * ((1 - b) + b * (float64(s.docLen.docLength(docID)) / s.docLen.averageDocumentLength())) + tf)
			if resultsIndex == -1 {
				// Document ID not yet in results.
				resList.ids = append(resList.ids, docID)
				resList.scores = append(resList.scores, score)
			} else {
				resList.scores[resultsIndex] += score
			}
		}
	}
	sort.Sort(resList)
	results = resList.ids
	return
}

// BuildIndices builds the inverted index and k-gram index
// from the document storage.
func (s *Searcher) BuildIndices() {
	s.storage.Apply(func(doc Document) {
		// Only take word count of Body.
		s.docLen.addDocumentLength(doc.Body)
		// Adds words in Title and Body to index.
		for _, token := range tokenize(doc.Title) {
			s.ii.addIDToPostingsList(token, doc.id)
			s.ki.addWordToPostingsList(token)
		}
		for _, token := range tokenize(doc.Body) {
			s.ii.addIDToPostingsList(token, doc.id)
			s.ki.addWordToPostingsList(token)
		}
	})
}