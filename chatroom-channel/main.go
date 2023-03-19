package main

import (
	"log"

	"github.com/hhow09/go-chatrooms/chatroom-channel/lib"
	"github.com/hhow09/go-chatrooms/chatroom-channel/server"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	redisc := lib.CreateRedisClient()
	defer redisc.Close()
	wsServer := server.NewWsServer(redisc)
	go wsServer.ListenToClientEvents()
	if err := wsServer.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
