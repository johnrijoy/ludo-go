package prompt

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
	"github.com/johnrijoy/ludo-go/frontend"
)

var (
	vlcPlayer     *app.VlcPlayer
	audioDb       *app.AudioDatastore
	isPipedSource bool = true
)

const defaultForwardRewind = 10

func Run() {
	exitSig := false

	log.SetOutput(io.Discard)
	err := app.Init()
	handleErrExit(err)
	defer app.Close()

	vlcPlayer = app.MediaPlayer()
	audioDb = app.AudioDb()

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

	case "setVol", "v":
		modifyVolume(arg)

	case "stop":
		resetPlayer()

	case "like":
		likeSong(arg)

	case "listSongs", "ls":
		fetchSongList(arg)

	case "setApi":
		modifyApi(arg)

	case "checkApi":
		fmt.Println("Piped Api: ", Green(app.Piped.GetPipedApi()))

	case "listApi":
		displayApiList()

	case "randApi":
		modifyApiRandom()

	case "setSource", "ss":
		modifySource(arg)

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
	fmt.Println("Api: ", Green(app.Piped.GetPipedApi()))
	fmt.Println("libVlc Binding Version: ", Green(app.Info().String()))
	fmt.Println("Vlc Runtime Version: ", Green(app.Info().Changeset()))
}

func modifyApiRandom() {
	apiList, err := app.Piped.GetPipedInstanceList()
	if err != nil {
		errorLog("Error in fetching Instance list:", err)
	} else {
		randIndex := rand.Intn(len(apiList))
		app.Piped.SetPipedApi(apiList[randIndex].ApiUrl)
	}
}

func displayApiList() {
	apiList, err := app.Piped.GetPipedInstanceList()
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
	app.Piped.SetPipedApi(newApi)
}

func modifyApi(arg string) {
	if arg != "" {
		app.Piped.SetPipedApi(arg)
		fmt.Println("Api changed from ", Gray(app.Piped.GetOldPipedApi()), " to ", Green(app.Piped.GetPipedApi()))
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
	var statusMsg string
	if vlcPlayer.IsPlaying() {
		statusMsg = Green(fmt.Sprintf("%-30s", "Now playing..."))
	} else if vlcPlayer.CheckMediaError() {
		statusMsg = Red(fmt.Sprintf("%-30s", "Error in playing media"))
	} else {
		statusMsg = Yellow(fmt.Sprintf("%-30s", vlcPlayer.FetchPlayerState()))
	}

	aud := vlcPlayer.GetAudioState().AudioBasic
	currPos, totPos := vlcPlayer.GetMediaPosition()

	scale := 50

	navMsg := Gray(strings.Repeat("-", scale))
	if totPos > currPos {
		scaledCurrPos := int(math.Round((float64(currPos) / float64(totPos)) * float64(scale)))
		restPosition := scale - scaledCurrPos

		navMsg = Magenta(strings.Repeat(">", scaledCurrPos)) + Gray(strings.Repeat("-", restPosition))
	}

	fmt.Printf("%s%10s%10s\n", statusMsg, app.GetFormattedTime(currPos), app.GetFormattedTime(totPos))
	fmt.Printf("%s\n", navMsg)
	fmt.Printf("%-30s%20s\n", safeTruncString(aud.Title, 30), safeTruncString(aud.Uploader, 20))
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
		audio, err = app.GetSong(isPipedSource)(arg, false)
		handleErrExit(err)

		err = vlcPlayer.ResetPlayer()
		handleErrExit(err)

		vlcPlayer.AppendAudio(audio)
		vlcPlayer.StartPlayback()
	}

	go func() {
		audioList, err := app.GetPlayList(isPipedSource)(audio.YtId, true, 1, 10)
		handleErrExit(err)

		for _, audio := range *audioList {
			vlcPlayer.AppendAudio(&audio)
		}
	}()
}

func appendPlay(arg string) {
	if arg != "" {
		audio, err := app.GetSong(isPipedSource)(arg, false)
		handleErrExit(err)

		err = vlcPlayer.AppendAudio(audio)
		handleErrExit(err)
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

	audioBasicList, err := app.GetSearchList(isPipedSource)(arg, 0, 10)
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

	audio, err := app.GetSong(isPipedSource)(audioBasic.YtId, true)
	if displayErr(err) {
		return
	}

	err = vlcPlayer.AppendAudio(audio)
	displayErr(err)

	if !vlcPlayer.IsPlaying() {
		vlcPlayer.StartPlayback()
	}
}

func modifyVolume(arg string) {

	if arg == "" {
		warnLog("No volume given")
		return
	}
	vol, err := strconv.Atoi(arg)
	if displayErr(err) {
		return
	}

	err = vlcPlayer.SetVol(vol)
	if !displayErr(err) {
		fmt.Println("volume set:", Green(vol))
	}
}

func fetchSongList(arg string) {
	if arg == "" {
		warnLog("No criteria given")
		return
	}

	var criteria app.AudioListCriteria
	switch arg {
	case "recent":
		fmt.Println("Recently Played")
		criteria = app.RecentlyPlayed
	case "plays":
		fmt.Println("Most Played")
		criteria = app.MostPlayed
	case "likes":
		fmt.Println("Most Liked")
		criteria = app.MostLikes
	}

	audDocs, err := audioDb.GetAudioList(criteria, 0, 10)
	if err != nil {
		errorLog(err)
	}

	for i, audDoc := range audDocs {
		fmt.Printf("%-2d - %-50s | %-20s | %20s\n", i+1, safeTruncString(audDoc.Title, 50), safeTruncString(audDoc.Uploader, 20), audDoc.GetFormattedDuration())
	}
}

func likeSong(arg string) {
	trackIndex := vlcPlayer.GetQueueIndex()
	if arg != "" {
		var err error
		trackIndex, err = strconv.Atoi(arg)
		displayErr(err)
		trackIndex -= 1
	}
	if trackIndex < 0 || trackIndex >= len(vlcPlayer.GetQueue()) {
		errorLog("invalid item index")
	}
	audioDb.UpdateLikes(vlcPlayer.GetQueue()[trackIndex].YtId)
}

func modifySource(arg string) {
	switch arg {
	case "youtube", "yt":
		isPipedSource = false
		fmt.Println("Source changed to ", Green("Youtube"))
	case "piped", "pp":
		isPipedSource = true
		fmt.Println("Source changed to ", Green("Piped"))
	case "":
		warnLog("no source given")
	default:
		errorLog("Invalid source")
	}
}

func showStartupMessage() {
	fmt.Println(Blue("==="), Magenta("LUDO GO"), Blue("==="))
	fmt.Println("Welcome to", Magenta("LudoGo"))
	fmt.Println("To start listening, enter " + Green("play <song name>"))
	fmt.Println("To show help, enter " + Green("help"))
}

func showHelp() {
	displayList := func(items []string) {
		for _, val := range items {
			splitVal := strings.Split(val, "-")
			cmd, help := splitVal[0], strings.Join(splitVal[1:], "-")
			cmd = strings.ReplaceAll(cmd, ",", ", ")
			helpSplit := strings.Split(help, "|")
			helpMsg := helpSplit[0]
			if len(helpSplit) > 1 {
				helpMsg = strings.TrimSpace(strings.Join(helpSplit[0:len(helpSplit)-1], "|"))
				usageMsg := Magenta(" | " + strings.TrimSpace(helpSplit[len(helpSplit)-1]))
				helpMsg += usageMsg
			}

			fmt.Printf("%-40s - %s\n", Green(cmd), helpMsg)
		}
	}
	fmt.Println("Commands")
	displayList(frontend.Commands)

	fmt.Println("\nProperties")
	displayList(frontend.Configs)
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
