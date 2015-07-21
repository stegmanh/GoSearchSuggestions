package models

import (
	"database/sql"
	"fmt"
)

type Article struct {
	Title  string `json:"title"`
	Date   string `json:"date"`
	Source string `json:"source"`
	Body   string `json:"body"`
}

type ArticleResponse struct {
	Articles []Article `json:"data"`
}

var ftsSearch string = "SELECT title, source, body, created_at FROM articles, to_tsvector(title) tvt, to_tsquery($1) tvq WHERE tvt @@ tvq ORDER BY ts_rank(tvt, tvq) DESC"

func GetArticle(id int) (Article, error) {
	var title, createdAt, source, body, ai string
	err := Db.QueryRow("SELECT * FROM articles WHERE id = $1", id).Scan(&title, &source, &body, &createdAt, &ai)
	article := Article{Title: title, Date: createdAt, Source: source, Body: body}
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("No Article with Id")
		return article, err
	case err != nil:
		fmt.Println(err)
		return article, err
	}
	return article, nil
}
