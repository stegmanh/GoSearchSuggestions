package routes

import (
	"GoSearchSuggestions/trie"
	"GoSearchSuggestions/websrc/models"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

type Suggestions struct {
	Term    string
	Results []string
}

type ErrorResponse struct {
	Err string `json:"error"`
}

var plainWord = regexp.MustCompile(`(^[a-zA-Z_]*$)`)
var trieTree *trie.Trie = nil
var searchHistory map[string]int

func AutoComplete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	searchTerm := vars["term"]
	if len(searchTerm) == 0 {
		fmt.Fprintf(w, "%#v", "Please send search results")
		return
	}
	searchResults := models.Trie.Find(searchTerm)
	responseJSON := Suggestions{Term: searchTerm, Results: searchResults}
	ReponseJSON(w, responseJSON, http.StatusOK)
}

func DbSearchHandler(w http.ResponseWriter, r *http.Request) {
	//Change to something other than title
	articleTitle := r.FormValue("title")
	if len(articleTitle) == 0 {
		response := ErrorResponse{Err: "You must input a title"}
		ReponseJSON(w, response, http.StatusNotAcceptable)
		return
	}
	articles, err := models.SearchArticles(articleTitle)
	if err != nil {
		fmt.Fprintf(w, "%#v", err)
		response := ErrorResponse{Err: "Error Reading from Database"}
		ReponseJSON(w, response, http.StatusInternalServerError)
		return
	}
	/*
		Disguesting
			for _, article := range articles.Articles {
				article.Body = bytes.Replace(article.Body, []byte("\\u003c"), []byte("<"), -1)
				article.Body = bytes.Replace(article.Body, []byte("\\u003e"), []byte(">"), -1)
				article.Body = bytes.Replace(article.Body, []byte("\\u0026"), []byte("&"), -1)
			}
	*/
	ReponseJSON(w, articles, http.StatusOK)
}

// func TitleSearch(w http.ResponseWriter, r *http.Request) {
// }

func ReponseJSON(w http.ResponseWriter, i interface{}, status int) {
	json, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		http.Error(w, "Error Encoding JSON", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)
}
