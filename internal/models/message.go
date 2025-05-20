package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// var DB *sql.DB
type Message struct {
	Message   string    `json:"message"`
	ID        uuid.UUID `json:"id"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	IsRead    bool      `json:"isRead"`
	CreatedAt time.Time `json:"timestamp"`
}

func (m *Message) MessageToDatabase() error {
	query := "INSERT INTO Message VALUES(?,?,?,?,?,?)"
	_, err := DB.Exec(query, m.ID, m.Sender, m.Receiver, m.Message, m.IsRead, m.CreatedAt)
	if err != nil {
		return fmt.Errorf("couldn't add message to database")
	}
	return nil
}

func (m *Message) MessageHistory(user1, user2 uuid.UUID) ([]Message, error) {
	query := `SELECT * FROM Messages 
	WHERE (sender=? AND receiver=?)
	OR(sender=? AND receiver=?)
	ORDER BY timestamp ASC
	LIMIT 10 OFFSET 0`

	messageRows, err := DB.Query(query, user1, user2, user2, user1)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("no messages shared yet")
		}
		return nil, fmt.Errorf("Database Issue")
	}
	var messages []Message
	for messageRows.Next() {
		var msg Message
		err := messageRows.Scan(&msg.ID, &msg.Sender, &msg.Receiver, &msg.Message, &msg.IsRead, msg.CreatedAt.Truncate(time.Hour))
		if err != nil {
			return nil, fmt.Errorf("Database issue", err)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
