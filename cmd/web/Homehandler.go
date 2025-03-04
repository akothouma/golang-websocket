package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

type PostView struct {
	Post     models.Post
	Comments []string
	Likes    int
	Dislikes int
}

func (dep *Dependencies) HomeHandler(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		dep.ErrorLog.Println("Error getting the working directory:\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// constructing the correct path
	tmplPath := wd + "/ui/templates/"

	// Parse base, home, and post templates
	tmpl, err := template.ParseFiles(
		tmplPath+"base.html",
		tmplPath+"home.html",
		tmplPath+"postContent.html",
		tmplPath+"categories.html",
	)
	if err != nil {
		dep.ErrorLog.Println("Error loading template:\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// HANDLING THE FILTER FORM
	// if r.Method != http.MethodPost{
	// 	dep.ErrorLog.Println("Error during filtration of categories\n", err)
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	err = r.ParseForm()
	if err != nil {
		dep.ErrorLog.Println("Error parsing filter form:\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	categories := r.Form["category"]
	//DEBUG
	fmt.Println("This are the selected categories\n", categories)
	if len(categories) == 0 {
		// get all the posts data
		posts, err := models.RenderPostsPage()
		if err != nil {
			dep.ErrorLog.Println("Error Retrieving Post data\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		//DEBUG
		fmt.Println("This is the posts data >>", posts)
		// Execute the template with data
		err = tmpl.ExecuteTemplate(w, "base", posts)
		if err != nil {
			dep.ErrorLog.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		for _, cat := range categories {
			filteredPosts, err := dep.Forum.Filters(cat)
			if err != nil {
				dep.ErrorLog.Println("Error Filtering the posts:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			// DEBUG
			fmt.Println("This are the filtered posts:\n", filteredPosts)
		}
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
