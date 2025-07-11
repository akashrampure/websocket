package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // restrict in prod
	},
}

var onReceive MessageHandler

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	clientID := r.RemoteAddr // or generate UUID
	addClient(clientID, conn)
	defer removeClient(clientID)

	conn.SetReadLimit(1024)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		if onReceive != nil {
			onReceive(clientID, msg)
		}
	}
}
