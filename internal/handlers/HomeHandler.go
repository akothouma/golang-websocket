package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Welcome to Forum")

	homeTemplate, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		http.Error(w, "Error loading home template", http.StatusNotFound)
		return
	}
	posts, err := models.GetAllPosts()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve posts: %v", err), http.StatusInternalServerError)
		return
	}

	homeTemplate.Execute(w, posts)
}

