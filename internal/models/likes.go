package models

import (
	"fmt"

	"github.com/google/uuid"
)

// processLike handles the database operations for likes
func (f *ForumModel) ProcessLike(itemType, itemID, userID, likeType string) error {
	var tableName, idColumn string
	// Set table and column names based on item type
	if itemType == "post" {
		tableName = "post_likes"
		idColumn = "post_id"
	} else {
		tableName = "comment_likes"
		idColumn = "comment_id"
	}
	// First verify the item exists
	exists, err := f.checkItemExists(itemType, itemID)
	if err != nil {
		return fmt.Errorf("error checking item existence: %v", err)
	}
	if !exists {
		return fmt.Errorf("item not found")
	}
	// Check for existing like
	var existingType string
	query := fmt.Sprintf("SELECT type FROM %s WHERE user_id = ? AND %s = ?", tableName, idColumn)
	err = f.DB.QueryRow(query, userID, itemID).Scan(&existingType)
	if err == nil {
		// Update existing like
		updateQuery := fmt.Sprintf("UPDATE %s SET type = ? WHERE user_id = ? AND %s = ?", tableName, idColumn)
		_, err = f.DB.Exec(updateQuery, likeType, userID, itemID)
	} else {
		// Insert new like
		insertQuery := fmt.Sprintf("INSERT INTO %s (id, user_id, %s, type) VALUES (?, ?, ?, ?)", tableName, idColumn)
		_, err = f.DB.Exec(insertQuery, uuid.New().String(), userID, itemID, likeType)
	}
	return err
}

// checkItemExists verifies if the post or comment exists
func (f *ForumModel) checkItemExists(itemType, itemID string) (bool, error) {
	var exists bool
	var query string
	if itemType == "post" {
		query = "SELECT EXISTS(SELECT 1 FROM posts WHERE post_id = ?)"
	} else {
		query = "SELECT EXISTS(SELECT 1 FROM comments WHERE id = ?)"
	}
	err := f.DB.QueryRow(query, itemID).Scan(&exists)
	return exists, err
}

func PostCommentLikesDislikes(itemType, id string) (int, int, error) {
	var tableName, idColumn string

	// Validate and set table name
	switch itemType {
	case "post":
		tableName = "post_likes"
		idColumn = "post_id"
	case "comment":
		tableName = "comment_likes"
		idColumn = "comment_id"
	default:
		return 0, 0, fmt.Errorf("invalid item type: %s", itemType)
	}

	var likes, dislikes int

	// Fetch likes count
	err := DB.QueryRow("SELECT COUNT(*) FROM "+tableName+" WHERE "+idColumn+" = ? AND type = 'like'", id).Scan(&likes)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch likes: %v", err)
	}

	// Fetch dislikes count
	err = DB.QueryRow("SELECT COUNT(*) FROM "+tableName+" WHERE "+idColumn+" = ? AND type = 'dislike'", id).Scan(&dislikes)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch dislikes: %v", err)
	}

	return likes, dislikes, nil
}
