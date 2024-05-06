package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients CLientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(CLientList),
	}
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

func (m *Manager) addClient(client Client) {
	m.Lock()

	defer m.Unlock()
	m.clients[&client] = true


}

func (m *Manager) removeClient(client Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[&client]; ok {
		client.connection.Close()
		delete(m.clients, &client)
	}
}