package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

// struct to hold application-wide dependencies
type Dependencies struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Forum    *models.ForumModel
}

func main() {
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

	// Initialize new instance of the dependency struct
	// containing dependencies
	dep := &Dependencies{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		Forum:    &models.ForumModel{DB: db},
	}

	serv := &http.Server{
		Handler:  dep.Routes(),
		Addr:     *addr,
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on port %v:", *addr)
	err = serv.ListenAndServe()
	errorLog.Fatal(err)
}
