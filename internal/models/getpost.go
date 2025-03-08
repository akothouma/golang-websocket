package models

import (
	"database/sql"
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
		mediaBase64 := MediaToBase64(media)
	

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

	data := make(map[string]interface{})

	username := ""

	if username, err = LogedInUser(r); err != nil{
		fmt.Println(err)
	}else{
		data["UserName"] = username
		data["Initial"] = string(username[0])
	}

	data["Posts"] = posts
	data["Categories"] = categories


	// fmt.Println("categories:", categories)
	


	RenderTemplates(w, "index.html", data)
}


func LogedInUser(r *http.Request)(string, error){

	session, err := r.Cookie("session_id")

	if err != nil {
		return "", fmt.Errorf("Session cookie not found")
	}

	sessionID := session.Value

	// Query to check if session is valid and fetch the username
	query := `
		SELECT u.username 
		FROM users u
		JOIN sessions s ON u.user_uuid = s.user_uuid
		WHERE s.id = ? AND s.expires_at > CURRENT_TIMESTAMP`
	
	var username string

	err = DB.QueryRow(query, sessionID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No valid session found for session ID: %s, err: %v", sessionID, err)
			return "", fmt.Errorf("no valid session found for session ID: %s", sessionID)
		} else if err != nil {
			log.Printf("Database error: %v", err)
			return "", fmt.Errorf("database error: %w", err)
		}
		
	}

	return username, nil
}