package server

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
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
	conn     *websocket.Conn // The actual websocket connection.
	wsServer *WsServer
	send     chan []byte
}

func newClient(conn *websocket.Conn, wsServer *WsServer) *Client {
	return &Client{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
	}

}

func (client *Client) disconnect() {
	client.wsServer.unregister <- client
	close(client.send)
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
		_, msg, err := c.conn.ReadMessage()
		fmt.Printf("received message: %v\n", string(msg))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("unexpected close error: %v", err)
			}
			break
		}
		c.wsServer.broadcast <- msg
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
