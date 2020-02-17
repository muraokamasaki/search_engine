package main

import (
    "bufio"
    "log"
    "os"
    "sort"
)

// Implementation of a inverted index.
type InvertedIndex struct {
    postingsLists map[string][]int
}

func NewInvertedIndex() *InvertedIndex {
    return &InvertedIndex{postingsLists: make(map[string][]int)}
}

// Builds the inverted index from a text file where each document exists on a single line.
func (ii InvertedIndex) BuildFromTextFile(filename string) {
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
        for _, word := range tokenize(line) {
            ii.addIDToPostingsList(word, docID)
        }
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

// Adds the ID to the postings list of a term in the inverted index.
func (ii InvertedIndex) addIDToPostingsList(term string, docID int) {
    if len(term) > 0 {
        pList := ii.PostingsList(term)
        if len(pList) == 0 || pList[len(pList) - 1] != docID {
            ii.postingsLists[term] = append(pList, docID)
        }
    }
}

func (ii InvertedIndex) PostingsList(term string) (plist []int) {
    return ii.postingsLists[term]
}

// Returns the IDs of documents that contain all the terms.
func (ii InvertedIndex) Intersect(terms []string) (result []int) {
    sort.Slice(terms, func(i, j int) bool { return len(ii.postingsLists[terms[i]]) < len(ii.postingsLists[terms[j]]) })
    result = ii.PostingsList(terms[0])
    pointer := 1
    for pointer < len(terms) && len(result) != 0 {
        result = IntersectPosting(result, ii.PostingsList(terms[pointer]))
        pointer++
    }
    return
}

// Returns the intersection of two postings lists.
func IntersectPosting(plist1 []int, plist2 []int) (result []int) {
    // plist1 and plist 2 are assumed to be sorted.
    result = []int{}
    pointer1, pointer2 := 0, 0
    for pointer1 < len(plist1) && pointer2 < len(plist2) {
        docID1 := plist1[pointer1]
        docID2 := plist2[pointer2]
        if docID1 == docID2 {
            result = append(result, docID1)
            pointer1++
            pointer2++
        } else if docID1 < docID2 {
            pointer1++
        } else {
            pointer2++
        }
    }
    return
}

// Returns the IDs of documents that contains at least one term.
func (ii InvertedIndex) Union(terms []string) (result []int) {
    set := make(map[int]bool)
    for _, term := range terms {
        for _, id := range ii.PostingsList(term) {
            set[id] = true
        }
    }
    result = make([]int, len(set))
    i := 0
    for k := range set {
        result[i] = k
        i++
    }
    sort.Ints(result)
    return
}

// Returns the union of two postings lists.
func UnionPosting(plist1 []int, plist2 []int) (result []int) {
    // plist1 and plist2 are assumed to be sorted.
    set := make(map[int]bool)
    for _, id := range plist1 {
        set[id] = true
    }
    for _, id := range plist2 {
        set[id] = true
    }
    result = make([]int, len(set))
    i := 0
    for k := range set {
        result[i] = k
        i++
    }
    sort.Ints(result)
    return
}