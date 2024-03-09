package frontend

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ludo-go/app"
)

var vlcPlayer app.VlcPlayer

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

	case "showq":
		audList := vlcPlayer.GetQueue()
		qIndex := vlcPlayer.GetQueueIndex()

		for i, audio := range audList {
			indexMsg := strconv.Itoa(i + 1)
			if qIndex == i {
				indexMsg = "**" + indexMsg
			}
			fmt.Println(indexMsg, " - ", audio.Title)
		}

	case "curr":
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

	case "stop":
		vlcPlayer.StopPlayback()

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
		fmt.Println("Run into Error: {}", err)
		os.Exit(1)
	}
}

func safeTruncString(label string, max int) string {
	var result string
	max = max + 3
	if len(label) <= max {
		result = label
	} else {

	}

	return result
}
