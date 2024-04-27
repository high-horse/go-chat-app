package main

import "net/http"

func main() {
	setupAPI()

	http.ListenAndServe(":8080", nil)
}

func setupAPI() {
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
}
