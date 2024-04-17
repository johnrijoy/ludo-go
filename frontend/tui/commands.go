package tui

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/johnrijoy/ludo-go/app"
	"github.com/johnrijoy/ludo-go/frontend"
)

var isPiped bool

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

	case "setSource", "ss":
		setSource(arg, m)

	case "version":
		displayVersion(m)

	case "help":
		setHelpMode(m)

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

// Audio Search

func appendPlay(arg string, m *mainModel) {
	if arg != "" {

		if strings.HasPrefix(arg, "/a") && isPiped {
			app.SetPipedAllFilterType(true)
			aft, _ := strings.CutPrefix(arg, "/a")
			arg = strings.TrimSpace(aft)
		}

		audio, err := app.GetSong(isPiped)(arg, false)
		if handleErr(err, m) {
			return
		}

		if isPiped && app.GetPipedAllFilterType() {
			app.SetPipedAllFilterType(false)
		}

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
		audio, err = app.GetSong(isPiped)(arg, false)
		handleErr(err, m)

		err = app.MediaPlayer().ResetPlayer()
		handleErr(err, m)

		app.MediaPlayer().AppendAudio(audio)
		app.MediaPlayer().StartPlayback()
	}

	go func() {
		audioList, err := app.GetPlayList(isPiped)(audio.YtId, true, 1, 10)
		if handleErr(err, m) {
			return
		}

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

	if strings.HasPrefix(arg, "/a") && isPiped {
		app.SetPipedAllFilterType(true)
		aft, _ := strings.CutPrefix(arg, "/a")
		arg = strings.TrimSpace(aft)
	}

	audioBasicList, err := app.GetSearchList(isPiped)(arg, 0, 10)
	if handleErr(err, m) {
		return
	}

	if isPiped && app.GetPipedAllFilterType() {
		app.SetPipedAllFilterType(false)
	}

	m.searchList = make([]string, len(*audioBasicList))
	m.highlightIndices = []int{}
	for i, audio := range *audioBasicList {
		m.searchList[i] = fmt.Sprintf("%-2d - %-30s | %-20s | %s", i+1, safeTruncString(audio.Title, 30), safeTruncString(audio.Uploader, 20), audio.GetFormattedDuration())
	}

	setInteractiveListMode(m, "> Enter index number (q to escape): ")

	m.postSearchFunc = func(index int, m *mainModel) {
		audioBasic := (*audioBasicList)[index]

		audio, err := app.GetSong(isPiped)(audioBasic.YtId, true)
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
		trackIndex--
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

// media playback control
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

func resetPlayer(m *mainModel) {
	err := app.MediaPlayer().ResetPlayer()
	handleErr(err, m)
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

// Info commands
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
	m.highlightIndices = []int{}
	for i, inst := range apiList {
		m.searchList[i] = fmt.Sprintf("%-2d - %s\n", i+1, inst)
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

func displayQueue(m *mainModel) {
	audList := app.MediaPlayer().GetQueue()
	qIndex := app.MediaPlayer().GetQueueIndex()

	if len(audList) == 0 {
		handleErr(Warn("no songs in queue"), m)
		return
	}

	m.searchList = make([]string, len(audList))
	m.highlightIndices = []int{}
	for i, audio := range audList {
		if qIndex == i {
			m.highlightIndices = []int{i}
		}
		m.searchList[i] = fmt.Sprintf("%-2d - %-50s | %s", i+1, safeTruncString(audio.Title, 50), safeTruncString(audio.Uploader, 50))
	}

	setListMode(m)
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
	m.highlightIndices = []int{}
	for i, audDoc := range audDocs {
		m.searchList[i] = fmt.Sprintf("%-2d - %-30s | %-20s | %s", i+1, safeTruncString(audDoc.Title, 30), safeTruncString(audDoc.Uploader, 20), audDoc.GetFormattedDuration())
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

func setSource(arg string, m *mainModel) {

	if arg == "" {
		if isPiped {
			m.resultMsg = "Source is Piped"
		} else {
			m.resultMsg = "Source is Youtube"
		}
		return
	}

	switch arg {
	case "youtube", "yt":
		isPiped = false
		m.resultMsg = "Source changed to Youtube"
	case "piped", "pp":
		isPiped = true
		m.resultMsg = "Source changed to Piped"
	default:
		handleErr(Warn("Source not valid (youtube/yt, piped/pp)"), m)
		return
	}
}

func showStartupMessage(m *mainModel) {
	fmt.Println("Welcome to", Magenta("LudoGo"))
	fmt.Println("To start listening, enter " + Green("play <song name>"))
	fmt.Println("To show help, enter " + Green("help"))
}

func showHelp() string {
	var help strings.Builder

	displayList := func(items []string) string {
		s := ""
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

			s += fmt.Sprintf("%-40s - %s\n", Green(cmd), helpMsg)

		}
		return s
	}
	help.WriteString("Commands\n")
	help.WriteString(displayList(frontend.Commands))

	help.WriteString("\nProperties\n")
	help.WriteString(displayList(frontend.Configs))

	return help.String()
}
