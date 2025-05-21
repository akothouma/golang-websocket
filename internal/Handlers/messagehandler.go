package handlers

import (
	"database/sql"
	"encoding/json"
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

type IncomingMessage struct{
	Receiver string `json:"reciever"`

	Content string `json:"content"`
}

type BroadcastMessage struct{
	SenderID string
	ReceiverID string
	Message models.Message
}


var (
	clients    = make(map[string]*websocket.Conn)
	broadcast  = make(chan BroadcastMessage)
	upgrader   = websocket.Upgrader{}
	db         *sql.DB
	clientsMux sync.Mutex
)

func chatHandler(w http.ResponseWriter, r *http.Request) {
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
	defer conn.Close()

	go handleClientConnections(r,conn)
	go broadcastToClients();
}

func handleClientConnections(r *http.Request, conn *websocket.Conn) {
	defer conn.Close()
	userID := r.Context().Value("user_uuid").(string)
	// var mess models.Message

	clientsMux.Lock()
	clients[userID] = conn
	clientsMux.Unlock()

	for{
        _,msg,err:=conn.ReadMessage()
		if err != nil {
			clientsMux.Lock()
			delete(clients, userID)
			clientsMux.Unlock()
			break
		}

		var incoming IncomingMessage
		err=json.Unmarshal(msg,&incoming)
		if err !=nil{
			continue
		}

		//message
		mess:=models.Message{
			ID: uuid.New(),
			Sender: userID,
			Receiver: incoming.Receiver,
			Message: incoming.Content,
			IsRead: false,
			CreatedAt: time.Now(),
			

		}
		_=mess.MessageToDatabase()
		

		//Send to broadcast channel
		broadcast<-BroadcastMessage{
			SenderID: userID,
			ReceiverID: incoming.Receiver,
			Message: mess,

		}
	}
}
func broadcastToClients(){
	for{
		select{
		case msg:=<-broadcast:
			clientsMux.Lock()
			if receiverConnection,ok:=clients[msg.ReceiverID];ok{
				receiverConnection.WriteJSON(msg.Message)
			}
			clientsMux.Unlock();
			clientsMux.Lock()
			if senderConnection,ok:=clients[msg.SenderID];ok{
				senderConnection.WriteJSON(msg.Message);
			}
			clientsMux.Unlock();
		}
	}

}