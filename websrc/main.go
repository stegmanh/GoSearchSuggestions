package main

import (
	"bufio"
	"log"
	"net/http"
	"os"

	"GoSearchSuggestions/trie"
	"GoSearchSuggestions/websrc/app/routes"
	"GoSearchSuggestions/websrc/models"
)

func main() {
	models.Init()

	//Read words.txt
	file, err := os.Open("./websrc/static/text/words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	trieTree := &trie.Trie{}
	trieTree.Initialize()
	scanner := bufio.NewScanner(file)
	trieTree.BuildTrie(scanner)

	http.ListenAndServe(":8080", routes.InitHandler())
}
