package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/mazzegi/mint"
)

func main() {
	prompt := func() { fmt.Print("> ") }
	state := mint.NewState()
	scr := bufio.NewScanner(os.Stdin)
	prompt()
	for scr.Scan() {
		txt := scr.Text()
		switch txt {
		case "":
		case "exit", "q", "quit", "bye":
			os.Exit(0)
		default:
			t0 := time.Now()
			res, err := state.Eval(txt)
			if err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println(fmt.Sprintf("<< (%s)", time.Since(t0)), res)
			}
		}
		prompt()
	}
}
