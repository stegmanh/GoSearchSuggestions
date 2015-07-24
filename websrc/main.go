package main

import (
	"net/http"

	"GoSearchSuggestions/websrc/app/routes"
	"GoSearchSuggestions/websrc/models"
)

//TODO: Move crawler information to models section

func main() {
	models.Init()

	http.ListenAndServe(":8080", routes.InitHandler())
}
