package prompts

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

const (
	GROUP_CHAT   string = "group chat"
	PRIVATE_CHAT string = "private chat"
)

var q1 = []*survey.Question{
	{
		Name: "program",
		Prompt: &survey.Select{
			Message: "Choose a program:",
			Options: []string{GROUP_CHAT, PRIVATE_CHAT},
			Default: GROUP_CHAT,
		},
	},
	{
		Name:     "username",
		Prompt:   &survey.Input{Message: "please input username"},
		Validate: survey.Required,
		Transform: survey.TransformString(func(s string) string {
			transformed := strings.Replace(s, " ", "-", -1) // replace whitespace with dash
			return transformed
		}),
	},
}

type Ans1 struct {
	Program  string `survey:"program"`
	UserName string `survey:"username"`
}

func ChooseProgram() Ans1 {
	answers := Ans1{}
	err := survey.Ask(q1, &answers)
	if err != nil {
		panic(err)
	}
	return answers
}
