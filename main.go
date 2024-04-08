package main

import (
	"flag"

	"github.com/johnrijoy/ludo-go/frontend/prompt"
	"github.com/johnrijoy/ludo-go/frontend/tui"
)

func main() {
	isTui := flag.Bool("tui", false, "Start in TUI mode")
	flag.Parse()

	if *isTui {
		tui.Run()
	} else {
		prompt.Run()
	}
}
