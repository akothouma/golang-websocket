package main

import (
	"net/http"
	// middleware "learn.zone01kisumu.ke/git/clomollo/forum/Middleware"
)

func (dep *Dependencies) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", dep.CSRFMiddleware(http.HandlerFunc(dep.HomeHandler)))
	mux.Handle("/post", dep.CSRFMiddleware(http.HandlerFunc(dep.PostHandler)))
	mux.Handle("/register", dep.CSRFMiddleware(http.HandlerFunc(dep.RegisterHandler)))
	mux.Handle("/logout", http.HandlerFunc(dep.LogoutHandler))
	mux.Handle("/login", dep.CSRFMiddleware(http.HandlerFunc(dep.LoginHandler)))
	mux.Handle("/styling/", http.StripPrefix("/styling/", http.FileServer(http.Dir("../../ui/static/styling"))))
	return mux
}
