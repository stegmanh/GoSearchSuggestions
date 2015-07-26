package main

import (
	"net/http"
	"os"

	"GoSearchSuggestions/websrc/app/routes"
	"GoSearchSuggestions/websrc/models"
	"github.com/gorilla/handlers"
)

func main() {
	models.Init()
	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	http.Handle("/", handlers.LoggingHandler(logFile, routes.GetContext(routes.InitHandler())))

	http.ListenAndServe(":8080", nil)
}
