package router

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
func requestGeneration(path string, t *testing.T, method string, data io.Reader) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, path, data)
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr, err
}

func TestHealthCheck_OK(t *testing.T) {
	rr, err := requestGeneration("/healthcheck", t, "GET", nil)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, http.StatusOK, rr.Code, "Wrong status code")
	assert.Contains(t, rr.Body.String(), "DBHealhty?: true", "Wrong body content")
}

func TestNotFound_OK(t *testing.T) {
	rr, err := requestGeneration("/b", t, "GET", nil)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, http.StatusNotFound, rr.Code, "Wrong status code")
}

func TestIndex_OK(t *testing.T) {
	rr, err := requestGeneration("/", t, "GET", nil)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, http.StatusOK, rr.Code, "Wrong status code")
	assert.Contains(t, rr.Body.String(), "Customer Index page hosted from the instance", "Wrong body content")
}
