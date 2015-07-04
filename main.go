package main

import (
	"./trie"
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Suggestions struct {
	Term    string
	Results []string
}

var plainWord = regexp.MustCompile(`(^[a-zA-Z_]*$)`)
var trieTree *trie.Trie = nil
var searchHistory map[string]int

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
	rows, err := db.Query("SELECT * FROM articles WHERE title = $1", articleTitle)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var title, createdAt, source, body string
		err = rows.Scan(&title, &createdAt, &source, &body)
		if err != nil {
			fmt.Println("Got an error here..")
			continue
		}
		fmt.Fprintf(w, "%v, %v, %v, %v", title, createdAt, source, body)
		break
	}
}

func main() {
	//Init DB, panic if fails
	db, err := sql.Open("postgres", "user=postgres password=qwert12345 dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	//Read words.txt
	file, err := os.Open("words.txt")
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
