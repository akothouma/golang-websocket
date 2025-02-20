package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitializeDB() error {
	dataBase, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}

	// defer dataBase.Close()

	// create the database queries tables
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
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

	CREATE TABLE IF NOT EXISTS posts(
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		category TEXT NOT NULL, 
		title TEXT NOT NULL,
		content TEXT NOT NULL, 
		media BLOB, 
		content_type TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (category) REFERENCES post_categories(category)
	);

	CREATE TABLE IF NOT EXISTS post_categories(
		post_id TEXT NOT NULL,
		category_id TEXT NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		PRIMARY KEY (post_id, category_id)
	);

	CREATE TABLE IF NOT EXISTS comments(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id TEXT,
		parent_comment_id INTEGER,
		user_id TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (parent_comment_id) REFERENCES comments(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS likes(
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		post_id TEXT,
		type TEXT CHECK(type IN ('like', 'dislike')),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);

	`

	if _, err := dataBase.Exec(query); err != nil {
		dataBase.Close()
		return err
	}

	DB = dataBase
	return nil
}
