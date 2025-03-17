package main

import (
	"log"
	"net/http"
	"encoding/json"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

// LikeHandler handles likes/dislikes for both posts and comments
func (dep *Dependencies) LikeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Likes handler called with ID: %s, Type: %s, Action: %s", 
        r.FormValue("id"), 
        r.FormValue("item_type"), 
        r.FormValue("type"))
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_uuid").(string)

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	// Get form values
	itemID := r.FormValue("id")          // This could be either post_id or comment_id
	itemType := r.FormValue("item_type") // "post" or "comment"
	likeType := r.FormValue("type")      // "like" or "dislike"
	// Validate inputs
	if itemID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	if itemType != "post" && itemType != "comment" {
		http.Error(w, "Invalid item type", http.StatusBadRequest)
		return
	}
	if likeType != "like" && likeType != "dislike" {
		http.Error(w, "Invalid like type", http.StatusBadRequest)
		return
	}

	// Process the like/dislike based on item type
	err := dep.Forum.ProcessLike(itemType, itemID, userID, likeType)
	if err != nil {
		http.Error(w, "Failed to process like/dislike", http.StatusInternalServerError)
		return
	}

	// Get updated likes and dislikes counts
	likes, dislikes, err := models.PostLikesDislikes(itemID)
	if err != nil {
		http.Error(w, "Failed to get updated counts", http.StatusInternalServerError)
		return
	}

	// Return success response with updated counts
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"success":  true,
		"likes":    likes,
		"dislikes": dislikes,
	}
	
	// Convert the response map to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create response", http.StatusInternalServerError)
		return
	}
	
	w.Write(jsonResponse)
}