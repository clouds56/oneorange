package main

import (
	"fmt"
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func initRouter() *mux.Router {
	var router = mux.NewRouter()
	router.HandleFunc("/{author}", handler)
	return router
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", initRouter()))
}
