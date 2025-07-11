package ws

import (
	"log"
	"net/http"
	"sync"
	"time"
)

var serverOnce sync.Once

func StartServer() {
	serverOnce.Do(func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/ws", wsHandler)
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status":"ok"}`))
		})

		srv := &http.Server{
			Addr:         ":8080",
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		log.Println("WebSocket server started on :8080")
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				log.Println("server closed:", err)
			}
		}()
	})
}

func ReceiveMessage(handler MessageHandler) {
	onReceive = handler
}

func SendMessage(clientID string, msg []byte) error {
	return sendToClient(clientID, msg)
}
