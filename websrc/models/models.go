package models

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"

	"GoSearchSuggestions/trie"
)

var db *sql.DB
var Trie *trie.Trie
var pool *redis.Pool
var err error

type dbConfig struct {
	User     string
	Password string
	Db_name  string
	Ssl_mode string
}

type redisConfig struct {
	Status   string
	Exists   string
	Queue    string
	Port     string
	Protocol string
}

type config struct {
	Db    dbConfig
	Redis redisConfig
}

func Init() {
	var c config
	loadConfig("config.json", &c)
	connString := fmt.Sprintf("user=%v password=%ord=%v dbname=%v sslmode=%v", c.Db.User, c.Db.Password, c.Db.Db_name, c.Db.Ssl_mode)
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

	pool, err = MakePool(c)
	if err != nil {
		panic(err)
	}
}

func loadConfig(path string, c *config) {
	content, err := ioutil.ReadFile(path)
	//Panic because reading config is going to be useful
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, c)
}

func MakePool(c config) (*redis.Pool, error) {
	pool := redis.Pool{
		MaxIdle:   50,
		MaxActive: 0,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(c.Redis.Protocol, c.Redis.Port)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
	return &pool, nil
}
