// internal/models/message.go

// Package models contains the data structures and database interaction logic for the application.
// This file specifically handles all operations related to chat messages.
package models

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
)

// Message represents the Go struct for a single chat message.
// It maps directly to the `messages` table in the database and is used for JSON serialization.
type Message struct {
	// The content of the message itself. JSON tag matches the `content` field from the client.
	Message string `json:"message"`

	// A unique identifier for the message, generated as a UUID.
	ID uuid.UUID `json:"id"`

	// The user_uuid of the person who sent the message.
	Sender string `json:"sender"`

	// The user_uuid of the intended recipient of the message.
	Receiver string `json:"receiver"`

	// A boolean flag to indicate if the message has been seen by the receiver.
	IsRead bool `json:"isRead"`

	// The exact time the message was created on the server. JSON tag `timestamp` is used for client-side consistency.
	CreatedAt time.Time `json:"timestamp"`
}

// MessageToDatabase inserts a new message record into the `messages` table.
// It explicitly names each column in the INSERT statement for safety and clarity.
func MessageToDatabase(m *Message) error {
	// The query specifies columns to ensure the order of values matches the database schema,
	// preventing errors if the table structure is ever changed.
	query := "INSERT INTO messages (messageID, sender, receiver, messageText, isRead, createdAt) VALUES(?,?,?,?,?,?)"
	_, err := DB.Exec(query, m.ID, m.Sender, m.Receiver, m.Message, m.IsRead, m.CreatedAt)
	if err != nil {
		return fmt.Errorf("couldn't add message to database %s", err)
	}
	return nil
}

// GetMessageHistory fetches a paginated conversation history between two users.
// It retrieves all messages where the two users are either the sender or receiver.
// This function supports infinite scrolling via the `lastTimestamp` and `limit` parameters.
func GetMessageHistory(user1, user2 string, lastTimestamp time.Time, limit int) ([]Message, error) {
	var messages []Message
	var query string
	var args []interface{}

	// The base query selects all relevant columns and messages for a two-way conversation.
	query = `SELECT messageID, sender, receiver, messageText, isRead, createdAt FROM messages 
	WHERE ((sender=? AND receiver=?) OR (sender=? AND receiver=?))`
	args = append(args, user1, user2, user2, user1)

	// If `lastTimestamp` is provided, it's used as a cursor to fetch messages older than it.
	// This is the core of the pagination logic.
	if !lastTimestamp.IsZero() {
		query += " AND createdAt < ?"
		args = append(args, lastTimestamp)
	}

	// Order by `createdAt DESC` to get the most recent messages first, up to the specified `limit`.
	query += " ORDER BY createdAt DESC LIMIT ?"
	args = append(args, limit)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query issue: %w", err)
	}
	defer rows.Close()

	// Scan each row from the database result into a Message struct.
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Sender, &msg.Receiver, &msg.Message, &msg.IsRead, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("database scan issue: %w", err)
		}
		messages = append(messages, msg)
	}

	// The messages were fetched in reverse-chronological order (newest first).
	// We reverse the slice in-place to return them in chronological order (oldest first),
	// which is easier for the frontend to render correctly when prepending history.
	// for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
	// 	messages[i], messages[j] = messages[j], messages[i]
	// }

	return messages, nil
}

// GetAllLastMessages fetches the single most recent message for every conversation in the database.
// This is used to populate the user list in the UI, showing a preview of the latest communication.
func GetAllLastMessages() ([]Message, error) {
	// This advanced SQL query uses a window function (ROW_NUMBER) to achieve this efficiently.
	// 1. `PARTITION BY`: It groups all messages into conversations. A conversation key is created
	//    by concatenating the sender and receiver IDs in alphabetical order (e.g., 'userA:userB').
	//    This ensures that messages from A->B and B->A are treated as the same conversation.
	// 2. `ORDER BY createdAt DESC`: Within each conversation partition, it orders messages from newest to oldest.
	// 3. `ROW_NUMBER() ... as rn`: It assigns a rank number (`rn`) to each message within its conversation group.
	//    The newest message in each conversation will always have `rn = 1`.
	// 4. `WHERE rn = 1`: The outer query then selects only those messages with a rank of 1, effectively
	//    giving us the last message from every unique conversation in a single, efficient query.
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

// GetConversationID is a helper function that creates a standardized, order-independent key
// for a conversation between two users. By sorting the user IDs before joining them,
// it ensures that the conversation between userA and userB has the same key as the one
// between userB and userA (e.g., "userA_id:userB_id"). This is crucial for lookup maps.
func GetConversationID(user1, user2 string) string {
	users := []string{user1, user2}
	sort.Strings(users)
	return users[0] + ":" + users[1]
}