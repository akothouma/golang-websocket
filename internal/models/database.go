package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type ForumModel struct {
	DB *sql.DB
}

func InitializeDB() (*sql.DB, error) {
	dataBase, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, err
	}

	// defer dataBase.Close()
	
	// create the database queries tabCategory: category,les
	query := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT UNIQUE NOT NULL,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        image_path TEXT
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id INTEGER,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
    );
	
	CREATE TABLE IF NOT EXISTS posts(
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		category TEXT,
		title TEXT NOT NULL,
		content TEXT NOT NULL, 
		user_id INTEGER,
		category TEXT,
		media BLOB, 
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (category) REFERENCES post_categories(id)
	);
	
	
	CREATE TABLE IF NOT EXISTS comments(
		post_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
		
		CREATE TABLE IF NOT EXISTS likes(
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			post_id TEXT,
		type TEXT CHECK(type IN ('like', 'dislike')),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id)
		);
		
		CREATE TABLE IF NOT EXISTS categories( 
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_id TEXT NOT NULL,
		category_value TEXT NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS post_categories(
			post_id TEXT NOT NULL,
			category_id TEXT NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			PRIMARY KEY (post_id, category_id)
		);


	`

	if _, err := dataBase.Exec(query); err != nil {
		dataBase.Close()
		return nil, err
	}

	return dataBase, nil
}
