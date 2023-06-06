package main

type audioPlayer interface {
	// struct control
	init()
	closePlayer()

	// media control
	appendAudio(AudioDetails)
	skipAudio()
	removeAudioFromIndex(int)
	removeLastAudio()

	// audio playback control
	startPlayback()
	stopPlayback()
	pauseResume()
	setVol()
	endPlayer()
}
