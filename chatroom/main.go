package main

import (
	"database/sql"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/hhow09/go-chatrooms/chatroom/lib"
	"github.com/hhow09/go-chatrooms/chatroom/server"
	"github.com/hhow09/go-chatrooms/chatroom/util"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	var redisc *redis.Client
	var db *sql.DB
	if util.IsPubsubEnv() {
		redisc = lib.CreateRedisClient()
		defer redisc.Close()
		db = lib.InitDB()
		defer db.Close()
	}

	wsServer := server.NewWsServer(redisc, db)
	go wsServer.ListenToClientEvents()
	if err := wsServer.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
