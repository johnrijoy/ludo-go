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
	handleErrExit(err)
	defer vlcPlayer.ClosePlayer()

	showStartupMessage()

	for !exitSig {
		command := StringPrompt(">>")

		exitSig = runCommand(command)
	}

	silentLog("Exiting player...")
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
		fmt.Println("Piped Api: ", Green(app.GetPipedApi()))

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
		warnLog("Invalid command")
	}

	return exitSig
}

// Commands //

func displayVersion() {
	fmt.Println("Ludo version: ", Green(app.Version))
	fmt.Println("Api: ", Green(app.GetPipedApi()))
	fmt.Println("libVlc Binding Version: ", Green(app.Info().String()))
	fmt.Println("Vlc Runtime Version: ", Green(app.Info().Changeset()))
}

func modifyApiRandom() {
	apiList, err := app.GetPipedInstanceList()
	if err != nil {
		errorLog("Error in fetching Instance list:", err)
	} else {
		randIndex := rand.Intn(len(apiList))
		app.SetPipedApi(apiList[randIndex].ApiUrl)
	}
}

func displayApiList() {
	apiList, err := app.GetPipedInstanceList()
	if err != nil {
		errorLog("Error in fetching Instance list:", err)
		return
	}

	for i, inst := range apiList {
		fmt.Printf("%-2d - %s\n", i+1, inst)
	}

	cmd := StringPrompt("> Enter index number to change api (q to escape): ")

	if cmd == "q" {
		silentLog("exiting api list...")
		return
	}

	index, err := strconv.Atoi(cmd)
	if err != nil {
		errorLog("please enter a number to play song or q to exit search")
		return
	}

	index -= 1

	if index < 0 || index >= len(apiList) {
		errorLog("Please enter a valid number")
		return
	}

	newApi := apiList[index].ApiUrl
	app.SetPipedApi(newApi)
}

func modifyApi(arg string) {
	if arg != "" {
		app.SetPipedApi(arg)
		fmt.Println("Api changed from ", Gray(app.GetOldPipedApi()), " to ", Green(app.GetPipedApi()))
	} else {
		errorLog("no arguments")
	}
}

func resetPlayer() {
	err := vlcPlayer.ResetPlayer()
	handleErrExit(err)
}

func audioRewind(arg string) {
	duration := defaultForwardRewind
	if arg != "" {
		var err error
		duration, err = strconv.Atoi(arg)
		handleErrExit(err)
	}
	err := vlcPlayer.RewindBySeconds(duration)
	handleErrExit(err)
}

func audioForward(arg string) {
	duration := defaultForwardRewind
	if arg != "" {
		var err error
		duration, err = strconv.Atoi(arg)
		handleErrExit(err)
	}
	err := vlcPlayer.ForwardBySeconds(duration)
	handleErrExit(err)
}

func removeAllIndex(arg string) {
	trackIndex := vlcPlayer.GetQueueIndex() + 1
	if arg != "" {
		var err error
		trackIndex, err = strconv.Atoi(arg)
		handleErrExit(err)
		trackIndex -= 1
	}

	err := vlcPlayer.RemoveAllAudioFromIndex(trackIndex)
	handleErrExit(err)
}

func removeIndex(arg string) {
	trackIndex := len(vlcPlayer.GetQueue()) - 1
	if arg != "" {
		var err error
		trackIndex, err = strconv.Atoi(arg)
		handleErrExit(err)
		trackIndex -= 1
	}

	err := vlcPlayer.RemoveAudioFromIndex(trackIndex)
	handleErrExit(err)
}

func skipIndex(arg string) {
	trackIndex := vlcPlayer.GetQueueIndex() + 1
	if arg != "" {
		var err error
		trackIndex, err = strconv.Atoi(arg)
		if displayErr(err) {
			return
		}
	}

	err := vlcPlayer.SkipToIndex(trackIndex)
	displayErr(err)
}

func skipPrevious() {
	err := vlcPlayer.SkipToPrevious()
	handleErrExit(err)
}

func skipNext() {
	err := vlcPlayer.SkipToNext()
	handleErrExit(err)
}

func displayCurrentSong() {
	if vlcPlayer.IsPlaying() {
		fmt.Println(Green("Now playing..."))
	} else if vlcPlayer.CheckMediaError() {
		fmt.Println(Red("Error in playing media"))
	} else {
		fmt.Println(Yellow(vlcPlayer.FetchPlayerState()))
	}

	audState := vlcPlayer.GetAudioState()
	currPos, totPos := (&audState).GetPositionDetails()

	scale := len((&audState).String())

	if totPos > currPos {
		scaledCurrPos := int(math.Round((float64(currPos) / float64(totPos)) * float64(scale)))
		restPosition := scale - scaledCurrPos

		navMsg := Magenta(strings.Repeat(">", scaledCurrPos)) + Gray(strings.Repeat("-", restPosition))
		fmt.Println(navMsg)
	}

	fmt.Println(&audState)
}

func displayQueue() {
	audList := vlcPlayer.GetQueue()
	qIndex := vlcPlayer.GetQueueIndex()

	for i, audio := range audList {
		msg := fmt.Sprintf("%-2d - %-50s | %-50s", i+1, safeTruncString(audio.Title, 50), safeTruncString(audio.Uploader, 50))
		if qIndex == i {
			msg = Magenta(msg)
		}
		fmt.Println(msg)
	}
}

func radioPlay(arg string) {
	if arg == "" {
		warnLog("please enter a search query")
		return
	}

	var audio *app.AudioDetails
	if arg == "." {
		audioD := vlcPlayer.GetAudioState().AudioDetails
		audio = &audioD
		removeAllIndex("")
	} else {
		var err error
		audio, err = app.GetYtSong(arg, false)
		handleErrExit(err)

		err = vlcPlayer.ResetPlayer()
		handleErrExit(err)

		vlcPlayer.AppendAudio(audio)
		vlcPlayer.StartPlayback()
	}

	go func() {
		audioList, err := app.GetYtRadioList(audio.YtId, true, 1, 10)
		handleErrExit(err)

		for _, audio := range *audioList {
			vlcPlayer.AppendAudio(&audio)
		}
	}()
}

func appendPlay(arg string) {
	if arg != "" {
		audio, err := app.GetPipedSong(arg, false)
		handleErrExit(err)

		vlcPlayer.AppendAudio(audio)
	}
	if len(vlcPlayer.GetQueue()) < 1 {
		warnLog("no song in queue for playback")
		return
	}
	err := vlcPlayer.StartPlayback()
	handleErrExit(err)
}

func searchPlay(arg string) {
	if arg == "" {
		warnLog("please enter a search query")
		return
	}

	audioBasicList, err := app.SearchPipedSong(arg, 0, 10)
	handleErrExit(err)

	for i, audio := range *audioBasicList {
		fmt.Printf("%-2d - %-50s | %-20s | %s\n", i+1, safeTruncString(audio.Title, 50), safeTruncString(audio.Uploader, 20), audio.GetFormattedDuration())
	}

	cmd := StringPrompt("> Enter index number (q to escape): ")

	if cmd == "q" {
		silentLog("exiting search...")
		return
	}
	index, err := strconv.Atoi(cmd)
	if err != nil {
		errorLog("please enter a number to play song or q to exit search")
		return
	}

	index -= 1

	if index < 0 || index >= len(*audioBasicList) {
		errorLog("please enter a valid number")
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
	fmt.Println(Blue("==="), Magenta("LUDO GO"), Blue("==="))
	fmt.Println("Welcome to", Magenta("LudoGo"))
	fmt.Println("To start listening, enter " + Green("play <song name>"))
	fmt.Println("To show help, enter " + Green("help"))
}

func showHelp() {
	for _, val := range commands {
		splitVal := strings.Split(val, "-")
		cmd, helpMsg := splitVal[0], splitVal[1]
		cmd = strings.ReplaceAll(cmd, ",", ", ")

		fmt.Printf("%-40s - %s\n", GreenH(cmd), helpMsg)
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

func handleErrExit(err error) {
	if err != nil {
		errorLog(err)
		vlcPlayer.ClosePlayer()
		os.Exit(1)
	}
}

func displayErr(err error) bool {
	if err != nil {
		errorLog(err)
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
