package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

const (
	MaxFileSize = 20 * 1024 * 1024 // 20MB to allow for some buffer
	ChunkSize   = 4096             // Read/write in 4KB chunks
)

// /home/clomollo/forum/ui/html/posts.html
func (dep *Dependencies) PostHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	dep.ErrorLog.Println("Error: Method not allowed:")
	// 	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	fmt.Println(">>> Method: ", r.Method)

	postTempl, err := models.InitTemplates() // initialize the post template
	if err != nil {
		dep.ErrorLog.Println("Error initializing post template:\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sessionId := r.Context().Value("session_id")
	sess1, err := r.Cookie("session_id")
	if err != nil {
		log.Println("error biggy")
		http.Error(w, "User not logged in: ", http.StatusUnauthorized)
		return
	}
	if sess1.Value != sessionId {
		log.Println("sess1.Value", sess1.Value, sessionId)
		log.Println("sessioId", sessionId)
		http.Error(w, "User not logged in: ", http.StatusUnauthorized)
		return
	}

	// Extract form data
	err = r.ParseForm()
	if err != nil {
		dep.ErrorLog.Println("Error parsing post form:\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	postContent := r.FormValue("post_content")
	postId := uuid.New().String()
	categories := r.Form["category"]
	fmt.Println("categories: ", categories)
	title := r.FormValue("post_title")
	userId := r.Context().Value("user_uuid").(string)

	post := models.Post{
		PostId:      postId,
		UserId:      userId,
		Category:    categories,
		Title:       title,
		PostContent: postContent,
	}
	fmt.Println("Categories1", post.Category)

	if err := dep.Forum.CreatePost(&post); err != nil {
		log.Println("Error while quering post db: ", err)
	} else {

		postTempl.ExecuteTemplate(w, "posts.html", nil)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Post created successfully"))
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

func (dep *Dependencies) PostsByFilters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
		return
	}

	FilteredTemplate, err := template.ParseFiles("./ui/html/allposts.html")
	if err != nil {
		http.Error(w, "Failed to parse file", http.StatusInternalServerError)
		return
	}

	categories := r.Form["category"]
	filteredPosts, err := dep.Forum.FilterCategories(categories)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get all posts", http.StatusInternalServerError)
		return
	}
	FilteredTemplate.ExecuteTemplate(w, "allposts.html", &filteredPosts)
}
