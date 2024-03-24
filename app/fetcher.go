package app

import (
	"log"
)

func GetSong(songName string, loadRelated bool) (*AudioDetails, error) {

	log.Println("Fetching song: ", songName)

	musicId, err := getPipedApiMusicId(songName)
	if err != nil {
		return nil, err
	}

	audio, err := getPipedApiSong(musicId, loadRelated)
	if err != nil {
		return nil, err
	}

	log.Println("audio: ", audio)

	return &audio, nil
}
