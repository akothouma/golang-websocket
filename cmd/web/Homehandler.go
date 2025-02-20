package main

import (
	"html/template"
	"net/http"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

func (dep *Dependencies) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse base, home, and post templates
	tmpl, err := template.ParseFiles(
		"../../ui/templates/base.html",
		"../../ui/templates/home.html",
		"../../ui/templates/postContent.html",
		"../../ui/templates/categories.html",
	)
	if err != nil {
		dep.ErrorLog.Println("Error loading template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Dummy Data for a post
	
	posts := []models.Post{
		{
			PostId:      "1",
			UserId:      "user123",
			Category:    []string{"Tech", "GoLang"},
			Title:       "Go Template Issue",
			PostContent: "Fixing struct and template mismatches.",
			Media:       "https://example.com/image.jpg",
		}, 
		{
			PostId:      "2",
			UserId:      "user456",
			Category:    []string{"Programming"},
			Title:       "Understanding Golang",
			PostContent: "Golang is powerful for backend development.",
			Media:       "",
		}, 
	}

	// Execute the template with data
	err = tmpl.ExecuteTemplate(w, "base", posts) 
	if err != nil {
		dep.ErrorLog.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
