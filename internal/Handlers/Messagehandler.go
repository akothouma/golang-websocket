package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

// type Client struct {
// 	Connection *websocket.Conn
// 	UserID     string
// 	IsOnline  bool
// }

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
	Content    string
}

type ErrorObject struct {
	Error string `json:"error"`
}

type BroadcastMessage struct {
	Message models.Message
}

var (
	Clients = make(map[string]*websocket.Conn)
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
	go dep.broadcastToClients(conn)
}

func (dep *Dependencies) handleClientConnections(r *http.Request, conn *websocket.Conn) {
	defer conn.Close()
	userID := r.Context().Value("user_uuid").(string)

	ClientsMux.Lock();
	Clients[userID]=conn;
	ClientsMux.Unlock()
	fmt.Println(Clients)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			ClientsMux.Lock()
			delete(Clients, userID);
			conn.Close();
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
			dep.getConnectedUsers(conn);
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

func (dep *Dependencies) broadcastToClients(sender *websocket.Conn) {
	fmt.Println("broadcast", len(Clients))
	for {
		select {
		case msg := <-broadcast:
			ClientsMux.Lock()
			// Send to receiver
			receiverConn, ok := Clients[msg.Message.Receiver] 
				if ok  {
					receiverConn.WriteJSON(msg.Message)
				}
			
			// Send back to sender
			sender.WriteJSON(msg.Message)
			ClientsMux.Unlock()
		}
	}
}

func (dep *Dependencies) getConnectedUsers(conn *websocket.Conn) {
	userid := []string{}
	for userID := range Clients {
		userid = append(userid, userID)
	}
	allConnectedUsers, err := dep.Forum.GetAllConnectedUsers(userid)
	if err != nil {
		conn.WriteJSON(ErrorObject{Error: "Something went wrong retrieving connected users"})
		return
	}
	conn.WriteJSON(map[string]any{
		"message": "connected_client_list",
		"value":   allConnectedUsers,
	})
}
