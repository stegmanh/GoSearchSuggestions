package routes

import (
	"fmt"
	"net/http"

	"GoSearchSuggestions/websrc/app/auth"
	ctx "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func InitHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", Use(Home, GetContext))
	router.HandleFunc("/login", Use(Login, GetContext))
	router.HandleFunc("/dashboard", Use(Dashboard, Authenticated, GetContext)).Methods("GET")
	router.HandleFunc("/articles", DbSearchHandler).Methods("GET")
	router.HandleFunc("/autocomplete", SearchHandler).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./websrc/static/")))).Methods("GET")
	return router
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Name is %s ", "Holden")))
}

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./websrc/views/index.html")
}

func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}

func Authenticated(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if u := ctx.Get(r, "user"); u != nil {
			handler.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/", 302)
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	session := ctx.Get(r, "session").(*sessions.Session)
	success, err := auth.Login(r)
	if err != nil {
		fmt.Println(err)
	}
	if success {
		session.Save(r, w)
		http.Redirect(w, r, "/dashboard", 302)
	} else {
		http.Redirect(w, r, "/", http.StatusForbidden)
	}
}

func GetContext(handler http.Handler) http.HandlerFunc {
	// Set the context here
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request form
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing request", http.StatusInternalServerError)
		}
		// Set the context appropriately here.
		// Set the session
		session, _ := auth.Store.Get(r, "admin")
		// Put the session in the context so that
		ctx.Set(r, "session", session)
		if _, ok := session.Values["id"]; ok {
			if err != nil {
				ctx.Set(r, "user", nil)
			} else {
				ctx.Set(r, "user", "holden")
			}
		} else {
			ctx.Set(r, "user", nil)
		}
		handler.ServeHTTP(w, r)
		// Remove context contents
		ctx.Clear(r)
	}
}
