package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)


func AddPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	userID := r.FormValue("user_id")
	category := r.FormValue("category")

	if title == "" || content == "" || userID == "" || category == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()

	postID, err := models.AddPost(id, userID, content)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add post: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"post_id": postID, "message": "Post added successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetAllPosts()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve posts: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
