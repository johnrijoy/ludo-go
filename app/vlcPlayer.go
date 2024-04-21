package app

import (
	"errors"
	"fmt"
	"io"
	"log"

	vlc "github.com/adrg/libvlc-go/v3"
	uuid "github.com/satori/go.uuid"
)

var vlcPlayer VlcPlayer

type VlcPlayer struct {
	player       *vlc.ListPlayer
	mediaList    *vlc.MediaList
	audioQueue   []AudioDetails
	audioState   AudioState
	eventIDs     EventIdList
	isMediaError bool
}

type EventIdList struct {
	player     []vlc.EventID
	listPlayer []vlc.EventID
}

var vlcLog = log.New(io.Discard, "vlc: ", log.LstdFlags|log.Lmsgprefix)
var eventLog = log.New(io.Discard, "vlcEvent: ", log.LstdFlags|log.Lmsgprefix)

var playerStateMap = map[int]string{
	0: "Nothing Special",
	1: "Media Opening",
	2: "Media Buffering",
	3: "Media Playing",
	4: "Media Paused",
	5: "Media Stopped",
	6: "Media Ended",
	7: "Media Error",
}

func PlayerStateString(i int) (string, bool) {
	val, ok := playerStateMap[i]
	return val, ok
}

// display information regarding libVlc version
func Info() vlc.VersionInfo {
	return vlc.Version()
}

// Creates and initialises a new vlc player
func (vlcPlayer *VlcPlayer) InitPlayer() error {
	err := vlc.Init("--no-video", "--quiet")
	if err != nil {
		return err
	}

	vlc.SetAppName(fmt.Sprintf("%s v%s", "LudoGo", Version), "")

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
	vlcPlayer.isMediaError = false
	vlcPlayer.audioState.currentTrackIndex = -1

	return vlcPlayer.attachEvents()
}

// Stops and releases the creates vlc player
func (vlcPlayer *VlcPlayer) ClosePlayer() error {
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
		return err
	}
	vlcLog.Println("VLC Player closed")
	return nil
}

func (vlcPlayer *VlcPlayer) ResetPlayer() error {
	vlcPlayer.ClosePlayer()
	return vlcPlayer.InitPlayer()
}

//////////////////////
// Playback Control //
//////////////////////

func (vlcPlayer *VlcPlayer) StartPlayback() error {

	mediaState, err := vlcPlayer.getPlayerState()
	if err != nil {
		return err
	}
	trackIndex := vlcPlayer.audioState.currentTrackIndex
	vlcLog.Println("Current Index:", trackIndex)

	if trackIndex < 0 {
		trackIndex = 0
	}

	if *mediaState == vlc.MediaEnded {
		return vlcPlayer.player.PlayAtIndex(uint(trackIndex + 1))
	}

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

func (vlcPlayer *VlcPlayer) SetVol(vol int) error {

	if vol < 0 || vol > 100 {
		return errors.New("invalid volume input")
	}

	player, err := vlcPlayer.player.Player()
	if err != nil {
		return errors.Join(errors.New("error in accessing player"), err)
	}

	return player.SetVolume(vol)
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

func (vlcPlayer *VlcPlayer) FetchPlayerState() int {
	vlcLog.Println("Getting player state")
	mediaState, err := vlcPlayer.player.MediaState()
	if err != nil {
		return 99
	}

	vlcLog.Printf("%v", mediaState)

	return int(mediaState)
}

func (vlcPlayer *VlcPlayer) CheckMediaError() bool {
	return vlcPlayer.isMediaError
}

func (vlcPlayer *VlcPlayer) GetMediaPosition() (int, int) {
	player, err := vlcPlayer.player.Player()
	if err != nil {
		return 0, 0
	}
	currTime, err := player.MediaTime()
	if err != nil {
		return 0, 0
	}
	totalTime, err := player.MediaLength()
	if err != nil {
		return 0, 0
	}
	currTime = currTime / 1000
	totalTime = totalTime / 1000
	return currTime, totalTime
}

///////////////////
// media control //
///////////////////

func (vlcPlayer *VlcPlayer) AppendAudio(audio *AudioDetails) error {
	audio.uid = uuid.NewV1().String()
	vlcLog.Println("Audio UUID:", audio.uid)
	err := vlcPlayer.addSongToQueue(audio)
	return err
}

func (vlcPlayer *VlcPlayer) RemoveAudioFromIndex(removeIndex int) error {
	currIndex := vlcPlayer.audioState.currentTrackIndex
	queueLen, err := vlcPlayer.mediaList.Count()
	if err != nil {
		return err
	}

	if removeIndex <= currIndex || removeIndex > queueLen {
		errString := fmt.Sprintf("%d, %d, %d", currIndex, removeIndex, queueLen)
		return errors.New("Invalid remove index: " + errString)
	}

	if err = vlcPlayer.mediaList.Lock(); err != nil {
		return err
	}

	if err = vlcPlayer.mediaList.RemoveMediaAtIndex(uint(removeIndex)); err != nil {
		vlcPlayer.mediaList.Unlock()
		return err
	}

	vlcPlayer.mediaList.Unlock()

	vlcPlayer.audioQueue = append(vlcPlayer.audioQueue[:removeIndex], vlcPlayer.audioQueue[removeIndex+1:]...)

	return nil
}

func (vlcPlayer *VlcPlayer) RemoveLastAudio(removeIndex int) error {
	return vlcPlayer.RemoveAudioFromIndex(len(vlcPlayer.audioQueue) - 1)
}

func (vlcPlayer *VlcPlayer) RemoveAllAudioFromIndex(removeIndex int) error {
	currIndex := vlcPlayer.audioState.currentTrackIndex
	queueLen, err := vlcPlayer.mediaList.Count()
	if err != nil {
		return err
	}

	if removeIndex <= currIndex || removeIndex > queueLen {
		errString := fmt.Sprintf("%d, %d, %d", currIndex, removeIndex, queueLen)
		return errors.New("Invalid remove index: " + errString)
	}

	if err = vlcPlayer.mediaList.Lock(); err != nil {
		return err
	}

	i := removeIndex
	for ; i < queueLen; i++ {
		if err = vlcPlayer.mediaList.RemoveMediaAtIndex(uint(removeIndex)); err != nil {
			vlcPlayer.mediaList.Unlock()
			break
		}
	}

	vlcPlayer.mediaList.Unlock()

	vlcPlayer.audioQueue = append(vlcPlayer.audioQueue[:removeIndex], vlcPlayer.audioQueue[i:]...)

	return nil
}

func (vlcPlayer *VlcPlayer) SkipToNext() error {
	return vlcPlayer.player.PlayNext()
}

func (vlcPlayer *VlcPlayer) SkipToPrevious() error {
	return vlcPlayer.player.PlayPrevious()
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

////////////////////////
// Internal Functions //
////////////////////////

func (vlcPlayer *VlcPlayer) addSongToQueue(audio *AudioDetails) error {
	var media *vlc.Media
	mediaCreated := false

	if mediaPath, ok := audioCache.LookupCache(audio.AudioBasic); ok {
		vlcLog.Println("Playing Cached audio:", audio.Title, ",", mediaPath)
		newMedia, err := vlc.NewMediaFromPath(mediaPath)
		if err == nil {
			media = newMedia
			mediaCreated = true
		}
	}
	if !mediaCreated {
		newMedia, err := vlc.NewMediaFromURL(audio.AudioStreamUrl)
		if err != nil {
			return err
		}
		media = newMedia
		mediaCreated = true
	}

	if !mediaCreated {
		return errors.New("could not create Media")
	}

	err := media.SetUserData(audio.uid)
	if err != nil {
		return err
	}

	vlcPlayer.audioQueue = append(vlcPlayer.audioQueue, *audio)
	return vlcPlayer.mediaList.AddMedia(media)
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
	return trackIndex >= 0 && trackIndex < len(vlcPlayer.audioQueue)
}

func (vlcPlayer *VlcPlayer) attachEvents() error {

	mediaChangedCallback := func(event vlc.Event, userData interface{}) {
		eventLog.Println("MediaChange Event")

		vlcPlayer, ok := userData.(*VlcPlayer)
		if !ok {
			eventLog.Println("!! [mediaChangedCallback] no vlc data")
			return
		}

		player, err := vlcPlayer.player.Player()
		if err != nil {
			eventLog.Println("!! [mediaChangedCallback] error:", err)
			return
		}
		media, err := player.Media()
		if err != nil {
			eventLog.Println("!! [mediaChangedCallback] error:", err)
			return
		}
		uData, err := media.UserData()
		if err != nil {
			eventLog.Println("!! [mediaChangedCallback] error:", err)
			return
		}
		currUid, ok := uData.(string)
		if !ok {
			eventLog.Println("!! [mediaChangedCallback] error: could not convert to uid")
			return
		}
		eventLog.Println("Curr UID:", currUid)

		currInd := -1
		for i, aud := range vlcPlayer.audioQueue {
			if aud.uid == currUid {
				currInd = i
				break
			}
		}

		vlcPlayer.audioState.currentTrackIndex = currInd
		trackIndex := vlcPlayer.audioState.currentTrackIndex
		eventLog.Println(trackIndex)
		if trackIndex < 0 || trackIndex > len(vlcPlayer.audioQueue) {
			eventLog.Println("!! [mediaChangedCallback] invalid track index")
			return
		}

		vlcPlayer.audioState.updateAudioState(&vlcPlayer.audioQueue[trackIndex])
		eventLog.Println(vlcPlayer.audioState.String())

		err = audioDb.SaveOrIncrementAudioDoc(vlcPlayer.audioState.AudioBasic)

		if err != nil {
			eventLog.Println("!! [mediaChangedCallback] error in saving to db")
			eventLog.Println(err)
		}

		audioCache.CacheAudio(vlcPlayer.audioState.AudioDetails)
	}

	/*
		positionChangedCallback := func(event vlc.Event, userData interface{}) {

			vlcPlayer, ok := userData.(*VlcPlayer)
			if !ok {
				eventLog.Println("!! [positionChangedCallback] could not vlc user instance")
				return
			}

			//eventLog.Println("PositionChange Event")
			player, err := vlcPlayer.player.Player()
			if err != nil {
				eventLog.Println("!! [positionChangedCallback] could not fetch player")
				return
			}

			currPos, err := player.MediaTime()
			if err != nil {
				eventLog.Println("!! [positionChangedCallback] could not media curr time")
				return
			}
			totPos, err := player.MediaLength()
			if err != nil {
				eventLog.Println("!! [positionChangedCallback] could not media total length")
				return
			}

			//eventLog.Printf("currPoss: %d, totalPos: %d\n", currPos, totPos)

			vlcPlayer.audioState.currentPos = currPos / 1000
			vlcPlayer.audioState.totalLength = totPos / 1000

			//vlcLog.Println(vlcPlayer.audioState.String())

			if vlcPlayer.isMediaError {
				vlcPlayer.isMediaError = false
			}
		}
	*/

	encounteredErrorCallback := func(event vlc.Event, userData interface{}) {
		eventLog.Println("List player encountered error")

		vlcPlayer, ok := userData.(*VlcPlayer)
		if !ok {
			eventLog.Println("!! [encounteredErrorCallback] could not vlc user instance")
			return
		}

		mediaState, err := vlcPlayer.player.MediaState()
		if err != nil {
			eventLog.Println("!! [encounteredErrorCallback] error: ", err)
			return
		}

		eventLog.Println("MediaState: ", playerStateMap[int(mediaState)])

		vlcPlayer.isMediaError = true
	}

	player, err := vlcPlayer.player.Player()
	if err != nil {
		return err
	}

	// Retrieve player event manager.
	listPlayerMan, err := player.EventManager()
	if err != nil {
		return err
	}

	// Retrieve List Player event manager. temporarily player event manager
	_, err = vlcPlayer.player.EventManager()
	if err != nil {
		return err
	}

	eventID1, err := listPlayerMan.Attach(vlc.MediaPlayerMediaChanged, mediaChangedCallback, vlcPlayer)
	if err != nil {
		return err
	}

	/*
		eventID2, err := listPlayerMan.Attach(vlc.MediaPlayerPositionChanged, positionChangedCallback, vlcPlayer)
		if err != nil {
			return err
		}
	*/

	eventID3, err := listPlayerMan.Attach(vlc.MediaPlayerEncounteredError, encounteredErrorCallback, vlcPlayer)
	if err != nil {
		return err
	}

	playerEventID := []vlc.EventID{eventID1, eventID3}
	vlcPlayer.eventIDs.player = playerEventID

	lPlayerEventID := []vlc.EventID{}
	vlcPlayer.eventIDs.listPlayer = lPlayerEventID

	return nil
}
