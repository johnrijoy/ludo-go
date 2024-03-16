package frontend

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/johnrijoy/ludo-go/app"
)

var vlcPlayer app.VlcPlayer

const defaultForwardRewind = 10

func RunPrompt() {
	exitSig := false

	log.SetOutput(io.Discard)
	err := vlcPlayer.Init()
	checkErr(err)

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
			audio, err := app.GetSong(arg)
			checkErr(err)

			vlcPlayer.AppendSong(audio)
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
		vlcPlayer.StopPlayback()

	case "setApi":
		if arg != "" {
			app.SetPipedApi(arg)
			fmt.Println("Api changed from ", app.GetOldPipedApi(), " to ", app.GetPipedApi())
		} else {
			fmt.Println("Run into Error:", "No arguments")
		}

	case "checkApi":
		fmt.Println("Api: ", app.GetPipedApi())

	case "quit":
		fmt.Println("Exiting player...")
		exitSig = true

	default:
		fmt.Println("Invalid command")
	}

	return exitSig
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
