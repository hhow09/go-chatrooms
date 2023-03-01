package model

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hhow09/go-chatrooms/chatroom-channel/util"
)

const (
	maxMessageSize = 10000
	pongWait       = 15 * time.Second    // short connection
	pingPeriod     = (pongWait * 9) / 10 // pingPeriod should less than pongWait
	writeWait      = 10 * time.Second    // Max wait time when writing message to peer
)

var (
	newline = []byte{'\n'}
)

// Client represents the websocket client at the server
type Client struct {
	conn         *websocket.Conn // The actual websocket connection.
	send         chan []byte
	unregister   chan *Client // send unregister notification to server
	broadcast    chan Message // send broadcast notification to server
	roomActions  chan Message // send room action to server
	ServerNotify chan string  // receive notification from server
	Name         string
	Room         *Room // a client can only join one room
}

func NewClient(conn *websocket.Conn, unregister chan *Client, broadcast chan Message, name string, roomActions chan Message) *Client {
	client := &Client{
		Name:         name,
		conn:         conn,
		unregister:   unregister,
		broadcast:    broadcast,
		send:         make(chan []byte, 256),
		ServerNotify: make(chan string),
		roomActions:  roomActions,
	}
	go client.readMessage()
	go client.writePump()
	return client
}

func (client *Client) disconnect() {
	client.unregister <- client                // unregister client from server
	client.Room.UnregisterClientInRoom(client) // unregister  client from room
	client.send = nil                          // deactivate the channel
	client.conn.Close()
}

func (c *Client) readMessage() {
	defer func() {
		c.disconnect()
	}()
	c.conn.SetReadLimit(maxMessageSize)              // maximum size in bytes for a message read
	c.conn.SetReadDeadline(time.Now().Add(pongWait)) // heartbeat wait time
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait)) // update the wait time
		return nil
	})
	for {
		util.Log("client readMessage")
		_, msg, err := c.conn.ReadMessage()
		util.Log(fmt.Sprintf("received message: %s\n", string(msg)))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("unexpected close error: %v", err)
			}
			break
		}
		c.handleNewMessage(msg)
	}
}

// write message to ws from channel c.send
// also keep the heartbeat with ticker
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // timeout
			if !ok {
				fmt.Println("channel already closed")
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println("error getting writer", err)
				return
			}
			w.Write(msg)

			// Attach queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil { // heeartbeat
				return
			}
		}

	}
}

func (c *Client) Send(msg []byte) {
	c.send <- msg
}

func (c *Client) GetName() string {
	return c.Name
}

func (c *Client) handleNewMessage(rawMessage []byte) {
	message, err := Decode(rawMessage)
	if err != nil {
		fmt.Println("decode error", err)
		return
	}

	message.Sender = c.Name

	switch message.Action {
	case SendMessageAction:
		c.broadcast <- message
	case JoinRoomAction:
		c.handleJoinRoomMessage(message)
	}
}

func (c *Client) handleJoinRoomMessage(msg Message) {
	c.roomActions <- msg
	util.Log(<-c.ServerNotify)
}
