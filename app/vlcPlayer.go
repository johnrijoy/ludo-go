package app

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	vlc "github.com/adrg/libvlc-go/v3"
)

type VlcPlayer struct {
	player     *vlc.ListPlayer
	mediaList  *vlc.MediaList
	audioQueue []AudioDetails
	audioState AudioState
	eventIDs   EventIdList
	Quit       chan struct{}
}

type EventIdList struct {
	player     []vlc.EventID
	listPlayer []vlc.EventID
}

var vlcLog = log.New(io.Discard, "vlc :", log.LstdFlags)

// display information regarding libVlc version
func Info() string {
	versionInfo := vlc.Version()

	var sb strings.Builder

	sb.WriteString(fmt.Sprintln("libVlc Binding Version: ", versionInfo.String()))
	sb.WriteString(fmt.Sprint("Vlc Runtime Version: ", versionInfo.Changeset()))

	return sb.String()
}

// Creates and initialises a new vlc player
func (vlcPlayer *VlcPlayer) Init() error {
	err := vlc.Init("--no-video", "--quiet")
	if err != nil {
		return err
	}

	// Create a new list player.
	player, err := vlc.NewListPlayer()
	if err != nil {
		return err
	}
	vlcLog.Println("List Player created")

	mediaList, err := vlc.NewMediaList()
	if err != nil {
		return err
	}

	player.SetMediaList(mediaList)
	vlcLog.Println("MediaList created")

	vlcPlayer.mediaList = mediaList
	vlcPlayer.player = player
	vlcPlayer.audioQueue = make([]AudioDetails, 0)
	vlcPlayer.audioState = AudioState{}
	vlcPlayer.Quit = make(chan struct{})
	vlcPlayer.audioState.currentTrackIndex = -1

	vlcPlayer.attachEvents()

	return nil
}

// Stops and releases the creates vlc player
func (vlcPlayer *VlcPlayer) Close() {
	vlcLog.Println("VLC Player closing...")
	vlcPlayer.player.Stop()
	vlcPlayer.mediaList.Release()

	player, err := vlcPlayer.player.Player()
	if err == nil {
		// Retrieve player event manager.
		manager, err := player.EventManager()
		if err == nil {
			vlcLog.Println("player events detached")
			manager.Detach(vlcPlayer.eventIDs.player...)
		}
	}
	vlcLog.Println("Reached here")

	manager, err := vlcPlayer.player.EventManager()
	if err == nil {
		vlcLog.Println("List player event detached")
		manager.Detach(vlcPlayer.eventIDs.listPlayer...)
	} else {
		vlcLog.Println(err)
	}

	err = vlcPlayer.player.Release()
	if err != nil {
		vlcLog.Println(err)
	}
	vlcLog.Println("VLC Player closed")
}

func (vlcPlayer *VlcPlayer) ResetPlayer() error {
	vlcPlayer.Close()
	return vlcPlayer.Init()
}

//////////////////////
// Playback Control //
//////////////////////

func (vlcPlayer *VlcPlayer) StartPlayback() error {
	return vlcPlayer.player.Play()
}

func (vlcPlayer *VlcPlayer) StopPlayback() error {

	err := vlcPlayer.player.Stop()
	if err != nil {
		return err
	}

	vlcPlayer.audioState.currentTrackIndex = -1
	return nil
}

func (vlcPlayer *VlcPlayer) PauseResume() error {
	mediaState, err := vlcPlayer.getPlayerState()
	if err != nil {
		return err
	}

	if *mediaState != vlc.MediaEnded {
		return vlcPlayer.player.TogglePause()
	}
	return nil
}

func (vlcPlayer *VlcPlayer) SkipToNext() error {
	return vlcPlayer.player.PlayNext()
}

func (vlcPlayer *VlcPlayer) SkipToPrevious() error {
	err := vlcPlayer.player.PlayPrevious()
	if err != nil {
		return err
	}

	return vlcPlayer.updateCurrentMedia(vlcPlayer.audioState.currentTrackIndex - 1)
}

func (vlcPlayer *VlcPlayer) SkipToIndex(trackIndex int) error {
	if !vlcPlayer.validateTrackIndex(trackIndex) {
		vlcLog.Println("!! [mediaChangedCallback] invalid track index")
		return errors.New("invalid track index")
	}

	err := vlcPlayer.player.PlayAtIndex(uint(trackIndex))
	if err != nil {
		return err
	}

	return vlcPlayer.updateCurrentMedia(trackIndex)
}

func (vlcPlayer *VlcPlayer) ForwardBySeconds(duration int) error {
	if duration < 0 {
		return errors.New("negative duration")
	}

	player, err := vlcPlayer.player.Player()
	if err != nil {
		return err
	}

	totalTime, err := player.MediaLength()
	if err != nil {
		return err
	}

	currTime, err := player.MediaTime()
	if err != nil {
		return err
	}

	newTime := currTime + duration*1000
	if newTime >= totalTime {
		newTime = totalTime
	}
	return player.SetMediaTime(newTime)
}

func (vlcPlayer *VlcPlayer) RewindBySeconds(duration int) error {
	if duration < 0 {
		return errors.New("negative duration")
	}

	player, err := vlcPlayer.player.Player()
	if err != nil {
		return err
	}

	currTime, err := player.MediaTime()
	if err != nil {
		return err
	}

	newTime := currTime - duration*1000
	if newTime <= 0 {
		newTime = 0
	}

	return player.SetMediaTime(newTime)
}

////////////////////
// info functions //
////////////////////

func (vlcPlayer *VlcPlayer) IsPlaying() bool {
	return vlcPlayer.player.IsPlaying()
}

func (vlcPlayer *VlcPlayer) GetAudioState() AudioState {
	return vlcPlayer.audioState
}

func (vlcPlayer *VlcPlayer) GetQueueIndex() int {
	return vlcPlayer.audioState.currentTrackIndex
}

func (vlcPlayer *VlcPlayer) GetQueue() []AudioDetails {
	return vlcPlayer.audioQueue
}

func (vlcPlayer *VlcPlayer) FetchPlayerState() vlc.MediaState {
	vlcLog.Println("Getting player state")
	mediaState, err := vlcPlayer.player.MediaState()
	if err != nil {
		return 99
	}

	vlcLog.Printf("%v", mediaState)

	return mediaState
}

///////////////////
// media control //
///////////////////

func (vlcPlayer *VlcPlayer) AppendSong(audio *AudioDetails) error {
	mediaState, err := vlcPlayer.getPlayerState()
	if err != nil {
		return err
	}

	if *mediaState == vlc.MediaEnded {
		err = vlcPlayer.ResetPlayer()
		if err != nil {
			return err
		}
	}
	vlcPlayer.addSongToQueue(audio)

	return nil
}

////////////////////////
// Internal Functions //
////////////////////////

func (vlcPlayer *VlcPlayer) addSongToQueue(audio *AudioDetails) error {
	vlcPlayer.audioQueue = append(vlcPlayer.audioQueue, *audio)
	err := vlcPlayer.mediaList.AddMediaFromURL(audio.AudioStreamUrl)

	return err
}

func (vlcPlayer *VlcPlayer) updateCurrentMedia(trackIndex int) error {
	if !vlcPlayer.validateTrackIndex(trackIndex) {
		vlcLog.Println("!! [mediaChangedCallback] invalid track index")
		return errors.New("invalid track index")
	}

	vlcPlayer.audioState.currentTrackIndex = trackIndex
	audio := vlcPlayer.audioQueue[trackIndex]
	vlcPlayer.audioState.updateAudioState(&audio)

	return nil
}

// Deprecated: not used, instead use Close and Init a new vlc player
func (vlcPlayer *VlcPlayer) resetMediaList() error {
	vlcPlayer.mediaList.Release()

	mediaList, err := vlc.NewMediaList()
	if err != nil {
		return err
	}

	err = vlcPlayer.player.SetMediaList(mediaList)
	return err
}

func (vlcPlayer *VlcPlayer) getPlayerState() (*vlc.MediaState, error) {
	vlcLog.Println("Getting player state")
	mediaState, err := vlcPlayer.player.MediaState()
	if err != nil {
		return nil, err
	}

	vlcLog.Printf("%v", mediaState)

	return &mediaState, nil
}

func (vlcPlayer *VlcPlayer) validateTrackIndex(trackIndex int) bool {
	return trackIndex < 0 || trackIndex >= len(vlcPlayer.audioQueue)
}

func (vlcPlayer *VlcPlayer) attachEvents() error {

	mediaChangedCallback := func(event vlc.Event, userData interface{}) {
		vlcLog.Println("MediaChange Event")

		vlcPlayer, ok := userData.(*VlcPlayer)
		if !ok {
			vlcLog.Println("!! [mediaChangedCallback] no vlc data")
			return
		}

		vlcPlayer.audioState.currentTrackIndex += 1
		trackIndex := vlcPlayer.audioState.currentTrackIndex
		vlcLog.Println(trackIndex)
		if trackIndex < 0 || trackIndex > len(vlcPlayer.audioQueue) {
			vlcLog.Println("!! [mediaChangedCallback] invalid track index")
			return
		}

		vlcPlayer.audioState.updateAudioState(&vlcPlayer.audioQueue[trackIndex])
		vlcLog.Println(vlcPlayer.audioState.String())
	}

	positionChangedCallback := func(event vlc.Event, userData interface{}) {

		vlcPlayer, ok := userData.(*VlcPlayer)
		if !ok {
			vlcLog.Println("!! [positionChangedCallback] could not vlc user instance")
			return
		}

		vlcLog.Println("PositionChange Event")
		player, err := vlcPlayer.player.Player()
		if err != nil {
			vlcLog.Println("!! [positionChangedCallback] could not fetch player")
			return
		}

		currPos, err := player.MediaTime()
		if err != nil {
			vlcLog.Println("!! [positionChangedCallback] could not media curr time")
			return
		}
		totPos, err := player.MediaLength()
		if err != nil {
			vlcLog.Println("!! [positionChangedCallback] could not media total length")
			return
		}

		//vlcLog.Printf("currPoss: %d, totalPos: %d\n", currPos, totPos)

		vlcPlayer.audioState.currentPos = currPos / 1000
		vlcPlayer.audioState.totalLength = totPos / 1000

		//vlcLog.Println(vlcPlayer.audioState.String())
	}

	mediaListEndedCallback := func(event vlc.Event, userData interface{}) {
		vlcLog.Println("Media List Ended Event")

		vlcPlayer, ok := userData.(*VlcPlayer)
		if !ok {
			vlcLog.Println("!! [mediaListEndedCallback] no vlc data")
			return
		}

		close(vlcPlayer.Quit)
	}

	player, err := vlcPlayer.player.Player()
	if err != nil {
		return err
	}

	// Retrieve player event manager.
	manager, err := player.EventManager()
	if err != nil {
		return err
	}

	// Retrieve List Player event manager.
	manager2, err := vlcPlayer.player.EventManager()
	if err != nil {
		return err
	}

	eventID1, err := manager.Attach(vlc.MediaPlayerMediaChanged, mediaChangedCallback, vlcPlayer)
	if err != nil {
		return err
	}

	eventID2, err := manager.Attach(vlc.MediaPlayerTimeChanged, positionChangedCallback, vlcPlayer)
	if err != nil {
		return err
	}

	eventID3, err := manager2.Attach(vlc.MediaListPlayerPlayed, mediaListEndedCallback, vlcPlayer)
	if err != nil {
		return err
	}

	var playerEventID []vlc.EventID
	playerEventID = append(playerEventID, eventID1, eventID2)
	vlcPlayer.eventIDs.player = playerEventID

	var lPlayerEventID []vlc.EventID
	lPlayerEventID = append(lPlayerEventID, eventID3)
	vlcPlayer.eventIDs.listPlayer = lPlayerEventID

	return nil
}
