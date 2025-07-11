package ws

import (
	"log"
)

var startOnce = make(chan struct{}, 1)
var handlerFunc func(string, []byte)

func init() {
	startOnce <- struct{}{}
}

func StartServer() {
	select {
	case <-startOnce:
		go runServer()
	default:
		log.Println("Server already started")
	}
}

func ReceiveMessage(handler func(clientID string, msg []byte)) {
	handlerFunc = handler
}

func SendMessage(clientID string, msg []byte) error {
	return sendToClient(clientID, msg)
}
