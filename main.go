package main

import (
	"fmt"
	"htmx-learning/lib"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello, World!")
		fmt.Fprintf(w, "Hello, World!")
	})

	var (
		entering = make(chan lib.Client)
		leaving  = make(chan lib.Client)
		messages = make(chan []byte) // all incoming client messages
	)

	// broadcaster is a goroutine that broadcasts messages to all clients
	broadcaster := lib.NewBroadcaster()
	go broadcaster.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// handle websocket
		broadcaster.HandleWebsocket(w, r, entering, leaving, messages)
	})

	http.ListenAndServe(":8080", nil)
}
