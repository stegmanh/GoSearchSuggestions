package main

import (
	"./trie"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Suggestions struct {
	Term    string
	Results []string
}

var plainWord = regexp.MustCompile(`(^[a-zA-Z]*$)`)
var trieTree *trie.Trie = nil
var searchHistory map[string]int

func searchHandler(w http.ResponseWriter, r *http.Request, t *trie.Trie) {
	searchTerm := r.FormValue("q")
	if len(searchTerm) == 0 {
		fmt.Fprintf(w, "%#v", "Please send search results")
		return
	}
	searchHistory[searchTerm]++
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

//Helper function to clear a map if searched for less than 2 times
func trimMap(m *map[string]int) {
	for k, v := range *m {
		if v < 2 {
			delete(*m, k)
		}
	}
}

func main() {
	//Init the searchhistory
	searchHistory = make(map[string]int)
	//Timer to keep search history and stuff
	go func() {
		c := time.Tick(30 * time.Second)
		for now := range c {
			fmt.Printf("%v, %v\n", now, searchHistory)
			trimMap(&searchHistory)
		}
	}()
	file, err := os.Open("titles.txt")
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
	http.ListenAndServe(":8080", nil)
}
