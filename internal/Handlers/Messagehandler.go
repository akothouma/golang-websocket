// Package handlers contains the HTTP handlers for the application, including the WebSocket logic.
package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

// UserChatInfo represents the data structure for a single user in the chat's sidebar list.
// It contains all the necessary information for the frontend to render the user's status and last message.
type UserChatInfo struct {
	UserID             string    `json:"userID"`             // The unique identifier for the user.
	Username           string    `json:"username"`           // The display name for the user.
	IsOnline           bool      `json:"isOnline"`           // A boolean flag indicating if the user is currently connected via WebSocket.
	LastMessageContent string    `json:"lastMessageContent"` // The content of the most recent message for sorting and display.
	LastMessageTime    time.Time `json:"lastMessageTime"`    // The timestamp of the most recent message, used for sorting.
	UnreadCount        int       `json:"unreadCount"`        // The count of unread messages from this user.
	IsMe               bool      `json:"isMe,omitempty"`     // A flag set to true only for the user receiving this payload, used to identify self.
}

// WebSocketMessage is the universal struct for all communication over the WebSocket.
// The `Type` field acts as a router, telling the recipient how to interpret the payload.
type WebSocketMessage struct {
	Type            string           `json:"type"`                      // e.g., "get_user_list", "private_message", "get_message_history".
	Content         string           `json:"content,omitempty"`         // The text content of a message being sent.
	Target          string           `json:"target,omitempty"`          // The userID of the intended recipient.
	Sender          string           `json:"sender,omitempty"`          // The userID of the sender, used for events like typing.
	Messages        []models.Message `json:"messages,omitempty"`        // A slice of messages, used for sending chat history.
	UserList        []UserChatInfo   `json:"userList,omitempty"`        // A slice of user info, used for updating the online users list.
	LastMessageTime time.Time        `json:"lastMessageTime,omitempty"` // The timestamp of the oldest message received, used for paginating history (infinite scroll).
}

// ErrorObject is a simple struct for sending error messages back to the client.
type ErrorObject struct {
	Error string `json:"error"`
}

// Global variables for the chat system.
var conn *websocket.Conn

var (
	// Clients is a map that stores active WebSocket connections, with the user's UUID as the key.
	// It acts as the central registry of all currently online users.
	Clients = make(map[string]*websocket.Conn)

	// broadcast is a buffered channel that acts as a message queue for all incoming messages.
	// It decouples message reception from message processing and broadcasting.
	broadcast = make(chan models.Message)

	// upgrader is a Gorilla WebSocket helper that upgrades an HTTP connection to a WebSocket connection.
	upgrader = websocket.Upgrader{}

	// ClientsMux is a mutex to prevent race conditions when multiple goroutines access the Clients map.
	// Any read or write to the Clients map must be protected by this lock.
	ClientsMux sync.Mutex

	// broadcastOnce ensures that the goroutine listening on the broadcast channel is started only once.
	broadcastOnce sync.Once
)

// ChatHandler is the HTTP handler for the /ws endpoint. It upgrades the connection and starts the client handler.
func (dep *Dependencies) ChatHandler(w http.ResponseWriter, r *http.Request) {
	// A user must be logged in (have a valid session) to connect to the chat.
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Allow all origins for WebSocket connections during development.
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to establish connection", err)
		return
	}

	// Retrieve the user's ID from the context, set by the authentication middleware.
	userID := dep.UserIdFromContext(r)

	// Start a new goroutine to handle this specific client's connection.
	go dep.handleClientConnections(userID, conn)
	// Ensure the global message broadcaster is started, but only once.
	broadcastOnce.Do(func() {
		dep.StartChatBroadcastHandler()
	})
}

// handleClientConnections runs in its own goroutine for each connected user.
// It listens for incoming messages from the client and handles cleanup on disconnect.
func (dep *Dependencies) handleClientConnections(userID string, conn *websocket.Conn) {
	// 1. Register the new client.
	ClientsMux.Lock()
	Clients[userID] = conn
	ClientsMux.Unlock()
	dep.broadcastUserListUpdate()

	// 2. Announce the new user's status to all other clients.

	// 3. Defer the cleanup logic to run when the function exits (i.e., when the client disconnects).
	defer func() {
		ClientsMux.Lock()
		delete(Clients, userID)
		ClientsMux.Unlock()
		// Announce the user's offline status to all remaining clients.
		dep.broadcastUserListUpdate()
		conn.Close()
	}()

	// 4. Start an infinite loop to read messages from the client.
	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			// Handle WebSocket close errors gracefully.
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			}
			break // Exit loop on error to trigger the deferred cleanup.
		}

		// Process the received message based on its type.
		switch msg.Type {
		case "get_user_list":
			// Client is explicitly requesting the user list (e.g., on reconnect).
			dep.broadcastUserListUpdate()

		case "private_message":
			// Client is sending a private message to another user.
			if msg.Target == "" || msg.Content == "" {
				continue
			}
			messageModel := models.Message{
				ID:        uuid.New(),
				Sender:    userID,
				Receiver:  msg.Target,
				Message:   msg.Content,
				IsRead:    false, // Default to unread.
				CreatedAt: time.Now(),
			}
			// Send the structured message to the central broadcast channel for processing.
			broadcast <- messageModel

		case "get_message_history":
			// Client is requesting the chat history with another user.
			if msg.Target == "" {
				continue
			}
			dep.sendMessageHistory(conn, userID, msg.Target, msg.LastMessageTime)
		case "mark_messages_as_read":
			// Client has opened a chat and is marking messages as read.
			if msg.Target == "" {
				continue
			}

			// Mark messages from the target (sender) to the current user (receiver) as read.
			err := models.MarkMessagesAsRead(msg.Target, userID)
			if err != nil {
				log.Printf("Error marking messages as read for user %s from %s: %v", userID, msg.Target, err)
			}
			// After marking as read, broadcast an updated user list so the unread count badge disappears.
			dep.broadcastUserListUpdate()

			// ---- START: NEW TYPING CASES ----
		case "start_typing":
			dep.relayTypingEvent("typing_started", userID, msg.Target)

		case "stop_typing":
			dep.relayTypingEvent("typing_stopped", userID, msg.Target)

		default:
			log.Printf("Unknown message type '%s' from user %s", msg.Type, userID)
		}
	}
}

// relayTypingEvent forwards a typing notification from a sender to a target user.
// It does not persist anything and acts as a simple, real-time relay.
func (dep *Dependencies) relayTypingEvent(eventType, senderID, targetID string) {
	if targetID == "" {
		return
	}

	ClientsMux.Lock()
	defer ClientsMux.Unlock()

	// Prepare the event payload to be sent to the recipient.
	response := WebSocketMessage{
		Type:   eventType, // "typing_started" or "typing_stopped"
		Sender: senderID,  // The ID of the user who is typing
	}

	// Find the recipient's connection and send them the event.
	if targetConn, ok := Clients[targetID]; ok {
		if err := targetConn.WriteJSON(response); err != nil {
			log.Printf("Error relaying typing event to user %s: %v", targetID, err)
		}
	}
}

// sendMessageHistory fetches a paginated chunk of messages and sends it back to the requesting client.
// It supports infinite scrolling by using the timestamp of the oldest message as a cursor.
func (dep *Dependencies) sendMessageHistory(conn *websocket.Conn, senderID, targetID string, lastTimestamp time.Time) {
	messages, err := models.GetMessageHistory(senderID, targetID, lastTimestamp, 10)
	if err != nil {
		log.Printf("Error getting message history for %s<->%s: %v", senderID, targetID, err)
		return
	}

	// Construct and send the response.
	response := WebSocketMessage{
		Type:     "history_response",
		Target:   targetID, // Let the frontend know which chat this history belongs to.
		Messages: messages,
	}
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending message history to %s: %v", senderID, err)
	}
}

// broadcastUserListUpdate compiles a fresh list of all users and broadcasts it to all connected clients.
// Each client receives a tailored list where the `IsMe` flag is set correctly.
// This is the single source of truth for the user list UI on the frontend.
func (dep *Dependencies) broadcastUserListUpdate() {
	// 1. Fetch all users and all last messages from the database.
	allDBUsers, err := dep.Forum.GetAllUsers()
	if err != nil {
		// log.Printf("broadcastUserListUpdate: Error getting all users: %v", err)
		return
	}
	lastMessages, err := models.GetAllLastMessages()
	if err != nil {
		// log.Printf("broadcastUserListUpdate: Error getting last messages: %v", err)
		return
	}

	// 2. Create a map for efficient lookup of the last message in any conversation.
	lastMessageMap := make(map[string]models.Message)
	for _, msg := range lastMessages {
		key := models.GetConversationID(msg.Sender, msg.Receiver)
		lastMessageMap[key] = msg
	}

	ClientsMux.Lock()
	defer ClientsMux.Unlock()

	// 3. Iterate through each connected client to send them a personalized user list.
	for userID, clientConn := range Clients {
		// Fetch the unread message counts specifically for the user receiving this update
		unreadCounts, err := models.GetUnreadMessageCounts(userID)
		if err != nil {
			log.Printf("Error getting unread counts for user %s: %v", userID, err)
			// Continue with an empty map if there's an error
			unreadCounts = make(map[string]int)
		}
		var usersWithStatus []UserChatInfo // Create a new list for each recipient.

		// 4. Build the list of users with their online status and "IsMe" flag.
		for _, dbUser := range allDBUsers {
			_, isOnline := Clients[dbUser.UserID]
			userInfo := UserChatInfo{
				UserID:      dbUser.UserID,
				Username:    dbUser.Username,
				IsOnline:    isOnline,
				IsMe:        (dbUser.UserID == userID), // Tailor this flag for the recipient.
				UnreadCount: unreadCounts[dbUser.UserID],
			}

			if !(userInfo.IsMe) {
				key := models.GetConversationID(userID, dbUser.UserID)
				if msg, ok := lastMessageMap[key]; ok {
					// if msg.CreatedAt.After(userInfo.LastMessageTime) {
					userInfo.LastMessageContent = msg.Message
					userInfo.LastMessageTime = msg.CreatedAt
					//}
				}
			}
			usersWithStatus = append(usersWithStatus, userInfo)
		}
		// 6. Send the final, enriched, and tailored list to the client.
		response := WebSocketMessage{
			Type:     "user_list_update",
			UserList: usersWithStatus,
		}
		if err := clientConn.WriteJSON(response); err != nil {
			log.Printf("Error broadcasting tailored user list to user %s: %v", userID, err)
		}
	}
}

// StartChatBroadcastHandler runs as a single, long-lived goroutine.
// It listens on the `broadcast` channel and processes messages sequentially.
// This prevents race conditions and ensures messages are handled in order.
func (dep *Dependencies) StartChatBroadcastHandler() {
	go func() {
		log.Println("Global broadcast listener goroutine STARTED.")
		for msg := range broadcast {
			// Step 1: Persist the message to the database.
			if err := models.MessageToDatabase(&msg); err != nil {
				log.Printf("DATABASE ERROR: Failed to save message %s: %v", msg.ID, err)
				continue
			}
			log.Printf("Message %s saved to DB.", msg.ID)

			// Step 2: Relay the message in real-time to the involved clients.
			dep.relayMessage(&msg)

			// Step 3: Broadcast an updated user list so everyone's sidebar re-sorts with the new "last message".
			dep.broadcastUserListUpdate()
		}
	}()
}

// relayMessage sends a given message to the sender (for UI confirmation) and the intended receiver.
func (dep *Dependencies) relayMessage(msg *models.Message) {
	ClientsMux.Lock()
	defer ClientsMux.Unlock()

	response := WebSocketMessage{
		Type:     "private_message",
		Content:  msg.Message,
		Messages: []models.Message{*msg}, // Embed the full message object for easy rendering on the frontend.
	}

	// Send to the receiver if they are online.
	if receiverConn, ok := Clients[msg.Receiver]; ok {
		if err := receiverConn.WriteJSON(response); err != nil {
			log.Printf("Error sending private message to receiver %s: %v", msg.Receiver, err)
		} else {
			log.Printf("Relayed message to receiver %s", msg.Receiver)
		}
	}

	// Send a copy to the sender so their UI updates instantly.
	if senderConn, ok := Clients[msg.Sender]; ok {
		if err := senderConn.WriteJSON(response); err != nil {
			log.Printf("Error sending private message copy to sender %s: %v", msg.Sender, err)
		} else {
			log.Printf("Relayed message copy to sender %s", msg.Sender)
		}
	}
}

func (dep *Dependencies) UserIdFromContext(r *http.Request) string {
	userIDFromContext := r.Context().Value("user_uuid")
	if userIDFromContext == nil {
		log.Println("Error: user_uuid not found in context for WebSocket connection")
		conn.Close()
		return " "
	}
	userID := userIDFromContext.(string)
	return userID
}
