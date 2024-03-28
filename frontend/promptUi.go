package frontend

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/johnrijoy/ludo-go/app"
)

var vlcPlayer app.VlcPlayer

const defaultForwardRewind = 10

var commands = []string{
	"play,add-play the song | play <song name>",
	"search,s-search the song and display search result | search <song name>",
	"radio-start radio for song | radio <song name>",
	"pause,resume,p-toggle pause/resume",
	"showq,q-display song queue",
	"curr,c-display current song",
	"skipn,n-skip to next song",
	"skipb,b-skip to previous song",
	"skip-skip to the specified index, default is 1 | skip <index>",
	"remove,rem-remove song at specified index, default is last | remove <index>",
	"removeAll,reml-remove all songs stating from at specified index, default is current+1 | removeAll <index>",
	"forward,f-forwads playback by 10s **",
	"rewind,r-reqinds playback by 10s **",
	"stop-resets the player",
	"checkApi-check the current piped api",
	"setApi-set new piped api | setApi <piped api>",
	"listApi-display all available instances",
	"randApi-randomly select an piped instance",
	"version-display application details",
	"quit-quit application",
}

func RunPrompt() {
	exitSig := false

	log.SetOutput(io.Discard)
	err := vlcPlayer.InitPlayer()
	checkErr(err)
	defer vlcPlayer.ClosePlayer()

	showStartupMessage()

	for !exitSig {
		command := StringPrompt(">>")

		exitSig = runCommand(command)
	}

	fmt.Println("Exiting player...")
}

func runCommand(command string) bool {
	exitSig := false
	command, arg := parseCommand(command)

	switch command {
	case "play", "add":
		appendPlay(arg)

	case "search", "s":
		searchPlay(arg)

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
		exitSig = true

	default:
		fmt.Println("Invalid command")
	}

	return exitSig
}

// Commands //

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

	scale := len((&audState).String())

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
		audio, err := app.GetYtSong(arg, false)
		checkErr(err)

		err = vlcPlayer.ResetPlayer()
		checkErr(err)

		vlcPlayer.AppendAudio(audio)

		go func() {
			audioList, err := app.GetYtRadioList(audio.YtId, true, 1, 10)
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
		audio, err := app.GetPipedSong(arg, false)
		checkErr(err)

		vlcPlayer.AppendAudio(audio)
	}
	vlcPlayer.StartPlayback()
}

func searchPlay(arg string) {
	if arg == "" {
		fmt.Println("Please enter a search query")
		fmt.Println()
		return
	}

	audioBasicList, err := app.SearchPipedSong(arg, 0, 10)
	checkErr(err)

	for i, audio := range *audioBasicList {
		fmt.Printf("%-2d - %-50s | %-20s | %s\n", i+1, safeTruncString(audio.Title, 50), safeTruncString(audio.Uploader, 20), audio.GetFormattedDuration())
	}

	fmt.Println("Enter index number (q to escape): ")
	cmd := StringPrompt(":> ")

	if cmd == "q" {
		fmt.Println("exiting search...")
		return
	}
	index, err := strconv.Atoi(cmd)
	if err != nil {
		fmt.Println("Error: Please enter a number to play song or q to exit search")
		return
	}

	index -= 1

	if index < 0 || index >= len(*audioBasicList) {
		fmt.Println("Error: Please enter a valid number")
		return
	}

	audioBasic := (*audioBasicList)[index]

	audio, err := app.GetPipedSong(audioBasic.YtId, true)
	if displayErr(err) {
		return
	}

	err = vlcPlayer.AppendAudio(audio)
	displayErr(err)

	if !vlcPlayer.IsPlaying() {
		vlcPlayer.StartPlayback()
	}
}

func showStartupMessage() {
	fmt.Println("=== LUDO GO ===")
	fmt.Println("Welcome to LudoGo")
	fmt.Println("To start listening, enter \"play <song name>\"")
	fmt.Println("To show help, enter \"help\"")
}

func showHelp() {
	for _, val := range commands {
		splitVal := strings.Split(val, "-")
		cmd, helpMsg := splitVal[0], splitVal[1]
		cmd = strings.ReplaceAll(cmd, ",", ", ")

		fmt.Printf("%-20s - %s\n", cmd, helpMsg)
	}
}

//////////////////////
// Helper functions //
//////////////////////

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

func displayErr(err error) bool {
	if err != nil {
		fmt.Println("Run into Error: ", err)
		return true
	}

	return false
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
