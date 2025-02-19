package main

import (
	"flag"
	"log"
	"net/http"

	middleware "learn.zone01kisumu.ke/git/clomollo/forum/Middleware"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/database"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/handlers"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")
	if err := database.InitializeDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Initialize dependencies
	deps := models.NewDependencies()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.Handle("/register", middleware.CSRFMiddleware(http.HandlerFunc(handlers.RegisterHandler)))
	mux.Handle("/logout", http.HandlerFunc(handlers.LogoutHandler))
	mux.Handle("/login", middleware.CSRFMiddleware(http.HandlerFunc(handlers.LoginHandler)))

	serv := &http.Server{
		Handler:  mux,
		Addr:     *addr,
		ErrorLog: deps.ErrorLog,
	}

	deps.InfoLog.Printf("Starting server on port %v:", *addr)
	err := serv.ListenAndServe()
	deps.ErrorLog.Fatal(err)
}
