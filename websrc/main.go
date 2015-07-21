package main

import (
	"GoSearchSuggestions/trie"
	"GoSearchSuggestions/websrc/app/routes"
	"bufio"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

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

	http.HandleFunc("/", routes.IndexHandler)
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		routes.SearchHandler(w, r, trieTree)
	})
	http.HandleFunc("/dbsearch", func(w http.ResponseWriter, r *http.Request) {
		routes.DbSearchHandler(w, r, db)
	})
	http.ListenAndServe(":8080", nil)
}
