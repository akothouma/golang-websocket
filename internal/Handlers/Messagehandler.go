// internal/handlers/messagehandler.go

package handlers

import (
	// "encoding/json"
	"log"
	"net/http"
	// "sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

// FIX: Renamed and enriched this struct. This is what represents a user in the chat list.
type UserChatInfo struct {
	UserID             string    `json:"userID"`
	Username           string    `json:"username"`
	IsOnline           bool      `json:"isOnline"`
	LastMessageContent string    `json:"lastMessageContent"`
	LastMessageTime    time.Time `json:"lastMessageTime"`
	IsMe               bool      `json:"isMe,omitempty"`
}

// FIX: Simplified the message structure for both incoming and outgoing websocket messages.
type WebSocketMessage struct {
	Type            string      `json:"type"` // e.g., "get_user_list", "private_message", "get_message_history"
	Content         string      `json:"content,omitempty"`
	Target          string      `json:"target,omitempty"`      // The userID of the recipient
	Messages        []models.Message `json:"messages,omitempty"`    // For sending history
	UserList        []UserChatInfo   `json:"userList,omitempty"`    // For sending user list
	LastMessageTime time.Time        `json:"lastMessageTime,omitempty"` // For pagination
}

type ErrorObject struct {
	Error string `json:"error"`
}

var (
	Clients         = make(map[string]*websocket.Conn)
	broadcast       = make(chan models.Message) // FIX: Channel will just carry the message model
	upgrader        = websocket.Upgrader{}
	ClientsMux      sync.Mutex
	broadcastOnce   sync.Once
)

func (dep *Dependencies) ChatHandler(w http.ResponseWriter, r *http.Request) {
    // ... your existing ChatHandler code (no changes needed here)
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to establish connection", err)
		return
	}

	userIDFromContext := r.Context().Value("user_uuid")
	if userIDFromContext == nil {
		log.Println("Error: user_uuid not found in context for WebSocket connection")
		conn.Close()
		return
	}
	userID := userIDFromContext.(string)

	go dep.handleClientConnections(userID, conn)
	broadcastOnce.Do(dep.StartChatBroadcastHandler)
}

func (dep *Dependencies) handleClientConnections(userID string, conn *websocket.Conn) {
	// Register client
	ClientsMux.Lock()
	Clients[userID] = conn
	ClientsMux.Unlock()
	log.Printf("User %s connected. Total clients: %d", userID, len(Clients))
	
	// Announce to everyone that a user's status changed
	dep.broadcastUserListUpdate()
	
	defer func() {
		ClientsMux.Lock()
		delete(Clients, userID)
		ClientsMux.Unlock()
		log.Printf("User %s disconnected. Remaining clients: %d", userID, len(Clients))
		// Announce to everyone that a user's status changed
		dep.broadcastUserListUpdate()
		conn.Close()
	}()

	for {
		// FIX: Use the new simplified WebSocketMessage struct for reading
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error for user %s: %v", userID, err)
			}
			break
		}
		
		log.Printf("Parsed message from user %s: Type='%s', Target='%s', ContentLength=%d", userID, msg.Type, msg.Target, len(msg.Content))

		switch msg.Type {
		case "get_user_list":
			// The user list is now broadcast automatically, but we can send it on demand if needed.
			dep.sendFullUserListToUser(userID)

		case "private_message":
			if msg.Target == "" || msg.Content == "" {
				log.Printf("Invalid private message from %s: missing target or content", userID)
				continue
			}
			// Create the message model
			messageModel := models.Message{
				ID:        uuid.New(),
				Sender:    userID,
				Receiver:  msg.Target,
				Message:   msg.Content,
				IsRead:    false, // Will be false until the receiver actually opens the chat
				CreatedAt: time.Now(),
			}
			// Send to broadcast channel to be relayed and saved
			broadcast <- messageModel

		case "get_message_history":
			if msg.Target == "" {
				log.Printf("Invalid history request from %s: missing target", userID)
				continue
			}
			dep.sendMessageHistory(conn, userID, msg.Target, msg.LastMessageTime)
		
		default:
			log.Printf("Unknown message type '%s' from user %s", msg.Type, userID)
		}
	}
}


// ADDED: New function to send history with pagination
func (dep *Dependencies) sendMessageHistory(conn *websocket.Conn, senderID, targetID string, lastTimestamp time.Time) {
	messages, err := models.GetMessageHistory(senderID, targetID, lastTimestamp, 10)
	if err != nil {
		log.Printf("Error getting message history for %s<->%s: %v", senderID, targetID, err)
		// Optionally send an error message to the client
		return
	}

	response := WebSocketMessage{
		Type:     "history_response",
		Target:   targetID, // The user whose history was requested
		Messages: messages,
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending message history to %s: %v", senderID, err)
	}
}

// REPLACED `sendFullUserList` and `broadcastUserStatusUpdate` with this unified function

// FIX: This is the critical function to change. We will iterate and send a custom list to each user.
func (dep *Dependencies) broadcastUserListUpdate() {
	allDBUsers, err := dep.Forum.GetAllUsers()
	if err != nil {
		log.Printf("broadcastUserListUpdate: Error getting all users: %v", err)
		return
	}

	lastMessages, err := models.GetAllLastMessages()
	if err != nil {
		log.Printf("broadcastUserListUpdate: Error getting last messages: %v", err)
		return
	}

	lastMessageMap := make(map[string]models.Message)
	for _, msg := range lastMessages {
		key := models.GetConversationID(msg.Sender, msg.Receiver)
		lastMessageMap[key] = msg
	}

	ClientsMux.Lock()
	defer ClientsMux.Unlock()

	log.Printf("Broadcasting tailored user lists to %d clients.", len(Clients))
	
	// For each connected client...
	for userID, clientConn := range Clients {
		
		var usersWithStatus []UserChatInfo // Create a fresh list for this specific user

		// Build the base list of all users from the DB
		for _, dbUser := range allDBUsers {
			_, isOnline := Clients[dbUser.UserID]

			userInfo := UserChatInfo{
				UserID:   dbUser.UserID,
				Username: dbUser.Username,
				IsOnline: isOnline,
				// Add the IsMe flag here!
				IsMe:     (dbUser.UserID == userID),
			}
			usersWithStatus = append(usersWithStatus, userInfo)
		}
		
		// Enrich with last message data
		for i, u := range usersWithStatus {
			for _, otherUser := range usersWithStatus {
				if u.UserID == otherUser.UserID { continue }

				key := models.GetConversationID(u.UserID, otherUser.UserID)
				if msg, ok := lastMessageMap[key]; ok {
					// We need to find the most recent message across all conversations for a user
					if msg.CreatedAt.After(usersWithStatus[i].LastMessageTime) {
						usersWithStatus[i].LastMessageContent = msg.Message
						usersWithStatus[i].LastMessageTime = msg.CreatedAt
					}
				}
			}
		}

		// Now send this tailored list to the current client in the loop
		response := WebSocketMessage{
			Type:     "user_list_update",
			UserList: usersWithStatus,
		}
		if err := clientConn.WriteJSON(response); err != nil {
			log.Printf("Error broadcasting tailored user list to user %s: %v", userID, err)
		}
	}
}


// ADDED: Helper to send the full list to just one user on demand
func (dep *Dependencies) sendFullUserListToUser(userID string) {
	// This function body would be almost identical to broadcastUserListUpdate, 
	// but only sends to the specific user's connection. 
	// For simplicity, we can just trigger a full broadcast which is okay for this app scale.
	dep.broadcastUserListUpdate()
}


func (dep *Dependencies) StartChatBroadcastHandler() {
	go func() {
		log.Println("Global broadcast listener goroutine STARTED.")
		for msg := range broadcast {
			// Save to database
			if err := models.MessageToDatabase(&msg); err != nil {
				log.Printf("DATABASE ERROR: Failed to save message %s: %v", msg.ID, err)
				continue
			}
			log.Printf("Message %s saved to DB.", msg.ID)
			
			// Relay to clients
			dep.relayMessage(&msg)
			
			// After a message is successfully sent, broadcast an updated user list
			// so the order changes in everyone's sidebar.
			dep.broadcastUserListUpdate()
		}
	}()
}


// ADDED: New function to just relay messages to relevant parties.
func (dep *Dependencies) relayMessage(msg *models.Message) {
	ClientsMux.Lock()
	defer ClientsMux.Unlock()

	response := WebSocketMessage{
		Type:    "private_message",
		Content: msg.Message,
		// This response needs to carry the full message object for rendering
		Messages: []models.Message{*msg},
	}
	
	// Send to receiver
	if receiverConn, ok := Clients[msg.Receiver]; ok {
		if err := receiverConn.WriteJSON(response); err != nil {
			log.Printf("Error sending private message to receiver %s: %v", msg.Receiver, err)
		} else {
			log.Printf("Relayed message to receiver %s", msg.Receiver)
		}
	}

	// Send back to sender so they can see their own message appear in the chat
	if senderConn, ok := Clients[msg.Sender]; ok {
		if err := senderConn.WriteJSON(response); err != nil {
			log.Printf("Error sending private message copy to sender %s: %v", msg.Sender, err)
		} else {
			log.Printf("Relayed message copy to sender %s", msg.Sender)
		}
	}
}

// DELETE or comment out the old `broadcastToClients` and `handleMessageBroadcast` as they are now replaced.