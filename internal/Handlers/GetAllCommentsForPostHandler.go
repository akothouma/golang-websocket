package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

func GetAllCommentsForPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	comments, err := models.GetAllCommentsForPost(postID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve comments: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
