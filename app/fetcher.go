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

func GetYtSong(songName string) (*AudioDetails, error) {
	musicId, err := GetYtMusicId(songName)
	if err != nil {
		return nil, err
	}

	audio, err := getPipedApiSong(musicId, false)
	if err != nil {
		return nil, err
	}

	log.Println("audio: ", audio)

	return &audio, nil
}

func GetYtRadioList(searchString string, radioLen int, skipFirst bool, isVideoID bool) (*[]AudioDetails, error) {
	var musicId string
	if isVideoID {
		musicId = searchString
	} else {
		var err error
		musicId, err = GetYtMusicId(searchString)
		if err != nil {
			return nil, err
		}
	}

	audioBasicList, err := GetYtPlaylist(musicId)
	if err != nil {
		return nil, err
	}

	// first song the song itself
	if skipFirst {
		audioBasicList = audioBasicList[1:]
	}

	if len(audioBasicList) > radioLen {
		audioBasicList = audioBasicList[:radioLen]
	}

	audioList := make([]AudioDetails, 0)
	for _, audio := range audioBasicList {
		log.Println("[GetRadioList] ", audio)
	}

	for i := 0; i < len(audioBasicList); i++ {

		audio, err := getPipedApiSong(audioBasicList[i].YtId, false)
		if err == nil {
			audioList = append(audioList, audio)
		}

	}

	return &audioList, nil
}
