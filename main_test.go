package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server *httptest.Server
	url    string
)

func init() {
	router, db = initRouter()
	server = httptest.NewServer(router)
	url = server.URL
	fmt.Println(url)
}

func TestTest(t *testing.T) {
	assert.True(t, true, "Canary test")
	assert.Contains(t, "a", "a")
}

func TestAuthorGet(t *testing.T) {
	resp, err := http.Get(url + "/Articles/Clouds")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.NoError(t, err)
	assert.Contains(t, string(body), "Clouds")
}
