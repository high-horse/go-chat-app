package main

import "github.com/gorilla/websocket"

type ClientList map[*Client]bool


type Client struct {
	connection *websocket.Conn
	manager *Manager
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:  manager,
	}
}