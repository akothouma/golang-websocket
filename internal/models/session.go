package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
}


func CreateSession(userID int) (string, error) {
	querry := `DELETE FROM Sessions WHERE user_id=?`
	_, err := DB.Exec(querry, userID)
	if err != nil {
		return "", fmt.Errorf("failed to delete existing sessions: %w", err)
	}
	SessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	query := "INSERT INTO Sessions(id, user_id, expires_at) VALUES(?, ?, ?)"
	_, err = DB.Exec(query, SessionID, userID, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to insert session: %w", err)
	}

	return SessionID, nil
}
