package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/hhow09/go-chatrooms/chatroom/model"
	"github.com/hhow09/go-chatrooms/chatroom/util"
)

type WsServer struct {
	clientMap   map[string]*model.Client // name -> *Client
	register    chan *model.Client
	unregister  chan *model.Client
	broadcast   chan model.Message
	roomActions chan model.Message
	router      *gin.Engine
	roomMap     map[string]model.Room
	redisClient *redis.Client
}

func NewWsServer(r *redis.Client) *WsServer {
	router := gin.New()

	s := &WsServer{
		clientMap:   map[string]*model.Client{},
		register:    make(chan *model.Client),
		unregister:  make(chan *model.Client),
		broadcast:   make(chan model.Message),
		roomActions: make(chan model.Message),
		router:      router,
		roomMap:     map[string]model.Room{},
		redisClient: r,
	}
	s.router.GET("/ws", func(ctx *gin.Context) {
		conn, name := WsHandler(ctx)
		client := model.NewClient(conn, s.unregister, s.broadcast, name, s.roomActions)
		if client != nil {
			s.register <- client
		}
	})
	s.router.GET("/rooms", func(ctx *gin.Context) {
		rooms := make([]string, 0, len(s.roomMap))
		for roomName := range s.roomMap {
			rooms = append(rooms, roomName)
		}

		ctx.JSON(200, rooms)
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
			s.broadcastToRoom(msg)
		case msg := <-s.roomActions:
			s.handleRoomActions(msg)
		}

	}
}

func (s *WsServer) registerClient(client *model.Client) {
	msg := fmt.Sprintf("new client joined: %s", client.GetName())
	if _, ok := s.clientMap[client.Name]; ok {
		// TODO handle duplicate client error
		client.CloseConn()
		return
	}
	util.Log(msg)
	s.clientMap[client.Name] = client
}

func (s *WsServer) unregisterClient(client *model.Client) {
	util.Log("new client levaed")
	delete(s.clientMap, client.Name)
}

func (s *WsServer) broadcastToRoom(message model.Message) {
	room, ok := s.roomMap[message.Target]
	if !ok {
		fmt.Println("cannot find room: ", message.Target)
	}
	_, ok = s.clientMap[message.Sender]
	if !ok {
		fmt.Println("cannot find sender (client): ", message.Sender)
	}
	room.GetBroadcastChan() <- message
}

func (s *WsServer) handleRoomActions(message model.Message) {
	util.Log("wsServer handleRoomActions")
	switch message.Action {
	case model.JoinRoomAction:
		client, ok := s.clientMap[message.Sender]
		if !ok {
			util.Log("client not exist", client.Name)
			return
		}
		if room, ok := s.roomMap[message.Target]; ok {
			// room exist
			util.Log("JoinRoomAction, existing room, client", client.Name)
			room.RegisterClientInRoom(client, false) // add client to room
			client.Room = room
		} else {
			// create room
			util.Log("wsServer JoinRoomAction, creating room")
			room := s.createRoom(message.Target, false)
			room.RegisterClientInRoom(client, true)
			client := s.clientMap[message.Sender]
			client.Room = room
			util.Log("wsServer create room sucess")
		}
		// notify client joined
		client.ServerNotify <- "join room success"
	}
}

func (s *WsServer) createRoom(name string, private bool) model.Room {
	room := model.NewRoom(name, private, s.redisClient)
	room.Setup()
	go room.Run()
	s.roomMap[room.GetName()] = room
	return room
}
