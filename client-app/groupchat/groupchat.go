package groupchat

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gorilla/websocket"
	"github.com/hhow09/go-chatrooms/client-app/input"
)

func GroupChatProgram(username string) {
	// init keyboard reader
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	flag.Parse()
	log.SetFlags(0)

	// exit sig
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("localhost:%s", os.Getenv("WEB_HOST")), Path: "/ws", RawQuery: fmt.Sprintf("name=%s", username)}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(fmt.Sprintf("dial:%v", err))
	}
	defer c.Close()

	done := make(chan struct{})

	// receive message from ws
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("unexpected close error: %v", err)
				}
				fmt.Println("err:", err)
				return
			}
			fmt.Printf("recv: %s\n", message)
		}
	}()

	// input reader
	ichan := input.NewInput(interrupt, []keyboard.Key{keyboard.KeyEsc, keyboard.KeyCtrlC})

	ticker := time.NewTicker(time.Second) // heartbeat timer
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case input := <-ichan:
			err := c.WriteMessage(websocket.TextMessage, []byte(input))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.PingMessage, []byte(t.String()))
			if err != nil {
				log.Println("heartbeat error:", err)
				return
			}
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
