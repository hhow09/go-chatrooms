package model

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hhow09/go-chatrooms/chatroom-channel/util"
)

type RoomChannel struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	clients   map[*Client]bool
	Broadcast chan Message
	Private   bool
}

func (room *RoomChannel) Setup() {}

// RunRoom runs our room, accepting various requests
func (room *RoomChannel) Run() {
	for {
		message := <-room.Broadcast
		room.BroadcastToClientsInRoom(message)
	}
}

// add client to the room
// then client can receivent the room broadcast
func (room *RoomChannel) RegisterClientInRoom(client *Client, isNewRoom bool) {
	util.Log("RoomChannel.registerClientInRoom")
	room.clients[client] = true
	room.NotifyClientJoined(client, isNewRoom)
}

// remove client from the room
func (room *RoomChannel) UnregisterClientInRoom(client *Client) {
	delete(room.clients, client)
	room.BroadcastToClientsInRoom(Message{Message: fmt.Sprintf(leavRoomMessage, client.Name, room.Name, len(room.clients)), Action: LeaveRoomAction})
}

// broadcast to all client in room
func (room *RoomChannel) BroadcastToClientsInRoom(message Message) {
	util.Log("RoomChannel.BroadcastToClientsInRoom", message.Message)
	for client := range room.clients {
		client.Send(message.encode())
	}
}

// send notification to all clients in room that new client has joined.
func (room *RoomChannel) NotifyClientJoined(client *Client, isNewRoom bool) {
	content := fmt.Sprintf(welcomeMessage, client.GetName(), room.Name, len(room.clients))
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

func (room *RoomChannel) GetName() string {
	return room.Name
}

func (room *RoomChannel) GetId() string {
	return room.ID.String()
}

func (room *RoomChannel) GetBroadcastChan() chan Message {
	return room.Broadcast
}
