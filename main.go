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

	tmpl.Execute(w, map[string]string{"Author":r.URL.Path[1:]})
}

func initRouter() *mux.Router {
	var router = mux.NewRouter()
	router.HandleFunc("/{author}", handler)
	return router
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", initRouter()))
}
