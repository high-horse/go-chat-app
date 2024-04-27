package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager

	// egress -> used to avoid concurrent writes on the websocket
	egress chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,

		egress: make(chan Event),
	}
}

func (c *Client) readMessages() {
	defer func() {
		// Cleanup connection
		c.manager.removeClient(c)
	}()

	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling event: %v", err)
			break
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("error handling message: ", err)
		}
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed :", err)
				}
				return

			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			log.Println("message sent")
		}
	}
}
