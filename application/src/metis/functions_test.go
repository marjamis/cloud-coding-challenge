package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

/*
Global variables for testing
*/
var (
	router *mux.Router
)

/*
Initialise the system for testing
*/
func init() {
	router = NewRouter()
}

/*
Helper script to make HTTPConnections to the router
*/
func HTTPConnection(path string, t *testing.T, method string, data io.Reader) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, data)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

/*
Test that the HealthCheck returns the expected results when there aren't any issues.
*/
func TestHealthCheck_OK(t *testing.T) {
	assert := assert.New(t)

	rr := HTTPConnection("/healthcheck", t, "GET", nil)
	assert.Equal(http.StatusOK, rr.Code, "Wrong status code")
	assert.Contains(rr.Body.String(), "DBHealhty?: true", "Wrong body content")
}

/*
Test that when routing match isn't found the right HTTP code is returned
TODO: expand this to content as well.
*/
func TestNotFound_OK(t *testing.T) {
	rr := HTTPConnection("/b", t, "GET", nil)
	assert.Equal(t, http.StatusNotFound, rr.Code, "Wrong status code")
}

/*
Test the the index will return the correct page and with the right HTTP code
TODO: this should also be quite heavily expanded
*/
func TestIndex_OK(t *testing.T) {
	assert := assert.New(t)

	rr := HTTPConnection("/", t, "GET", nil)

	assert.Equal(http.StatusOK, rr.Code, "Wrong status code")
	assert.Contains(rr.Body.String(), "Customer Index page hosted from the instance", "Wrong body content")
}
