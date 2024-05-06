package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn

	manager *Manager
	egress  chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}

func (client *Client) readMessages() {
	defer func() {
		client.manager.removeClient(client)
	}()

	for {
		messagetype, payload, err := client.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message : %v \n", err)
			}
			break
		}
		log.Println(messagetype)
		log.Println(string(payload))

		for wsclient := range client.manager.clients {
			wsclient.egress <- payload
		}
	}
}

func (client *Client) writeMessage() {
	defer func() {
		client.manager.removeClient(client)
	}()

	for {
		select {
		case message, ok := <-client.egress:
			if !ok {
				if err := client.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed :", err)
				}
				return
			}
			if err := client.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println(err)
			}
			log.Println("Message sent.")
		}
	}
}
