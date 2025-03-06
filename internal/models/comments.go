package models

import (
	"fmt"
	"log"
	"time"
)

// Comment struct
type Comment struct {
	ID              int       `json:"id"`
	PostID          string    `json:"post_id,omitempty"`
	ParentCommentID int       `json:"parent_comment_id,omitempty"`
	UserUuiD          string    `json:"user_id"`
	UserName        string    `json:"username"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
}

// AddComment adds a new comment (post or reply)
func AddComment(postID, UserUuiD string, content string) (int64, error) {
	//query to get the user name
	var username string
	err := DB.QueryRow("SELECT username FROM users WHERE user_uuid = ?", UserUuiD).Scan(&username)
	if err != nil {
		log.Printf("Failed to fetch username: %v", err)
		return 0, fmt.Errorf("failed to fetch username: %w", err)
	}

	query := `INSERT INTO comments (post_id, user_uuid, username, content, created_at) 
			  VALUES (?, ?, ?, ?, ?)`
	result, err := DB.Exec(query, postID, UserUuiD, username, content, time.Now())
	if err != nil {
		log.Printf("Failed to add comment: %v", err)
		return 0, fmt.Errorf("failed to add comment: %w", err)
	}

	commentID, _ := result.LastInsertId()
	return commentID, nil
}

// AddReply adds a new comment (post or reply)
func AddReply(parentCommentID, UserUuiD string, content string) (int64, error) {
	//query to get the user name
	var username string
	err := DB.QueryRow("SELECT username FROM users WHERE user_uuid = ?", UserUuiD).Scan(&username)
	if err != nil {
		log.Printf("Failed to fetch username: %v", err)
		return 0, fmt.Errorf("failed to fetch username: %w", err)
	}

	query := `INSERT INTO comments (parent_comment_id, username, user_uuid, content, created_at) 
			  VALUES (?, ?, ?, ?, ?)`
	result, err := DB.Exec(query, parentCommentID, username, UserUuiD, content, time.Now())
	if err != nil {
		log.Printf("Failed to add Reply: %v", err)
		return 0, fmt.Errorf("failed to add Reply: %w", err)
	}

	commentID, _ := result.LastInsertId()
	return commentID, nil
}

// GetAllCommentsForPost retrieves all top-level comments for a post
func GetAllCommentsForPost(postID string) ([]Comment, error) {
	query := `SELECT id, post_id, username, user_uuid, content, created_at 
			  FROM comments WHERE post_id = ? AND parent_comment_id IS NULL ORDER BY created_at DESC`

	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserName, &c.UserUuiD, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, c)
	}
	return comments, nil
}

// GetRepliesForComment retrieves all replies to a specific comment
func GetAllRepliesForComment(commentID int) ([]Comment, error) {
	query := `SELECT id, parent_comment_id, username, user_uuid, content, created_at 
			  FROM comments WHERE parent_comment_id = ? ORDER BY created_at ASC`

	rows, err := DB.Query(query, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get replies: %w", err)
	}
	defer rows.Close()

	var replies []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.ParentCommentID, &c.UserName, &c.UserUuiD, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan reply: %w", err)
		}
		replies = append(replies, c)
	}
	return replies, nil
}
