package main

import (
	"html/template"
	"net/http"
)

func (dep *Dependencies) HomeHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate, err := template.ParseFiles("./ui/templates/home.html")
	if (err == nil) {
		homeTemplate.Execute(w, nil)
	}else {
		http.Error(w, "Error loading home template", http.StatusNotFound)
		return
	}
	
}
