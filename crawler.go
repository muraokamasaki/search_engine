package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"
)

// readAPI issues a GET request to a url, reading into JSON
// and storing it in the target interface. Returns an error
// if it fails to get the url or fails to read its content.
func readAPI(url string, target *interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

// getWikiContents returns a url to get the introduction paragraph
// for given title.
func getWikiContents(title string) string {
	w, err := url.Parse("https://en.wikipedia.org/w/api.php")
	if err != nil {
		log.Println(err)
	}
	q := w.Query()
	q.Set("titles", title)
	q.Set("action", "query")
	q.Set("prop", "extracts")
	q.Set("format", "json")
	q.Set("exlimit", "1")
	q.Set("explaintext", "1")
	q.Set("exintro", "1")
	q.Set("formatversion", "2")
	w.RawQuery = q.Encode()

	return w.String()
}

// getWikiContents returns a url to get up to 3 out-going
// links from the Wikipedia article of the given title.
func getWikiLinks(title string) string {
	w, err := url.Parse("https://en.wikipedia.org/w/api.php")
	if err != nil {
		log.Println(err)
	}
	q := w.Query()
	q.Set("titles", title)
	q.Set("action", "query")
	q.Set("prop", "links")
	q.Set("format", "json")
	q.Set("pllimit", "3")
	q.Set("formatversion", "2")
	w.RawQuery = q.Encode()

	return w.String()
}

// scrapeWikiContents uses the WikiAPI to retrieve the introductory
// paragraph for the given title. Returns the contents as a string
// and a bool indicating if the scraping was successful.
func scrapeWikiContents(title string) (string, bool) {
	var webPage interface{}
	err := readAPI(getWikiContents(title), &webPage)
	if err != nil {
		log.Println(err)
		return "", false
	}
	ext, ok := webPage.(map[string]interface{})["query"].(map[string]interface{})["pages"].
		([]interface{})[0].(map[string]interface{})["extract"].(string)
	return ext, ok
}

// scrapeWikiLinks uses the WikiAPI to retrieve up to 3 out-going
// links from the Wikipedia article. Returns the links as a string
// slice and a bool indicating if scraping was successful.
func scrapeWikiLinks(title string) (links []string, ok bool) {
	var webPage interface{}
	err := readAPI(getWikiLinks(title), &webPage)
	if err != nil {
		log.Println(err)
		return
	}
	ext, ok := webPage.(map[string]interface{})["query"].(map[string]interface{})["pages"].
		([]interface{})[0].(map[string]interface{})["links"].([]interface{})
	if ok {
		for _, i := range ext {
			t, ok := i.(map[string]interface{})["title"].(string)
			if !ok {
				continue
			}
			links = append(links, t)
		}
	}
	return
}

// getWikiURL returns the address to a Wikipedia article for a
// given title. Does not check if the article exists.
func getWikiURL(title string) string {
	u, err := url.Parse("https://en.wikipedia.org/wiki/")
	if err != nil {
		log.Println(err)
	}
	u.Path = path.Join(u.Path, title)
	return u.String()
}

// LinkMap keeps track of links that have been scraped.
type LinkMap struct {
	seenLinks map[string]bool
	mux       sync.Mutex
}

func (m *LinkMap) Length() int {
	m.mux.Lock()
	defer m.mux.Unlock()
	return len(m.seenLinks)
}

func (m *LinkMap) Value(key string) bool {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.seenLinks[key]
}

func (m *LinkMap) AddLink(key string) {
	m.mux.Lock()
	m.seenLinks[key] = true
	m.mux.Unlock()
}

// crawlWikiContents will scrape the article at the given URL,
// create a Document and send it through the channel.
func crawlWikiContents(link string, ch chan Document) {
	title := strings.ReplaceAll(link, " ", "_")
	contents, ok := scrapeWikiContents(title)
	if !ok {
		log.Println("Cannot retrieve contents for", link)
		return
	}
	ch <- Document{
		Title: link,
		Body:  contents,
		URL:   getWikiURL(title),
	}
}

// crawlWikiLinks will scrape the out-going links from a given
// URL and send each link though the channel.
func crawlWikiLinks(link string, ch chan string) {
	outlinks, ok := scrapeWikiLinks(link)
	if !ok {
		log.Println("Cannot retrieve outlinks for", link)
		return
	}
	for _, outlink := range outlinks {
		ch <- outlink
	}
}

// CrawlWiki crawls Wikipedia articles, starting from a given seed of articles,
// and saves the introductory paragraph in the DocumentSaver. Will scrape up
// to a given capacity of documents, or continue forever if -1 is passed. Will
// wait a given duration after each article is scraped for politeness.
// (Note) Article titles in seed must be capitalization properly as it can
// cause articles to not be retrieved.
func CrawlWiki(seed []string, docSaver DocumentSaver, capacity int, duration time.Duration) {
	linkMap := LinkMap{seenLinks: make(map[string]bool)}
	linkCh := make(chan string, len(seed))
	docCh := make(chan Document)
	documentsAdded := 0
	for _, s := range seed {
		linkCh <- s
	}
	for capacity == -1 || documentsAdded < capacity {
		select {
		case link := <- linkCh:
			if !linkMap.Value(link) && (capacity == -1 || linkMap.Length() < capacity) {
				linkMap.AddLink(link)
				go crawlWikiLinks(link, linkCh)
				go crawlWikiContents(link, docCh)
				time.Sleep(duration)
			}
		case document := <-docCh:
			docSaver.Save(document)
			if capacity != -1 {
				documentsAdded++
			}
		}
	}
}
