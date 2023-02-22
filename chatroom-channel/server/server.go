package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hhow09/go-chatrooms/chatroom-channel/model"
)

type WsServer struct {
	clients    map[*model.Client]bool
	register   chan *model.Client
	unregister chan *model.Client
	broadcast  chan []byte
	router     *gin.Engine
}

func NewWsServer() *WsServer {
	router := gin.New()

	s := &WsServer{
		clients:    make(map[*model.Client]bool),
		register:   make(chan *model.Client),
		unregister: make(chan *model.Client),
		broadcast:  make(chan []byte),
		router:     router,
	}
	s.router.GET("/ws", func(ctx *gin.Context) {
		conn, name := WsHandler(ctx)
		client := model.NewClient(conn, s.unregister, s.broadcast, name)
		if client != nil {
			s.register <- client
		}
	})
	return s
}

func (s *WsServer) Run() error {
	err := s.router.Run(fmt.Sprintf(":%v", os.Getenv("WEB_HOST")))
	if err != nil {
		return err
	}
	return nil
}

func (s *WsServer) ListenToClientEvents() {
	for {
		select {
		case client := <-s.register:
			s.registerClient(client)
		case client := <-s.unregister:
			s.unregisterClient(client)
		case msg := <-s.broadcast:
			s.broadcastToClients(msg)
		}

	}
}

func (s *WsServer) registerClient(client *model.Client) {
	msg := fmt.Sprintf("new client joined: %s", client.GetName())
	fmt.Println(msg)
	s.broadcastToClients([]byte(msg))
	s.clients[client] = true
}

func (s *WsServer) unregisterClient(client *model.Client) {
	fmt.Println("new client levaed")
	delete(s.clients, client)
}

func (s *WsServer) broadcastToClients(message []byte) {
	for client := range s.clients {
		client.Send(message)
	}
}
