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

type ClientConnection struct {
	Event   string  `json:"event"`
	Payload Payload `json:"payload"`
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
	// defer conn.Close()

	go dep.handleClientConnections(r, conn)
	// 5. Start ChatBroadcastHandler (via sync.Once)
	chatBroadcastOnce.Do(func() {
		dep.StartChatBroadcastHandler()
	})
}

func (dep *Dependencies) handleClientConnections(r *http.Request, conn *websocket.Conn) {
	defer conn.Close()

	userID := r.Context().Value("user_uuid").(string)
	ClientsMux.Lock()
	Clients[userID] = conn
	ClientsMux.Unlock()
	fmt.Println(Clients)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			ClientsMux.Lock()
			delete(Clients, userID)
			conn.Close()
			ClientsMux.Unlock()
			break

		}

		var incoming ClientConnection

		err = json.Unmarshal(msg, &incoming)
		if err != nil {
			continue
		}
		fmt.Println(incoming)
		messageType := incoming.Payload.Msg

		switch messageType {
		case "get_online_users":
			dep.getConnectedUsers(conn)
			break // Won't create message for this instance now
		case "chat_message":
			// We only creating message for actual chat messages now
			mess := models.Message{
				ID:        uuid.New(),
				Sender:    userID,
				Receiver:  incoming.Payload.ReceiverID,
				Message:   incoming.Payload.Content,
				IsRead:    false,
				CreatedAt: time.Now(),
			}
			_ = mess.MessageToDatabase()

			// Send to broadcast channel
			broadcast <- BroadcastMessage{
				Message: mess,
			}
		}

	}
}

// This function starts the single, global broadcast handler for chat messages
func (dep *Dependencies) StartChatBroadcastHandler() {
	go func() {
		log.Println(">>> Global broadcast listener goroutine STARTED. Listening on broadcast channel. <<<")
		for {
			select {
			case msg := <-broadcast:
				log.Printf(">>> Broadcast listener: Received message ID %s from broadcast channel for sender %s to receiver %s. <<<", msg.Message.ID, msg.Message.Sender, msg.Message.Receiver)
				dep.broadcastToClients(msg)
			}
		}
	}()
}

func (dep *Dependencies) broadcastToClients(msg BroadcastMessage) {
	ClientsMux.Lock()
	defer ClientsMux.Unlock()
	fmt.Println("broadcast", len(Clients))

	// senderID := msg.Message.Sender
	// receiverID := msg.Message.Receiver
	for {
		select {
		case msg := <-broadcast:
			ClientsMux.Lock()
			// Send to receiver
			receiverConn, ok := Clients[msg.Message.Receiver]
			if ok {
				receiverConn.WriteJSON(msg.Message)
			}

			// Send back to sender
			// sender.WriteJSON(msg.Message)
			ClientsMux.Unlock()
		}
	}
}

func (dep *Dependencies) getConnectedUsers(conn *websocket.Conn) {
	connectedUserList := []string{}
	for userID := range Clients {
		connectedUserList = append(connectedUserList, userID)
	}
	allConnectedUsers, err := dep.Forum.GetAllConnectedUsers(connectedUserList)
	if err != nil {
		conn.WriteJSON(ErrorObject{Error: "Something went wrong retrieving connected users"})
		return
	}

	response := map[string]interface{}{
		"message": "connected_client_list", // Frontend expects this key
		"value":   allConnectedUsers,       // Frontend expects this key
	}
	conn.WriteJSON(response)
}
