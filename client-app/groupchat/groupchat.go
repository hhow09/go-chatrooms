package groupchat

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gorilla/websocket"
	"github.com/hhow09/go-chatrooms/client-app/api"
	"github.com/hhow09/go-chatrooms/client-app/input"
	"github.com/hhow09/go-chatrooms/client-app/model"
	"github.com/hhow09/go-chatrooms/client-app/prompts"
)

type rtn_code int

const (
	ERROR_EXIT rtn_code = 1
	DONE_EXIT  rtn_code = 0
)

// reconnect
func Run(username string) {
	code := ERROR_EXIT
	for code != DONE_EXIT {
		code = GroupChatProgram(username)
		if code == DONE_EXIT {
			return
		}
		fmt.Println("restarting...")
	}
}

func GroupChatProgram(username string) rtn_code {
	// init keyboard reader
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		keyboard.Close()
	}()

	flag.Parse()
	log.SetFlags(0)

	// exit sig
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	wc := api.NewWSClient(username)
	defer wc.Close()

	done := make(chan rtn_code)
	defer close(done)
	room := chooseRoom()
	// join room
	resp, err := joinRoom(wc.Conn, username, room)
	if err != nil {
		panic(fmt.Sprintf("join room failed %v", err))
	}
	// input reader
	inputi, ichan, doneCh := input.NewInput(interrupt, []keyboard.Key{keyboard.KeyEsc, keyboard.KeyCtrlC})
	defer func() {
		// prevent input from blocking infinitely
		select {
		case doneCh <- struct{}{}:
		case <-time.After(time.Second):
		}
		close(doneCh)
	}()
	handleReceiveMessage(inputi, resp)

	ticker := wc.HeartbeatSetup()
	defer ticker.Stop()

	go readPump(wc.Conn, done, inputi)
	for {
		select {
		case code := <-done:
			return code
		case input := <-ichan: // input
			err := wc.Conn.WriteMessage(websocket.TextMessage, model.NewTextMessage(username, input, room).Encode())
			if err != nil {
				log.Println("write:", err)
				return ERROR_EXIT
			}
		case t := <-ticker.C: //heartbeat
			err := wc.Conn.WriteMessage(websocket.PingMessage, []byte(t.String()))
			if err != nil {
				log.Println("heartbeat error:", err)
				return ERROR_EXIT
			}
		case <-interrupt: //os interrupt
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := wc.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return DONE_EXIT
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return DONE_EXIT
		}
	}
}

// receive message from ws
func readPump(c *websocket.Conn, done chan rtn_code, inputi *input.Input) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("unexpected close error: ", err)
				done <- ERROR_EXIT
				return
			} else if err, ok := err.(net.Error); ok && err.Timeout() {
				// handle timeout error
				// set by c.SetReadDeadline in heartbeatSetup
				fmt.Println("timeout error:", err)
				done <- ERROR_EXIT
				return
			}
			fmt.Println("err:", err)
			done <- ERROR_EXIT
			return
		}
		handleReceiveMessage(inputi, message)
	}
}

// send join room action to server and wait for response
func joinRoom(c *websocket.Conn, username, room string) ([]byte, error) {
	c.WriteMessage(websocket.TextMessage, model.NewJoinRoomMessage(username, room).Encode())

	fmt.Println("wait for server response...")
	_, resp, err := c.ReadMessage()
	if err != nil {
		fmt.Println("err: ", err, resp)
		return nil, err
	}
	msg, err := model.Decode(resp)
	if err != nil {
		fmt.Println("error decoding message", err)
	}
	if msg.Action == model.JoinRoomSuccessAction {
		return resp, nil
	}
	return nil, fmt.Errorf("unexpcted response %v", msg)
}

// display default message on screen
func handleReceiveMessage(inputi *input.Input, rawMessage []byte) {
	msg, err := model.Decode(rawMessage)
	if err != nil {
		fmt.Println("error decoding message", err)
	}
	if inputi.BufLen() > 0 {
		input.ClearLine()
	}
	fmt.Printf("From %s: %s\n", msg.Sender, msg.Message)
	inputi.ResumeBuffer()
}

// choose room or create room with prompt
func chooseRoom() (room string) {
	roomList, _ := api.GetRoomList()
	ans1 := prompts.ChooseRoom(roomList)
	if ans1.Room == prompts.OPTION_CREATE_ROOM || ans1.Room == "" {
		room = prompts.CreateRoom().Room
	} else {
		room = ans1.Room
	}
	return room
}
