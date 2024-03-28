package app

import (
	"strings"

	"github.com/raitonoberu/ytmusic"
)

func getYtMusicId(searchString string) (string, error) {
	search := ytmusic.Search(searchString)
	result, err := search.Next()
	if err != nil {
		return "", err
	}

	return result.Tracks[0].VideoID, nil
}

func getYtPlaylist(musicId string) ([]AudioBasic, error) {
	trackItems, err := ytmusic.GetWatchPlaylist(musicId)
	if err != nil {
		return []AudioBasic{}, err
	}

	return trackItemToAudioBasic(trackItems), nil
}

func getYtSearchList(searchString string, offset int, limit int) (*[]AudioBasic, error) {
	search := ytmusic.Search(searchString)
	result, err := search.Next()
	if err != nil {
		return nil, err
	}

	trackList := trimList(result.Tracks, offset, limit)

	audioBasicList := trackItemToAudioBasic(trackList)

	return &audioBasicList, nil
}

// Helper functions //

func trackItemToAudioBasic(trackItems []*ytmusic.TrackItem) []AudioBasic {
	audioList := make([]AudioBasic, len(trackItems))

	for i := 0; i < len(trackItems); i++ {
		audioList[i].Title = trackItems[i].Title
		audioList[i].Duration = trackItems[i].Duration
		audioList[i].YtId = trackItems[i].VideoID
		audioList[i].Uploader = getArtistString(trackItems[i].Artists)
	}

	return audioList
}

func getArtistString(artists []ytmusic.Artist) string {
	var artistsSlice []string
	for _, artist := range artists {
		artistsSlice = append(artistsSlice, artist.Name)
	}

	var artistString string
	if len(artistsSlice) > 2 {
		artistString = strings.Join(artistsSlice[:len(artistsSlice)-2], ", ")
	}
	if len(artistsSlice) > 1 {
		artistString = artistString + " & " + artistsSlice[len(artistsSlice)-1]
	} else if len(artistsSlice) > 0 {
		artistString = artistsSlice[0]
	}

	return artistString
}
