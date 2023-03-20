package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/hhow09/go-chatrooms/client-app/groupchat"
	"github.com/hhow09/go-chatrooms/client-app/prompts"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	// recover from failed program
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("program error happened:", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		// the answers will be written to this struct
		fmt.Println("Welcome to chatrooms program!")
		answers := prompts.ChooseProgram()
		if answers.Program == prompts.GROUP_CHAT {
			groupchat.Run(answers.UserName)
		}
		if answers.Program == prompts.PRIVATE_CHAT {
			fmt.Println("WIP")
			continue
		}
	}

}
