package main

import (
	"database/sql"
	"encoding/csv"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Document struct {
	id    int
	Title string
	Body  string
	URL   string
}

// In-memory struct for fast access of lengths of document needed for document length normalization.
type DocumentLengths struct {
	lengths []int
	totalLength int
}

// Adds a document to the document length list.
func (docLen *DocumentLengths) addDocumentLength(document string) {
	docLength := wordCount(document)
	docLen.lengths = append(docLen.lengths, docLength)
	docLen.totalLength += docLength
	return
}

// Returns the length (the number of words) of a document.
func (docLen DocumentLengths) DocLength(docID int) int {
	return docLen.lengths[docID-1]  // -1 since documents are indexed from 1.
}

// Returns the average length of documents.
func (docLen DocumentLengths) averageDocumentLength() float64 {
	return float64(docLen.totalLength) / float64(len(docLen.lengths))
}

// Returns the number of words in a document.
func wordCount(document string) int {
	return len(strings.Fields(document))
}

// Functions that reads documents from external sources.

// Defines functions that consumes documents.
type documentFn func(document Document)

// Reads through a CSV file of columns 'id', 'title', 'body' and 'URL' and applies a function to each document.
func readDocumentFromCSV(filename string, fn documentFn) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	_, err = r.Read()  // Ignores the header.
	if err != nil {
		log.Fatal(err)
	}
	for {
		doc, err := r.Read()
		if err == io.EOF {
			break
		}
		id, err := strconv.Atoi(doc[0])
		if err != nil {
			log.Fatal(err)
		}
		fn(Document{
			id:    id,
			Title: doc[1],
			Body:  doc[2],
			URL: doc[3],
		})
	}
}

// Retrieve a list of documents from a CSV with columns 'id', 'title', 'body' and 'URL.
func getDocumentFromCSV(filename string, ids []int) (resultsList []Document) {
	resultsList = make([]Document, len(ids))

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	_, err = r.Read()  // Ignores the header.
	if err != nil {
		log.Fatal(err)
	}

	idsCopy := make([]int, len(ids))
	copy(idsCopy, ids)
	sort.Ints(idsCopy)
	idPointer := 0

	for idPointer < len(idsCopy) {
		doc, err := r.Read()
		if err == io.EOF {
			// Some documents cannot be found.
			break
		}
		id, err := strconv.Atoi(doc[0])
		if err != nil {
			log.Fatal(err)
		}
		if id == idsCopy[idPointer] {
			idx := 0
			for id != ids[idx] {
				idx++
			}
			resultsList[idx] = Document{
				id:    id,
				Title: doc[1],
				Body:  doc[2],
				URL: doc[3],
			}
			idPointer++
		}
	}
	return
}

// Read from a SQL database table 'documents' with attributes 'id', 'title', 'body' and 'URL'
func readDocumentFromSQL(filename string, fn documentFn) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, body, URL FROM documents")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var document Document
	for rows.Next() {
		if err = rows.Scan(&document.id, &document.Title, &document.Body, &document.URL); err != nil {
			log.Fatal(err)
		}
		fn(document)
	}
}

// Retrieve a list of documents from a SQL database table 'documents' with attributes 'id', 'title', 'body' and 'URL'
func getDocumentFromSQL(filename string, ids []int) (resultsList []Document) {
	resultsList = make([]Document, len(ids))

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for idx, id := range ids {
		row := db.QueryRow("SELECT id, title, body, URL FROM documents WHERE id=?", id)
		row.Scan(&resultsList[idx].id, &resultsList[idx].Title, &resultsList[idx].Body, &resultsList[idx].URL) // Ignore error, skip idx if no document can be found.
	}
	return
}