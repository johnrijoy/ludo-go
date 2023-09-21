package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func getSong(songName string) (*AudioDetails, error) {

	log.Println("Fetching song: ", songName)

	musicId, err := getPipedApiMusicId(songName)
	if err != nil {
		return nil, err
	}

	target := getPipedApi() + "streams/" + musicId

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
		audio.title = title
	}
	if uploader, ok := getValue(response, path{"uploader"}).(string); ok {
		audio.uploader = uploader
	}
	if duration, ok := getValue(response, path{"duration"}).(float64); ok {
		audio.duration = int(duration)
	}
	if trackUrl, ok := getValue(response, path{"audioStreams", 0, "url"}).(string); ok {
		audio.audioStreamUrl = trackUrl
	}

	log.Println("audio: ", audio)

	return &audio, nil
}
