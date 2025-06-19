package models

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	// "learn.zone01kisumu.ke/git/clomollo/forum/internal/database"
)

type User struct {
	ID             int
	UserID         string
	FirstName      string `json:"firstName"` // Add if you want to receive and store
	LastName       string `json:"lastName"`  // Add if you want to receive and store
	Email          string `json:"email"`
	Username       string `json:"username"`
	Password       string `json:"password"`         // This will store the HASH
	Age            int    `json:"age,omitempty"`    // Add if you want to receive and store
	Gender         string `json:"gender,omitempty"` // Add if you want to receive and store
	ProfilePicture string `json:"profilePicture,omitempty"`
	ContentType    string `json:"contentType,omitempty"`
}

func (f *ForumModel) CreateUser(
	userUuid string,
	firstName string, // New parameter
	lastName string, // New parameter
	email string,
	username string,
	plainPassword string, // Renamed for clarity
	age int, // New parameter (pass 0 or a specific value if not provided/optional)
	gender string, // New parameter (pass empty string if not provided/optional)
) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Make sure the query statement matches your updated schema column names
	// and the order of parameters.
	queryStatement := `
        INSERT INTO users (
            user_uuid, 
            first_name, 
            last_name, 
            email, 
            username, 
            password, 
            age, 
            gender
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	// log.Printf("Executing query: %s with params: %s, %s, %s, %s, %s, %s, %d, %s",
	// 	queryStatement, userUuid, firstName, lastName, email, username, " HASHED_PASSWORD ", age, gender)

	_, err = f.DB.Exec(queryStatement,
		userUuid,
		firstName,
		lastName,
		email,
		username,
		string(hashedPassword),
		age,    // If age can be optional and your schema allows NULL, you might need sql.NullInt32
		gender, // If gender can be optional and your schema allows NULL, you might need sql.NullString
	)

	if err != nil {
		// You might want more specific error handling here to check for UNIQUE constraint violations
		// (e.g., for email or username) to return more user-friendly messages.
		return fmt.Errorf("failed to insert user into database: %w", err)
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

func (f *ForumModel) GetAllConnectedUsers(usersID []string) ([]User, error) {
	var users []User
	query := "SELECT username,profile_picture FROM users WHERE user_uuid=?"
	for _, userID := range usersID {
		row, err := f.DB.Query(query, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("This user doesn't exist")
			}
			return nil, fmt.Errorf(err.Error())
		}

		var u User
		for row.Next() {
			var profilepic []byte

			err = row.Scan(&u.Username, &profilepic)
			if err != nil {
				return nil, fmt.Errorf(err.Error())
			}
			u.ProfilePicture = MediaToBase64(profilepic)
			u.UserID = userID
		}
		users = append(users, u)
	}
	return users, nil
}

func (f *ForumModel) GetAllUsers() ([]User, error) {
	var users []User

	query := "SELECT user_uuid, username, profile_picture FROM users"

	rows, err := f.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		var profilepic []byte

		err := rows.Scan(&u.UserID, &u.Username, &profilepic)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %v", err)
		}

		u.ProfilePicture = MediaToBase64(profilepic)
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return users, nil
}
