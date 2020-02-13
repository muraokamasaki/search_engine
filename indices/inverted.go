package indices

import (
    "bufio"
    "log"
    "os"
    "sort"
    "strings"
)

type InvertedIndex struct {
    PostingsLists map[string][]int
}

func NewInvertedIndex() *InvertedIndex {
    return &InvertedIndex{PostingsLists: make(map[string][]int)}
}

func (ii InvertedIndex) BuildFromTextFile(filename string) {
    // Builds the inverted index from a text file where each document exists on a single line.
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
        for _, word := range Tokenize(line) {
            ii.addIDToPostingsList(word, docID)
        }
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func (ii InvertedIndex) addIDToPostingsList(word string, docID int) {
    if len(word) > 0 {
        word = strings.ToLower(word)
        pList := ii.PostingsLists[word]
        if len(pList) == 0 || pList[len(pList) - 1] != docID {
            ii.PostingsLists[word] = append(pList, docID)
        }
    }
}

func (ii InvertedIndex) Intersect(terms []string) (result []int) {
    sort.Slice(terms, func(i, j int) bool { return len(ii.PostingsLists[terms[i]]) < len(ii.PostingsLists[terms[j]]) })
    result = ii.PostingsLists[terms[0]]
    pointer := 1
    for pointer < len(terms) && len(result) != 0 {
        result = IntersectPosting(result, ii.PostingsLists[terms[pointer]])
        pointer++
    }
    return
}

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

func (ii InvertedIndex) Union(terms []string) (result []int) {
    set := make(map[int]bool)
    for _, term := range terms {
        for _, id := range ii.PostingsLists[term] {
            set[id] = true
        }
    }
    result = make([]int, len(set))
    i := 0
    for k := range set {
        result[i] = k
        i++
    }
    return
}

func UnionPosting(plist1 []int, plist2 []int) (result []int) {
    // plist1 and plist 2 are assumed to be sorted.
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