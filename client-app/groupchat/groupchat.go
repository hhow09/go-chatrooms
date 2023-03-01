package groupchat

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gorilla/websocket"
	"github.com/hhow09/go-chatrooms/client-app/input"
	"github.com/hhow09/go-chatrooms/client-app/model"
)

func getRoomList() ([]string, error) {
	u := url.URL{Scheme: "http", Host: fmt.Sprintf("localhost:%s", os.Getenv("WEB_HOST")), Path: "/rooms"}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Non-OK HTTP status: %v", res.StatusCode)
	}
	defer res.Body.Close()
	var rommlist []string
	err = json.NewDecoder(res.Body).Decode(&rommlist)
	if err != nil {
		return nil, err
	}
	return rommlist, nil
}

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

	// setup ws connection
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("localhost:%s", os.Getenv("WEB_HOST")), Path: "/ws", RawQuery: fmt.Sprintf("name=%s", username)}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(fmt.Sprintf("dial:%v", err))
	}
	defer c.Close()

	var room string
	roomList, _ := getRoomList()
	ans1 := ChooseRoomPrompt(roomList)
	if ans1.Room == OPTION_CREATE_ROOM || ans1.Room == "" {
		room = CreateRoomPrompt().Room
	} else {
		room = ans1.Room
	}

	done := make(chan struct{})

	// join room
	err = joinRoom(c, username, room)
	if err != nil {
		panic(fmt.Sprintf("join room failed %v", err))
	}
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
			handleReceiveMessage(message)
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
		case input := <-ichan: // input
			err := c.WriteMessage(websocket.TextMessage, model.NewTextMessage(username, input, room).Encode())
			if err != nil {
				log.Println("write:", err)
				return
			}
		case t := <-ticker.C: //heartbeat
			err := c.WriteMessage(websocket.PingMessage, []byte(t.String()))
			if err != nil {
				log.Println("heartbeat error:", err)
				return
			}
		case <-interrupt: //os interrupt
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

// send join room action to server and wait for response
func joinRoom(c *websocket.Conn, username, room string) error {
	c.WriteMessage(websocket.TextMessage, model.NewJoinRoomMessage(username, room).Encode())

	fmt.Println("wait for server response...")
	_, resp, err := c.ReadMessage()
	if err != nil {
		fmt.Println("err: ", err, resp)
		return err
	}
	msg, err := model.Decode(resp)
	if err != nil {
		fmt.Println("error decoding message", err)
	}
	if msg.Action == model.JoinRoomSuccessAction {
		handleReceiveMessage(resp)
		return nil
	}
	return fmt.Errorf("unexpcted response %v", msg)
}

// display default message on screen
func handleReceiveMessage(rawMessage []byte) {
	msg, err := model.Decode(rawMessage)
	if err != nil {
		fmt.Println("error decoding message", err)
	}
	fmt.Printf("From %s: %s\n", msg.Sender, msg.Message)
}
