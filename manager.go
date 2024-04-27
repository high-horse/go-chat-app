package main

import (
	// "errors"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	webSocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)


var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

type Manager struct {
	clients ClientList

	// Using a syncMutex here to be able to lcok state before editing clients
	// Could also use Channels to block
	sync.RWMutex
	// handlers are functions that are used to handle Events
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = func(e Event, c *Client) error {
		fmt.Println(e)
		return nil
	}
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

func (m *Manager) serveWs(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")

	// upgrage http to websocket
	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panicln(err)
		return
	}

	// conn.Close()
	client := NewClient(conn, m)

	m.addClient(client)

	// start client processes
	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}
