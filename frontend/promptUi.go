package frontend

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/johnrijoy/ludo-go/app"
)

var vlcPlayer app.VlcPlayer

const defaultForwardRewind = 10

var commands = map[string]string{
	"0,play,add":       "play the song | play <song name>",
	"1,radio":          "start radio for song | radio <song name>",
	"2,pause,resume,p": "toggle pause/resume",
	"3,showq,q":        "display song queue",
	"4,curr,c":         "display current song",
	"5,skipn,n":        "skip to next song",
	"6,skipb,b":        "skip to previous song",
	"7,skip":           "skip to the specified index, default is 1 | skip <index>",
	"8,forward,f":      "forwads playback by 10s",
	"9,rewind,r":       "reqinds playback by 10s",
	"10,stop":          "resets the player",
	"11,checkApi":      "check the current piped api",
	"12,setApi":        "set new piped api | setApi <piped api>",
	"13,version":       "display application details",
	"14,quit":          "quit application",
}

func RunPrompt() {
	exitSig := false

	log.SetOutput(io.Discard)
	err := vlcPlayer.Init()
	checkErr(err)
	defer vlcPlayer.Close()

	for !exitSig {
		command := StringPrompt(">>")

		exitSig = runCommand(command)
	}
}

func runCommand(command string) bool {
	exitSig := false
	cmdList := strings.Split(command, " ")
	command = cmdList[0]
	arg := ""
	if len(cmdList) > 1 {
		arg = strings.Join(cmdList[1:], " ")
	}

	switch command {
	case "play", "add":
		if arg != "" {
			audio, err := app.GetSong(arg, false)
			checkErr(err)

			vlcPlayer.AppendSong(audio)
		}
		vlcPlayer.StartPlayback()

	case "radio":
		if arg != "" {
			audio, err := app.GetSong(arg, false)
			checkErr(err)

			err = vlcPlayer.ResetPlayer()
			checkErr(err)

			vlcPlayer.AppendSong(audio)

			go func() {
				audioList, err := app.GetYtRadioList(arg, 10, true)
				checkErr(err)

				for _, audio := range *audioList {
					vlcPlayer.AppendSong(&audio)
				}
			}()

		} else {
			fmt.Println("Run into Error:", "No arguments")
		}
		vlcPlayer.StartPlayback()

	case "p", "pause", "resume":
		vlcPlayer.PauseResume()

	case "showq", "q":
		audList := vlcPlayer.GetQueue()
		qIndex := vlcPlayer.GetQueueIndex()

		for i, audio := range audList {
			indexMsg := strconv.Itoa(i + 1)
			if qIndex == i {
				indexMsg = "**" + indexMsg
			}
			fmt.Println(indexMsg, " - ", safeTruncString(audio.Title, 50), ", ", safeTruncString(audio.Uploader, 50))
		}

	case "curr", "c":
		const scale = 10
		audState := vlcPlayer.GetAudioState()
		currPos, totPos := (&audState).GetPositionDetails()

		if totPos > currPos {
			scaledCurrPos := int(math.Round((float64(currPos) / float64(totPos)) * scale))
			restPosition := scale - scaledCurrPos

			navMsg := strings.Repeat(">", scaledCurrPos) + strings.Repeat("-", restPosition)
			fmt.Println(navMsg)
		}

		fmt.Println(&audState)

	case "skipn", "n":
		err := vlcPlayer.SkipToNext()
		checkErr(err)

	case "skipb", "b":
		err := vlcPlayer.SkipToPrevious()
		checkErr(err)

	case "skip":
		trackIndex := vlcPlayer.GetQueueIndex() + 1
		if arg != "" {
			var err error
			trackIndex, err = strconv.Atoi(arg)
			checkErr(err)
		}

		err := vlcPlayer.SkipToIndex(trackIndex)
		checkErr(err)

	case "forward", "f":
		duration := defaultForwardRewind
		if arg != "" {
			var err error
			duration, err = strconv.Atoi(arg)
			checkErr(err)
		}
		err := vlcPlayer.ForwardBySeconds(duration)
		checkErr(err)

	case "rewind", "r":
		duration := defaultForwardRewind
		if arg != "" {
			var err error
			duration, err = strconv.Atoi(arg)
			checkErr(err)
		}
		err := vlcPlayer.RewindBySeconds(duration)
		checkErr(err)

	case "stop":
		err := vlcPlayer.ResetPlayer()
		checkErr(err)

	case "setApi":
		if arg != "" {
			app.SetPipedApi(arg)
			fmt.Println("Api changed from ", app.GetOldPipedApi(), " to ", app.GetPipedApi())
		} else {
			fmt.Println("Run into Error:", "No arguments")
		}

	case "checkApi":
		fmt.Println("Piped Api: ", app.GetPipedApi())

	case "version":
		fmt.Println("Ludo version: ", app.Version)
		fmt.Println("Api: ", app.GetPipedApi())
		fmt.Println(app.Info())

	case "help":
		showHelp()

	case "quit":
		fmt.Println("Exiting player...")
		exitSig = true

	default:
		fmt.Println("Invalid command")
	}

	return exitSig
}

func showHelp() {
	var keys []string
	for key := range commands {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		cmd := strings.Join(strings.Split(key, ",")[1:], ", ")
		//info := strings.Split(commands[key], " | ")
		info := commands[key]
		fmt.Printf("%-20s - %s\n", cmd, info)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Run into Error: ", err)
		vlcPlayer.Close()
		os.Exit(1)
	}
}

func safeTruncString(label string, max int) string {
	var result string

	if len(label) <= max {
		result = label
	} else {
		result = label[0:(max-3)] + "..."
	}

	return result
}
