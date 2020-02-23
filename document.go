package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type DocumentList struct {
	nextDocID int
	lengths []int
	totalLength int
}

// Adds a document to the document list. Returns the ID of the document.
func (docList *DocumentList) addToDocumentList(document string) (docID int) {
	docID = docList.nextDocID
	docList.nextDocID++
	docLength := wordCount(document)
	docList.lengths = append(docList.lengths, docLength)
	docList.totalLength +=  docLength
	return
}

// Gets length of a document.
func (docList DocumentList) DocLength(docID int) int {
	return docList.lengths[docID]
}

// Returns the average length of documents.
func (docList DocumentList) averageDocumentLength() float64 {
	return float64(docList.totalLength) / float64(len(docList.lengths))
}

// Returns the number of words in a document.
func wordCount(document string) int {
	return len(strings.Fields(document))
}

// Functions that reads documents from external sources.

// Defines functions that consumes documents.
type documentFn func(string)

// Read documents from a text file, where each document exists on a line.
func readLinesFromTextFile(filename string, fn documentFn) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		document := scanner.Text()
		fn(document)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}