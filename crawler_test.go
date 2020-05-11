package main

import (
	"strings"
	"testing"
	"time"
)

type TestSaver struct {
	docs []Document
}

func (s *TestSaver) Save(document Document) {
	s.docs = append(s.docs, document)
}

func TestCrawlWiki(t *testing.T) {
	s := &TestSaver{}
	capacity := 2
	seed := []string{"Pet door"}
	expectedResults := []Document{
		{
			id:    1,
			Title: "Pet door",
			Body:  "A pet door or pet flap (also referred to in more specific terms",
			URL:   "https://en.wikipedia.org/wiki/Pet_door",
		},
		{
			id:    2,
			Title: "Ancient Egypt",
			Body:  "Ancient Egypt was a civilization of ancient North Africa",
			URL:   "https://en.wikipedia.org/wiki/Ancient_Egypt",
		},
	}
	CrawlWiki(seed, s, capacity, time.Millisecond)  // Only wait 1 millisecond since we scrape 2 links.
	if len(s.docs) != capacity {
		t.Errorf("Wrong number of documents retrieved. Got %v, Wanted %v.", len(s.docs), capacity)
	}
	for idx, doc := range s.docs {
		if expectedResults[idx].Title != doc.Title {
			t.Errorf("Did not get expected title. Got %v, Wanted %v.", expectedResults[idx].Title, doc.Title)
		}
		if !strings.HasPrefix(doc.Body, expectedResults[idx].Body) {
			// Body is long and might change so we
			t.Logf("Did not get expected document. This might be because the article was edited." +
				"Got %v, Wanted %v as the prefix.", expectedResults[idx].Body, doc.Body)
		}

		if expectedResults[idx].URL != doc.URL {
			t.Errorf("Did not get expected URL. Got %v, Wanted %v.", expectedResults[idx].URL, doc.URL)
		}
	}
}
