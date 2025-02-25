package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

// struct to hold application-wide dependencies
type Dependencies struct {
	ErrorLog  *log.Logger
	InfoLog   *log.Logger
	Forum     *models.ForumModel
	Templates *template.Template
}

func main() {
	// DEBUG
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Current working directory:", cwd)

	addr := flag.String("addr", ":8000", "HTTP network address")
	// leveled logging, informational messages output to standard out(stdout)
	// Error messages output to standard error(stderr)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t",
		log.Ldate|log.Ltime|log.Lshortfile)
	db, err := models.InitializeDB()
	if err != nil {
		errorLog.Fatalf("Failed to initialize database: %v", err)
	}

	defer db.Close()

	// initializing dependencies
	dep := &Dependencies{
		ErrorLog:  errorLog,
		InfoLog:   infoLog,
		Forum:     &models.ForumModel{DB: db},
	}

	// Creating a server
	serv := &http.Server{
		Handler:  dep.Routes(),
		Addr:     *addr,
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on port %v:", *addr)
	err = serv.ListenAndServe()
	errorLog.Fatal(err)
}
