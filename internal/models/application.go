package models

import (
	"log"
	"os"
)

// struct to hold application-wide dependencies
type Dependencies struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

// NewDependencies initializes and returns a Dependencies struct.
func NewDependencies() *Dependencies {
	return &Dependencies{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
