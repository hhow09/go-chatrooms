package groupchat

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

const OPTION_CREATE_ROOM = "create new room"

// the questions to ask
func getQ1(roomList []string) []*survey.Question {
	opts := append(roomList, OPTION_CREATE_ROOM)
	fmt.Println(opts)
	return []*survey.Question{{
		Name: "room",
		Prompt: &survey.Select{
			Message: "Choose a room:",
			Options: opts,
		},
	},
	}
}

type Ans1 struct {
	Room string `survey:"room"`
}

func ChooseRoomPrompt(roomList []string) Ans1 {
	if len(roomList) == 0 {
		return Ans1{}
	}
	answers := Ans1{}
	err := survey.Ask(getQ1(roomList), &answers)
	if err != nil {
		fmt.Println(err)
		return Ans1{}
	}
	return answers
}

// the questions to ask
func getQ2() []*survey.Question {
	return []*survey.Question{{
		Name: "room",
		Prompt: &survey.Input{
			Message: "Create a room:",
		},
		Validate: func(val interface{}) error {
			if str, ok := val.(string); !ok || str == "" {
				return errors.New("room name is required")
			}
			if val.(string) == OPTION_CREATE_ROOM {
				return errors.New("room name is restricted")
			}

			return nil
		},
	},
	}
}

type Ans2 struct {
	Room string `survey:"room"`
}

func CreateRoomPrompt() Ans2 {
	answers := Ans2{}
	err := survey.Ask(getQ2(), &answers)
	if err != nil {
		fmt.Println(err)
		return Ans2{}
	}
	return answers
}
