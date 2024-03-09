package app

import (
	"fmt"
	"log"
)

type AudioDetails struct {
	YtId           string
	Title          string
	Uploader       string
	Duration       int
	AudioStreamUrl string
}

func (audioDetails *AudioDetails) GetFormattedDuration() string {
	formattedDuartion := fmt.Sprintf("%dm%ds", audioDetails.Duration/60, audioDetails.Duration%60)

	return formattedDuartion
}

func (audioDetails *AudioDetails) String() string {
	formattedAudioDetails := fmt.Sprintf("< song: '%s'  duration: %s >", audioDetails.Title, audioDetails.GetFormattedDuration())

	return formattedAudioDetails
}

type AudioState struct {
	AudioDetails
	currentTrackIndex int
	currentPos        int
	totalLength       int
}

func (audioState *AudioState) updateAudioState(audioDetails *AudioDetails) {
	log.Println("Updating Audio State...")
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
