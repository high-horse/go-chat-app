package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	// websocketUpgrader = websocket.Upgrader{
	// 	ReadBUddderSize: 1024,
	// 	WriteBuffersize: 1024,
	// }
)

type Manager struct {}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) ServerWS(w http.ResponseWriter, r *http.Request) {
	println("new connection")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn.Close()
}