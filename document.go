package main

import (
	"database/sql"
	"encoding/csv"
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

// DocumentLengths stores the lengths of document and total length
// of the documents.
type DocumentLengths struct {
	lengths []int
	totalLength int
}

// addDocumentLength stores the length of the document.
func (docLen *DocumentLengths) addDocumentLength(document string) {
	docLength := wordCount(document)
	docLen.lengths = append(docLen.lengths, docLength)
	docLen.totalLength += docLength
	return
}

// docLength returns the length (word count) of a document.
func (docLen *DocumentLengths) docLength(docID int) int {
	return docLen.lengths[docID-1]  // -1 since documents are indexed from 1.
}

// averageDocumentLength returns the average length of stored documents.
func (docLen *DocumentLengths) averageDocumentLength() float64 {
	return float64(docLen.totalLength) / float64(len(docLen.lengths))
}

// wordCount returns the number of words in a document.
func wordCount(document string) int {
	return len(strings.Fields(document))
}

// Functions that reads documents from external sources.

// documentFn defines functions that consumes documents.
type documentFn func(document Document)

// DocumentStorage is an interface that supports
// retrieving and applying functions to documents.
type DocumentStorage interface {
	// Apply applies a function to each document in DocumentStorage.
	Apply(fn documentFn)
	// Get returns a list of documents from a given list of ids.
	Get(ids []int) []Document
}

// DocumentSaver is an interface that supports
// adding documents to a collection.
type DocumentSaver interface {
	// Save stores the given document.
	// Used when crawling.
	Save(document Document)
}

// CSVStorage contains a csv file storing documents with columns
// 'id', 'title', 'body' and 'URL'.
type CSVStorage struct {
	filename string
}

func NewCSVStorage(filename string) *CSVStorage {
	return &CSVStorage{filename:filename}
}

func (store *CSVStorage) Apply(fn documentFn) {
	f, err := os.Open(store.filename)
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

func (store *CSVStorage) Get(ids []int) (resultsList []Document) {
	resultsList = make([]Document, len(ids))

	f, err := os.Open(store.filename)
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

// SQLStorage holds a database handle which contains a
// connection to a database storing documents with
// attributes 'id', 'title', 'body' and 'URL'
type SQLStorage struct {
	*sql.DB
}

func NewSQLStorage(db *sql.DB) *SQLStorage {
	return &SQLStorage{db}
}

func (store *SQLStorage) Apply(fn documentFn) {
	rows, err := store.Query("SELECT id, title, body, URL FROM documents")
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

func (store *SQLStorage) Get(ids []int) (resultsList []Document) {
	resultsList = make([]Document, len(ids))
	for idx, id := range ids {
		row := store.QueryRow("SELECT id, title, body, URL FROM documents WHERE id=?", id)
		row.Scan(&resultsList[idx].id, &resultsList[idx].Title, &resultsList[idx].Body, &resultsList[idx].URL) // Ignore error, skip idx if no document can be found.
	}
	return
}

func (store *SQLStorage) Save(document Document) {
	statement, err := store.Prepare("INSERT INTO documents (title, body, url) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(document.Title, document.Body, document.URL)
}