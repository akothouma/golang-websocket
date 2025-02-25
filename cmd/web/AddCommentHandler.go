package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)


func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.FormValue("post_id")
	parentCommentIDStr := r.FormValue("parent_comment_id")
	userIDStr := r.FormValue("user_id")
	content := r.FormValue("content")

	if userIDStr == "" || content == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	var postID, parentCommentID *int
	if postIDStr != "" {
		id, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		postID = &id
	}

	if parentCommentIDStr != "" {
		id, err := strconv.Atoi(parentCommentIDStr)
		if err != nil {
			http.Error(w, "Invalid parent comment ID", http.StatusBadRequest)
			return
		}
		parentCommentID = &id
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	commentID, err := models.AddComment(postID, parentCommentID, userID, content)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add comment: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":   "Comment added successfully",
		"comment_id": commentID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}