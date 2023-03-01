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
	JoinRoomSuccessAction Action = "join-room-success"
)

type Message struct {
	Action  Action `json:"action"`
	Message string `json:"message"`
	Target  string `json:"target"`
	Sender  string `json:"sender"`
}

func (m *Message) Encode() []byte {
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

func NewTextMessage(username, content, room string) *Message {
	return &Message{
		Action:  SendMessageAction,
		Message: content,
		Target:  room,
		Sender:  username,
	}
}

func NewJoinRoomMessage(username, room string) *Message {
	return &Message{
		Action:  JoinRoomAction,
		Target:  room,
		Sender:  username,
		Message: "",
	}
}
