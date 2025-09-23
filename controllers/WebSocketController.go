package controllers

import (
	"bougette-backend/common"
	"bougette-backend/configs"
	"strconv"

	"github.com/labstack/echo/v4"
)

func HandleWebSocket(c echo.Context) error {
	userID := c.Param("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid user ID")
	}

	conn, err := configs.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Failed to upgrade to WebSocket")
	}
	defer conn.Close()

	configs.WebsocketConnections.Lock()
	configs.WebsocketConnections.Connections[uint(id)] = &configs.Connection{
		Conn: conn,
	}
	configs.WebsocketConnections.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			configs.WebsocketConnections.Lock()
			delete(configs.WebsocketConnections.Connections, uint(id))
			configs.WebsocketConnections.Unlock()
			break
		}
	}

	return nil
}
