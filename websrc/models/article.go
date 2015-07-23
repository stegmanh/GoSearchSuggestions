package models

import (
	"fmt"
	"strings"
)

type Article struct {
	Title  string `json:"title"`
	Date   string `json:"date"`
	Source string `json:"source"`
	Body   string `json:"body"`
	Id     int    `json:"-"`
}

type ArticleResponse struct {
	Articles []Article `json:"data"`
}

var ftsSearch string = "SELECT title, source, body, created_at, id FROM articles, to_tsvector(title) tvt, to_tsquery($1) tvq WHERE tvt @@ tvq ORDER BY ts_rank(tvt, tvq) DESC LIMIT 10"

func GetArticle(id int) (Article, error) {
	var title, createdAt, source, body string
	var ai int
	err := db.QueryRow("SELECT * FROM articles WHERE id = $1", id).Scan(&title, &source, &body, &createdAt, &ai)
	article := Article{Title: title, Date: createdAt, Source: source, Body: body, Id: ai}
	switch {
	case err != nil:
		fmt.Println(err)
		return article, err
	}
	return article, nil
}

func SearchArticles(query string) (ArticleResponse, error) {
	var response = ArticleResponse{Articles: make([]Article, 0)}
	query = strings.Join(strings.Split(query, " "), " | ")
	rows, err := db.Query(ftsSearch, query)
	if err != nil {
		return response, err
	}
	defer rows.Close()
	for rows.Next() {
		var title, createdAt, source, body string
		var ai int
		err = rows.Scan(&title, &source, &body, &createdAt, &ai)
		if err != nil {
			fmt.Println("Got an error here..", err)
			continue
		}
		article := Article{Title: title, Date: createdAt, Source: source, Body: body, Id: ai}
		response.Articles = append(response.Articles, article)
	}
	return response, nil
}
