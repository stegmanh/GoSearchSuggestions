package main

import (
	"net/http"

	"GoSearchSuggestions/websrc/app/routes"
	"GoSearchSuggestions/websrc/models"
)

func main() {
	models.Init()

	http.ListenAndServe(":8080", routes.InitHandler())
}
