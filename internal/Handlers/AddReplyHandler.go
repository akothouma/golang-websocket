package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

func AddReplyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Error biggy")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	sessionId := r.Context().Value("session_id")
	sess1, err := r.Cookie("session_id")
	if err != nil {
		return
	}
	if sess1.Value != sessionId {
		log.Println("sess1.Value", sess1.Value, sessionId)
		log.Println("sessioId", sessionId)
		return
	}

	parent_comment_id := r.FormValue("parent_comment_id")
	userID := r.Context().Value("user_uuid").(string)


	content := r.FormValue("content")


	if userID == "" || content == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	commentID, err := models.AddReply(parent_comment_id, userID, content)
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
