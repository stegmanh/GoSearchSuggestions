package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"

	_ "github.com/lib/pq"
)

var Db *sql.DB
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
	Db, err = sql.Open("postgres", connString)
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
