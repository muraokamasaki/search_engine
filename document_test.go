package main

import (
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

func setUpDocumentList() (docList *DocumentList) {
	docList = &DocumentList{}
	docList.addToDocumentList("My name is John.")
	docList.addToDocumentList("  to be  or not    to be")
	docList.addToDocumentList("Document A: This is a hat. This is a cat.")
	return
}

func TestDocumentList_DocLength(t *testing.T) {
	pairs := []struct{
		idx int
		length int
	}{
		{0, 4},
		{1, 6},
		{2, 10},
	}
	docList := setUpDocumentList()
	for _, pair := range pairs {
		length := docList.DocLength(pair.idx)
		if length != pair.length {
			t.Errorf("Wrong length for index %d: Got %d, Wanted %d.", pair.idx, length, pair.length)
		}
	}
}

func TestDocumentList_AverageDocumentLength(t *testing.T) {
	docList := setUpDocumentList()
	avgLen := (4 + 6 + 10) / 3.0
	if docList.averageDocumentLength() != avgLen {
		t.Errorf("Wrong average length: Got %f, Wanted %f.", docList.averageDocumentLength(), avgLen)
	}
}