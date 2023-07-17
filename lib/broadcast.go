package lib

import (
	"encoding/json"
	"htmx-learning/models"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Broadcaster interface {
	Run()
	HandleWebsocket(w http.ResponseWriter, r *http.Request, entering chan Client, leaving chan Client, messages chan []byte)
}

type broadcaster struct {
	ns       NameStore
	entering chan Client
	leaving  chan Client
	messages chan []byte
	// so we can send the last 100 messages to new clients
	tempstore Tempstore
}

type NameStore struct {
	mu        sync.Mutex
	nameStore []string
}

var nameStore = []string{
	"Ziggyzorp",
	"Quibblefink",
	"Muddlechops",
	"Gobbledygook",
	"Blunderbuss",
	"Wobblebottom",
	"Snorklewhip",
	"Jabberwocky",
	"Flibbertigibbet",
	"Woozleflap",
	"Zonkers",
	"Wobbleflop",
	"Noodlebop",
	"Fluffernutter",
	"Bamboozle",
	"Gizmo",
	"Wobblesnatch",
	"Jellybean",
	"Squiggles",
	"Snickerdoodle",
	"Kookaburra",
	"Zigzag",
	"Wobblegobble",
	"Squeegee",
	"Gigglemuffin",
	"Whippersnapper",
	"Muffinhead",
	"Skedaddle",
	"Blunderbluss",
	"Squizzle",
	"Zoinks",
	"Wiggleworm",
	"Dingledorf",
	"Flapdoodle",
	"Quizzlestick",
	"Banjo",
	"Wobblewhisk",
	"Bumblebee",
	"Gobbledygoo",
	"Fuddlewump",
	"Schnookums",
	"Woobleflop",
	"Noodlebrain",
	"Blubberbutt",
	"Whimsy",
	"Goober",
	"Dingleberry",
	"Bafflegab",
	"Flibberflabber",
	"Jibberjabber",
	"Wobblekins",
	"Fuzzbucket",
}

func NewBroadcaster() Broadcaster {
	return &broadcaster{
		ns: NameStore{
			nameStore: nameStore,
		},
		entering:  make(chan Client),
		leaving:   make(chan Client),
		messages:  make(chan []byte),
		tempstore: *NewTempstore(50),
	}
}

func (b *broadcaster) Run() {
	clients := make(map[Client]bool) // all connected clients
	for {
		select {
		case msg := <-b.messages: // incoming message
			// broadcast incoming message to all clients' outgoing message channels
			for cli := range clients {
				cli.send <- msg
			}
		case cli := <-b.entering: // incoming client
			clients[cli] = true
		case cli := <-b.leaving: // leaving client
			delete(clients, cli)
			close(cli.send)
		}
	}
}

func (b *broadcaster) GetRandomName() string {
	randInt := rand.Intn(len(b.ns.nameStore) - 1)
	b.ns.mu.Lock()
	name := b.ns.nameStore[randInt]
	b.ns.nameStore = append(b.ns.nameStore[:randInt], b.ns.nameStore...)
	b.ns.mu.Unlock()
	return name
}

func (b *broadcaster) AddBackName(name string) {
	b.ns.mu.Lock()
	b.ns.nameStore = append(b.ns.nameStore, name)
	b.ns.mu.Unlock()
}

func (b *broadcaster) HandleWebsocket(w http.ResponseWriter, r *http.Request, entering chan Client, leaving chan Client, messages chan []byte) {
	// upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// create client
	client := &Client{name: b.GetRandomName(), conn: conn, send: make(chan []byte, 256), broadcaster: b}
	// assign client a name
	// register client
	b.entering <- *client

	// send last 100 messages to new client
	for _, msg := range b.tempstore.GetMessages() {
		client.send <- msg
	}

	enterMessage := &models.Message{
		User:      client.name,
		Body:      conn.RemoteAddr().String() + " entered",
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Type:      "enter",
	}

	jsonBytes, _ := json.Marshal(enterMessage)
	b.messages <- jsonBytes

	// handle all reads
	go client.readPump()

	// handle all writes
	go client.writePump()

}
