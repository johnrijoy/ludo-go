package main

import (
	"log"
)

func getPipedApi() string {
	return "https://piapi.ggtyler.dev/"
}

func main() {
	playMusicVlc2("3 sec")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
