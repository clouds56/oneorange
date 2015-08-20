package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

var (
	router *mux.Router
	db     *sql.DB
)

type Author struct {
	Id          int
	Name        string
	Description string
}

func getAuthor(name string) (*Author, error) {
	var author Author
	err := db.QueryRow("SELECT id, name, description FROM authors WHERE name=$1", name).Scan(&author.Id, &author.Name, &author.Description)
	return &author, err
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/author.html")
	if err != nil {
		log.Fatal(err)
	}

	params := mux.Vars(r)
	author, err := getAuthor(params["Author"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
	} else {
		tmpl.Execute(w, map[string]interface{}{"Author": author})
	}
}

func initRouter() (*mux.Router, *sql.DB) {
	db, err := sql.Open("postgres", "port=9456 dbname=orangez sslmode=disable")
	if err != nil {
		panic("Open postgres failed")
	}
	router := mux.NewRouter()
	sub := router.PathPrefix("/Articles").Subrouter()
	sub.HandleFunc("/{Author}", handler)
	return router, db
}

func main() {
	router, db = initRouter()
	//log.Fatal(db.Ping())
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	log.Fatal(http.ListenAndServe(":8080", router))
}
