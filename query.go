package main

import (
	"fmt"
	"strings"
)

type Searcher struct {
	ii InvertedIndex
	ki KGramIndex
}

func NewSearcher(k int) *Searcher {
	return &Searcher{ii: *NewInvertedIndex(), ki: *NewKGramIndex(k)}
}

func (s Searcher) TermQuery(query string) (results []int) {
	// Tokenize the query term and look for a document which contains all of the tokens.
	results = s.ii.Intersect(Tokenize(query))
	return
}

func (s Searcher) BooleanQuery(query string) (results []int, err error) {
	// Allows boolean operations between terms. Terms should only consist of a single word.
	var queryTerms []string
	intersectFlag := strings.Contains(query, "&&")
	unionFlag := strings.Contains(query, "||")
	// If the query does not contain both && or ||, we do not need to create postfix tree.
	if intersectFlag && unionFlag {
		queryTerms = shuntingYard(parseInfix(query))
		var stack [][]int
		for i := 0; i < len(queryTerms); i++ { // Parse postfix expression tree.
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
				stack = append(stack, s.ii.Intersect([]string{queryTerms[i]}))
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

func shuntingYard(tokens []string) (output []string) {
	// Partial implementation of Shunting-yard algorithm, which only considers left associative, binary operators.
	orderOfOperations := map[string]int{"||": 1, "&&": 2} // Larger value implies higher precedence; Values are arbitrary.
	var operators []string
	for _, token := range tokens {
		if _, ok := orderOfOperations[token]; !ok {
			output = append(output, token)
		} else {
			for i := len(operators) - 1; i>= 0; i-- {
				if orderOfOperations[operators[i]] >= orderOfOperations[token] { // If equal, assumes token is left associative.
					output = append(output, operators[i])
					operators = operators[:i]
				} else {
					break
				}
			}
			operators = append(operators, token)
		}
	}
	for i := len(operators) - 1; i>= 0; i-- {
		output = append(output, operators[i])
	}
	return
}

func splitTrimToLower(str string, split string) (out []string) {
	out = strings.Split(str, split)
	for i := range out {
		out[i] = strings.TrimSpace(strings.ToLower(out[i]))
	}
	return
}

func (s Searcher) FuzzyQuery(query string) (results []int) {
	var fuzziness int
	for _, queryToken := range Tokenize(query) {
		if len(queryToken) <= 2 {
			fuzziness = 0
		} else if len(queryToken) <= 5 {
			fuzziness = 1
		} else {
			fuzziness = 2
		}
		terms := s.ki.GetCloseTerms(queryToken, fuzziness)
		var partialResult []string
		for _, res := range terms {
			if EditDistance(queryToken, res) <= fuzziness {
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

func (s Searcher) WildcardQuery(query string) (results []int) {
	for _, queryToken := range TokenizeWildcard(query) {
		terms := s.ki.KGramMatch(queryToken)
		var partialResult  []string
		for _, res := range terms {
			if WildcardMatch(queryToken, res) {
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