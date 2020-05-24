package model

import (
	"time"
)

// Message is a model to store messages
type Message struct {
	Base

	PeerID string `json:"peer_id"`
	Data   string `json:"data"`
}

// NewMessage returns instatiated Message struct
func NewMessage() *Message {
	message := new(Message)
	message.CreatedAt = time.Now().UTC()
	message.UpdatedAt = time.Now().UTC()

	return message
}
