package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool


type Client struct {
	connection *websocket.Conn
	manager *Manager

	// egress -> used to avoid concurrent writes on the websocket
	egress chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:  manager,

		egress: make(chan []byte),
	}
}

func (c *Client) readMessages() {
	defer func ()  {
		// Cleanup connection
		c.manager.removeClient(c)
	}()

	for {
		messagetype, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		for wsclient := range c.manager.clients{
			wsclient.egress <- payload
		}
		log.Println(messagetype)
		log.Println(string(payload))
	}
}

func (c *Client) writeMessages() {
	defer func ()  {
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
			
			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			log.Println("message sent")
		}
	}
}