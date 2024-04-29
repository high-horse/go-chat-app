package main

import (
	"log"
	"net/http"
)

func main() {
	println("starting server at port :8000")
	setupAPI()

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func setupAPI() {
	manager := NewManager()

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWs)
}
