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

func GetSession(sessionID string)(*Session,error){
	query:=`SELECT id, user_id, expires_at FROM Sessions WHERE id=?`
	row:=DB.QueryRow(query,sessionID)
	var session Session

	err:=row.Scan(&session.ID,&session.UserID,&session.ExpiresAt)
	if err!=nil{
		return nil,fmt.Errorf("failed to get session: %w",err)
	}
	return &session,nil
}

func DeleteSession(sessionID string)error{
	query:="DELETE FROM Sessions WHERE id=?"
	_,err:=DB.Exec(query,sessionID)
	if err!=nil{
		return fmt.Errorf("failed to delete session: %w",err)
	}
	return nil
}

