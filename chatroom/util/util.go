package util

import (
	"fmt"
	"os"
)

func Log(s ...string) {
	if os.Getenv("DEBUG") == "true" {
		fmt.Println("[Log]", s)
	}
}
