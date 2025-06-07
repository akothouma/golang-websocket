// internal/models/message.go
package models

import (
	// "database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Message   string    `json:"message"`
	ID        uuid.UUID `json:"id"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	IsRead    bool      `json:"isRead"`
	CreatedAt time.Time `json:"timestamp"`
}


func MessageToDatabase(m *Message) error {
    // FIX: Column names now match your schema exactly (messageText, createdAt)
	query := "INSERT INTO messages (messageID, sender, receiver, messageText, isRead, createdAt) VALUES(?,?,?,?,?,?)"
	_, err := DB.Exec(query, m.ID, m.Sender, m.Receiver, m.Message, m.IsRead, m.CreatedAt)
	if err != nil {
		return fmt.Errorf("couldn't add message to database %s", err)
	}
	return nil
}


func GetMessageHistory(user1, user2 string, lastTimestamp time.Time, limit int) ([]Message, error) {
	var messages []Message
	var query string
	var args []interface{}

    // FIX: Using your exact column names: messageID, messageText, createdAt
	query = `SELECT messageID, sender, receiver, messageText, isRead, createdAt FROM messages 
	WHERE (sender=? AND receiver=?) OR (sender=? AND receiver=?)`
	args = append(args, user1, user2, user2, user1)

	if !lastTimestamp.IsZero() {
		query += " AND createdAt < ?" // Use the correct column name here too
		args = append(args, lastTimestamp)
	}

	query += " ORDER BY createdAt DESC LIMIT ?"
	args = append(args, limit)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query issue: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Sender, &msg.Receiver, &msg.Message, &msg.IsRead, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("database scan issue: %w", err)
		}
		messages = append(messages, msg)
	}
	
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
        messages[i], messages[j] = messages[j], messages[i]
    }
	
	return messages, nil
}


func GetAllLastMessages() ([]Message, error) {
    // FIX: Using your exact column names: messageID, messageText, createdAt
    query := `
        SELECT messageID, Sender, Receiver, messageText, IsRead, createdAt FROM (
            SELECT *, ROW_NUMBER() OVER(PARTITION BY
                CASE WHEN Sender < Receiver THEN Sender || ':' || Receiver
                     ELSE Receiver || ':' || Sender END
                ORDER BY createdAt DESC
            ) as rn
            FROM messages
        ) WHERE rn = 1`

    rows, err := DB.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to query last messages: %w", err)
    }
    defer rows.Close()

    var messages []Message
    for rows.Next() {
        var msg Message
        if err := rows.Scan(&msg.ID, &msg.Sender, &msg.Receiver, &msg.Message, &msg.IsRead, &msg.CreatedAt); err != nil {
            return nil, fmt.Errorf("failed to scan last message: %w", err)
        }
        messages = append(messages, msg)
    }
    return messages, nil
}


func GetConversationID(user1, user2 string) string {
	users := []string{user1, user2}
	sort.Strings(users)
	return users[0] + ":" + users[1]
}