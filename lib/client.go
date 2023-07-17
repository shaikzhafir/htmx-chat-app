package lib

import (
	"encoding/json"
	"fmt"
	"htmx-learning/models"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	name        string
	conn        *websocket.Conn
	send        chan []byte
	broadcaster *broadcaster
}

// readPump pumps messages from the websocket connection to the broadcaster
// user types something
// message first received in ReadMessage in readPump

func (c *Client) readPump() {
	defer func() {
		c.broadcaster.leaving <- *c
		leavingMessage := &models.Message{
			User:      c.name,
			Body:      c.conn.RemoteAddr().String() + " left the chat",
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			Type:      "leave",
		}
		// clear the name
		c.broadcaster.AddBackName(c.name)

		jsonBytes, _ := json.Marshal(leavingMessage)
		c.broadcaster.messages <- jsonBytes
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		// do transformation into message struct here
		var htmxMessage models.HTMXMessage
		err = json.Unmarshal(message, &htmxMessage)
		if err != nil {
			fmt.Println(err)
			break
		}
		var msg models.Message
		msg.User = c.name
		msg.Body = htmxMessage.Message
		msg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		msg.Type = "message"
		jsonBytes, _ := json.Marshal(msg)
		// send message to broadcaster
		c.broadcaster.messages <- jsonBytes
		c.broadcaster.tempstore.AddMessage(jsonBytes)
	}
}

// writePump pumps messages from the broadcaster to the websocket connection
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		if !ok {
			// channel closed
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		// add the users name to the message
		var msg models.Message
		err := json.Unmarshal(message, &msg)
		if err != nil {
			log.Println(err)
			return
		}

		username := ""
		// if its a message, add the user's name to the message
		if msg.Type == "message" {
			username = `<span class="font-mono text-sm text-green-600">` + msg.User + `</span>`
		}
		if msg.Type == "enter" {
			username = `<span class="font-mono text-sm text-green-600">` + msg.User + ` entered </span>`
			msg.Body = ""
		}
		if msg.Type == "leave" {
			username = `<span class="font-mono text-sm text-red-600">` + msg.User + ` left </span>`
			msg.Body = ""
		}

		htmlString := `<div id="messages" class="flex" hx-swap-oob="beforeend">
			<div class="w-full flex flex-col">` +
			`<div><span class="font-mono text-sm">` + msg.Timestamp + `:</span>` +
			username +
			`</div>` +
			`<div class="flex flex-col"><span class="break-words">` + msg.Body +
			`</span></div></div></div>`

		c.conn.WriteMessage(websocket.TextMessage, []byte(htmlString))
	}

}
