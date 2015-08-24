package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	server *httptest.Server
)

func init() {
	router, db = initRouter()
	server = httptest.NewServer(router)
	log.Println(server.URL)
}

type IT struct {
	t       *testing.T
	message string
	resp    *http.Response
	body    string
	parsed  bool
	failed  bool
}

func I(t *testing.T, message string) *IT {
	var it IT
	it.Settings(t, message)
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
		resp, err = http.PostForm(server.URL+url, data)
	} else if method == "GET" {
		resp, err = http.Get(server.URL + url)
	} else {
		it.t.Errorf("Unkown http method at %s", assert.CallerInfo())
		it.t.FailNow()
		return nil
	}
	if assert.NoError(it.t, err) && assert.NotNil(it.t, resp) {
		it.resp = resp
		return it
	}
	return it.FailNow("Failed : %s at %s", it.message, assert.CallerInfo())
}

func (it *IT) Redirect(url string) *IT {
	if it.failed {
		return it
	}
	// if it.resp.StatusCode >= http.StatusMultipleChoices && it.resp.StatusCode <= http.StatusTemporaryRedirect {
	// 	return it
	// } else {
	// 	assert.Fail(it.t, fmt.Sprintf("Expected StatusCode between %d and %d but get %d", http.StatusMultipleChoices, http.StatusTemporaryRedirect, it.resp.StatusCode))
	// 	it.t.Fail()
	// }
	if assert.Equal(it.t, it.resp.Request.URL.Path, url) {
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
	rows, err := db.Query(query, args...)
	if assert.NoError(dt.t, err) {
		dt.rows = rows
		return dt
	}
	return dt.FailNow("Failed : %s at %s", dt.message, assert.CallerInfo())
}

func (dt *DT) Exec(query string, args ...interface{}) *DT {
	_, err := db.Exec(query, args...)
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
	assert.NotNil(t, router)
}

func TestSignup(t *testing.T) {
	I(t, "should be a sign up html").Method("GET", "/Articles/Sign-Up", nil).Contains(`type="submit"`)
	if D(t, "clear author").Exec("DELETE FROM authors WHERE name=$1", "test1").PASS() {
		D(t, "shouldn't have test1").Query("SELECT * FROM authors WHERE name=$1", "test1").Empty()
		I(t, "should be a success page").Method("POST", "/Articles/Sign-Up", url.Values{"username": {"test1"}, "password": {"123456"}, "description": {"lazy and nothing"}}).Redirect("/Articles/Success")
		D(t, "shouldn't have test1").Query("SELECT * FROM authors WHERE name=$1", "test1").NonEmpty()
		I(t, "should be a 500 page").Method("POST", "/Articles/Sign-Up", url.Values{"username": {"test1"}, "password": {"123456"}, "description": {"lazy and nothing"}}).
			Contains(`pq: duplicate key value violates unique constraint "authors_name_key"`)
		D(t, "clear author again").Exec("DELETE FROM authors WHERE name=$1", "test1").PASS()
	}
}

func TestAuthorGet(t *testing.T) {
	I(t, "should contains Clouds").Method("GET", "/Articles/Clouds", nil).Contains("Clouds")
}
