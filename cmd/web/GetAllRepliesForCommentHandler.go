package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)


func GetAllRepliesForCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	commentIDStr := r.URL.Query().Get("comment_id")
	if commentIDStr == "" {
		http.Error(w, "Comment ID is required", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	replies, err := models.GetAllRepliesForComment(commentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve replies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(replies)
}