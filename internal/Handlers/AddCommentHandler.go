package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	sessionId := r.Context().Value("session_id")
	sess1, err := r.Cookie("session_id")
	if err != nil {
		log.Println("error biggy", err)
		return
	}
	if sess1.Value != sessionId {
		log.Println("sess1.Value", sess1.Value, sessionId)
		log.Println("sessionId", sessionId)
		return
	}

	postID := r.FormValue("post_id")
	userID := r.Context().Value("user_uuid").(string)


	content := r.FormValue("content")
	content = strings.TrimSpace(content)
	userID = strings.TrimSpace(userID)

	if userID == "" || content == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	commentID, err := models.AddComment(postID, userID, content)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add comment: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":    "Comment added successfully",
		"comment_id": commentID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
