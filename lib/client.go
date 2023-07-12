package lib

import (
	"encoding/json"
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
	conn        *websocket.Conn
	send        chan []byte
	broadcaster *broadcaster
}

// readPump pumps messages from the websocket connection to the broadcaster
func (c *Client) readPump() {
	defer func() {
		c.broadcaster.leaving <- *c
		leavingMessage := &models.Message{
			User:      c.conn.RemoteAddr().String(),
			Body:      c.conn.RemoteAddr().String() + " left the chat",
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			Type:      "leave",
		}

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
		c.broadcaster.messages <- message
		c.broadcaster.tempstore.AddMessage(message)
	}
}

// writePump pumps messages from the broadcaster to the websocket connection
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// channel closed
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// add the users name to the message
			var msg models.Message
			var untyped map[string]interface{}

			err := json.Unmarshal(message, &msg)
			if err != nil {
				log.Println(err)
				return
			}
			if msg.Body == "" {
				err := json.Unmarshal(message, &untyped)
				if err != nil {
					log.Println(err)
					return
				}
				msg.Body = untyped["message"].(string)
				msg.User = c.conn.RemoteAddr().String()
				msg.Type = "message"
			}
			msg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			htmlString := `<div id="notifications" class="flex" hx-swap-oob="beforeend">
			<div>` + msg.Timestamp + ": " + msg.Body + `</div>
		   </div>`
			c.conn.WriteMessage(websocket.TextMessage, []byte(htmlString))
		}
	}
}
