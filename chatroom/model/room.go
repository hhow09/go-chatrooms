package model

import (
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hhow09/go-chatrooms/chatroom/util"
)

const welcomeMessage = "%s joined the room [%s]."
const leavRoomMessage = "%s leaved the room [%s]."

type Room interface {
	Setup()
	Run()
	RegisterClientInRoom(client *Client, isNewRoom bool)
	UnregisterClientInRoom(client *Client)
	BroadcastToClientsInRoom(message Message)
	NotifyClientJoined(client *Client, isNewRoom bool)
	GetName() string
	GetId() string
	GetBroadcastChan() chan Message
	GetPrivate() bool
}

func NewRoom(name string, private bool, redisClient *redis.Client) Room {
	if util.IsPubsubEnv() {
		return &RoomPubsub{
			ID:          uuid.New(),
			Name:        name,
			clients:     make(map[*Client]bool),
			Broadcast:   make(chan Message),
			Private:     private,
			redisClient: redisClient,
		}
	}
	return &RoomBasic{
		ID:        uuid.New(),
		Name:      name,
		clients:   make(map[*Client]bool),
		Broadcast: make(chan Message),
		Private:   private,
	}
}
