package main

import (
	"htmx-learning/lib"
	"log"
	"net/http"
	"os"
)

func main() {

	// set up logger
	fileName := "chatapp.log"
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	defer logFile.Close()

	log.SetOutput(logFile)

	// serve static files like index.html
	http.Handle("/", http.FileServer(http.Dir("./")))

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
	isDev := os.Getenv("DEV")
	if isDev == "true" {
		log.Fatal(http.ListenAndServe(":8080", nil))

	} else {
		log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "cert.key", nil))
	}
}
