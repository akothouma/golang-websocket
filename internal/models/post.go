package models

import (
	"database/sql"
	"fmt"
)

type Post struct {
	PostId      string 
	UserId      string
	Title       string  
	PostContent string  
	Media       string 
}

func (f *ForumModel) CreatePost(p *Post) error {
	fmt.Println("here")
	fmt.Println(p.PostId,p.UserId,p.Title,p.PostContent)
	query := "INSERT INTO posts(post_id,title,content) VALUES(?,?,?)"
	_, err := f.DB.Exec(query, p.PostId,p.UserId, p.Title, p.PostContent)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to insert a post")
	}
	fmt.Println("Your post has ben succesfully created")
	return nil
}

func (f *ForumModel) FindPostById(id string) (*Post, error) {
	query := "SELECT id,category,title,content FROM posts WHERE id=?"
	row := f.DB.QueryRow(query, id)
	post := Post{}
	err := row.Scan(&post.PostId, &post.Title, &post.PostContent)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &post, nil
}

