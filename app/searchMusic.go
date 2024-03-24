package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func getPipedApiMusicId(search string) (string, error) {
	escapedSearch := url.QueryEscape(search)
	target := GetPipedApi() + "/search?q=" + escapedSearch + "&filter=music_songs"

	log.Println("target: ", target)

	resp, err := http.Get(target)
	if err != nil {
		return "", err
	}

	log.Println("Resp status: ", resp.Status)

	if resp.StatusCode != 200 {
		err = errors.New("[getPipedApiMusicId] bad response from api")
		return "", err
	}

	var response interface{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	if trackUrl, ok := getValue(response, path{"items", 0, "url"}).(string); ok {
		musicId := strings.Split(trackUrl, "=")[1]
		return musicId, nil
	}

	return "", errors.New("could not fetch music Id")
}

func getPipedApiSong(musicId string, loadRelated bool) (AudioDetails, error) {
	target := GetPipedApi() + "/streams/" + musicId

	log.Println("target: ", target)

	resp, err := http.Get(target)
	if err != nil {
		return AudioDetails{}, err
	}

	log.Println("Resp status: ", resp.Status)

	if resp.StatusCode != 200 {
		err = errors.New("[GetSong] bad response from api")
		return AudioDetails{}, err
	}

	var response interface{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return AudioDetails{}, err
	}

	var audio AudioDetails
	audio.YtId = musicId

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
	if loadRelated {
		audio.RelatedAudioList = getPipedApiRelatedSongs(response)
	}

	return audio, nil
}

func getPipedApiRelatedSongs(response interface{}) []AudioBasic {

	relatedList, ok := getValue(response, path{"relatedStreams"}).([]interface{})
	if !ok || len(relatedList) <= 0 {
		return []AudioBasic{}
	}

	var audioList []AudioBasic

	for _, relatedItem := range relatedList {
		streamType := false
		if audioType, ok := getValue(relatedItem, path{"type"}).(string); ok && audioType == "stream" {
			streamType = true
		}

		if streamType {
			var audio AudioBasic

			if trackUrl, ok := getValue(relatedItem, path{"url"}).(string); ok {
				audio.YtId = strings.Split(trackUrl, "=")[1]
			}
			if title, ok := getValue(relatedItem, path{"title"}).(string); ok {
				audio.Title = title
			}
			if uploader, ok := getValue(relatedItem, path{"uploader"}).(string); ok {
				audio.Uploader = uploader
			}
			if duration, ok := getValue(relatedItem, path{"duration"}).(float64); ok {
				audio.Duration = int(duration)
			}

			if audio.Duration < 500 {
				audioList = append(audioList, audio)
			}
		}
	}

	return audioList
}

// Json parser

type path []interface{}

func getValue(source interface{}, path path) interface{} {
	value := source
	for _, element := range path {
		mustBreak := false
		switch element.(type) {
		case string:
			if val, ok := value.(map[string]interface{})[element.(string)]; ok {
				value = val
			} else {
				value = nil
				mustBreak = true
			}
		case int:
			if len(value.([]interface{})) > element.(int) {
				value = value.([]interface{})[element.(int)]
			} else {
				value = nil
				mustBreak = true
			}
		}
		if mustBreak {
			break
		}
	}
	return value
}
