package main

import (
	"log"

	vlc "github.com/adrg/libvlc-go/v3"
	"github.com/johnrijoy/ludo-go/app"
)

func playMusicVlc(target string) {
	err := vlc.Init("--no-video", "--quiet")
	checkErr(err)
	defer vlc.Release()

	player, err := vlc.NewPlayer()
	checkErr(err)

	defer func() {
		player.Stop()
		player.Release()
	}()

	// target := "https://pipedproxy-bom-2.kavin.rocks/videoplayback?expire=1685488136&ei=qC12ZLO0IOj74-EPjP68wAM&ip=140.238.251.167&id=o-ADbyf3kMeJOr9Bodc9Jcv1w436P8kB-jwe2QSDL-3SjG&itag=139&source=youtube&requiressl=yes&mh=Nv&mm=31%2C29&mn=sn-cvh76nl7%2Csn-cvh7knzr&ms=au%2Crdu&mv=m&mvi=5&pl=26&gcr=in&initcwndbps=22951250&spc=qEK7BzlfF7rP8VHIk3tp_IU1BO2rL8M&vprv=1&svpuc=1&mime=audio%2Fmp4&gir=yes&clen=1803991&dur=295.557&lmt=1680796357291454&mt=1685466053&fvip=2&keepalive=yes&fexp=24007246&beids=24350018&c=ANDROID&txp=4532434&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cgcr%2Cspc%2Cvprv%2Csvpuc%2Cmime%2Cgir%2Cclen%2Cdur%2Clmt&sig=AOq0QJ8wRQIgaUyEbjr3Qwti0HCW92-IKy6X7p5JuZDuf-KFVXkeBhYCIQDG37aW-NZMarB3Ie7_qrGfb6bkNakVomTzVhS2a8fbXA%3D%3D&lsparams=mh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Cinitcwndbps&lsig=AG3C_xAwRQIgFpvKT6kl3a7eCAh5SIDdfmPDKzet9rRTF2LmRHaS3zACIQDWM-EUfeX9UepZDRqbqVcFtveNdQaaNGZ1wZy-S5XHPQ%3D%3D&cpn=kTNjgr_Lw9FztwiC&host=rr5---sn-cvh76nl7.googlevideo.com"

	media, err := player.LoadMediaFromURL(target)
	checkErr(err)
	defer media.Release()

	manager, err := player.EventManager()
	checkErr(err)

	quit := make(chan struct{})
	eventCallback := func(event vlc.Event, userData interface{}) {
		close(quit)
	}

	eventCallback2 := func(event vlc.Event, userData interface{}) {
		log.Println("Position Changed...")
	}

	eventID, err := manager.Attach(vlc.MediaPlayerEndReached, eventCallback, nil)
	checkErr(err)
	defer manager.Detach(eventID)

	eventID2, err := manager.Attach(vlc.MediaPlayerTimeChanged, eventCallback2, nil)
	checkErr(err)
	defer manager.Detach(eventID2)

	err = player.Play()
	checkErr(err)
	<-quit
}

func playMusicVlc2(songName string) {
	log.Println("Hello Start")
	// target := musicSearch("Hello")
	// playMusicVlc(target)

	audio, err := app.GetSong(songName)
	checkErr(err)

	log.Printf("%s\n", audio)
	log.Printf("%s\n", audio.Uploader)
	log.Printf("%d\n", audio.Duration)
	//fmt.Println("{}", audio.audioStreamUrl)

	var vlcPlayer app.VlcPlayer
	err = vlcPlayer.Init()
	checkErr(err)
	defer vlcPlayer.Close()

	vlcPlayer.AppendSong(audio)
	vlcPlayer.StartPlayback()

	<-vlcPlayer.Quit
	log.Println("Control reached back")
}

func playMusicVlc3(songList ...string) {
	var vlcPlayer app.VlcPlayer
	err := vlcPlayer.Init()
	checkErr(err)
	defer vlcPlayer.Close()

	log.Println("Before mediastate: {}")

	for _, songName := range songList {
		log.Println(songName)
		// target := musicSearch("Hello")
		// playMusicVlc(target)

		audio, err := app.GetSong(songName)
		checkErr(err)

		log.Printf("%s\n", audio)
		log.Printf("%s\n", audio.Uploader)
		log.Printf("%d\n", audio.Duration)
		//fmt.Println("{}", audio.audioStreamUrl)

		vlcPlayer.AppendSong(audio)
	}

	vlcPlayer.StartPlayback()

	<-vlcPlayer.Quit
	log.Println("Control reached back")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
