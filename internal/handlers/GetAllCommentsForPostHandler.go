package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)


func GetAllCommentsForPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
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

