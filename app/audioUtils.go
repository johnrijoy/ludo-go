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
	formattedDuartion := fmt.Sprintf("%dm%ds", audioBasic.Duration/60, audioBasic.Duration%60)

	return formattedDuartion
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
	currentPos        int
	totalLength       int
}

func (audioState *AudioState) updateAudioState(audioDetails *AudioDetails) {
	audioUtilsLog.Println("Updating Audio State...")
	audioState.currentPos = 0
	audioState.YtId = audioDetails.YtId
	audioState.Title = audioDetails.Title
	audioState.Duration = audioDetails.Duration
	audioState.AudioStreamUrl = audioDetails.AudioStreamUrl
}

func (audioState *AudioState) GetFormattedPosition() string {
	formattedPos := fmt.Sprintf("%dm%ds", audioState.currentPos/60, audioState.currentPos%60)

	return formattedPos
}

func (audioState *AudioState) String() string {
	formattedAudioDetails := fmt.Sprintf("< song: '%s' Pos: %s  duration: %s >", audioState.Title, audioState.GetFormattedPosition(), audioState.GetFormattedDuration())

	return formattedAudioDetails
}

func (audioState *AudioState) GetPositionDetails() (int, int) {
	return audioState.currentPos, audioState.totalLength
}

// Helpers
