package models

import (
	"database/sql"
	"fmt"
	"log"
)

type Post struct {
	Id          int
	PostId      string
	UserId      string
	Media       []byte
	Category    []string
	Title       string
	Content string
	ContentType string
	TimeStamp   string
}


func (f *ForumModel) CreatePost(p *Post) error {
	var username string
	err := DB.QueryRow("SELECT username FROM users WHERE user_uuid = ?", p.UserId).Scan(&username)
	if err != nil {
		log.Printf("Failed to fetch username: %v", err)
		return fmt.Errorf("failed to fetch username: %w", err)
	}

	query := "INSERT INTO posts(post_id, user_uuid, username, title, content,  media, content_type) VALUES(?, ?, ?, ?, ?, ?, ?)"
	_, err = f.DB.Exec(query, p.PostId, p.UserId, username, p.Title, p.Content, p.Media, p.ContentType)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to insert a post")
	}

	// Insert categories
	fmt.Println("categories",p.Category)
	
	for _, categoryNames := range p.Category {
		_, err = DB.Exec(`
            INSERT INTO post_categories (post_id, category_id)
            VALUES (?, ?)`,
			p.PostId, categoryNames,
		)
		if err != nil {
			fmt.Println(err)

			return err
		}
	}
	fmt.Println("Your post has been succesfully created")
	return nil
}

func (f *ForumModel) FindPostById(id string) (*Post, error) {
	query := "SELECT id,user_uuid,category,title,content FROM posts WHERE id=?"
	row := f.DB.QueryRow(query, id)
	post := Post{}
	err := row.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content)
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
		err := rows.Scan(&p.Id, &p.PostId, &p.UserId, &p.Title, &p.Content, &p.TimeStamp)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (f *ForumModel) FilterCategories(categories []string) ([]Post, error) {
	posts := []Post{}
	for _, categoryID := range categories {
		query := "SELECT post_id, user_uuid, title, content,  media, content_type FROM posts p JOIN post_categories pc ON p.post_id=pc.post_id AND pc.category_id=?"
		rows, err := f.DB.Query(query, categoryID)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var p Post
			err := rows.Scan(&p.Id, &p.PostId, &p.UserId, &p.Media, &p.Category, &p.Title, &p.Content, &p.ContentType, &p.TimeStamp)
			if err != nil {
				return nil, err
			}
			posts = append(posts, p)
		}
	}
	return posts, nil
}