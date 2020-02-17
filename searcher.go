package main

import (
	"fmt"
	"strings"
)

// The Searcher type represents a search engine.
type Searcher struct {
	ii InvertedIndex
	ki KGramIndex
}

func NewSearcher(k int) *Searcher {
	return &Searcher{ii: *NewInvertedIndex(), ki: *NewKGramIndex(k)}
}

// Query methods that our Searcher implements.

// Filters documents that contain all of the provided terms.
func (s Searcher) TermsQuery(query string) (results []int) {
	results = s.ii.Intersect(tokenize(query))
	return
}

// Filters documents based on the boolean retrieval model.
// Only supports AND (&&) and OR (||).
func (s Searcher) BooleanQuery(query string) (results []int, err error) {
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
					err = fmt.Errorf("not enough elements in stack: %v", stack)
					break
				}
			} else if queryTerms[i] == "||" {
				if len(stack) >= 2 {
					stack[len(stack)-2] = UnionPosting(stack[len(stack)-1], stack[len(stack)-2])
					stack = stack[:len(stack)-1]
				} else {
					err = fmt.Errorf("not enough elements in stack: %v", stack)
					break
				}
			}  else {
				stack = append(stack, s.ii.PostingsList(queryTerms[i]))
			}
		}
		if len(stack) == 1 {
			results = stack[0]
		} else if len(stack) > 1 {
			err = fmt.Errorf("stack is not empty: %v", stack)
		}
	} else if unionFlag {
		results = s.ii.Union(splitTrimToLower(query, "||"))
	} else {
		results = s.ii.Intersect(splitTrimToLower(query, "&&"))
	}
	return
}

// Parses the expression into a set of tokens. Assumes that the expression is written in the infix notation.
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

// Reorders the tokens from infix notation to a postfix notation.
// Partial implementation of Shunting-yard algorithm, which only considers left associative, binary operators.
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

// Splits the input string with the provided token. Trim and lowercase the output tokens.
func splitTrimToLower(str string, split string) (out []string) {
	out = strings.Split(str, split)
	for i := range out {
		out[i] = strings.TrimSpace(strings.ToLower(out[i]))
	}
	return
}

// Filters documents that contain all of the provided terms.
// Each term permits a spelling correction to terms within a certain edit distance.
func (s Searcher) FuzzyQuery(query string) (results []int) {
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

// Edit distance for each term is based on the length of the term.
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

// Filters documents that contain all of the provided terms.
// Terms can contain the characters `?` and/or `*` which expands into one or more terms.
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