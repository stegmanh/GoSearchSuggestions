package main

import (
	"GoSearchSuggestions/trie"
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Suggestions struct {
	Term    string
	Results []string
}

type Article struct {
	Title  string `json:"title"`
	Date   string `json:"date"`
	Source string `json:"source"`
	Body   string `json:"body"`
}

type ArticleResponse struct {
	Articles []Article `json:"data"`
}

var plainWord = regexp.MustCompile(`(^[a-zA-Z_]*$)`)
var trieTree *trie.Trie = nil
var searchHistory map[string]int

var ftsSearch string = "SELECT title, source, body, created_at FROM articles, to_tsvector(title) tvt, to_tsquery($1) tvq WHERE tvt @@ tvq ORDER BY ts_rank(tvt, tvq) DESC LIMIT 5"

func searchHandler(w http.ResponseWriter, r *http.Request, t *trie.Trie) {
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1>", "Hello Noah!")
}

func dbSearchHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	articleTitle := r.FormValue("title")
	if len(articleTitle) == 0 {
		fmt.Fprintf(w, "%v", "Please input an article title")
		return
	}
	articleTitle = strings.Join(strings.Split(articleTitle, " "), " | ")
	rows, err := db.Query(ftsSearch, articleTitle)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	defer rows.Close()
	toReturn := ArticleResponse{Articles: make([]Article, 0)}
	for rows.Next() {
		var title, createdAt, source, body string
		err = rows.Scan(&title, &source, &body, &createdAt)
		if err != nil {
			fmt.Println("Got an error here..", err)
			continue
		}
		articleAdd := Article{Title: title, Date: createdAt, Source: source, Body: body}
		toReturn.Articles = append(toReturn.Articles, articleAdd)
	}
	js, err := json.Marshal(toReturn)
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

func main() {
	//Init DB, panic if fails
	db, err := sql.Open("postgres", "user=postgres password=qwert12345 dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	//Read words.txt
	file, err := os.Open("./websrc/static/text/titles.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	trieTree := &trie.Trie{}
	trieTree.Initialize()
	scanner := bufio.NewScanner(file)
	trieTree.BuildTrie(scanner)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		searchHandler(w, r, trieTree)
	})
	http.HandleFunc("/dbsearch", func(w http.ResponseWriter, r *http.Request) {
		dbSearchHandler(w, r, db)
	})
	http.ListenAndServe(":8080", nil)
}
