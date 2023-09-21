package main

import (
	"log"
)

func getPipedApi() string {
	return "https://piapi.ggtyler.dev/"
}

func main() {
	playMusicVlc2("Hello Adele")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
