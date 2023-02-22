package main

import (
	"log"

	"github.com/hhow09/chatroom-channel/server"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	wsServer := server.NewWsServer()
	go wsServer.ListenToClientEvents()
	if err := wsServer.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
