package ws

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[string]*websocket.Conn)
	clientsMu sync.RWMutex
)

func addClient(id string, conn *websocket.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	clients[id] = conn
}

func removeClient(id string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	delete(clients, id)
}

func sendToClient(id string, msg []byte) error {
	clientsMu.RLock()
	defer clientsMu.RUnlock()
	conn, ok := clients[id]
	if !ok {
		return errors.New("client not found")
	}
	return conn.WriteMessage(websocket.TextMessage, msg)
}
