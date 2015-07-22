package routes

import (
	"GoSearchSuggestions/trie"
	"GoSearchSuggestions/websrc/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type Suggestions struct {
	Term    string
	Results []string
}

var plainWord = regexp.MustCompile(`(^[a-zA-Z_]*$)`)
var trieTree *trie.Trie = nil
var searchHistory map[string]int

func SearchHandler(w http.ResponseWriter, r *http.Request, t *trie.Trie) {
	searchTerm := r.FormValue("q")
	if len(searchTerm) == 0 {
		fmt.Fprintf(w, "%#v", "Please send search results")
		return
	}
	searchResults := t.Find(searchTerm)
	responseJSON := Suggestions{Term: searchTerm, Results: searchResults}
	js, err := json.Marshal(responseJSON)
	if err != nil {
		fmt.Fprintf(w, "%#v", "Error encoding json")
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1>", "Hello Noah!")
}

func DbSearchHandler(w http.ResponseWriter, r *http.Request) {
	articleTitle := r.FormValue("title")
	if len(articleTitle) == 0 {
		fmt.Fprintf(w, "%v", "Please input an article title")
		return
	}
	articles, err := models.SearchArticles(articleTitle)
	if err != nil {
		fmt.Fprintf(w, "%#v", err)
	}
	js, err := json.Marshal(articles)
	if err != nil {
		fmt.Fprintf(w, "%#v", "Error encoding json")
		fmt.Println(err)
		return
	}
	//Unencode.. This is disguesting but I can't find much else...
	js = bytes.Replace(js, []byte("\\u003c"), []byte("<"), -1)
	js = bytes.Replace(js, []byte("\\u003e"), []byte(">"), -1)
	js = bytes.Replace(js, []byte("\\u0026"), []byte("&"), -1)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
