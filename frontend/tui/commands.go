package tui

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/johnrijoy/ludo-go/app"
)

func doCommand(cmd string, m *mainModel) {
	setCommandMode(m)

	cmd, arg := parseCommand(cmd)
	switch cmd {
	case "play":
		appendPlay(arg, m)
	case "nothing":
		m.err = errors.New("not supported")
	case "search":
		m.mode = interactiveListMode
		m.cmdInput.Prompt = "> Please select a song: "
		m.searchList = []string{"item1", "item2", "item3"}
	case "show":
		m.mode = listMode
		m.searchList = []string{"item4", "item5", "item6"}
	default:
		handleErr(Warn("Invalid command"), m)
	}

}

func doInterativeList(ind string, m *mainModel) {
	defer setCommandMode(m)

	i, err := strconv.Atoi(ind)
	if err != nil {
		m.err = errors.New("invalid index")
		return
	}

	i--

	if i < 0 || i >= len(m.searchList) {
		m.err = errors.New("index out of bounds")
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
	handleErr(err, m)

	m.searchList = make([]string, len(*audioBasicList))
	for i, audio := range *audioBasicList {
		m.searchList[i] = fmt.Sprintf("%-50s | %-20s | %s\n", i+1, safeTruncString(audio.Title, 50), safeTruncString(audio.Uploader, 20), audio.GetFormattedDuration())
	}

	setInteractiveListMode(m, "> Enter index number (q to escape): ")

	m.postSearchFunc = func(index int, m *mainModel) {
		audioBasic := (*audioBasicList)[index]

		audio, err := app.GetSong(true)(audioBasic.YtId, true)
		handleErr(err, m)

		err = app.MediaPlayer().AppendAudio(audio)
		handleErr(err, m)

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
