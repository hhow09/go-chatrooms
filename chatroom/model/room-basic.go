package model

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hhow09/go-chatrooms/chatroom/util"
)

type RoomBasic struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	clients   map[*Client]bool
	Broadcast chan Message
	Private   bool
}

func (room *RoomBasic) Setup() {}

// RunRoom runs our room, accepting various requests
func (room *RoomBasic) Run() {
	for {
		message := <-room.Broadcast
		room.BroadcastToClientsInRoom(message)
	}
}

// add client to the room
// then client can receivent the room broadcast
func (room *RoomBasic) RegisterClientInRoom(client *Client, isNewRoom bool) {
	util.Log("RoomBasic.registerClientInRoom")
	room.clients[client] = true
	room.NotifyClientJoined(client, isNewRoom)
}

// remove client from the room
func (room *RoomBasic) UnregisterClientInRoom(client *Client) {
	delete(room.clients, client)
	room.BroadcastToClientsInRoom(Message{Message: fmt.Sprintf(leavRoomMessage, client.Name, room.Name), Action: LeaveRoomAction})
}

// broadcast to all client in room
func (room *RoomBasic) BroadcastToClientsInRoom(message Message) {
	util.Log("RoomBasic.BroadcastToClientsInRoom", message.Message)
	for client := range room.clients {
		client.Send(message.encode())
	}
}

// send notification to all clients in room that new client has joined.
func (room *RoomBasic) NotifyClientJoined(client *Client, isNewRoom bool) {
	content := fmt.Sprintf(welcomeMessage, client.GetName(), room.Name)
	if isNewRoom {
		content = "new room created. \n" + content
	}
	msg := Message{
		Action:  JoinRoomSuccessAction,
		Target:  room.Name,
		Message: content,
		Sender:  SenderServer,
	}
	room.BroadcastToClientsInRoom(msg)
}

func (room *RoomBasic) GetName() string {
	return room.Name
}

func (room *RoomBasic) GetId() string {
	return room.ID.String()
}

func (room *RoomBasic) GetBroadcastChan() chan Message {
	return room.Broadcast
}

func (room *RoomBasic) GetPrivate() bool {
	return room.Private
}
