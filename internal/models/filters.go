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

func (f *ForumModel) Filters(categoryName string) ([]FilteredPosts, error) {
	query := `
		SELECT posts.id, posts.title, posts.content, posts.created_at
		FROM posts
		JOIN post_categories ON posts.post_id = post_categories.post_id
		WHERE category_id = ?;
		`

	rows, err := f.DB.Query(query, categoryName)
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
		//DEBUG
		fmt.Println("This are the filtered posts\n", post)
		filteredPosts = append(filteredPosts, post)
	}

	return filteredPosts, nil
}
