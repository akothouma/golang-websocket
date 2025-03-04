package models

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"time"
	// "learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

// var database *models.ForumModel
var DB *sql.DB

// var f *ForumModel

func RenderPostsPage() (map[string]interface{}, error) {
	// if r.Method == http.MethodGet {
	var categories []struct {
		ID   string
		Name string
	}
	categoryRows, err := DB.Query("SELECT id, name FROM categories ORDER BY name")
	if err != nil {
		// Error("Failed to load categories", http.StatusInternalServerError)
		return nil, errors.New("Failed to load categories")
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
		// http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return nil, errors.New("Failed to load posts")
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, username, title, content string
		var createdAt time.Time
		var media []byte
		var contentType *string

		if err := rows.Scan(&id, &username, &title, &content, &media, &contentType, &createdAt); err != nil {
			// fmt.Println(err)
			// http.Error(w, "Failed to parse posts", http.StatusInternalServerError)
			return nil, errors.New("Failed to parse posts")
		}

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
			// http.Error(w, "Failed to fetch post categories", http.StatusInternalServerError)
			return nil, errors.New("Failed to fetch post categories")
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
			// http.Error(w, "Failed to fetch likes", http.StatusInternalServerError)
			// return
			return nil, errors.New("Failed to fetch likes")
		}

		err = DB.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND type = 'dislike'", id).Scan(&dislikes)
		if err != nil {
			// http.Error(w, "Failed to fetch dislikes", http.StatusInternalServerError)
			// return
			return nil, errors.New("Failed to fetch dislikes")
		}

		commentRows, err := DB.Query(`
                SELECT c.id, c.content, c.created_at,
                (SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id AND type = 'like') as likes,
                (SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id AND type = 'dislike') as dislikes
                FROM comments c
                WHERE c.post_id = ?
                ORDER BY c.created_at DESC`, id)
		if err != nil {
			// http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
			// return
			return nil, errors.New("Failed to fetch comments")
		}
		defer commentRows.Close()

		var comments []map[string]interface{}
		for commentRows.Next() {
			var commentID, commentContent string
			var commentCreatedAt time.Time
			var commentLikes, commentDislikes int

			if err := commentRows.Scan(&commentID, &commentContent, &commentCreatedAt, &commentLikes, &commentDislikes); err != nil {
				// http.Error(w, "Failed to parse comment", http.StatusInternalServerError)
				return nil, errors.New("Failed to parse comments")
			}

			comments = append(comments, map[string]interface{}{
				"ID":        commentID,
				"Content":   commentContent,
				"CreatedAt": commentCreatedAt,
				"Likes":     commentLikes,
				"Dislikes":  commentDislikes,
			})
		}

		posts = append(posts, map[string]interface{}{
			"ID":          id,
			"Title":       title,
			"Content":     content,
			"Likes":       likes,
			"Dislikes":    dislikes,
			"Comments":    comments,
			"Username":    username,
			"Categories":  postCategories,
			"Media":       mediaBase64,
			"ContentType": contentTypeStr,
			"CreatedAt":   createdAt,
		})
	}
	var data map[string]interface{}
	data = map[string]interface{}{
		"Posts":      posts,
		"Categories": categories,
	}

	// // RenderTemplates(w, "posts.html", data)
	// postHtml, err := template.ParseFiles("./ui/html/posts.html")
	// if err != nil {
	// 	// fmt.Println("Error loading posts html\n",err )
	// 	fmt.Println("")
	// }
	// postHtml.Execute(w, data)
	// }
	return data, nil
}
