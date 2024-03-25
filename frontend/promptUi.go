package frontend

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
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
	"8,remove,rem":     "remove song at specified index, default is last | remove <index>",
	"9,removeAll,reml": "remove all songs stating from at specified index, default is current+1 | removeAll <index>",
	"10,forward,f":     "forwads playback by 10s",
	"11,rewind,r":      "reqinds playback by 10s",
	"12,stop":          "resets the player",
	"13,checkApi":      "check the current piped api",
	"14,setApi":        "set new piped api | setApi <piped api>",
	"15,listApi":       "display all available instances",
	"16,randApi":       "randomly select an piped instance",
	"17,version":       "display application details",
	"18,quit":          "quit application",
}

func RunPrompt() {
	exitSig := false

	log.SetOutput(io.Discard)
	err := vlcPlayer.InitPlayer()
	checkErr(err)
	defer vlcPlayer.ClosePlayer()

	for !exitSig {
		command := StringPrompt(">>")

		exitSig = runCommand(command)
	}
}

func runCommand(command string) bool {
	exitSig := false
	command, arg := parseCommand(command)

	switch command {
	case "play", "add":
		appendPlay(arg)

	case "radio":
		radioPlay(arg)

	case "p", "pause", "resume":
		vlcPlayer.PauseResume()

	case "showq", "q":
		displayQueue()

	case "curr", "c":
		displayCurrentSong()

	case "skipn", "n":
		skipNext()

	case "skipb", "b":
		skipPrevious()

	case "skip":
		skipIndex(arg)

	case "remove", "rem":
		removeIndex(arg)

	case "removeAll", "reml":
		removeAllIndex(arg)

	case "forward", "f":
		audioForward(arg)

	case "rewind", "r":
		audioRewind(arg)

	case "stop":
		resetPlayer()

	case "setApi":
		modifyApi(arg)

	case "checkApi":
		fmt.Println("Piped Api: ", app.GetPipedApi())

	case "listApi":
		displayApiList()

	case "randApi":
		modifyApiRandom()

	case "version":
		displayVersion()

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

//////////////////////
// Helper functions //
//////////////////////

func displayVersion() {
	fmt.Println("Ludo version: ", app.Version)
	fmt.Println("Api: ", app.GetPipedApi())
	fmt.Println(app.Info())
}

func modifyApiRandom() {
	apiList, err := app.GetPipedInstanceList()
	if err != nil {
		fmt.Println("Error in fetching Instance list: ", err)
	} else {
		randIndex := rand.Intn(len(apiList))
		app.SetPipedApi(apiList[randIndex].ApiUrl)
	}
}

func displayApiList() {
	apiList, err := app.GetPipedInstanceList()
	if err != nil {
		fmt.Println("Error in fetching Instance list: ", err)
	} else {
		for i, inst := range apiList {
			fmt.Printf("%-2d - %s\n", i, inst)
		}
	}
}

func modifyApi(arg string) {
	if arg != "" {
		app.SetPipedApi(arg)
		fmt.Println("Api changed from ", app.GetOldPipedApi(), " to ", app.GetPipedApi())
	} else {
		fmt.Println("Run into Error:", "No arguments")
	}
}

func resetPlayer() {
	err := vlcPlayer.ResetPlayer()
	checkErr(err)
}

func audioRewind(arg string) {
	duration := defaultForwardRewind
	if arg != "" {
		var err error
		duration, err = strconv.Atoi(arg)
		checkErr(err)
	}
	err := vlcPlayer.RewindBySeconds(duration)
	checkErr(err)
}

func audioForward(arg string) {
	duration := defaultForwardRewind
	if arg != "" {
		var err error
		duration, err = strconv.Atoi(arg)
		checkErr(err)
	}
	err := vlcPlayer.ForwardBySeconds(duration)
	checkErr(err)
}

func removeAllIndex(arg string) {
	trackIndex := vlcPlayer.GetQueueIndex() + 1
	if arg != "" {
		var err error
		trackIndex, err = strconv.Atoi(arg)
		checkErr(err)
		trackIndex -= 1
	}

	err := vlcPlayer.RemoveAllAudioFromIndex(trackIndex)
	checkErr(err)
}

func removeIndex(arg string) {
	trackIndex := len(vlcPlayer.GetQueue()) - 1
	if arg != "" {
		var err error
		trackIndex, err = strconv.Atoi(arg)
		checkErr(err)
		trackIndex -= 1
	}

	err := vlcPlayer.RemoveAudioFromIndex(trackIndex)
	checkErr(err)
}

func skipIndex(arg string) {
	trackIndex := vlcPlayer.GetQueueIndex() + 1
	if arg != "" {
		var err error
		trackIndex, err = strconv.Atoi(arg)
		checkErr(err)
	}

	err := vlcPlayer.SkipToIndex(trackIndex)
	checkErr(err)
}

func skipPrevious() {
	err := vlcPlayer.SkipToPrevious()
	checkErr(err)
}

func skipNext() {
	err := vlcPlayer.SkipToNext()
	checkErr(err)
}

func displayCurrentSong() {
	audState := vlcPlayer.GetAudioState()
	currPos, totPos := (&audState).GetPositionDetails()

	scale := len(audState.String())

	if totPos > currPos {
		scaledCurrPos := int(math.Round((float64(currPos) / float64(totPos)) * float64(scale)))
		restPosition := scale - scaledCurrPos

		navMsg := strings.Repeat(">", scaledCurrPos) + strings.Repeat("-", restPosition)
		fmt.Println(navMsg)
	}

	fmt.Println(&audState)
}

func displayQueue() {
	audList := vlcPlayer.GetQueue()
	qIndex := vlcPlayer.GetQueueIndex()

	for i, audio := range audList {
		indexMsg := strconv.Itoa(i + 1)
		if qIndex == i {
			indexMsg = "**" + indexMsg
		}
		fmt.Println(indexMsg, " - ", safeTruncString(audio.Title, 50), ", ", safeTruncString(audio.Uploader, 50))
	}
}

func radioPlay(arg string) {
	if arg != "" {
		audio, err := app.GetYtSong(arg)
		checkErr(err)

		err = vlcPlayer.ResetPlayer()
		checkErr(err)

		vlcPlayer.AppendAudio(audio)

		go func() {
			audioList, err := app.GetYtRadioList(audio.YtId, 10, true, true)
			checkErr(err)

			for _, audio := range *audioList {
				vlcPlayer.AppendAudio(&audio)
			}
		}()

	} else {
		fmt.Println("Run into Error:", "No arguments")
	}
	vlcPlayer.StartPlayback()
}

func appendPlay(arg string) {
	if arg != "" {
		audio, err := app.GetSong(arg, false)
		checkErr(err)

		vlcPlayer.AppendAudio(audio)
	}
	vlcPlayer.StartPlayback()
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

func parseCommand(command string) (string, string) {
	cmdList := strings.Split(command, " ")
	command = cmdList[0]
	arg := ""
	if len(cmdList) > 1 {
		arg = strings.Join(cmdList[1:], " ")
	}

	return command, arg
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Run into Error: ", err)
		vlcPlayer.ClosePlayer()
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
