package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var (
	server *httptest.Server
)

func init() {
	app = initRouter()
	server = httptest.NewServer(app.Router)
	log.Println(server.URL)
}

type IT struct {
	t       *testing.T
	message string
	client  http.Client
	resp    *http.Response
	body    string
	parsed  bool
	failed  bool
}

func I(t *testing.T, message string) *IT {
	var it IT
	it.Settings(t, message)
	jar, err := cookiejar.New(nil)
	if !assert.NoError(t, err) {
		it.failed = true
	}
	it.client = http.Client{Jar: jar}
	return &it
}

func (it *IT) Settings(t *testing.T, message string) *IT {
	it.t = t
	it.message = message
	it.failed = false
	return it
}

func (it *IT) FailNow(format string, args ...interface{}) *IT {
	it.t.Errorf(format, args...)
	it.t.Fail()
	it.failed = true
	return it
}

func (it *IT) Method(method, url string, data url.Values) *IT {
	if it.failed {
		return it
	}
	var resp *http.Response
	var err error
	if method == "POST" {
		resp, err = it.client.Post(server.URL+url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	} else if method == "GET" {
		resp, err = it.client.Get(server.URL + url)
	} else {
		return it.FailNow("Unkown http method at %s", assert.CallerInfo())
	}
	if assert.NoError(it.t, err) && assert.NotNil(it.t, resp) {
		it.resp = resp
		it.parsed = false
		return it
	}
	return it.FailNow("Failed : %s at %s", it.message, assert.CallerInfo())
}

func (it *IT) Redirect(url string) *IT {
	if it.failed {
		return it
	}
	if assert.Equal(it.t, url, it.resp.Request.URL.Path) {
		return it
	}
	return it.FailNow("Failed : %s at %s", it.message, assert.CallerInfo())
}

func (it *IT) HttpCode(code int) *IT {
	if it.failed {
		return it
	}
	if assert.Equal(it.t, code, it.resp.StatusCode) {
		return it
	}
	return it.FailNow("Failed : %s at %s", it.message, assert.CallerInfo())
}

func (it *IT) ParseBody() *IT {
	if it.failed {
		return it
	}
	if it.parsed {
		return it
	}
	body, err := ioutil.ReadAll(it.resp.Body)
	it.resp.Body.Close()
	if assert.NoError(it.t, err) && assert.NotNil(it.t, body) {
		it.parsed = true
		it.body = string(body)
		return it
	}
	return it.FailNow("Failed : %s at %s", it.message, assert.CallerInfo())
}

func (it *IT) Contains(str string) *IT {
	if it.failed {
		return it
	}
	it.ParseBody()
	if assert.Contains(it.t, it.body, str) {
		return it
	}
	return it.FailNow("Failed : %s at %s", it.message, assert.CallerInfo())
}

func (it *IT) NotContains(str string) *IT {
	if it.failed {
		return it
	}
	it.ParseBody()
	if assert.NotContains(it.t, it.body, str) {
		return it
	}
	return it.FailNow("Failed : %s at %s", it.message, assert.CallerInfo())
}

func (it *IT) PASS() bool {
	return !it.failed
}

type DT struct {
	t       *testing.T
	message string
	rows    *sql.Rows
	failed  bool
}

func D(t *testing.T, message string) *DT {
	var dt DT
	dt.Settings(t, message)
	return &dt
}

func (dt *DT) Settings(t *testing.T, message string) *DT {
	dt.t = t
	dt.message = message
	dt.failed = false
	return dt
}

func (dt *DT) FailNow(format string, args ...interface{}) *DT {
	dt.t.Errorf(format, args...)
	dt.t.Fail()
	dt.failed = true
	return dt
}

func (dt *DT) Query(query string, args ...interface{}) *DT {
	rows, err := app.DB.Query(query, args...)
	if assert.NoError(dt.t, err) {
		dt.rows = rows
		return dt
	}
	return dt.FailNow("Failed : %s at %s", dt.message, assert.CallerInfo())
}

func (dt *DT) Exec(query string, args ...interface{}) *DT {
	_, err := app.DB.Exec(query, args...)
	if assert.NoError(dt.t, err) {
		return dt
	}
	return dt.FailNow("Failed : %s at %s", dt.message, assert.CallerInfo())
}

func (dt *DT) Empty() *DT {
	if assert.False(dt.t, dt.rows.Next()) {
		return dt
	}
	return dt.FailNow("Failed : %s at %s", dt.message, assert.CallerInfo())
}

func (dt *DT) NonEmpty() *DT {
	if assert.True(dt.t, dt.rows.Next()) {
		return dt
	}
	return dt.FailNow("Failed : %s at %s", dt.message, assert.CallerInfo())
}

func (dt *DT) PASS() bool {
	return !dt.failed
}

func TestTest(t *testing.T) {
	assert.True(t, true, "Canary test")
	assert.Contains(t, "a", "a")
	assert.NotNil(t, app.Router)
	assert.NotNil(t, app.DB)
	assert.NotNil(t, app.Store)
	if !D(t, "shouldn't have testNonExist").Query("SELECT * FROM authors WHERE name=$1", "testNonExist").Empty().PASS() {
		D(t, "clear testNonExist").Exec("DELETE FROM authors WHERE name=$1", "testNonExist")
		D(t, "shouldn't have testNonExist").Query("SELECT * FROM authors WHERE name=$1", "testNonExist").Empty()
	}
	if !D(t, "should have testExist").Query("SELECT * FROM authors WHERE name=$1", "testExist").NonEmpty().PASS() {
		D(t, "create testExist").Exec("INSERT INTO authors (name, password, description) VALUES ($1, $2, $3)", "testExist", cryptoPassword("123"), "lazy and nothing")
		D(t, "should have testExist").Query("SELECT * FROM authors WHERE name=$1", "testExist").NonEmpty()
	}
	if !D(t, "should have Clouds").Query("SELECT * FROM authors WHERE name=$1", "Clouds").NonEmpty().PASS() {
		D(t, "clear author").Exec("INSERT INTO authors (name, password, description) VALUES ($1, $2, $3)", "Clouds", cryptoPassword("zxc"), "seven square and seven")
		D(t, "should have Clouds").Query("SELECT * FROM authors WHERE name=$1", "Clouds").NonEmpty()
	}
}

func TestSignup(t *testing.T) {
	I(t, "should be a sign up html").Method("GET", "/Articles/Sign-Up", nil).Contains("Sign up").Contains(`type="submit"`)
	I(t, "should be a sign up html").Method("GET", "/Articles/Sign-Up?err=authors_name_character&username=a.c", nil).
		Contains("Sign up").Contains(`type="submit"`).Contains("Invalid username").Contains("a.c")

	I(t, "should be a success page").Method("POST", "/Articles/Sign-Up/Submit", url.Values{"username": {"testNonExist"}, "password": {"123456"}, "description": {"lazy and nothing"}}).
		Redirect("/Articles/testNonExist")
	D(t, "should have testNonExist").Query("SELECT * FROM authors WHERE name=$1", "testNonExist").NonEmpty()
	D(t, "clear testNonExist").Exec("DELETE FROM authors WHERE name=$1", "testNonExist")

	I(t, "should show authors_name_key error").Method("POST", "/Articles/Sign-Up/Submit", url.Values{"username": {"testExist"}, "password": {"123456"}, "description": {"lazy and nothing"}}).
		Contains("Duplicate username")
	D(t, "should have testExist").Query("SELECT * FROM authors WHERE name=$1", "testExist").NonEmpty()

	I(t, "should show authors_name_character error").Method("POST", "/Articles/Sign-Up/Submit", url.Values{"username": {"a..b"}, "password": {"123456"}, "description": {"lazy and nothing"}}).
		Contains("Invalid username")
	I(t, "should show authors_name_character error").Method("POST", "/Articles/Sign-Up/Submit", url.Values{"username": {"a#bc"}, "password": {"123456"}, "description": {"lazy and nothing"}}).
		Contains("Invalid username")
	I(t, "should show authors_name_character error").Method("POST", "/Articles/Sign-Up/Submit", url.Values{"username": {"a1"}, "password": {"123456"}, "description": {"lazy and nothing"}}).
		Contains("Invalid username")
}

func TestSignin(t *testing.T) {
	I(t, "should be a sign in html").Method("GET", "/Articles/Sign-In", nil).Contains("Sign in").Contains(`type="submit"`)

	I(t, "should be a success page").Method("POST", "/Articles/Sign-In/Submit", url.Values{"username": {"testExist"}, "password": {"123"}}).
		Redirect("/Articles/testExist")

	I(t, "should show authors_name_nonexist error").Method("POST", "/Articles/Sign-In/Submit", url.Values{"username": {"testNonExist"}, "password": {"123"}}).
		Contains("Username not exists")
	I(t, "should show authors_password_notmatch error").Method("POST", "/Articles/Sign-In/Submit", url.Values{"username": {"testExist"}, "password": {"321"}}).
		Contains("Invalid password")

	I(t, "should be a 500 page").Method("POST", "/Articles/Sign-In/Submit", url.Values{"username": {"testExist"}}).
		Contains("Invalid password")
}

func TestAuthorGet(t *testing.T) {
	I(t, "should contains Clouds").Method("GET", "/Articles/Clouds", nil).Contains("Clouds").NotContains("ID")
	I(t, "should contains Clouds and ID").Method("POST", "/Articles/Sign-In/Submit", url.Values{"username": {"Clouds"}, "password": {"zxc"}}).
		Redirect("/Articles/Clouds").Contains("Clouds").Contains("ID").
		Method("GET", "/Articles/testExist", nil).NotContains("ID").
		Method("GET", "/Articles/Clouds", nil).Contains("ID")
	I(t, "should be a 404 page").Method("GET", "/Articles/testNonExist", nil).HttpCode(404)
}
