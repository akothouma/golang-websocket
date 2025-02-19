package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	middleware "learn.zone01kisumu.ke/git/clomollo/forum/Middleware"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/database"
)

// struct to hold application-wide dependencies
type Dependencies struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

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

	//Initialize new instance of the dependency struct
	//containing dependencies
	dep:=&Dependencies{
		ErrorLog: errorLog,
		InfoLog: infoLog,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", dep.HomeHandler)
	mux.Handle("/register", middleware.CSRFMiddleware(http.HandlerFunc(dep.RegisterHandler)))
	mux.Handle("/logout", http.HandlerFunc(dep.LogoutHandler))
	mux.Handle("/login", middleware.CSRFMiddleware(http.HandlerFunc(dep.LoginHandler)))

	serv := &http.Server{
		Handler:  mux,
		Addr:     *addr,
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on port %v:", *addr)
	err := serv.ListenAndServe()
	errorLog.Fatal(err)
}
