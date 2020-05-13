package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"testing"
)

func TestWordCount(t *testing.T) {
	pairs := []struct{
		str string
		length int
	}{
		{"My name is John.", 4},
		{"  to be  or not    to be", 6},
		{"Document A: This is a hat. This is a cat.", 10},
	}
	for _, pair := range pairs {
		length := wordCount(pair.str)
		if length != pair.length {
			t.Errorf("Wrong length for %s: Got %d, Wanted %d.", pair.str, length, pair.length)
		}
	}
}

func setUpDocumentLengths() (docList *DocumentLengths) {
	docList = &DocumentLengths{}
	docList.addDocumentLength("My name is John.")
	docList.addDocumentLength("  to be  or not    to be")
	docList.addDocumentLength("Document A: This is a hat. This is a cat.")
	return
}

func TestDocumentList_DocLength(t *testing.T) {
	pairs := []struct{
		id int
		length int
	}{
		{1, 4},
		{2, 6},
		{3, 10},
	}
	docLen := setUpDocumentLengths()
	for _, pair := range pairs {
		length := docLen.docLength(pair.id)
		if length != pair.length {
			t.Errorf("Wrong length for index %d: Got %d, Wanted %d.", pair.id, length, pair.length)
		}
	}
}

func TestDocumentList_AverageDocumentLength(t *testing.T) {
	docLen := setUpDocumentLengths()
	avgLen := (4 + 6 + 10) / 3.0
	if docLen.averageDocumentLength() != avgLen {
		t.Errorf("Wrong average length: Got %f, Wanted %f.", docLen.averageDocumentLength(), avgLen)
	}
}

func TestCSVStorage_Apply(t *testing.T) {
	wanted := []string{"Cohen's kappa", "Latent semantic analysis", "Code-division multiple access"}
	csvStore := NewCSVStorage("example.csv")
	csvStore.Apply(func (d Document) {
		 if wanted[d.id - 1] != d.Title {
			 t.Errorf("Read wrong document. Got %v, Wanted %v.", d.Title, wanted[d.id - 1])
		 }
	})
}

func TestGetDocumentFromCSV(t *testing.T) {
	wanted := []string{"Latent semantic analysis", "Cohen's kappa", "Code-division multiple access"}
	csvStore := NewCSVStorage("example.csv")
	resultList := csvStore.Get([]int{2, 1, 3})
	if resultList[0].Title != wanted[0] && resultList[1].Title != wanted[1] && resultList[2].Title != wanted[2] {
		t.Errorf("Wrong document titles. Got %v, Wanted %v.", resultList, wanted)
	}
	// Test getting non-existent document.
	resultList = csvStore.Get([]int{4})
	for _, result := range resultList {
		if result.id != 0 && result.Title != "" && result.Body != "" && result.URL != "" {
			t.Errorf("Retrived wrong document. Got %v, Wanted empty document", result)
		}
	}
}

func TestReadDocumentFromSQL(t *testing.T) {
	wanted := []string{"Cohen's kappa", "Latent semantic analysis", "Code-division multiple access"}
	db, err := sql.Open("sqlite3", "example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbStore := NewSQLStorage(db)
	dbStore.Apply(func (d Document) {
		if wanted[d.id - 1] != d.Title {
			t.Errorf("Read wrong document. Got %v, Wanted %v.", d.Title, wanted[d.id - 1])
		}
	})
}

func TestGetDocumentFromSQL(t *testing.T) {
	wanted := []string{"Latent semantic analysis", "Cohen's kappa", "Code-division multiple access"}
	db, err := sql.Open("sqlite3", "example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbStore := NewSQLStorage(db)
	resultList := dbStore.Get([]int{2, 1, 3})
	if resultList[0].Title != wanted[0] && resultList[1].Title != wanted[1] && resultList[2].Title != wanted[2] {
		t.Errorf("Wrong document titles. Got %v, Wanted %v.", resultList, wanted)
	}
	// Test getting non-existent document.
	resultList = dbStore.Get([]int{4})
	for _, result := range resultList {
		if result.id != 0 && result.Title != "" && result.Body != "" && result.URL != "" {
			t.Errorf("Retrived wrong document. Got %v, Wanted empty document", result)
		}
	}
}