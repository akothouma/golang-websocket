package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

type UserWithStatus struct{
	UserID string `json:"userid"`
	Username string `json:"username"`
	IsOnline bool `json:"isOnline"`
}

type ClientMessage struct {
	Event   string  `json:"event"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	Msg        string `json:"messageType"`
	ReceiverID string `json:"receiverID"`
	Content    string `json:"content"`
}

var chatBroadcastOnce sync.Once

type ErrorObject struct {
	Error string `json:"error"`
}

type BroadcastMessage struct {
	Message models.Message
}

var (
	Clients    = make(map[string]*websocket.Conn)
	broadcast  = make(chan BroadcastMessage)
	upgrader   = websocket.Upgrader{}
	ClientsMux sync.Mutex
)

func (dep *Dependencies) ChatHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println("failed to establish connection", err)
		return
	}

	//Pass the userID from context setup by middleware
	userIDFromContext := r.Context().Value("user_uuid")
	if userIDFromContext==nil{
		log.Println("Error: user_uuid not found in context for WebSocket connection")
		conn.Close()
		return

	}
	userID:=userIDFromContext.(string)


	go dep.handleClientConnections(userID, conn)

	// Start ChatBroadcastHandler (via sync.Once)
	chatBroadcastOnce.Do(func() {
		dep.StartChatBroadcastHandler()
	})
}


//pass user id as arg now
func (dep *Dependencies) handleClientConnections(userID string, conn *websocket.Conn) {
	defer func() {
		ClientsMux.Lock()
		delete(Clients, userID)
		ClientsMux.Unlock()
		log.Printf("User %s disconnected. Remaining clients: %d", userID, len(Clients))
		dep.broadcastUserStatusUpdate(userID, false) // Broadcast offline status to others
		conn.Close()
	}()

	log.Printf("New WebSocket connection established for user: %s", userID)

	ClientsMux.Lock()
	Clients[userID] = conn
	ClientsMux.Unlock()
	log.Printf("Connected clients count: %d", len(Clients))
	//Broadcast to other connected clients that this user is now online
	// The newly connected user will request the full list themselves.

	dep.broadcastUserStatusUpdate(userID, true)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error for user %s: %v", userID, err)
			ClientsMux.Lock()
			delete(Clients, userID)
			ClientsMux.Unlock()
			log.Printf("User %s disconnected. Remaining clients: %d", userID, len(Clients))
			break
		}

		log.Printf("Raw message received from user %s: %s", userID, string(msg))

		var incoming ClientMessage
		err = json.Unmarshal(msg, &incoming)
		if err != nil {
			log.Printf("JSON unmarshal error from user %s: %v. Raw message: %s", userID, err, string(msg))
			continue
		}

		log.Printf("Parsed message from user %s: Event='%s', MessageType='%s', ReceiverID='%s', Content='%s'",
			userID, incoming.Event, incoming.Payload.Msg, incoming.Payload.ReceiverID, incoming.Payload.Content)

		messageType := incoming.Payload.Msg

		switch messageType {
			//Client requests the full list of users with status
		case "get_user_list":
			log.Printf("Processing get_user_list request from user %s", userID)
			dep.sendFullUserList(conn, userID)
		case "chat_message":
			log.Printf("Processing chat_message from user %s to user %s", userID, incoming.Payload.ReceiverID)
			dep.handleMessageBroadcast(conn, userID, incoming.Payload)
		case "get_message_history":
			dep.GetMessageHistory(conn, incoming.Payload.ReceiverID, userID)
		default:
			log.Printf("Unknown message type '%s' from user %s", messageType, userID)
		}
	}
}

//New function to send the full list of users with their online status
func (dep *Dependencies) sendFullUserList(conn *websocket.Conn, currentUserID string) {
	allDBUsers, err := dep.Forum.GetAllUsers() // Assuming this returns []models.User
	if err != nil {
		log.Printf("Error getting all users from database: %v", err)
		conn.WriteJSON(ErrorObject{Error: "Failed to retrieve user list"})
		return
	}

	usersWithStatus := []UserWithStatus{}
	ClientsMux.Lock() // Lock before accessing Clients map
	for _, dbUser := range allDBUsers {
		var dbUserIDStr string=dbUser.UserID
		_, isOnline := Clients[dbUserIDStr]
		usersWithStatus = append(usersWithStatus, UserWithStatus{
			UserID:   dbUserIDStr,
			Username: dbUser.Username,
			IsOnline: isOnline,
		})
	}
	ClientsMux.Unlock() // Unlock after accessing Clients map

	err = conn.WriteJSON(map[string]interface{}{
		"message":     "user_list_full", // New message type for the client
		"value":       usersWithStatus,
		"currentUser": currentUserID, // Send back the ID of the user making the request
	})
	if err != nil {
		log.Printf("Error sending full user list to %s: %v", currentUserID, err)
	} else {
		log.Printf("Sent full user list to %s. Count: %d", currentUserID, len(usersWithStatus))
	}
}

//New function to broadcast user online/offline status updates
func (dep *Dependencies) broadcastUserStatusUpdate(changedUserID string, isOnline bool) {
	ClientsMux.Lock()
	defer ClientsMux.Unlock()

	message := map[string]interface{}{
		"message": "user_status_update", // New message type for client
		"payload": map[string]interface{}{
			"userID":   changedUserID,
			"isOnline": isOnline,
		},
	}

	log.Printf("Broadcasting status update for user %s: isOnline=%t. Total clients: %d", changedUserID, isOnline, len(Clients))

	for userID, clientConn := range Clients {
		if userID == changedUserID && isOnline {
			log.Printf("Skipping broadcastUserStatusUpdate to self: %s", userID)
			continue
		}

		if err := clientConn.WriteJSON(message); err != nil {
			log.Printf("Error broadcasting status update to user %s for %s: %v", userID, changedUserID, err)
		} else {
			log.Printf("Successfully broadcasted status update for %s to %s", changedUserID, userID)
		}
	}
}


func (dep *Dependencies) StartChatBroadcastHandler() {
	go func() {
		log.Println(">>> Global broadcast listener goroutine STARTED. Listening on broadcast channel. <<<")
		for {
			select {
			case msg := <-broadcast:
				log.Printf(">>> Broadcasting message ID %s from %s to %s <<<", msg.Message.ID, msg.Message.Sender, msg.Message.Receiver)
				dep.broadcastToClients(msg)
			}
		}
	}()
}

func (dep *Dependencies) broadcastToClients(msg BroadcastMessage) {
	ClientsMux.Lock()
	defer ClientsMux.Unlock()

	log.Printf("Broadcasting to %d connected clients", len(Clients))

	// Send to receiver
	receiverId := msg.Message.Receiver
	senderID := msg.Message.Sender // This is the ORIGINAL sender of the message

	log.Println("Receiver ID:", receiverId) // Your debug log

	// Send to receiver
	if receiverConn, ok := Clients[receiverId]; ok {
		err := receiverConn.WriteJSON(map[string]any{
			"message":  "send_private_message",
			"value":    msg.Message.Message,
			"senderID": senderID, // Correct: payload contains the ORIGINAL sender's ID
		})
		if err != nil {
			log.Printf("Error sending to receiver %s: %v", receiverId, err)
			delete(Clients, receiverId)
		} else {
			log.Printf("Message sent to receiver %s", receiverId)
		}
	} else {
		log.Printf("Receiver %s not found in connected clients", receiverId)
	}

	// Send confirmation back to sender
	if senderConn, ok := Clients[senderID]; ok { // Check ORIGINAL sender
		err := senderConn.WriteJSON(map[string]any{
			"message":    "message_sent_confirmation",
			"value":      msg.Message.Message,
			"receiverID": receiverId, // Correct: confirmation includes who it was sent TO
		})
		if err != nil {
			log.Printf("Error sending confirmation to sender %s: %v", senderID, err)
			delete(Clients, senderID)
		}
		// It would be good to have an 'else' log here for successful confirmation sending too:
		// else {
		//	  log.Printf("Message confirmation sent to original sender %s", senderID)
		// }
	} else {
		log.Printf("Original sender %s for confirmation not found in connected clients", senderID)
	}
}

// func (dep *Dependencies) getConnectedUsers(conn *websocket.Conn, currentuser string) {
// 	connectedUserList := []string{}
// 	ClientsMux.Lock()
// 	for userID := range Clients {
// 		connectedUserList = append(connectedUserList, userID)
// 	}
// 	ClientsMux.Unlock()

// 	allConnectedUsers, err := dep.Forum.GetAllConnectedUsers(connectedUserList)
// 	if err != nil {
// 		conn.WriteJSON(ErrorObject{Error: "Something went wrong retrieving connected users"})
// 		return
// 	}

// 	conn.WriteJSON(map[string]any{
// 		"message":     "connected_client_list",
// 		"value":       allConnectedUsers,
// 		"currentUser": currentuser,
// 	})
// }

func (dep *Dependencies) handleMessageBroadcast(c *websocket.Conn, senderid string, p Payload) {
	log.Printf("Handling message broadcast - Sender: %s, Receiver: %s, Content: '%s'", senderid, p.ReceiverID, p.Content)

	if p.ReceiverID == "" {
		log.Printf("Error: Missing receiver ID")
		c.WriteJSON(ErrorObject{Error: "Invalid message: missing receiver"})
		return
	}

	if p.Content == "" {
		log.Printf("Error: Missing message content")
		c.WriteJSON(ErrorObject{Error: "Invalid message: missing content"})
		return
	}

	mess := models.Message{
		ID:        uuid.New(),
		Sender:    senderid,
		Receiver:  p.ReceiverID,
		Message:   p.Content,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
	// Send to broadcast channel FIRST, or make sure it's always sent.
	log.Printf("Queueing message for broadcast (ID: %s): %s -> %s: '%s'", mess.ID, senderid, p.ReceiverID, p.Content)
	// Send to broadcast channel
	broadcast <- BroadcastMessage{
		Message: mess,
	}
	err := models.MessageToDatabase(&mess)
	if err != nil {
		log.Printf("Database error details: %v", err)
		c.WriteJSON(ErrorObject{Error: "Failed to save message. Please try again."})
		return
	} else {
		log.Printf("Message saved to database: %s -> %s: '%s'", senderid, p.ReceiverID, p.Content)
	}
}

func (dep *Dependencies) GetMessageHistory(conn *websocket.Conn, receiver, sender string) {
	prevMess, err := models.MessageHistory(receiver, sender)
	if err != nil {
		fmt.Println("Error getting message history", err)
	}
	fmt.Println("messagesHistory", prevMess)
	conn.WriteJSON(map[string]any{
		"message":    "message_history",
		"value":      prevMess,
		"receiverID": receiver,
		"senderID":sender,
	})
}

// func (dep *Dependencies) GetAllUsers(conn *websocket.Conn) {
// 	allUsers, err := dep.Forum.GetAllUsers()
// 	if err != nil {
// 		fmt.Println("Error getting all users", err)
// 		return
// 	}
// 	conn.WriteJSON(map[string]any{
// 		"message": "all_users",
// 		"value":   allUsers,
// 	})
// }
