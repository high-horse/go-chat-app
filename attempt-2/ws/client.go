package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type Client struct {
	connection *websocket.Conn

	manager *Manager
	egress  chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (client *Client) readMessages() {
	defer func() {
		client.manager.removeClient(client)
	}()

	client.connection.SetReadLimit(512)

	if err := client.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Panicln(err)
		return
	}

	client.connection.SetPongHandler(client.pongHandler)

	for {
		// messagetype, payload, err := client.connection.ReadMessage()
		_, payload, err := client.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message : %v \n", err)
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("Error marshalling message : %v \n", err)
			break
		}

		if err := client.manager.routeEvent(request, client); err != nil {
			log.Panicln("error handling message : ", err)
		}
	}
}

func (client *Client) pongHandler(pongMsg string) error {
	log.Println("pong")
	return client.connection.SetReadDeadline(time.Now().Add(pongWait))
}

func (client *Client) writeMessage() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()

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
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			if err := client.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}
			log.Println("Message sent.")
		case <-ticker.C:
			log.Println("ping")
			if err := client.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("writemsg: ", err)
				return
			}
		}
	}
}
