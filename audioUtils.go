package main

import "fmt"

type AudioDetails struct {
	ytId           string
	title          string
	uploader       string
	duration       int
	audioStreamUrl string
}

func (audioDetails *AudioDetails) getFormattedDuration() string {
	formattedDuartion := fmt.Sprintf("%dm%ds", audioDetails.duration/60, audioDetails.duration%60)

	return formattedDuartion
}

func (audioDetails *AudioDetails) String() string {
	formattedAudioDetails := fmt.Sprintf("< song: '%s'  duration: %s >", audioDetails.title, audioDetails.getFormattedDuration())

	return formattedAudioDetails
}

type AudioState struct {
	AudioDetails
	currentTrackIndex int
	currentPos        int
	totalLength       int
}

func (audioState *AudioState) updateAudioState(audioDetails *AudioDetails) {
	audioState.ytId = audioDetails.ytId
	audioState.title = audioDetails.title
	audioState.duration = audioDetails.duration
	audioState.audioStreamUrl = audioDetails.audioStreamUrl
}

func (audioState *AudioState) String() string {
	formattedAudioDetails := fmt.Sprintf("< song: '%s'  duration: %s >", audioState.title, audioState.getFormattedDuration())

	return formattedAudioDetails
}
