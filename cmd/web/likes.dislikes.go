package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

// LikeHandler handles likes/dislikes for both posts and comments
func (dep *Dependencies) LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_uuid").(string)

	// Log the content type and body
	log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))

	// Parse form data - handle both form-data and URL-encoded forms
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
			log.Printf("Error parsing multipart form: %v", err)
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
	}
	// Log ALL headers and form data for debugging
	log.Printf("Request headers: %+v", r.Header)
	log.Printf("Form data: %+v", r.Form)
	log.Printf("PostForm data: %+v", r.PostForm)
	// Parse form data
	// if err := r.ParseForm(); err != nil {
	// 	http.Error(w, "Failed to parse form data", http.StatusBadRequest)
	// 	return
	// }
	// Get form values explicitly from PostForm for POST requests
	itemID := r.PostFormValue("id")
	itemType := r.PostFormValue("item_type")
	likeType := r.PostFormValue("type")

	log.Printf("Parsed values - ID: '%s', Item Type: '%s', Like Type: '%s'",
		itemID, itemType, likeType)
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

	fmt.Println("likes handler executed")
	// Process the like/dislike based on item type
	err := dep.Forum.ProcessLike(itemType, itemID, userID, likeType)
	if err != nil {
		log.Printf("ProcessLike error: %v (itemType: %s, itemID: %s, userID: %s, likeType: %s)",
			err, itemType, itemID, userID, likeType)
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
