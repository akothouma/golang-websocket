package main

import (
	"fmt"
	"net/http"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

func (dep *Dependencies) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// wd, err := os.Getwd()
	// if err != nil {
	// 	dep.ErrorLog.Println("Error getting the working directory:\n", err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// // constructing the correct path
	// tmplPath := wd + "/ui/html/"

	// // Parse base, home, and post templates
	// tmpl, err := template.ParseFiles(
	// 	tmplPath + "posts.html")
	// if err != nil {
	// 	dep.ErrorLog.Println("Error loading template:\n", err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }
	allTemplates, err := models.InitTemplates()
	if err != nil {
		dep.ErrorLog.Println("Error loading template:\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// HANDLING THE FILTER FORM
	// Getting the values of the filter form
	err = r.ParseForm()
	if err != nil {
		dep.ErrorLog.Println("Error parsing filter form:\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	var filteredPosts interface{}
	categories := r.Form["category"]
	// DEBUG
	fmt.Println("This are the selected categories\n", categories)
	if len(categories) == 0 {
		// get all the posts data
		posts, err := models.RenderPostsPage()
		if err != nil {
			dep.ErrorLog.Println("Error Retrieving Post data\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// DEBUG
		fmt.Println("This is the posts data >>", posts)
		// Execute the template with data
		
		err = allTemplates.ExecuteTemplate(w, "base", nil)
		err = allTemplates.ExecuteTemplate(w, "content", posts)
		if err != nil {
			dep.ErrorLog.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// RenderTemplates(w, "posts.html", posts)
	} else {
		for _, cat := range categories {
			filteredPosts, err = dep.Forum.Filters(cat)
			if err != nil {
				dep.ErrorLog.Println("Error Filtering the posts:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			// DEBUG
			fmt.Println("This are the filtered posts:\n", filteredPosts)
		}
		allTemplates.ExecuteTemplate(w, "base", nil)
		allTemplates.ExecuteTemplate(w, "content", filteredPosts)
	}
}

// This function gets the likes, dislikes and comments
// that concerns the post using the PostID
