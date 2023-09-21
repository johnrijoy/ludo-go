package main

import (
	"fmt"
	"log"
)

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
	log.Println("Updating Audio State...")
	audioState.ytId = audioDetails.ytId
	audioState.title = audioDetails.title
	audioState.duration = audioDetails.duration
	audioState.audioStreamUrl = audioDetails.audioStreamUrl
}

func (audioState *AudioState) getFormattedPosition() string {
	formattedPos := fmt.Sprintf("%dm%ds", audioState.currentPos/60, audioState.currentPos%60)

	return formattedPos
}

func (audioState *AudioState) String() string {
	formattedAudioDetails := fmt.Sprintf("< song: '%s' Pos: %s  duration: %s >", audioState.title, audioState.getFormattedPosition(), audioState.getFormattedDuration())

	return formattedAudioDetails
}
