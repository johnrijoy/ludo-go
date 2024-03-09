package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func GetSong(songName string) (*AudioDetails, error) {

	log.Println("Fetching song: ", songName)

	musicId, err := getPipedApiMusicId(songName)
	if err != nil {
		return nil, err
	}

	target := GetPipedApi() + "streams/" + musicId

	log.Println("target: ", target)

	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}

	log.Println("Resp status: ", resp.Status)

	if resp.StatusCode != 200 {
		err = errors.New("bad response from api")
		return nil, err
	}

	// bB, err := io.ReadAll(resp.Body)
	// checkErr(err)

	var response interface{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var audio AudioDetails
	if title, ok := getValue(response, path{"title"}).(string); ok {
		audio.Title = title
	}
	if uploader, ok := getValue(response, path{"uploader"}).(string); ok {
		audio.Uploader = uploader
	}
	if duration, ok := getValue(response, path{"duration"}).(float64); ok {
		audio.Duration = int(duration)
	}
	if trackUrl, ok := getValue(response, path{"audioStreams", 0, "url"}).(string); ok {
		audio.AudioStreamUrl = trackUrl
	}

	log.Println("audio: ", audio)

	return &audio, nil
}
