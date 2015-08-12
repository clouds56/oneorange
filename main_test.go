package main

import (
	"fmt"
	"testing"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
)

var (
	server *httptest.Server
	url string
)

func init() {
	server = httptest.NewServer(initRouter())
	url = server.URL
	fmt.Println(url)
}

func TestTest(t *testing.T) {
	assert.True(t, true, "Canary test")
	assert.Contains(t, "a", "a")
}

func TestAuthorGet(t *testing.T) {
	resp, err := http.Get(url+"/clouds")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.NoError(t, err)
	assert.Contains(t, string(body), "clouds")
}
