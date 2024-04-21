package main

import (
	"flag"

	"github.com/johnrijoy/ludo-go/frontend/prompt"
	"github.com/johnrijoy/ludo-go/frontend/tui"
)

func main() {
	isPrompt := flag.Bool("p", false, "Start in prompt mode")
	flag.Parse()

	if *isPrompt {
		prompt.Run()
	} else {
		tui.Run()
	}
}
