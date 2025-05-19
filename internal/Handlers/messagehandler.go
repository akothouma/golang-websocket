package handlers

import (
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

type Client struct {
	Connection *websocket.Conn
	UserID     string
	IsOnline  bool
}


var (
	clients    = make(map[Client]bool)
	broadcast  = make(chan string)
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

	go handleClients(r,conn)
}

func handleClients(r *http.Request, conn *websocket.Conn) {
	userID := r.Context().Value("user_uuid").(string)

	client := Client{
		Connection: conn,
		UserID:     userID,
		IsOnline: true,

	}

	var mess models.Message

	clientsMux.Lock()
	clients[client] = true
	clientsMux.Unlock()

	for{
        _,msg,err:=conn.ReadMessage()
		if err != nil {
			clientsMux.Lock()
			delete(clients, client)
			clientsMux.Unlock()
			break
		}
		mess.Message=string(msg)
		mess.CreatedAt=time.Now()
		mess.ID=uuid.New()
		mess.Sender=userID
		_=mess.MessageToDatabase()

	}
}
