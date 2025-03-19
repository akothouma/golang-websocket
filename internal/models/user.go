package models

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	// "learn.zone01kisumu.ke/git/clomollo/forum/internal/database"
)

type User struct {
	ID       int
	UserID   string
	Email    string
	Username string
	Password string
	ProfilePicture string // 🔹 This stores the Base64-encoded image
	// forum    database.ForumModel
}

func (f *ForumModel) CreateUser(userUuid, email, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	querystatement := "INSERT INTO users(user_uuid,email, username, password) VALUES(?,?,?,?)"
	_, err = f.DB.Exec(querystatement, userUuid, email, username, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (f *ForumModel) GetUserByEmail(email string) (*User, error) {
	query := "SELECT id, user_uuid,email, username, password FROM users WHERE email = ?"
	row := f.DB.QueryRow(query, email)
	user := User{}
	err := row.Scan(&user.ID, &user.UserID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (f *ForumModel) GetUserByID(userID int) (*User, error) {
	query := "SELECT id, email, username, password, image_path FROM users WHERE id=?"

	row := f.DB.QueryRow(query, userID)
	user := User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Failed to get user by ID: %v", err)
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}
func (f *ForumModel) GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, email, username, password, profile_picture FROM users WHERE username = ?"
	row := f.DB.QueryRow(query, username)

	user := User{}
	var profilePic []byte

	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &profilePic)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Println(" Failed to get user by username:", err)
		return nil, err
	}

	// Convert image BLOB to Base64
	user.ProfilePicture = MediaToBase64(profilePic)
	return &user, nil
}

