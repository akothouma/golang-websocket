package main

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)
func UploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("UploadProfilePictureHandler executed!")
    
    // Check request method
    if r.Method != "POST" {
        fmt.Println("Invalid request method:", r.Method)
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }
    
    // Parse form with higher size limit
    err := r.ParseMultipartForm(32 << 20) // 32MB limit
    if err != nil {
        fmt.Println("Error parsing form:", err)
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }
    
    // Print form field names for debugging
    fmt.Println("Form fields:", r.MultipartForm.Value)
    fmt.Println("File fields:", r.MultipartForm.File)
    
    // Retrieve file - ensure the field name matches exactly what's in the form
    file, header, err := r.FormFile("profile_picture")
    if err != nil {
        fmt.Println("Error retrieving file:", err)
        
        // More specific error handling
        if strings.Contains(err.Error(), "no such file") {
            fmt.Println("No file was uploaded or field name mismatch")
            http.Error(w, "Please select a file to upload", http.StatusBadRequest)
        } else {
            http.Error(w, "File upload error", http.StatusBadRequest)
        }
        return
    }
    defer file.Close()
    
    fmt.Println("File received:", header.Filename, "Size:", header.Size)
    
    // Validate file type
    ext := strings.ToLower(filepath.Ext(header.Filename))
    if !isValidFileType(ext) {
        fmt.Println("Invalid file type:", ext)
        http.Error(w, "Invalid file type", http.StatusBadRequest)
        return
    }
    
    // Read file into buffer more efficiently
    buffer, err := io.ReadAll(file)
    if err != nil {
        fmt.Println("Error reading file:", err)
        http.Error(w, "Error reading file", http.StatusInternalServerError)
        return
    }
    
    // Get logged-in user
    username, err := models.LogedInUser(r)
    if err != nil {
        fmt.Println("User not logged in")
        http.Error(w, "User not logged in", http.StatusUnauthorized)
        return
    }

	 // Get content type based on file extension
	 contentType := getContentType(ext)
    
    fmt.Println("User identified:", username)
    
    // Store binary data in the database
    // Fixed SQL statement - removed asterisks that appeared to be typos
    _, err = models.DB.Exec("UPDATE users SET profile_picture=?, content_type=? WHERE username=?", 
                           buffer, contentType, username)
    if err != nil {
        fmt.Println("Error saving profile picture in DB:", err)
        http.Error(w, "Error saving profile picture", http.StatusInternalServerError)
        return
    }
    
    fmt.Println("Profile picture uploaded successfully for user:", username)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
