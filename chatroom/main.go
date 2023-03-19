package main

import (
	"log"

	"github.com/hhow09/go-chatrooms/chatroom/lib"
	"github.com/hhow09/go-chatrooms/chatroom/server"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	redisc := lib.CreateRedisClient()
	defer redisc.Close()
	db := lib.InitDB()
	defer db.Close()
	wsServer := server.NewWsServer(redisc, db)
	go wsServer.ListenToClientEvents()
	if err := wsServer.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
