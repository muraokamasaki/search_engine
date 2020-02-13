package query

import (
	"fmt"
	"github.com/muraokamasaki/search_engine/indices"
	"strings"
)

func (s Searcher) TermQuery(query string) (results []int) {
	// Tokenize the query term and look for a document which contains all of the tokens.
	results = s.ii.Intersect(indices.Tokenize(query))
	return
}

func (s Searcher) BooleanQuery(query string) (results []int, err error) {
	// Allows boolean operations between terms. Terms should only consist of a single word.
	var queryTerms []string
	intersectFlag, unionFlag := false, false
	if strings.Contains(query, "&&") {
		intersectFlag = true
	}
	if strings.Contains(query, "||") {
		unionFlag = true
	}
	// If the query does not contain both && or ||, we do not need to create postfix tree.
	if intersectFlag && unionFlag {
		queryTerms = shuntingYard(parseInfix(query))
		var stack [][]int
		for i := 0; i < len(queryTerms); i++ { // Parse postfix expression tree.
			if queryTerms[i] == "&&" {
				if len(stack) >= 2 {
					stack[len(stack)-2] = indices.IntersectPosting(stack[len(stack)-1], stack[len(stack)-2])
					stack = stack[:len(stack)-1]

				} else {
					err = fmt.Errorf("not enough elements in stack: %v", stack)
					break
				}
			} else if queryTerms[i] == "||" {
				if len(stack) >= 2 {
					stack[len(stack)-2] = indices.UnionPosting(stack[len(stack)-1], stack[len(stack)-2])
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
		results = s.ii.Union(splitAndTrim(query, "||"))
	} else {
		results = s.ii.Intersect(splitAndTrim(query, "&&"))
	}
	return
}

func parseInfix(expr string) (output []string) {
	t1 := strings.Split(expr, "&&")
	for i := 0; i < len(t1); i++ {
		t2 := strings.Split(t1[i], "||")
		for j := 0; j < len(t2) - 1; j++ {
			output = append(output, strings.TrimSpace(t2[j]), "||")
		}
		output = append(output, strings.TrimSpace(t2[len(t2) - 1]), "&&")
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

func splitAndTrim(str string, split string) (out []string) {
	out = strings.Split(str, split)
	for i := range out {
		out[i] = strings.TrimSpace(out[i])
	}
	return
}

func FuzzyQuery(query string) (results []int) {
	return
}