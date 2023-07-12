package lib

type Tempstore struct {
	// so we can send the last 100 messages to new clients
	messages [][]byte
	capacity int
}

func NewTempstore(capacity int) *Tempstore {
	return &Tempstore{
		messages: [][]byte{},
		capacity: capacity,
	}
}

func (t *Tempstore) AddMessage(msg []byte) {
	// check if we need to overwrite the oldest message
	if t.capacity == len(t.messages) {
		// drop the oldest message
		t.messages = t.messages[1:]
		t.messages = append(t.messages, msg)
	} else {
		t.messages = append(t.messages, msg)
	}
}

func (t *Tempstore) GetMessages() [][]byte {
	return t.messages
}
