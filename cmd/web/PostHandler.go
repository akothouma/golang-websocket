package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

// /home/clomollo/forum/ui/html/posts.html
func (dep *Dependencies) PostHandler(w http.ResponseWriter, r *http.Request) {
	// Load the template file to use. ("posts")
	// PostTemplate, err := template.ParseFiles("./ui/html/allposts.html")
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, "NOT FOUND\nError parsing post templates", http.StatusNotFound)
	// 	return
	// }
	// if r.Method == http.MethodGet {
	// 	csrfToken,ok:= r.Context().Value("csrf_token").(string)
	// 	fmt.Println("here",csrfToken)
	// 	if !ok || csrfToken == "" {
	// 		http.Error(w, "Internal Server Error: CSRF token missing or invalid", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	PostTemplate.ExecuteTemplate(w, "allposts.html", map[string]interface{}{
	// 		"CSRFToken": csrfToken,
	// 	})
	// 	return
	// }
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	// if !dep.ValidateCSRFToken(r) {
	// 	dep.ClientError(w, http.StatusForbidden)
	// 	return
	// }

	println("hello")

	if r.Method == http.MethodPost {
		// log.Println("Method not allowed")

		sessionId := r.Context().Value("session_id")
		sess1,err:=r.Cookie("session_id")
		if err !=nil {
			log.Println("error biggy", err)
			return
		}
        if sess1.Value !=sessionId{
			log.Println("sess1.Value",sess1.Value,sessionId)
			log.Println("sessioId",sessionId)
			return
		}
	
		if err := r.ParseForm(); err != nil {
			log.Println("Error while parsing the form", err)
		}

		postContent := r.FormValue("post_content")
		fmt.Println("here1", postContent)
		postId := uuid.New().String()
		categories := r.Form["category[]"]
		title := r.FormValue("post_title")
		userId := r.Context().Value("user_uuid").(string)
		fmt.Println("here2", userId, categories)

		post := models.Post{
			PostId:      postId,
			UserId:      userId,
			Category:    categories,
			Title:       title,
			PostContent: postContent,
		}

		dep.Forum.CreatePost(&post)

		http.Redirect(w, r, "/allposts", http.StatusSeeOther)
	}
}

func (dep *Dependencies) AllPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	PostsTemplate, err := template.ParseFiles("./ui/html/allposts.html")
	if err != nil {
		http.Error(w, "NOT FOUND\nError parsing post templates", http.StatusNotFound)
		return
	}
	posts, err := dep.Forum.AllPosts()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get all posts", http.StatusInternalServerError)
		return
	}
	PostsTemplate.ExecuteTemplate(w, "allposts.html", &posts)
}
