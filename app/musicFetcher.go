package app

import (
	"io"
	"log"
)

var fetcherLog = log.New(io.Discard, "musicFetcher: ", log.LstdFlags)

type audioFetcher interface {
	GetSong(songName string, isVideoID bool) (*AudioDetails, error)
	GetPlayList(searchString string, radioLen int, skipFirst bool, isVideoID bool) (*[]AudioDetails, error)

	SearchSong(songName string, offset int, limit int) (*[]AudioBasic, error)
}

func GetPipedSong(searchString string, isVideoID bool) (*AudioDetails, error) {

	fetcherLog.Println("Fetching song: ", searchString)

	musicId, err := resolveMusicId(searchString, isVideoID, getPipedApiMusicId)
	if err != nil {
		return nil, err
	}

	audio, err := getPipedApiAudioStream(musicId, false)
	if err != nil {
		return nil, err
	}

	fetcherLog.Println("audio: ", audio)

	return &audio, nil
}

func GetPipedRadioList(searchString string, isVideoID bool, offset int, limit int) (*[]AudioDetails, error) {
	musicId, err := resolveMusicId(searchString, isVideoID, getPipedApiMusicId)
	if err != nil {
		return nil, err
	}

	audioDetails, err := getPipedApiAudioStream(musicId, true)
	if err != nil {
		return nil, err
	}
	audioBasicList := audioDetails.RelatedAudioList

	audioBasicList = trimList(audioBasicList, offset, limit)

	audioList := make([]AudioDetails, 0)
	for _, audio := range audioBasicList {
		fetcherLog.Println("[GetRadioList] ", audio)
	}

	for i := 0; i < len(audioBasicList); i++ {

		audio, err := getPipedApiAudioStream(audioBasicList[i].YtId, false)
		if err == nil {
			audioList = append(audioList, audio)
		}

	}

	return &audioList, nil
}

func SearchPipedSong(searchString string, offset int, limit int) (*[]AudioBasic, error) {
	return getPipedSearchList(searchString, offset, limit)
}

func GetYtSong(searchString string, isVideoID bool) (*AudioDetails, error) {
	musicId, err := resolveMusicId(searchString, isVideoID, getYtMusicId)
	if err != nil {
		return nil, err
	}

	audio, err := getPipedApiAudioStream(musicId, false)
	if err != nil {
		return nil, err
	}

	fetcherLog.Println("audio: ", audio)

	return &audio, nil
}

func GetYtRadioList(searchString string, isVideoID bool, offset int, limit int) (*[]AudioDetails, error) {
	musicId, err := resolveMusicId(searchString, isVideoID, getYtMusicId)
	if err != nil {
		return nil, err
	}

	audioBasicList, err := getYtPlaylist(musicId)
	if err != nil {
		return nil, err
	}

	audioBasicList = trimList(audioBasicList, offset, limit)

	audioList := make([]AudioDetails, 0)
	for _, audio := range audioBasicList {
		fetcherLog.Println("[GetRadioList] ", audio)
	}

	for _, audioBasic := range audioBasicList {
		audio, err := getPipedApiAudioStream(audioBasic.YtId, false)
		if err == nil {
			audioList = append(audioList, audio)
		}
	}

	return &audioList, nil
}

func SearchYtSong(searchString string, offset int, limit int) (*[]AudioBasic, error) {
	return getYtSearchList(searchString, offset, limit)
}

// Helper Funcs //

func resolveMusicId(searchStr string, isVideoID bool, fetchMusicId func(string) (string, error)) (string, error) {
	if isVideoID {
		return searchStr, nil
	}

	musicId, err := fetchMusicId(searchStr)
	if err != nil {
		return "", err
	}

	return musicId, nil
}
