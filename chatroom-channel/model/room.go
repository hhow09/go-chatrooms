package model

import (
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const welcomeMessage = "%s joined the room [%s]. currently %d user in room."
const leavRoomMessage = "%s leaved the room [%s]. currently %d user in room."

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
}

// NewRoom creates a new RoomChannel
func NewRoom(name string, private bool, redisClient *redis.Client) Room {
	if os.Getenv("REDIS_PUBSUB") == "true" {
		return &RoomPubsub{
			ID:          uuid.New(),
			Name:        name,
			clients:     make(map[*Client]bool),
			Broadcast:   make(chan Message),
			Private:     private,
			redisClient: redisClient,
		}
	}
	return &RoomChannel{
		ID:        uuid.New(),
		Name:      name,
		clients:   make(map[*Client]bool),
		Broadcast: make(chan Message),
		Private:   private,
	}
}
