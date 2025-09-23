package utilities

import (
	"bougette-backend/configs"
	"log"
)

type SendWebsocketMessagePayload struct {
	UserID  uint   `json:"user_id"`
	Message string `json:"message"`
	IsRead  bool   `json:"is_read"`
}

func SendWebsocketMessage(userID uint, payload SendWebsocketMessagePayload) {
	configs.WebsocketConnections.RLock()
	defer configs.WebsocketConnections.RUnlock()

	conn, ok := configs.WebsocketConnections.Connections[userID]
	if ok {
		conn.Mutex.Lock()
		err := conn.Conn.WriteJSON(payload)
		conn.Mutex.Unlock()
		if err != nil {
			log.Println("Error sending websocket message to user", userID, err)
			configs.WebsocketConnections.Lock()
			delete(configs.WebsocketConnections.Connections, userID)
			configs.WebsocketConnections.Unlock()
		}
	} else {
		log.Println("Websocket connection not found for user", userID)
	}
}
