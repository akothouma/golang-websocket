package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")

	// DEBUG
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Current working directory:", cwd)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t",
		log.Ldate|log.Ltime|log.Lshortfile)

	// Creating a server
	serv := &http.Server{
		Handler:  Routes(),
		Addr:     *addr,
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on port %v:", *addr)
	fmt.Println("server started")
	err = serv.ListenAndServe()
	errorLog.Fatal(err)
}
