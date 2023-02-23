package input

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/eiannone/keyboard"
)

type Input struct {
	lock     sync.Mutex
	buf      []rune
	resultCh chan string
}

func NewInput(interrupt chan os.Signal, exitKeys []keyboard.Key) chan string {
	resultCh := make(chan string)
	i := &Input{
		resultCh: resultCh,
	}

	go i.readStdin(interrupt, exitKeys)
	return resultCh
}

func (i *Input) ResetBuffer() {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.resetBuffer()
}

func (i *Input) readStdin(interrupt chan os.Signal, exitKeys []keyboard.Key) {
	for {
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}

		i.lock.Lock()
		// exit conditions
		for _, ekeys := range exitKeys {
			if key == ekeys {
				fmt.Println("exit")
				interrupt <- os.Interrupt
				return
			}
		}
		fmt.Printf("%v", string(char)) // print out
		if char == 0 {                 // enter key
			fmt.Println("") // start new line
			i.resultCh <- strings.TrimSpace(string(i.buf))
			i.resetBuffer()
		} else {
			i.buf = append(i.buf, char)
		}
		i.lock.Unlock()
	}
}

func (i *Input) resetBuffer() {
	i.buf = []rune{}
}
