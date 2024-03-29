package input

import (
	"fmt"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
)

// restore the cursor position and clear the line
func ClearLine() {
	fmt.Print("\033[u\033[K")
}

func (i *Input) BufLen() int {
	return len(i.buf)
}

type Input struct {
	buf      []rune
	resultCh chan string
}

func NewInput(interrupt chan os.Signal, exitKeys []keyboard.Key) (*Input, chan string, chan interface{}) {
	resultCh := make(chan string)
	i := &Input{
		resultCh: resultCh,
	}
	doneCh := make(chan interface{})

	go i.readStdin(interrupt, exitKeys, doneCh)
	return i, resultCh, doneCh
}

func (i *Input) readStdin(interrupt chan os.Signal, exitKeys []keyboard.Key, doneCh chan interface{}) {
	for {
		select {
		case <-doneCh:
			return
		default:
			char, key, err := keyboard.GetSingleKey()
			if err != nil {
				fmt.Println(err)
				return
			}

			// exit conditions
			for _, ekeys := range exitKeys {
				if key == ekeys {
					fmt.Println("exit")
					interrupt <- os.Interrupt
					return
				}
			}
			switch key {
			case keyboard.KeyEnter:
				ClearLine()
				i.resultCh <- strings.TrimSpace(string(i.buf))
				i.resetBuffer()
			case keyboard.KeyBackspace2, keyboard.KeyBackspace:
				if len(i.buf) > 0 {
					i.buf = i.buf[:len(i.buf)-1]
				}
			case keyboard.KeySpace:
				i.buf = append(i.buf, rune(' '))
			default:
				i.buf = append(i.buf, char)
			}
			// print out
			ClearLine()
			fmt.Printf("%v", string(i.buf))
		}
	}
}

func (i *Input) resetBuffer() {
	i.buf = []rune{}
}

func (i *Input) ResumeBuffer() {
	fmt.Printf("%v", string(i.buf))
}
