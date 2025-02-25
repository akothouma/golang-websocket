package main

import (
	"fmt"
	"net/http"
	"time"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

var database models.ForumModel

// Display posts
func RenderPostsPage(w http.ResponseWriter, r *http.Request) {
	// if r.Method == http.MethodPost {
	// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
	// 	return
	// }
	if r.Method == http.MethodGet {
		rows, err := database.DB.Query("SELECT id, user_id, category, title, content, media, content_type, created_at FROM posts ORDER BY created_at DESC ")
		if err != nil {
			http.Error(w, "Failed to load posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var posts []map[string]interface{}
		var id, title, content, user_id, category, contentType string
		var createdAt time.Time
		var media []byte
		for rows.Next() {
			// scan through the database table to get the specified data
			if err := rows.Scan(&id, &user_id, &category, &title, &content, &media, &contentType, &createdAt); err != nil {
				http.Error(w, "Failed to parse posts", http.StatusInternalServerError)
				fmt.Println(err)
				return
			}
			// Fetch total likes, dislikes, and comments for the post
			var likes, dislikes int
			err = database.DB.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND type = 'like'", id).Scan(&likes)
			if err != nil {
				http.Error(w, "Failed to fetch likes", http.StatusInternalServerError)
				return
			}
			err = database.DB.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND type = 'dislike'", id).Scan(&dislikes)
			if err != nil {
				http.Error(w, "Failed to fetch dislikes", http.StatusInternalServerError)
				return
			}
			// Fetch comments for this post
			commentRows, err := database.DB.Query(`
                SELECT c.id, c.content, c.created_at,
                    (SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id AND type = 'like') as likes,
                    (SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id AND type = 'dislike') as dislikes
                FROM comments c
                WHERE c.post_id = ?
                ORDER BY c.created_at DESC`, id)
			if err != nil {
				http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
				return
			}
			defer commentRows.Close()
			// Append post data with likes, dislikes, and comments
			var comments []map[string]interface{}
			for commentRows.Next() {
				var commentID string
				var commentContent string
				var commentCreatedAt time.Time
				var commentLikes, commentDislikes int
				err := commentRows.Scan(&commentID, &commentContent, &commentCreatedAt, &commentLikes, &commentDislikes)
				if err != nil {
					http.Error(w, "Failed to parse comment", http.StatusInternalServerError)
					return
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
				"ID":       id,
				"Title":    title,
				"Content":  content,
				"Likes":    likes,
				"Dislikes": dislikes,
				"Comments": comments,
				"Username": user_id,
				"Category": category,
				"Media":    media,
			})
		}
		RenderTemplates(w, "posts.html", posts)
	}
}
