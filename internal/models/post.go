package models

import (
	"database/sql"
	"fmt"
)

type Post struct {
	Id          int
	PostId      string
	UserId      string
	Media       []byte
	Category    []string
	Title       string
	PostContent string
	ContentType string
	TimeStamp   string
}

func (f *ForumModel) CreatePost(p *Post) error {
	query := "INSERT INTO posts(post_id, user_uuid, title, content,  media, content_type) VALUES(?, ?, ?, ?, ?, ?)"
	_, err := f.DB.Exec(query, p.PostId, p.UserId, p.Title, p.PostContent, p.Media, p.ContentType)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to insert a post")
	}

	// Insert categories
	for _, categoryID := range p.Category {
		_, err = DB.Exec(`
            INSERT INTO post_categories (post_id, category_id)
            VALUES (?, ?)`,
			p.PostId, categoryID,
		)
		if err != nil {
			return err
		}
	}

	fmt.Println("Your post has ben succesfully created")
	return nil
}

func (f *ForumModel) FindPostById(id string) (*Post, error) {
	query := "SELECT id,user_uuid,category,title,content FROM posts WHERE id=?"
	row := f.DB.QueryRow(query, id)
	post := Post{}
	err := row.Scan(&post.PostId, &post.UserId, &post.Title, &post.PostContent)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &post, nil
}

func (f *ForumModel) AllPosts() ([]Post, error) {
	query := "SELECT * FROM posts"
	rows, err := f.DB.Query(query)
	posts := []Post{}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.Id, &p.PostId, &p.UserId, &p.Title, &p.PostContent, &p.TimeStamp)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
