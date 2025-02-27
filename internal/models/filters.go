package models

import (
	"fmt"
)

type FilteredPosts struct {
	PostID           int    `json:"postid"`
	PostTitle        string `json:"post_title"`
	PostContent      string `json:"post_content"`
	PostCreationDate string `json:"post_creation_date"`
}

func (fp *FilteredPosts) Filters(categoryName string, db *ForumModel) ([]FilteredPosts, error) {
	query := `
		SELECT posts.id, posts.title, posts.content, posts.created_at
		FROM posts
		JOIN post_categories ON posts.id = post_categories.post_id
		JOIN categories ON post_categories.category_id = categories.id
		WHERE categories.name = ?;
	`

	rows, err := db.DB.Query(query, categoryName)
	if err != nil {
		return nil, fmt.Errorf("error executing filters: %w", err)
	}
	defer rows.Close()

	var filteredPosts []FilteredPosts
	for rows.Next() {
		var post FilteredPosts
		err := rows.Scan(
			&post.PostID,
			&post.PostTitle,
			&post.PostContent,
			&post.PostCreationDate,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning the filtration rows: %w", err)
		}
		filteredPosts = append(filteredPosts, post)
	}

	return filteredPosts, nil
}
