package model

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	SendMessageAction = "send-message"
)

type Message struct {
	Action  string  `json:"action"`
	Message string  `json:"message"`
	Target  *Room   `json:"target"`
	Sender  *Client `json:"sender"`
}

func (m *Message) encode() []byte {
	json, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	return json
}

func Decode(m []byte) (Message, error) {
	var message Message
	if err := json.Unmarshal(m, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return Message{}, err
	}
	return message, nil
}
