package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	middleware "learn.zone01kisumu.ke/git/clomollo/forum/Middleware"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/database"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/handlers"
)

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")
	if err := database.InitializeDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// leveled logging, informational messages output to standard out(stdout)
	// Error messages output to standard error(stderr)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t",
log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.Handle("/register", middleware.CSRFMiddleware(http.HandlerFunc(handlers.RegisterHandler)))
	mux.Handle("/logout", http.HandlerFunc(handlers.LogoutHandler))
	mux.Handle("/login", middleware.CSRFMiddleware(http.HandlerFunc(handlers.LoginHandler)))

	serv := &http.Server{
		Handler: mux,
		Addr:    *addr,
	}

	infoLog.Printf("Starting server on port %v:", *addr)
	err := serv.ListenAndServe()
	errorLog.Fatal(err)
}
