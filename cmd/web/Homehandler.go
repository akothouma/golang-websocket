package main

import (
	"html/template"
	"net/http"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

type PostView struct {
	Post     models.Post
	Comments []string
	Likes    int
	Dislikes int
}

func (dep *Dependencies) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse base, home, and post templates
	tmpl, err := template.ParseFiles(
		"./ui/templates/base.html",
		"./ui/templates/home.html",
		"./ui/templates/postContent.html",
		"./ui/templates/categories.html",
	)
	if err != nil {
		dep.ErrorLog.Println("Error loading template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// get all the posts data
	posts, err := models.AllPosts()
	if err != nil {
		dep.ErrorLog.Println("Error Retrieving Post data\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	p := PostProcessor(posts)

	// Execute the template with data
	err = tmpl.ExecuteTemplate(w, "base", p)
	if err != nil {
		dep.ErrorLog.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// This function gets the likes, dislikes and comments
// that concerns the post using the PostID

func PostProcessor(post []models.Post) *[]PostView {
	var postViews []PostView

	for _, p := range post {
		// comments := GetComments(p.PostId)
		// likes := GetLikes(p.PostId)
		// dislikes := GetDislikes(p.PostId)

		postViews = append(postViews, PostView{
			Post: p,
		})
	}

	return &postViews
}
