package app

import (
	"fmt"
	"io"
	"log"
)

var audioUtilsLog = log.New(io.Discard, "audioUtils: ", log.LstdFlags)

// BasicAudio
type AudioBasic struct {
	YtId     string
	Title    string
	Uploader string
	Duration int
}

func (audioBasic *AudioBasic) GetFormattedDuration() string {
	return GetFormattedTime(audioBasic.Duration)
}

func (audioBasic *AudioBasic) String() string {
	formattedAudioDetails := fmt.Sprintf("< song: '%s' uploader: '%s'  duration: %s >", audioBasic.Title, audioBasic.Uploader, audioBasic.GetFormattedDuration())

	return formattedAudioDetails
}

func (audioBasic *AudioBasic) validate() bool {
	isValid := false
	if audioBasic.YtId != "" && audioBasic.Title != "" {
		isValid = true
	}
	return isValid
}

// AudioDetails

type AudioDetails struct {
	AudioBasic
	AudioStreamUrl   string
	RelatedAudioList []AudioBasic
	uid              string
}

func (audioDetails *AudioDetails) validate() bool {
	isValid := false
	if audioDetails.AudioStreamUrl != "" {
		isValid = true
	}
	return isValid
}

// AudioState

type AudioState struct {
	AudioDetails
	currentTrackIndex int
}

func (audioState *AudioState) updateAudioState(audioDetails *AudioDetails) {
	audioUtilsLog.Println("Updating Audio State...")
	audioState.AudioDetails = *audioDetails
}

// Helpers
