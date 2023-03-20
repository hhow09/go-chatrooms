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

func IsPubsubEnv() bool {
	return os.Getenv("REDIS_PUBSUB") == "true"
}
