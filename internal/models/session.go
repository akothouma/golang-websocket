package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
	// Forum *database.ForumModel
}

func (f *ForumModel) CreateSession(userID string) (string, error) {
	querry := `DELETE FROM Sessions WHERE user_uuid=?`
	_, err := f.DB.Exec(querry, userID)
	if err != nil {
		return "", fmt.Errorf("failed to delete existing sessions: %w", err)
	}
	SessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	query := "INSERT INTO Sessions(id, user_uuid, expires_at) VALUES(?, ?, ?)"
	_, err = f.DB.Exec(query, SessionID, userID, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to insert session: %w", err)
	}

	return SessionID, nil
}

func (f *ForumModel) GetSession(sessionID string) (*Session, error) {
	query := `SELECT id, user_uuid, expires_at FROM Sessions WHERE id=?`
	row := f.DB.QueryRow(query, sessionID)
	var session Session

	err := row.Scan(&session.ID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &session, nil
}

func (f *ForumModel) DeleteSession(sessionID string) error {
	query := "DELETE FROM Sessions WHERE id=?"
	_, err := f.DB.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
