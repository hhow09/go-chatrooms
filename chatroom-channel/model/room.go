package model

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hhow09/go-chatrooms/chatroom-channel/util"
)

const welcomeMessage = "%s joined the room [%s]. currently %d user in room."
const leavRoomMessage = "%s leaved the room [%s]. currently %d user in room."

type Room struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	clients   map[*Client]bool
	Broadcast chan Message
	Private   bool
}

// NewRoom creates a new Room
func NewRoom(name string, private bool) *Room {
	return &Room{
		ID:        uuid.New(),
		Name:      name,
		clients:   make(map[*Client]bool),
		Broadcast: make(chan Message),
		Private:   private,
	}
}

// RunRoom runs our room, accepting various requests
func (room *Room) Run() {
	for {
		message := <-room.Broadcast
		room.BroadcastToClientsInRoom(message)
	}
}

// add client to the room
// then client can receivent the room broadcast
func (room *Room) RegisterClientInRoom(client *Client, isNewRoom bool) {
	util.Log("Room.registerClientInRoom")
	room.clients[client] = true
	room.NotifyClientJoined(client, isNewRoom)
}

// remove client from the room
func (room *Room) UnregisterClientInRoom(client *Client) {
	delete(room.clients, client)
	room.BroadcastToClientsInRoom(Message{Message: fmt.Sprintf(leavRoomMessage, client.Name, room.Name, len(room.clients)), Action: LeaveRoomAction})
}

// broadcast to all client in room
func (room *Room) BroadcastToClientsInRoom(message Message) {
	util.Log("Room.BroadcastToClientsInRoom", message.Message)
	for client := range room.clients {
		client.Send(message.encode())
	}
}

// send notification to all clients in room that new client has joined.
func (room *Room) NotifyClientJoined(client *Client, isNewRoom bool) {
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

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetId() string {
	return room.ID.String()
}
