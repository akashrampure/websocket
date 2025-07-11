package ws

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	clientConn           *websocket.Conn
	clientMu             sync.Mutex
	clientReceiveHandler func([]byte)
	maxReconnectAttempts = 5
	reconnectDelay       = 2 * time.Second
)

func Subscribe(url string) {
	var err error
	clientConn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		clientConn, err = reconnectToServer(url)
		if err != nil {
			log.Fatalf("Failed to reconnect to server: %v", err)
		}
	}
	log.Println("Connected to server:", url)

	go func() {
		for {
			_, msg, err := clientConn.ReadMessage()
			if err != nil {
				log.Println("client read error:", err)
				if websocket.IsUnexpectedCloseError(err) {
					log.Println("Connection closed unexpectedly, attempting to reconnect...")
					reconnectToServer(url)
				}
				break
			}
			if clientReceiveHandler != nil {
				clientReceiveHandler(msg)
			}
		}
	}()
}

func OnMessageFromServer(handler func([]byte)) {
	clientReceiveHandler = handler
}

func SendToServer(msg []byte) error {
	clientMu.Lock()
	defer clientMu.Unlock()

	if clientConn == nil {
		return fmt.Errorf("client is not connected")
	}
	return clientConn.WriteMessage(websocket.TextMessage, msg)
}

func reconnectToServer(url string) (*websocket.Conn, error) {
	for i := 0; i < maxReconnectAttempts; i++ {
		if clientConn != nil {
			clientConn.Close()
		}
		var err error
		clientConn, _, err = websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Printf("Reconnect attempt %d failed: %v", i+1, err)
			time.Sleep(reconnectDelay)
		} else {
			log.Println("Connected to server:", url)
			return clientConn, nil
		}
	}
	return nil, fmt.Errorf("failed to reconnect to server after %d attempts", maxReconnectAttempts)
}
