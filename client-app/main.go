package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/hhow09/go-chatrooms/client-app/groupchat"
	"github.com/joho/godotenv"

	"github.com/AlecAivazis/survey/v2"
)

func init() {
	godotenv.Load()
}

const (
	GROUP_CHAT   string = "group chat"
	PRIVATE_CHAT string = "private chat"
)

// the questions to ask
var qs = []*survey.Question{
	{
		Name:     "username",
		Prompt:   &survey.Input{Message: "please input username"},
		Validate: survey.Required,
	},
	{
		Name: "program",
		Prompt: &survey.Select{
			Message: "Choose a program:",
			Options: []string{GROUP_CHAT, PRIVATE_CHAT},
			Default: GROUP_CHAT,
		},
	},
}

type Ans struct {
	Program  string `survey:"program"`
	UserName string `survey:"username"`
}

func main() {
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
		answers := Ans{}
		err := survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err)
			return
		}
		if answers.Program == GROUP_CHAT {
			groupchat.GroupChatProgram(answers.UserName)
		}
		if answers.Program == PRIVATE_CHAT {
			fmt.Println("WIP")
			continue
		}
	}

}
