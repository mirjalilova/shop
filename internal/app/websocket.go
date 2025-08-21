package app

import (
	"log"
	"net/http"
	"shop/internal/usecase"
	"shop/pkg/ws"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

func WebSocketHandler(db *usecase.UseCase, hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🔥 WebSocket handler called: %s", r.URL.Path)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}

		hub.AddClient(conn)
		log.Printf("✅ New connection. Active: %d", hub.Count())

		go func() {
			defer func() {
				hub.RemoveClient(conn)
				log.Printf("❌ Connection closed. Active: %d", hub.Count())
			}()

			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					break
				}
			}
		}()
	}
}
