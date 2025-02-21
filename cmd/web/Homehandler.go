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

	posts := []PostView{
		{
			Post: models.Post{
				PostId: "1",
				UserId: "user123",
				// Category:    []string{"Tech", "GoLang"},
				Title:       "Go Template Issue",
				PostContent: "Fixing struct and template mismatches.",
				Media:       "https://example.com/image.jpg",
			},
			Comments: []string{"Great Post!", "Very helpful"},
			Likes:    10,
			Dislikes: 2,
		},
		{
			Post: models.Post{
				PostId: "2",
				UserId: "user456",
				// Category:    []string{"Programming"},
				Title:       "Understanding Golang",
				PostContent: "Golang is powerful for backend development.",
				Media:       "",
			},
			Comments: []string{"Nice explanation", "Thanks for sharing"},
			Likes:    15,
			Dislikes: 0,
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

// This function gets the likes, dislikes and comments
// that concerns the post using the PostID

/*func PostProcessor(post []models.Post) *[]PostView {
	var postViews []PostView

	for _, p := range post{
		comments := GetComments(p.PostId)
		likes := GetLikes(p.PostId)
		dislikes := GetDislikes(p.PostId)

		postViews = append(postViews, PostView{
			Post: p,
			Comments: comments,
			Likes: likes,
			Dislikes: dislikes,
		})
	}

	return &postViews
}
*/
