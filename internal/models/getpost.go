package models

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

// var database *models.ForumModel
var DB *sql.DB

// var f *ForumModel

func RenderPostsPage(w http.ResponseWriter, r *http.Request) {
	// if r.Method == http.MethodGet {
	var categories []struct {
		ID   string
		Name string
	}
	categoryRows, err := DB.Query("SELECT id, name FROM categories ORDER BY name")
	if err != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}
	defer categoryRows.Close()

	for categoryRows.Next() {
		var cat struct {
			ID   string
			Name string
		}
		if err := categoryRows.Scan(&cat.ID, &cat.Name); err != nil {
			continue
		}
		categories = append(categories, cat)
	}

	rows, err := DB.Query(`
            SELECT p.post_id, u.username, p.title, p.content, p.media, p.content_type, p.created_at 
            FROM posts p 
            JOIN users u ON p.user_uuid = u.user_uuid
        `)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		// return
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, username, title, content string
		var createdAt time.Time
		var media []byte
		var contentType *string

		if err := rows.Scan(&id, &username, &title, &content, &media, &contentType, &createdAt); err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to parse posts", http.StatusInternalServerError)
			return
		}

		// Convert media to base64 if it exists
		// Convert media to base64 if it exists
		var mediaBase64 string
		if len(media) > 0 {
			mediaBase64 = base64.StdEncoding.EncodeToString(media)
		}

		// Handle contentType
		var contentTypeStr string
		if contentType != nil {
			contentTypeStr = *contentType // Dereference the pointer
		} else {
			contentTypeStr = "" // Or set a default value if needed
		}

		categoryRows, err := DB.Query(`
                SELECT c.id, c.name 
                FROM categories c 
                JOIN post_categories pc ON c.name = pc.category_id 
                WHERE pc.post_id = ?`, id)
		if err != nil {
			http.Error(w, "Failed to fetch post categories", http.StatusInternalServerError)
			return
		}
		defer categoryRows.Close()

		var postCategories []map[string]string
		for categoryRows.Next() {
			var catID, catName string
			if err := categoryRows.Scan(&catID, &catName); err != nil {
				continue
			}
			postCategories = append(postCategories, map[string]string{
				"ID":   catID,
				"Name": catName,
			})
		}

		var likes, dislikes int
		err = DB.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND type = 'like'", id).Scan(&likes)
		if err != nil {
			http.Error(w, "Failed to fetch likes", http.StatusInternalServerError)
			return
		}

		err = DB.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND type = 'dislike'", id).Scan(&dislikes)
		if err != nil {
			http.Error(w, "Failed to fetch dislikes", http.StatusInternalServerError)
			return
		}

		comments, err := GetAllCommentsForPost(id)
		if err != nil {
			http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
			return
		}		
		
		posts = append(posts, map[string]interface{}{
			"ID":             id,
			"Title":          title,
			"Content":        content,
			"Likes":          likes,
			"Dislikes":       dislikes,
			"Comments":       comments,
			"CommentsLenght": len(comments),
			"Username":       username,
			"Initial":        string(username[0]),
			"Categories":     postCategories,
			"Media":          mediaBase64,
			"ContentType":    contentTypeStr,
			"CreatedAt":      createdAt,
		})
	}

	for i, com :=range posts{
		fmt.Println("post",i,":", com["ID"])
	}


	data := map[string]interface{}{
		"Posts":      posts,
		"Categories": categories,
	}

	// Safely get userId from context
	userId, ok := r.Context().Value("user_uuid").(string)
	if !ok {
		log.Println("user_uuid not found in context")
	}

	// Fetch username from DB
	query := `SELECT username FROM users WHERE user_uuid=?`
	var username string
	err = DB.QueryRow(query, userId).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No username found for user_uuid:", userId)
		} else {
			log.Printf("Database error: %v", err)
		}
		username = "" // Ensure username is not nil
	}

	// Only add username if found
	if username != "" {
		data["UserName"] = username
		data["Initial"] = string(username[0])
	}

	RenderTemplates(w, "index.html", data)
}