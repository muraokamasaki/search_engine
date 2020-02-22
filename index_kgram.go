package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// Implementation of a k-gram index.
type KGramIndex struct {
	k             int
	postingsLists map[string][]string
}

func NewKGramIndex(k int) *KGramIndex {
	return &KGramIndex{k: k, postingsLists: make(map[string][]string)}
}

// Builds the k-gram index from a text file where each document exists on a single line.
func (ki KGramIndex) BuildFromTextFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		for _, token := range tokenize(line) {
			ki.addWordToPostingsList(token)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Builds the k-gram index for a term.
func (ki KGramIndex) addWordToPostingsList(term string) {
	m := buildKGrams(term, ki.k)
	for _, gram := range m {
		pList := ki.postingsLists[gram]
		if !matchInArray(pList, term) {
			ki.postingsLists[gram] = append(pList, term)
		}
	}
}

// Checks if a string exists in the array.
func matchInArray(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

// Generate k-grams (padded with `$`) from a given string.
func buildKGrams(str string, k int) (grams []string) {
	if len(str) < k - 1{
		grams = []string{str}
		return
	}
	grams = make([]string, len(str) + k - 1)
	for i := 0; i < len(str) - k + 1; i++ {
		grams[i] = str[i : i+k]
	}
	for i := 0; i < k - 1; i++ {
		padding := strings.Repeat("$", i + 1)
		grams[i + len(str)] = padding + str[: k - i - 1]
		grams[len(str) - i - 1] = str[len(str) - k + i + 1:] + padding
	}
	return
}

// Finds the number of overlapping k-grams between a string and terms in the postings list.
func (ki KGramIndex) KGramOverlap(str string) (count map[string]int) {
	count = make(map[string]int)
	grams := buildKGrams(str, ki.k)
	for _, g := range grams {
		for _, t := range ki.postingsLists[g] {
			count[t]++
		}
	}
	return
}

// Find all terms in the posting list that contains every k-grams generated by input string.
func (ki KGramIndex) KGramMatch(str string) (terms []string) {
	count := make(map[string]int)
	grams := buildKGrams(str, ki.k)
	wcGramCount := 0
	for _, g := range grams {
		if strings.Contains(g, "*") || strings.Contains(g, "?"){
			wcGramCount++
			continue
		}
		for _, t := range ki.postingsLists[g] {
			count[t]++
		}
	}
	for k, v := range count {
		if v == len(grams) - wcGramCount {
			terms = append(terms, k)
		}
	}
	return
}

// Finds the lower bound of matching k-gram terms between strings such that they are within a certain edit distance.
func lowerBoundKGramOverlap(s1 string, s2 string, maxEditDistance int, k int) int {
	return max(len(s1), len(s2)) - 1 - (maxEditDistance - 1) * k
}

// Returns terms that are within a certain edit distance from the input string.
func (ki KGramIndex) GetCloseTerms(str string, maxEditDistance int) (terms []string) {
	count := ki.KGramOverlap(str)
	for k, v := range count {
		// Only calculate edit distance if number of matching k-grams is above the lower bound.
		if v >= lowerBoundKGramOverlap(str, k, maxEditDistance, ki.k) && editDistance(str, k) <= maxEditDistance {
			terms = append(terms, k)
		}
	}
	return
}