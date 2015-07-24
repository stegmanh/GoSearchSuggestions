package models

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"

	"GoSearchSuggestions/trie"
)

var db *sql.DB
var Trie *trie.Trie
var err error

type dbConfig struct {
	User     string
	Password string
	Db_name  string
	Ssl_mode string
}

type config struct {
	Db dbConfig
}

func Init() {
	var c config
	loadConfig("config.json", &c)
	connString := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=%v", c.Db.User, c.Db.Password, c.Db.Db_name, c.Db.Ssl_mode)
	db, err = sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	file, err := os.Open("./text/words.txt")
	if err != nil {
		panic("Error Loading Trie")
	}
	defer file.Close()

	Trie = &trie.Trie{}
	Trie.Initialize()
	scanner := bufio.NewScanner(file)
	Trie.BuildTrie(scanner)
}

func loadConfig(path string, c *config) {
	content, err := ioutil.ReadFile(path)
	//Panic because reading config is going to be useful
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, c)
}
