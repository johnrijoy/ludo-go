package tui

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/johnrijoy/ludo-go/app"
	"github.com/johnrijoy/ludo-go/frontend"
)

func doCommand(cmd string, m *mainModel) {
	setCommandMode(m)

	cmd, arg := parseCommand(cmd)
	switch cmd {
	case "play":
		appendPlay(arg, m)
	case "search", "s":
		searchPlay(arg, m)

	case "radio":
		radioPlay(arg, m)

	case "p", "pause", "resume":
		app.MediaPlayer().PauseResume()

	case "showq", "q":
		displayQueue(m)

	case "curr", "c":
		displayCurrentSong(m)

	case "skipn", "n":
		skipNext(m)

	case "skipb", "b":
		skipPrevious(m)

	case "skip":
		skipIndex(arg, m)

	case "remove", "rem":
		removeIndex(arg, m)

	case "removeAll", "reml":
		removeAllIndex(arg, m)

	case "forward", "f":
		audioForward(arg, m)

	case "rewind", "r":
		audioRewind(arg, m)

	case "setVol", "v":
		modifyVolume(arg, m)

	case "stop":
		resetPlayer(m)

	case "like":
		likeSong(arg, m)

	case "listSongs", "ls":
		fetchSongList(arg, m)

	case "setApi":
		modifyApi(arg, m)

	case "checkApi":
		fmt.Println("Piped Api: ", Green(app.Piped.GetPipedApi()))

	case "listApi":
		displayApiList(m)

	case "version":
		displayVersion(m)

	case "help":
		showHelp(m)

	case "quit":
		m.quit = true
	default:
		handleErr(Warn("Invalid command"), m)
	}

}

func doInterativeList(ind string, m *mainModel) {
	defer setCommandMode(m)

	if ind == "q" {
		m.resultMsg = "exiting search..."
		return
	}

	i, err := strconv.Atoi(ind)
	if handleErr(err, m) {
		return
	}

	i--

	if i < 0 || i >= len(m.searchList) {
		handleErr(errors.New("index out of bounds"), m)
		return
	}

	m.postSearchFunc(i, m)
}

// Commands

func appendPlay(arg string, m *mainModel) {
	if arg != "" {
		audio, err := app.GetSong(true)(arg, false)
		handleErr(err, m)

		app.MediaPlayer().AppendAudio(audio)
	}
	if len(app.MediaPlayer().GetQueue()) < 1 {
		handleErr(Warn("No songs in queue"), m)
	}
	err := app.MediaPlayer().StartPlayback()
	handleErr(err, m)
}

func radioPlay(arg string, m *mainModel) {
	if arg == "" {
		handleErr(Warn("please enter a search query"), m)
		return
	}

	var audio *app.AudioDetails
	if arg == "." {
		audioD := app.MediaPlayer().GetAudioState().AudioDetails
		audio = &audioD
		removeAllIndex("", m)
	} else {
		var err error
		audio, err = app.GetSong(false)(arg, false)
		handleErr(err, m)

		err = app.MediaPlayer().ResetPlayer()
		handleErr(err, m)

		app.MediaPlayer().AppendAudio(audio)
		app.MediaPlayer().StartPlayback()
	}

	go func() {
		audioList, err := app.GetPlayList(false)(audio.YtId, true, 1, 10)
		handleErr(err, m)

		for _, audio := range *audioList {
			app.MediaPlayer().AppendAudio(&audio)
		}
	}()
}

func searchPlay(arg string, m *mainModel) {
	if arg == "" {
		handleErr(Warn("please enter a search query"), m)
		return
	}

	audioBasicList, err := app.GetSearchList(true)(arg, 0, 10)
	if handleErr(err, m) {
		return
	}

	m.searchList = make([]string, len(*audioBasicList))
	for i, audio := range *audioBasicList {
		m.searchList[i] = fmt.Sprintf("%-30s | %-20s | %s", safeTruncString(audio.Title, 30), safeTruncString(audio.Uploader, 20), audio.GetFormattedDuration())
	}

	setInteractiveListMode(m, "> Enter index number (q to escape): ")

	m.postSearchFunc = func(index int, m *mainModel) {
		audioBasic := (*audioBasicList)[index]

		audio, err := app.GetSong(true)(audioBasic.YtId, true)
		if handleErr(err, m) {
			return
		}

		err = app.MediaPlayer().AppendAudio(audio)
		if handleErr(err, m) {
			return
		}

		if !app.MediaPlayer().IsPlaying() {
			app.MediaPlayer().StartPlayback()
		}
	}
}

// media queue control

func removeAllIndex(arg string, m *mainModel) {
	trackIndex := app.MediaPlayer().GetQueueIndex() + 1
	if arg != "" {
		var err error
		trackIndex, err = strconv.Atoi(arg)
		handleErr(err, m)
		trackIndex -= 1
	}

	err := app.MediaPlayer().RemoveAllAudioFromIndex(trackIndex)
	handleErr(err, m)
}

// To be modified
func displayVersion(m *mainModel) {
	s := fmt.Sprintln("Ludo version: ", Green(app.Version))
	s += fmt.Sprintln("Api: ", Green(app.Piped.GetPipedApi()))
	s += fmt.Sprintln("libVlc Binding Version: ", Green(app.Info().String()))
	s += fmt.Sprintln("Vlc Runtime Version: ", Green(app.Info().Changeset()))

	m.resultMsg = s
}

func displayApiList(m *mainModel) {
	apiList, err := app.Piped.GetPipedInstanceList()
	if err != nil {
		handleErr(errors.Join(errors.New("error in fetching Instance list"), err), m)
		return
	}

	m.searchList = make([]string, len(apiList))

	for i, inst := range apiList {
		m.searchList[i] = fmt.Sprintf("%s\n", inst)
	}

	setInteractiveListMode(m, "> Enter index number to change api (q to escape): ")

	m.postSearchFunc = func(index int, m *mainModel) {
		newApi := apiList[index].ApiUrl
		app.Piped.SetPipedApi(newApi)
	}
}

func modifyApi(arg string, m *mainModel) {
	if arg == "" {
		handleErr(Warn("no arguments"), m)
		return
	}

	app.Piped.SetPipedApi(arg)
	m.resultMsg = fmt.Sprintln("Api changed from ", Gray(app.Piped.GetOldPipedApi()), " to ", Green(app.Piped.GetPipedApi()))
}

func resetPlayer(m *mainModel) {
	err := app.MediaPlayer().ResetPlayer()
	handleErr(err, m)
}

func audioRewind(arg string, m *mainModel) {
	duration := defaultForwardRewind
	if arg != "" {
		var err error
		duration, err = strconv.Atoi(arg)
		if handleErr(err, m) {
			return
		}
	}
	err := app.MediaPlayer().RewindBySeconds(duration)
	handleErr(err, m)
}

func audioForward(arg string, m *mainModel) {
	duration := defaultForwardRewind
	if arg != "" {
		var err error
		duration, err = strconv.Atoi(arg)
		if handleErr(err, m) {
			return
		}
	}
	err := app.MediaPlayer().ForwardBySeconds(duration)
	handleErr(err, m)
}

func removeIndex(arg string, m *mainModel) {
	trackIndex := len(app.MediaPlayer().GetQueue()) - 1
	if arg != "" {
		var err error
		trackIndex, err = strconv.Atoi(arg)
		if handleErr(err, m) {
			return
		}
		trackIndex -= 1
	}

	err := app.MediaPlayer().RemoveAudioFromIndex(trackIndex)
	handleErr(err, m)
}

func skipIndex(arg string, m *mainModel) {
	trackIndex := app.MediaPlayer().GetQueueIndex() + 1
	if arg != "" {
		var err error
		trackIndex, err = strconv.Atoi(arg)
		if handleErr(err, m) {
			return
		}
	}

	err := app.MediaPlayer().SkipToIndex(trackIndex)
	handleErr(err, m)
}

func skipPrevious(m *mainModel) {
	err := app.MediaPlayer().SkipToPrevious()
	handleErr(err, m)
}

func skipNext(m *mainModel) {
	err := app.MediaPlayer().SkipToNext()
	handleErr(err, m)
}

// not required
func displayCurrentSong(m *mainModel) {
	if app.MediaPlayer().IsPlaying() {
		fmt.Println(Green("Now playing..."))
	} else if app.MediaPlayer().CheckMediaError() {
		fmt.Println(Red("Error in playing media"))
	} else {
		val, ok := app.PlayerStateString(app.MediaPlayer().FetchPlayerState())
		if !ok {
			val = "Error"
		}
		fmt.Println(Yellow(val))
	}

	audState := app.MediaPlayer().GetAudioState()
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

func displayQueue(m *mainModel) {
	audList := app.MediaPlayer().GetQueue()
	qIndex := app.MediaPlayer().GetQueueIndex()

	if len(audList) == 0 {
		handleErr(Warn("no songs in queue"), m)
		return
	}

	m.searchList = make([]string, len(audList))
	for i, audio := range audList {
		msg := fmt.Sprintf("%-50s | %-50s", safeTruncString(audio.Title, 50), safeTruncString(audio.Uploader, 50))
		if qIndex == i {
			msg = Magenta(msg)
		}
		m.searchList[i] = msg
	}

	setListMode(m)
}

func modifyVolume(arg string, m *mainModel) {

	if arg == "" {
		handleErr(Warn("No volume given"), m)
		return
	}

	vol, err := strconv.Atoi(arg)
	if handleErr(err, m) {
		return
	}

	err = app.MediaPlayer().SetVol(vol)
	if !handleErr(err, m) {
		m.resultMsg = fmt.Sprintln("volume set:", Green(fmt.Sprintf("%d", vol)))
	}
}

func fetchSongList(arg string, m *mainModel) {
	if arg == "" {
		handleErr(Warn("No criteria given"), m)
		return
	}

	var criteria app.AudioListCriteria
	switch arg {
	case "recent":
		m.listTitle = "Recently Played"
		criteria = app.RecentlyPlayed
	case "plays":
		m.listTitle = "Most Played"
		criteria = app.MostPlayed
	case "likes":
		m.listTitle = "Most Liked"
		criteria = app.MostLikes
	}

	audDocs, err := app.AudioDb().GetAudioList(criteria, 0, 10)
	if handleErr(err, m) {
		return
	}

	m.searchList = make([]string, len(audDocs))
	for i, audDoc := range audDocs {
		m.searchList[i] = fmt.Sprintf("%-30s | %-20s | %s", safeTruncString(audDoc.Title, 30), safeTruncString(audDoc.Uploader, 20), audDoc.GetFormattedDuration())
	}

	setListMode(m)
}

func likeSong(arg string, m *mainModel) {
	trackIndex := app.MediaPlayer().GetQueueIndex()
	if arg != "" {
		var err error
		trackIndex, err = strconv.Atoi(arg)
		if handleErr(err, m) {
			return
		}
		trackIndex -= 1
	}
	if trackIndex < 0 || trackIndex >= len(app.MediaPlayer().GetQueue()) {
		handleErr(errors.New("invalid item index"), m)
		return
	}
	app.AudioDb().UpdateLikes(app.MediaPlayer().GetQueue()[trackIndex].YtId)
}

func showStartupMessage(m *mainModel) {
	fmt.Println(Blue("==="), Magenta("LUDO GO"), Blue("==="))
	fmt.Println("Welcome to", Magenta("LudoGo"))
	fmt.Println("To start listening, enter " + Green("play <song name>"))
	fmt.Println("To show help, enter " + Green("help"))
}

func showHelp(m *mainModel) {
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
