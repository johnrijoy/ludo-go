package main

import (
	"fmt"
	"log"
)

func main() {
	log.Println("Hello Start")
	// target := musicSearch("Hello")
	// playMusicVlc(target)

	audio, err := getSong("Hello Adele")
	checkErr(err)

	log.Printf("%s\n", audio)
	log.Printf("%s\n", audio.uploader)
	log.Printf("%d\n", audio.duration)
	//fmt.Println("{}", audio.audioStreamUrl)

	var vlcPlayer VlcPlayer
	err = vlcPlayer.init()
	checkErr(err)
	defer vlcPlayer.close()

	vlcPlayer.appendSong(audio)
	vlcPlayer.startPlayback()

	var inp string
	fmt.Scanln(&inp)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
