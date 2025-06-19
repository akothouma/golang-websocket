package handlers

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

func UploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	// Check request method
	if r.Method != "POST" {
		fmt.Println("Invalid request method:", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form with higher size limit
	err := r.ParseMultipartForm(32 << 20) // 32MB limit
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Retrieve file - ensure the field name matches exactly what's in the form
	file, header, err := r.FormFile("profile_picture")
	if err != nil {
		// More specific error handling
		if strings.Contains(err.Error(), "no such file") {
			http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect without an error
			return
		}
		http.Error(w, "File upload error", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isValidFileType(ext) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Read file into buffer more efficiently
	buffer, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Get logged-in user
	username, err := models.LogedInUser(r)
	if err != nil {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	// Get content type based on file extension
	contentType := getContentType(ext)

	// Store binary data in the database
	// Fixed SQL statement - removed asterisks that appeared to be typos
	_, err = models.DB.Exec("UPDATE users SET profile_picture=?, content_type=? WHERE username=?",
		buffer, contentType, username)
	if err != nil {
		http.Error(w, "Error saving profile picture", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
