package main

import (
	"flag"
	"log"
	"net/http"

	middleware "learn.zone01kisumu.ke/git/clomollo/forum/Middleware"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")
	if err := models.InitializeDB();err !=nil{
		log.Fatalf("Failed to initialize database: %v",err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeHandler)
	mux.Handle("/register", middleware.CSRFMiddleware(http.HandlerFunc(RegisterHandler)))
	mux.Handle("/logout", http.HandlerFunc(LogoutHandler))
	mux.Handle("/login", middleware.CSRFMiddleware(http.HandlerFunc(LoginHandler)))

	serv := &http.Server{
		Handler: mux,
		Addr:    *addr,
	}

	log.Printf("Starting server on port %v:", *addr)
	err := serv.ListenAndServe()
	log.Fatal(err)
}
