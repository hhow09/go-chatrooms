package api

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = 15 * time.Second    // short connection
	pingPeriod = (pongWait * 9) / 10 // pingPeriod should less than pongWait
	writeWait  = 10 * time.Second    // Max wait time when writing message to peer
)

type WSClient struct {
	Conn *websocket.Conn
}

// setup ws connection
func NewWSClient(username string) WSClient {
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("localhost:%s", os.Getenv("WEB_HOST")), Path: "/ws", RawQuery: fmt.Sprintf("name=%s", username)}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(fmt.Sprintf("dial:%v", err))
	}

	return WSClient{Conn: c}
}

func (wc *WSClient) Close() {
	wc.Conn.Close()
}

func (wc *WSClient) HeartbeatSetup() *time.Ticker {
	wc.Conn.SetPongHandler(func(appData string) error {
		wc.Conn.SetReadDeadline(time.Now().Add(pongWait)) // update the wait time
		return nil
	})
	wc.Conn.SetReadDeadline(time.Now().Add(pongWait))
	return time.NewTicker(pingPeriod) // heartbeat timer
}
