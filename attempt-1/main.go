package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	println("starting server at port :8000")
	setupAPI()

	log.Fatal(http.ListenAndServeTLS(":8000", "server.crt", "server.key", nil))
}

func setupAPI() {
	ctx := context.Background()
	manager := NewManager(ctx)

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWs)
	http.HandleFunc("/login", manager.loginHandler)
}
