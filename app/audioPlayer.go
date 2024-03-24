package app

type audioPlayer interface {
	// struct control
	InitPlayer()
	ClosePlayer()
	ResetPlayer()

	// media control
	AppendAudio(AudioDetails) error
	SkipToNext() error
	SkipToPrevious() error
	SkipToIndex(trackIndex int) error
	RemoveAudioFromIndex(int) error
	RemoveLastAudio()

	// audio playback control
	startPlayback()
	stopPlayback()
	PauseResume()
	ForwardBySeconds(int) error
	RewindBySeconds(int) error
	SetVol(int) error

	// Info Functions
	IsPlaying() bool
	GetQueue() []AudioDetails
	GetQueueIndex() int
}
