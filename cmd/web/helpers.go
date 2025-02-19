package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// ServerError helper writes an error message and stack trace to the
// Errolog
// Sends generic 500 internal server error response to the user
func (dep *Dependencies) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	dep.ErrorLog.Print(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and
// corresponding description to user
func (dep *Dependencies) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// A wrapper around clientError which sends a 404 not found response
// to user
func (dep *Dependencies) NotFound(w http.ResponseWriter) {
	dep.ClientError(w, http.StatusNotFound)
}
