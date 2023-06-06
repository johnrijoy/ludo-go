package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/raitonoberu/ytmusic"
)

func getSong(songName string) (*AudioDetails, error) {

	log.Println("Fetching song: ", songName)

	searchClient := ytmusic.Search(songName)
	result, err := searchClient.Next()
	checkErr(err)

	//jsonstr, _ := json.MarshalIndent(result, "", "    ")
	// fmt.Println(string(jsonstr))

	// for _, val := range result.Tracks {
	// 	arts := ""
	// 	for _, art := range val.Artists {
	// 		arts += art.Name
	// 	}
	// 	fmt.Printf("%v-%v %v\n", val.Title, arts, val.VideoID)
	// }

	target := "https://pipedapi.kavin.rocks/streams/" + result.Tracks[0].VideoID

	log.Println("target: ", target)

	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}

	log.Println("Resp status: ", resp.Status)

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
