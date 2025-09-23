package configs

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var Upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Connection struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

var WebsocketConnections = struct {
	sync.RWMutex
	Connections map[uint]*Connection
}{
	Connections: make(map[uint]*Connection),
}
