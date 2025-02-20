package models

import (
	"fmt"
	"log"
	"time"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/database"
)

// Comment struct
type Comment struct {
	ID              int        `json:"id"`
	PostID          *int       `json:"post_id,omitempty"`
	ParentCommentID *int       `json:"parent_comment_id,omitempty"`
	UserID          int        `json:"user_id"`
	Content         string     `json:"content"`
	CreatedAt       time.Time  `json:"created_at"`
}

// AddComment adds a new comment (post or reply)
func AddComment(postID, parentCommentID *int, userID int, content string) (int64, error) {
	query := `INSERT INTO comments (post_id, parent_comment_id, user_id, content, created_at) 
			  VALUES (?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, postID, parentCommentID, userID, content, time.Now())
	if err != nil {
		log.Printf("Failed to add comment: %v", err)
		return 0, fmt.Errorf("failed to add comment: %w", err)
	}

	commentID, _ := result.LastInsertId()
	return commentID, nil
}

// GetAllCommentsForPost retrieves all top-level comments for a post
func GetAllCommentsForPost(postID int) ([]Comment, error) {
	query := `SELECT id, post_id, parent_comment_id, user_id, content, created_at 
			  FROM comments WHERE post_id = ? AND parent_comment_id IS NULL ORDER BY created_at DESC`

	rows, err := database.DB.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParentCommentID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, c)
	}
	return comments, nil
}