package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

var (
	router http.Handler
	db     *sql.DB
	store  sessions.Store
)

type Author struct {
	Id          string
	Name        string
	Password    string
	Description string
}

func cryptoPassword(password string) string {
	if password == "" {
		return ""
	}
	return "SALT" + password
}

func checkUser(author *Author) error {
	password := author.Password
	err := db.QueryRow("SELECT password FROM authors WHERE name=$1", author.Name).Scan(&author.Password)
	if err != nil {
		return err
	}
	if author.Password != cryptoPassword(password) {
		return errors.New("db: password not match")
	}
	return nil
}

func getAuthor(name string, auth bool) (*Author, error) {
	var author Author
	var err error
	if auth {
		err = db.QueryRow("SELECT id, name, description FROM authors WHERE name=$1", name).Scan(&author.Id, &author.Name, &author.Description)
	} else {
		err = db.QueryRow("SELECT name, description FROM authors WHERE name=$1", name).Scan(&author.Name, &author.Description)
	}
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
	session, err := store.Get(r, "_session")
	auth := session.Values["logined"] == true && session.Values["username"] == params["Author"]
	author, err := getAuthor(params["Author"], auth)
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

func signupSubmitHandler(w http.ResponseWriter, r *http.Request) {
	var author Author
	author.Name, author.Password = r.PostFormValue("username"), r.PostFormValue("password")
	author.Description = r.PostFormValue("description")
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
	http.Redirect(w, r, fmt.Sprintf("/Articles/%s", author.Name), http.StatusSeeOther)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/signin.html")
	if err != nil {
		log.Panic(err)
	}
	params := mux.Vars(r)
	params["Username"] = r.FormValue("username")
	switch r.FormValue("err") {
	case "authors_name_nonexist":
		params["Error"] = "Username not exists"
	case "authors_password_notmatch":
		params["Error"] = "Invalid password"
	}
	tmpl.Execute(w, params)
}

func signinSubmitHandler(w http.ResponseWriter, r *http.Request) {
	var author Author
	author.Name, author.Password = r.PostFormValue("username"), r.PostFormValue("password")
	err := checkUser(&author)
	if err != nil {
		params := ""
		switch err.Error() {
		case "sql: no rows in result set":
			params = url.Values{"err": {"authors_name_nonexist"}, "username": {author.Name}}.Encode()
		case "db: password not match":
			params = url.Values{"err": {"authors_password_notmatch"}, "username": {author.Name}}.Encode()
		}
		http.Redirect(w, r, "/Articles/Sign-In?"+params, http.StatusFound)
		return
	}
	session, err := store.Get(r, "_session")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("%#v\n", session)
	session.Values["username"] = author.Name
	session.Values["logined"] = true
	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("/Articles/%s", author.Name), http.StatusSeeOther)
}

func initRouter() (http.Handler, *sql.DB, sessions.Store) {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	db, err := sql.Open("postgres", "port=9456 dbname=orangez sslmode=disable")
	if err != nil {
		panic("Open postgres failed")
	}
	router := mux.NewRouter()
	sub := router.PathPrefix("/Articles").Subrouter()
	sub.HandleFunc("/Sign-Up", signupHandler)
	sub.HandleFunc("/Sign-Up/Submit", signupSubmitHandler).Methods("POST")
	sub.HandleFunc("/Sign-In", signinHandler)
	sub.HandleFunc("/Sign-In/Submit", signinSubmitHandler).Methods("POST")
	sub.HandleFunc("/{Author}", authorHandler)
	store := pgstore.NewPGStore("port=9456 dbname=orangez sslmode=disable", []byte("something-secret"))
	return router, db, store
}

func forceHttps() http.Handler {
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
	router, db, store = initRouter()
	err := db.Ping()
	if err != nil {
		log.Panic(err)
	}
	go func() { log.Fatal(http.ListenAndServe(":8080", forceHttps())) }()
	log.Fatal(http.ListenAndServeTLS(":8443", "cert/orangez.cert.bundle.pem", "cert/orangez.key.pem", router))
}
