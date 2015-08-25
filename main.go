package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

var (
	router *mux.Router
	db     *sql.DB
)

type Author struct {
	Id          string
	Name        string
	Password    string
	Description string
}

func cryptoPassword(password string) string {
	return "SALT" + password
}

func getAuthor(name string) (*Author, error) {
	var author Author
	err := db.QueryRow("SELECT id, name, description FROM authors WHERE name=$1", name).Scan(&author.Id, &author.Name, &author.Description)
	return &author, err
}

func addAuthor(author *Author) error {
	_, err := db.Exec("INSERT INTO authors (name, password, description) VALUES ($1, $2, $3)", author.Name, cryptoPassword(author.Password), author.Description)
	return err
}

func authorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/author.html")
	if err != nil {
		log.Panic(err)
	}

	params := mux.Vars(r)
	author, err := getAuthor(params["Author"])
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.NotFound(w, r)
			return
		}
		log.Printf("%#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		tmpl.Execute(w, map[string]interface{}{"Author": author})
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		log.Panic(err)
	}
	params := mux.Vars(r)
	params["Username"] = r.FormValue("username")
	switch r.FormValue("err") {
	case "authors_name_key":
		params["Error"] = "Duplicate username"
	case "authors_name_character":
		params["Error"] = "Invalid username"
	}
	tmpl.Execute(w, params)
}

func newuserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var author Author
	if len(r.Form["username"]) != 1 {
		http.Error(w, "No username", http.StatusInternalServerError)
		return
	}
	author.Name = r.Form["username"][0]
	if len(r.Form["password"]) != 1 {
		http.Error(w, "No password", http.StatusInternalServerError)
		return
	}
	author.Password = r.Form["password"][0]
	switch len(r.Form["description"]) {
	case 0:
		author.Description = ""
	case 1:
		author.Description = r.Form["description"][0]
	default:
		http.Error(w, "Multiple description", http.StatusInternalServerError)
		return
	}
	err := addAuthor(&author)
	if err != nil {
		switch err := err.(type) {
		case *pq.Error:
			switch {
			case err.Code == "23505" && err.Constraint == "authors_name_key", err.Code == "23514" && err.Constraint == "authors_name_character":
				params := url.Values{"err": {err.Constraint}, "username": {author.Name}}.Encode()
				http.Redirect(w, r, "/Articles/Sign-Up?"+params, http.StatusFound)
				return
			}
		}
		log.Printf("%#v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "Success", http.StatusFound)
}

func initRouter() (*mux.Router, *sql.DB) {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	db, err := sql.Open("postgres", "port=9456 dbname=orangez sslmode=disable")
	if err != nil {
		panic("Open postgres failed")
	}
	router := mux.NewRouter()
	sub := router.PathPrefix("/Articles").Subrouter()
	sub.HandleFunc("/Sign-Up", newuserHandler).Methods("Post")
	sub.HandleFunc("/Sign-Up", signupHandler)
	sub.HandleFunc("/{Author}", authorHandler)
	return router, db
}

func forceHttps() *http.ServeMux {
	mux := http.NewServeMux()
	re := regexp.MustCompile(":8080$")
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		u := r.URL
		u.Scheme = "https"
		if r.Host != "" {
			u.Host = r.Host
		}
		u.Host = re.ReplaceAllString(u.Host, ":8443")
		log.Println(u)
		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	})
	return mux
}

func main() {
	router, db = initRouter()
	err := db.Ping()
	if err != nil {
		log.Panic(err)
	}
	go func() { log.Fatal(http.ListenAndServe(":8080", forceHttps())) }()
	log.Fatal(http.ListenAndServeTLS(":8443", "cert/orangez.cert.bundle.pem", "cert/orangez.key.pem", router))
}
