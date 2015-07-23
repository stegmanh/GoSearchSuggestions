package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func InitHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", Home)
	router.HandleFunc("/dashboard", Dashboard)
	router.HandleFunc("/articles", DbSearchHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./websrc/static/"))))
	return router
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Name is %s ", "Holden")))
}

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./websrc/views/index.html")
}
