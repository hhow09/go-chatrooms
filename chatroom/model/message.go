package model

import (
	"encoding/json"
	"fmt"
	"log"
)

type Action string

const (
	SendMessageAction     Action = "send-message"
	JoinRoomAction        Action = "join-room"
	LeaveRoomAction       Action = "leave-room"
	JoinRoomSuccessAction Action = "join-room-success"
	SenderServer                 = "server"
)

type Message struct {
	Action  Action `json:"action"`
	Message string `json:"message"`
	Target  string `json:"target"`
	Sender  string `json:"sender"`
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
