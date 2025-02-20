package models

import (
	"database/sql"
	"fmt"
)

type Post struct {
	PostId      string   `json:"id"`
	UserId      string   `json:"user_id"`
	Category    []string `json:"category"`
	Title       string   `json:"title"`
	PostContent string   `json:"content_type"`
	Media       string   `json:"videoLink"`
}

func (f *ForumModel) CreatePost(id, title, postContent string, category []string) error {
	query := "INSERT INTO posts(id,category,title,content) VALUES(?,?,?,?)"
	_, err := f.DB.Exec(query, id, title, postContent, category)
	if err != nil {
		return fmt.Errorf("failed to insert a post")
	}
	return nil
}

func (f *ForumModel) FindPostById(id string) (*Post, error) {
	query := "SELECT id,category,title,content FROM posts WHERE id=?"
	row := f.DB.QueryRow(query, id)
	post := Post{}
	err := row.Scan(&post.PostId, &post.Category, &post.Title, &post.PostContent)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &post, nil
}
