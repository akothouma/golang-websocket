package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

const (
	MaxFileSize = 20 * 1024 * 1024 // 20MB to allow for some buffer
	ChunkSize   = 4096             // Read/write in 4KB chunks
)

// /home/clomollo/forum/ui/html/posts.html
func (dep *Dependencies) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Fetch categories for the form
		rows, err := DB.Query("SELECT id, name FROM categories ORDER BY name")
		if err != nil {

			http.Error(w, "Failed to load categories", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var categories []struct {
			ID   string
			Name string
		}
		for rows.Next() {
			var cat struct {
				ID   string
				Name string
			}
			if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
				continue
			}
			categories = append(categories, cat)
		}
		RenderTemplates(w, "posts.html", map[string]interface{}{
			"Categories": categories,
		})
		return
	}
	if r.Method == http.MethodPost {
		// log.Println("Method not allowed")

		sessionId := r.Context().Value("session_id")
		sess1, err := r.Cookie("session_id")
		if err != nil {
			log.Println("error biggy", err)
			return
		}
		if sess1.Value != sessionId {
			log.Println("sess1.Value", sess1.Value, sessionId)
			log.Println("sessioId", sessionId)
			return
		}

		// Increase the maximum memory allocated for form parsing
		if err := r.ParseMultipartForm(MaxFileSize); err != nil {
			fmt.Println(err)
			http.Error(w, "File too large", http.StatusBadRequest)
			return
		}

		// Extract form data
		postContent := r.FormValue("post_content")
		postId := uuid.New().String()
		categories := r.Form["category[]"]
		title := r.FormValue("post_title")
		userId := r.Context().Value("user_uuid").(string)

		post := models.Post{
			PostId:      postId,
			UserId:      userId,
			Category:    categories,
			Title:       title,
			PostContent: postContent,
		}

		// Handle file upload
		file, header, err := r.FormFile("media")
		if err == nil {
			defer file.Close()

			// Validate file size
			if header.Size > MaxFileSize {
				http.Error(w, "File size exceeds maximum limit", http.StatusBadRequest)
				return
			}

			// Validate file type
			ext := strings.ToLower(filepath.Ext(header.Filename))
			if !isValidFileType(ext) {
				http.Error(w, "Invalid file type", http.StatusBadRequest)
				return
			}

			// Read file in chunks
			buffer := make([]byte, 0, header.Size)
			tempBuffer := make([]byte, ChunkSize)
			for {
				n, err := file.Read(tempBuffer)
				if err == io.EOF {
					break
				}
				if err != nil {
					http.Error(w, "Error reading file", http.StatusInternalServerError)
					return
				}
				buffer = append(buffer, tempBuffer[:n]...)
			}

			post.Media = buffer
			post.ContentType = getContentType(ext)
		}

		if err := dep.Forum.CreatePost(&post); err != nil {
			log.Println("Error while quering post db: ", err)
		}

		// http.Redirect(w, r, "/allposts", http.StatusSeeOther)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Post created successfully"))
	}
}

func (dep *Dependencies) AllPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	PostsTemplate, err := template.ParseFiles("./ui/html/allposts.html")
	if err != nil {
		http.Error(w, "NOT FOUND\nError parsing post templates", http.StatusNotFound)
		return
	}
	posts, err := dep.Forum.AllPosts()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get all posts", http.StatusInternalServerError)
		return
	}
	PostsTemplate.ExecuteTemplate(w, "allposts.html", &posts)
}

func isValidFileType(ext string) bool {
	validTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".mp4":  true,
		".mov":  true,
		".webm": true,
	}
	return validTypes[ext]
}

func getContentType(ext string) string {
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".mp4":  "video/mp4",
		".mov":  "video/quicktime",
		".webm": "video/webm",
	}
	return contentTypes[ext]
}
