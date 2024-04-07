package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
	http.HandleFunc("/play", playHandler)
}

func homeHandle(wr http.ResponseWriter, rq *http.Request) {
	restLog.Println("Reached home")
	json.NewEncoder(wr).Encode(resp{Item: "LudoGo"})
}

// media control
func playHandler(wr http.ResponseWriter, rq *http.Request) {
	restLog.Println("Reached play")
	q := rq.URL.Query()
	if q.Has("s") {
		searchStr := q.Get("s")
		restLog.Println("query:", searchStr)
		audio, err := app.GetSong(true)(searchStr, false)
		if err != nil {
			handleError(wr, err)
			return
		}

		app.MediaPlayer().AppendAudio(audio)
	}

	if len(app.MediaPlayer().GetQueue()) < 1 {
		handleWarn(wr, "No songs in queue")
		return
	}
	err := app.MediaPlayer().StartPlayback()
	if err != nil {
		handleError(wr, err)
	}

	status := struct{ Status string }{Status: "Playing"}

	json.NewEncoder(wr).Encode(resp{Item: status})
}

func searchHandler(wr http.ResponseWriter, rq *http.Request) {}

func radioHandler(wr http.ResponseWriter, rq *http.Request) {}

// playback control
func pauseHandler(wr http.ResponseWriter, rq *http.Request) {}

func showqHandler(wr http.ResponseWriter, rq *http.Request) {}

func currentHandler(wr http.ResponseWriter, rq *http.Request) {}

func stopHandler(wr http.ResponseWriter, rq *http.Request) {}

func quitHandler(wr http.ResponseWriter, rq *http.Request) {}

// Utils

func handleError(wr http.ResponseWriter, err error) {
	json.NewEncoder(wr).Encode(resp{Err: err})
	wr.WriteHeader(http.StatusInternalServerError)
}

func handleWarn(wr http.ResponseWriter, msg string) {
	warn := struct{ Warn string }{Warn: msg}
	json.NewEncoder(wr).Encode(resp{Item: warn})
	wr.WriteHeader(http.StatusNoContent)
}
