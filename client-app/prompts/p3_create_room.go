package prompts

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// the questions to ask
func Q2() []*survey.Question {
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

type Ans3 struct {
	Room string `survey:"room"`
}

func CreateRoom() Ans3 {
	answers := Ans3{}
	err := survey.Ask(Q2(), &answers)
	if err != nil {
		fmt.Println(err)
		return Ans3{}
	}
	return answers
}
