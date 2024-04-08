package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/johnrijoy/ludo-go/app"
)

var restLog = log.New(os.Stdout, "restLog: ", log.LstdFlags|log.Lmsgprefix)

func Run() {
	port := 8080
	restLog.Println("Starting server at: ", port)
	loadCommandHandlers()
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func loadCommandHandlers() {
	http.HandleFunc("/", homeHandle)

	// media control
	http.HandleFunc("/play", wrapperFunc(playHandler))
	http.HandleFunc("/search", wrapperFunc(searchHandler))
	http.HandleFunc("/radio", wrapperFunc(radioHandler))
	http.HandleFunc("/pause", wrapperFunc(pauseHandler))
	http.HandleFunc("/showq", wrapperFunc(showqHandler))
	http.HandleFunc("/curr", wrapperFunc(currentHandler))
	http.HandleFunc("/stop", wrapperFunc(stopHandler))
}

func homeHandle(wr http.ResponseWriter, rq *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	restLog.Println("Reached home")
	json.NewEncoder(wr).Encode(resp{Item: "LudoGo"})
}

//

func wrapperFunc(respHandler func(q url.Values) (any, error)) func(wr http.ResponseWriter, rq *http.Request) {
	return func(wr http.ResponseWriter, rq *http.Request) {
		wr.Header().Set("Content-Type", "application/json")
		q := rq.URL.Query()
		response, err := respHandler(q)
		if err != nil {
			handleError(wr, err)
		}

		json.NewEncoder(wr).Encode(resp{Item: response})
	}
}

// media control
func playHandler(q url.Values) (any, error) {
	restLog.Println("Reached play")

	if q.Has("s") {
		searchStr := q.Get("s")
		restLog.Println("query:", searchStr)
		audio, err := app.GetSong(true)(searchStr, false)
		if err != nil {
			return nil, err
		}

		app.MediaPlayer().AppendAudio(audio)
	}

	if len(app.MediaPlayer().GetQueue()) < 1 {
		return nil, Warn("No songs in queue")
	}
	err := app.MediaPlayer().StartPlayback()
	if err != nil {
		return nil, err
	}

	status := struct{ Status string }{Status: "Playing"}
	return status, nil
}

func searchHandler(q url.Values) (any, error) {
	restLog.Println("Reached Search")
	if !q.Has("s") {
		return nil, errors.New("no search query")
	}

	searchStr := q.Get("s")
	restLog.Println("query:", searchStr)

	audioBasicList, err := app.GetSearchList(true)(searchStr, 0, 10)
	if err != nil {
		return nil, err
	}

	resultList := struct{ SearchResult []app.AudioBasic }{SearchResult: *audioBasicList}

	return resultList, nil
}

func radioHandler(q url.Values) (any, error) {
	restLog.Println("Reached Search")

	var audio *app.AudioDetails

	if !q.Has("s") {
		if len(app.MediaPlayer().GetQueue()) < 1 {
			return nil, Warn("No songs in queue")
		}
		audioD := app.MediaPlayer().GetAudioState().AudioDetails
		audio = &audioD
		removeAllIndex("")
	} else {
		searchStr := q.Get("s")
		var err error
		audio, err = app.GetSong(false)(searchStr, false)
		if err != nil {
			return nil, err
		}

		err = app.MediaPlayer().ResetPlayer()
		if err != nil {
			return nil, err
		}

		app.MediaPlayer().AppendAudio(audio)
		app.MediaPlayer().StartPlayback()
	}

	go func() {
		audioList, _ := app.GetPlayList(false)(audio.YtId, true, 1, 10)

		for _, audio := range *audioList {
			app.MediaPlayer().AppendAudio(&audio)
		}
	}()

	status := struct{ msg string }{msg: "Started radio"}
	return status, nil
}

// playback control
func pauseHandler(q url.Values) (any, error) {
	restLog.Println("Reached pause")

	err := app.MediaPlayer().PauseResume()
	if err != nil {
		return nil, err
	}

	status := struct{ isPaused bool }{isPaused: app.MediaPlayer().IsPlaying()}
	return status, nil
}

func showqHandler(q url.Values) (any, error) {
	restLog.Println("Reached showq")

	audList := app.MediaPlayer().GetQueue()
	qIndex := app.MediaPlayer().GetQueueIndex()

	audQ := struct {
		Q     []app.AudioDetails
		Index int
	}{Q: audList, Index: qIndex}
	return audQ, nil
}

func currentHandler(q url.Values) (any, error) {
	restLog.Println("Reached current")

	stat := ""

	if app.MediaPlayer().IsPlaying() {
		stat = "Playing"
	} else if app.MediaPlayer().CheckMediaError() {
		stat = "Media Error"
	} else {
		stat = app.MediaPlayer().FetchPlayerState()
	}

	audState := app.MediaPlayer().GetAudioState()
	currPos, totPos := (&audState).GetPositionDetails()

	audStat := struct {
		Audio  app.AudioBasic
		Curr   int
		Total  int
		Status string
	}{Audio: audState.AudioBasic, Curr: currPos, Total: totPos, Status: stat}
	return audStat, nil
}

func stopHandler(q url.Values) (any, error) {
	restLog.Println("Reached current")

	err := app.MediaPlayer().ResetPlayer()
	if err != nil {
		return nil, err
	}

	status := struct{ isPaused bool }{isPaused: app.MediaPlayer().IsPlaying()}
	return status, nil
}

func quitHandler(q url.Values) (any, error) { return nil, nil }

// Utils

func removeAllIndex(arg string) error {
	trackIndex := app.MediaPlayer().GetQueueIndex() + 1
	if arg != "" {
		var err error
		trackIndex, err = strconv.Atoi(arg)
		if err != nil {
			return err
		}
		trackIndex -= 1
	}

	err := app.MediaPlayer().RemoveAllAudioFromIndex(trackIndex)
	return err
}

func handleError(wr http.ResponseWriter, err error) {

	if warn, ok := err.(ErrWarn); ok {
		warn := struct{ Warn string }{Warn: warn.msg}
		json.NewEncoder(wr).Encode(resp{Item: warn})
		wr.WriteHeader(http.StatusNoContent)
		return
	}

	json.NewEncoder(wr).Encode(resp{Err: err})
	wr.WriteHeader(http.StatusInternalServerError)
}
