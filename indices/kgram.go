package indices

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type KGramIndex struct {
	k int
	PostingsLists map[string][]string
}

func NewKGramIndex(k int) *KGramIndex {
	return &KGramIndex{k: k, PostingsLists: make(map[string][]string)}
}

func (ki KGramIndex) BuildFromTextFile(filename string) {
	// Builds the k-gram index from a text file where each document exists on a single line.
	// Generates k-grams from only the title (separated from body with a colon).
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	docID := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		docID++
		title := strings.Split(line, ":")[0] // Change this!
		ki.addWordToPostingsList(title)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (ki KGramIndex) addWordToPostingsList(word string) {
	m := buildKGrams(word, ki.k)
	for _, gram := range m {
		pList := ki.PostingsLists[gram]
		ki.PostingsLists[gram] = append(pList, word)
	}
}

func buildKGrams(str string, k int) (grams []string) {
	str = strings.ToLower(StripPunctuation(str))
	if len(str) < k {
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

func (ki KGramIndex) KGramOverlap(q string) (count map[string]int) {
	// Finds the number of overlapping k-grams between the query and terms in the postings list.
	count = make(map[string]int)
	grams := buildKGrams(q, ki.k)
	for _, g := range grams {
		for _, t := range ki.PostingsLists[g] {
			count[t]++
		}
	}
	return
}

func lowerBoundOverlap(s1 string, s2 string, editDistance int, k int) int {
	return max(len(s1), len(s2)) - 1 - (editDistance - 1) * k
}

func (ki KGramIndex) FuzzyQuery(q string) (terms []string) {
	count := ki.KGramOverlap(q)
	editDistance := 2
	for k, v := range count {
		if v >= lowerBoundOverlap(q, k, editDistance, ki.k) {
			terms = append(terms, k)
		}
	}
	return
}
