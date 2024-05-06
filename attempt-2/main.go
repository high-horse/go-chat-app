package main

import (
	"log"
	"net/http"

	"go-chat-app/ws"
)

func main() {
	setupApi()

	println("Server started at :3000")

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func setupApi() {
	manager := ws.NewManager()
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.ServerWS)
}