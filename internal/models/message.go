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

func MessageToDatabase(m *Message) error {
	query := "INSERT INTO messages VALUES(?,?,?,?,?,?)"
	_, err := DB.Exec(query, m.ID, m.Sender, m.Receiver, m.Message, m.IsRead, m.CreatedAt)
	if err != nil {
		return fmt.Errorf("couldn't add message to database %s",err)
	}
	return nil
}

func MessageHistory(user1, user2 string) ([]Message, error) {
	fmt.Println("The users",user1,user2);
	query := `SELECT * FROM messages 
	WHERE (sender=? AND receiver=?)
	OR(sender=? AND receiver=?)
	ORDER BY CreatedAt ASC`
	
	messageRows, err := DB.Query(query, user1, user2, user2, user1)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no messages shared yet")
		}
		return nil, fmt.Errorf("Database Issue %s",err)
	}
	var messages []Message
	for messageRows.Next() {
		var msg Message
		err := messageRows.Scan(&msg.ID, &msg.Sender, &msg.Receiver, &msg.Message, &msg.IsRead, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Database issue", err)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
