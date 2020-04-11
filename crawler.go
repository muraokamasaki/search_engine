package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
)

// Gets JSON from url and stores it in target interface.
func readAPI(url string, target *interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
	return json.Unmarshal(body, target)
}

// Uses the wiki API to get the introduction paragraph for the title topic.
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

// Uses the wiki API to get up to 3 out-going links in the wiki page.
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

// Parses the JSON the contents to retrieve contents.
func scrapeWikiContents(title string) (string, bool) {
	var webPage interface{}
	err := readAPI(getWikiContents(title), &webPage)
	if err != nil {
		log.Println(err)
	}
	ext, ok := webPage.(map[string]interface{})["query"].(map[string]interface{})["pages"].
		([]interface{})[0].(map[string]interface{})["extract"].(string)
	return ext, ok
}

// Parses the JSON the contents to retrieve out-going links.
func scrapeWikiLinks(title string) (links []string, ok bool) {
	var webPage interface{}
	err := readAPI(getWikiLinks(title), &webPage)
	if err != nil {
		log.Println(err)
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

func getWikiURL(title string) string {
	u, err := url.Parse("https://en.wikipedia.org/wiki/")
	if err != nil {
		log.Println(err)
	}
	u.Path = path.Join(u.Path, title)
	return u.String()
}