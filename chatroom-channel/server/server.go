package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	router     *gin.Engine
}

func NewWsServer() *WsServer {
	router := gin.New()

	s := &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		router:     router,
	}
	s.router.GET("/ws", func(ctx *gin.Context) {
		conn := WsHandler(ctx)
		client := newClient(conn, s)
		go client.readMessage()
		go client.writePump()
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

func (s *WsServer) registerClient(client *Client) {
	fmt.Println("new client joined")
	client.send <- []byte("welcome to the chatroom!")
	s.broadcastToClients([]byte("new client joined"))
	s.clients[client] = true
}

func (s *WsServer) unregisterClient(client *Client) {
	fmt.Println("new client levaed")
	delete(s.clients, client)
}

func (s *WsServer) broadcastToClients(message []byte) {
	for client := range s.clients {
		client.send <- message
	}
}
