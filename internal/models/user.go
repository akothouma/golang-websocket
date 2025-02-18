package models

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/database"
)

type User struct {
	ID       int
	Email    string
	Username string
	Password string
}

func CreateUser(email, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	querystatement := "INSERT INTO users(email, username, password) VALUES(?,?,?)"
	_, err = database.DB.Exec(querystatement, email, username, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	query := "SELECT id, email, username, password FROM users WHERE email = ?"
	row := database.DB.QueryRow(query, email)
	user := User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
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

func GetUserByID(userID int)(*User, error){
	query:="SELECT id, email, username, password, image_path FROM users WHERE id=?"

	row:=database.DB.QueryRow(query,userID)
	user:=User{}
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

func GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, email, username, password FROM users WHERE username = ?"
	row := database.DB.QueryRow(query, username)
	user := User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err==sql.ErrNoRows{
			return nil,nil
		}
		log.Printf("Failed to get user by username: %v", err)
        return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}