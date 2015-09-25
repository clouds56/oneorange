package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/lib/pq"
)

type Application struct {
	Router    http.Handler
	DB        *sql.DB
	Store     sessions.Store
	Templates map[string]*template.Template
}

var app Application

type Data map[string]interface{}
type TmplData struct {
	Request *http.Request
	Params  map[string]string
	Session *sessions.Session
	Data    Data
}

type Author struct {
	Id          string
	Name        string
	Password    string
	Description string
}

type Anthology struct {
	Id          string
	Name        string
	Author      *Author
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
	err := app.DB.QueryRow("SELECT password FROM authors WHERE name=$1", author.Name).Scan(&author.Password)
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
		err = app.DB.QueryRow("SELECT id, name, description FROM authors WHERE name=$1", name).Scan(&author.Id, &author.Name, &author.Description)
	} else {
		err = app.DB.QueryRow("SELECT name, description FROM authors WHERE name=$1", name).Scan(&author.Name, &author.Description)
	}
	return &author, err
}

func addAuthor(author *Author) error {
	_, err := app.DB.Exec("INSERT INTO authors (name, password, description) VALUES ($1, $2, $3)", author.Name, cryptoPassword(author.Password), author.Description)
	return err
}

func getAnthology(name string, author_name string, auth bool) (*Anthology, error) {
	var author Author
	var anthology Anthology
	anthology.Author = &author
	var err error
	query := `SELECT
			%s anthologies.name, anthologies.description, authors.name
		FROM
			anthologies, authors
		WHERE
			anthologies.name=$1 AND authors.name=$2 AND authors.id=anthologies.author_id`
	if auth {
		query = fmt.Sprintf(query, "anthologies.id,")
		err = app.DB.QueryRow(query, name, author_name).Scan(&anthology.Id, &anthology.Name, &anthology.Description, &anthology.Author.Name)
	} else {
		query = fmt.Sprintf(query, "")
		err = app.DB.QueryRow(query, name, author_name).Scan(&anthology.Name, &anthology.Description, &anthology.Author.Name)
	}
	return &anthology, err
}

func newSession(w http.ResponseWriter, r *http.Request, username string) *sessions.Session {
	session, err := app.Store.Get(r, "_session")
	if err != nil {
		log.Println(err.Error())
		// ignored and works
	}
	// log.Printf("%#v\n", session)
	session.Values["username"] = username
	session.Values["logined"] = true
	session.Save(r, w)
	// log.Printf("%#v\n", session)
	return session
}

func getSession(w http.ResponseWriter, r *http.Request) (*sessions.Session, string) {
	username := ""
	session, err := app.Store.Get(r, "_session")
	if err != nil || session.Values["logined"] != true {
		session.Options.MaxAge = -1
		sessions.Save(r, w)
		return session, username
	}
	if session.Values["logined"] == true {
		username = session.Values["username"].(string)
	}
	return session, username
}

func authorHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	session, err := app.Store.Get(r, "_session")
	// log.Printf("%#v\n", session)
	session, username := getSession(w, r)
	author, err := getAuthor(params["Author"], username == params["Author"])
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.NotFound(w, r)
			return
		}
		log.Printf("%#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		app.Templates["author"].Execute(w, TmplData{r, params, session, Data{"Author": author}})
	}
}

func anthologyHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	session, username := getSession(w, r)
	anthology, err := getAnthology(params["Anthology"], params["Author"], username == params["Author"])
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.NotFound(w, r)
			return
		}
		log.Printf("%#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		app.Templates["anthology"].Execute(w, TmplData{r, params, session, Data{"Anthology": anthology}})
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	data := Data{}
	switch r.FormValue("err") {
	case "authors_name_key":
		data["Error"] = "Duplicate username"
	case "authors_name_character":
		data["Error"] = "Invalid username"
	}
	app.Templates["signup"].Execute(w, TmplData{r, params, nil, data})
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
	params := mux.Vars(r)
	data := Data{}
	switch r.FormValue("err") {
	case "authors_name_nonexist":
		data["Error"] = "Username not exists"
	case "authors_password_notmatch":
		data["Error"] = "Invalid password"
	}
	app.Templates["signin"].Execute(w, TmplData{r, params, nil, data})
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
	newSession(w, r, author.Name)
	http.Redirect(w, r, fmt.Sprintf("/Articles/%s", author.Name), http.StatusSeeOther)
}

func initRouter() Application {
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
	sub.HandleFunc("/{Author}/{Anthology}", anthologyHandler)
	store := pgstore.NewPGStore("port=9456 dbname=orangez sslmode=disable", []byte("something-secret"))
	store.Cleanup(time.Minute * 5)
	tmpls := map[string]*template.Template{}
	for _, t := range []string{"author", "anthology", "signin", "signup"} {
		tmpls[t] = template.Must(template.ParseFiles(fmt.Sprintf("templates/%s.html", t)))
	}
	return Application{router, db, store, tmpls}
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
	app = initRouter()
	err := app.DB.Ping()
	if err != nil {
		log.Panic(err)
	}
	go func() { log.Fatal(http.ListenAndServe(":8080", forceHttps())) }()
	log.Fatal(http.ListenAndServeTLS(":8443", "cert/orangez.cert.bundle.pem", "cert/orangez.key.pem", app.Router))
}
