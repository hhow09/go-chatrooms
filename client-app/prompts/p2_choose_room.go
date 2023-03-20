package prompts

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

const OPTION_CREATE_ROOM = "create new room"

// the questions to ask
func getQ1(roomList []string) []*survey.Question {
	opts := append(roomList, OPTION_CREATE_ROOM)
	return []*survey.Question{{
		Name: "room",
		Prompt: &survey.Select{
			Message: "Choose a room:",
			Options: opts,
		},
	},
	}
}

type Ans2 struct {
	Room string `survey:"room"`
}

func ChooseRoom(roomList []string) Ans2 {
	if len(roomList) == 0 {
		return Ans2{}
	}
	answers := Ans2{}
	err := survey.Ask(getQ1(roomList), &answers)
	if err != nil {
		fmt.Println(err)
		return Ans2{}
	}
	return answers
}
