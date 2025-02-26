package main

import (
	"html/template"
	"log"
	"net/http"
)

// function to handle the Errors in the system all in one page
// prone to change if there exists a better one
func ErrorHandler(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	temp, err := template.ParseFiles("./ui/html/error.html")
	if err != nil {
		log.Println("Error while parsing the error page:", err)
		http.Error(w, "Page temporarily down", http.StatusInternalServerError)
		return
	}
	if err := temp.Execute(w, map[string]int{"Code": code}); err != nil {
		log.Println("Error while executing the error template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
