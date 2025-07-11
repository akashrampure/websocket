package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var clientConn *websocket.Conn
var clientMu sync.Mutex
var clientReceiveHandler func([]byte)

func Subscribe(url string) {
	var err error
	clientConn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("dial error: %v", err)
	}
	log.Println("Connected to server:", url)

	go func() {
		for {
			_, msg, err := clientConn.ReadMessage()
			if err != nil {
				log.Println("client read error:", err)
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
		return ErrClientNotConnected
	}
	return clientConn.WriteMessage(websocket.TextMessage, msg)
}
