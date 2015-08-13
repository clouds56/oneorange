package main

import (
	"html/template"
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/author.html")
	if err != nil {
		log.Fatal(err)
	}

	params := mux.Vars(r)
	tmpl.Execute(w, params)
}

func initRouter() *mux.Router {
	router := mux.NewRouter()
	sub := router.PathPrefix("/Articles").Subrouter()
	sub.HandleFunc("/{Author}", handler)
	return router
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", initRouter()))
}
