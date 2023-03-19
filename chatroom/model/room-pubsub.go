package model

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hhow09/go-chatrooms/chatroom/util"
)

type RoomPubsub struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	clients     map[*Client]bool
	Broadcast   chan Message
	Private     bool
	redisClient *redis.Client
	sub         *redis.PubSub
}

func (room *RoomPubsub) Setup() {
	ctx := context.Background()
	sub := room.redisClient.Subscribe(ctx, room.GetName())
	_, err := sub.Receive(ctx)
	if err != nil {
		util.Log("subscribe error", err.Error())
	}
	room.sub = sub
}

// RunRoom runs our room, accepting various requests
func (room *RoomPubsub) Run() {
	go room.subscribeToRoomMessages(room.sub.Channel())
	for {
		message := <-room.Broadcast
		room.BroadcastToClientsInRoom(message)
	}
}

// add client to the room
// then client can receivent the room broadcast
func (room *RoomPubsub) RegisterClientInRoom(client *Client, isNewRoom bool) {
	util.Log("RoomPubsub.registerClientInRoom")
	room.clients[client] = true
	room.NotifyClientJoined(client, isNewRoom)
}

// remove client from the room
func (room *RoomPubsub) UnregisterClientInRoom(client *Client) {
	delete(room.clients, client)
	room.BroadcastToClientsInRoom(Message{Message: fmt.Sprintf(leavRoomMessage, client.Name, room.Name), Action: LeaveRoomAction})
}

// broadcast to all client in room
func (room *RoomPubsub) BroadcastToClientsInRoom(message Message) {
	util.Log("RoomPubsub.BroadcastToClientsInRoom", message.Message)
	err := room.redisClient.Publish(context.Background(), room.GetName(), message.encode()).Err()
	if err != nil {
		log.Println(err)
	}

}

// send notification to all clients in room that new client has joined.
func (room *RoomPubsub) NotifyClientJoined(client *Client, isNewRoom bool) {
	util.Log("RoomPubsub.NotifyClientJoined", client.Name)
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

func (room *RoomPubsub) GetName() string {
	return room.Name
}

func (room *RoomPubsub) GetId() string {
	return room.ID.String()
}

func (room *RoomPubsub) GetBroadcastChan() chan Message {
	return room.Broadcast
}

func (room *RoomPubsub) subscribeToRoomMessages(ch <-chan *redis.Message) {
	util.Log("RoomPubsub subscribeToRoomMessages")
	for msg := range ch {
		m, err := Decode([]byte(msg.Payload))
		if err != nil {
			util.Log("RoomPubsub subscribeToRoomMessages decode error: ", err.Error())
		}
		for client := range room.clients {
			client.Send(m.encode())
		}
	}
}

func (room *RoomPubsub) GetPrivate() bool {
	return room.Private
}
