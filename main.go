package main

import (
	"github.com/johnrijoy/ludo-go/frontend/prompt"
	"github.com/johnrijoy/ludo-go/frontend/rest"
)

func main() {
	go rest.Run()
	prompt.Run()
}
