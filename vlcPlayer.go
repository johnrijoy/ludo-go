package main

import (
	"log"

	vlc "github.com/adrg/libvlc-go/v3"
)

type VlcPlayer struct {
	player     *vlc.ListPlayer
	mediaList  *vlc.MediaList
	audioQueue []AudioDetails
	audioState AudioState
	eventIDs   EventIdList
	quit       chan struct{}
}

type EventIdList struct {
	player     []vlc.EventID
	listPlayer []vlc.EventID
}

func (vlcPlayer *VlcPlayer) init() error {
	err := vlc.Init("--no-video", "--quiet")
	if err != nil {
		return err
	}

	// Create a new list player.
	player, err := vlc.NewListPlayer()
	if err != nil {
		return err
	}
	log.Println("List Player created")

	mediaList, err := vlc.NewMediaList()
	if err != nil {
		return err
	}

	player.SetMediaList(mediaList)
	log.Println("MediaList created")

	vlcPlayer.mediaList = mediaList
	vlcPlayer.player = player
	vlcPlayer.audioState.currentTrackIndex = -1
	vlcPlayer.quit = make(chan struct{})

	vlcPlayer.attachEvents()

	return nil
}

func (vlcPlayer *VlcPlayer) attachEvents() error {

	mediaChangedCallback := func(event vlc.Event, userData interface{}) {
		log.Println("MediaChange Event")

		vlcPlayer, ok := userData.(*VlcPlayer)
		if !ok {
			log.Println("!! [mediaChangedCallback] no vlc data")
			return
		}

		vlcPlayer.audioState.currentTrackIndex += 1
		trackIndex := vlcPlayer.audioState.currentTrackIndex
		log.Println(trackIndex)
		if trackIndex < 0 || trackIndex > len(vlcPlayer.audioQueue) {
			log.Println("!! [mediaChangedCallback] invalid track index")
			return
		}

		vlcPlayer.audioState.updateAudioState(&vlcPlayer.audioQueue[trackIndex])
		log.Println(vlcPlayer.audioState.String())
	}

	positionChangedCallback := func(event vlc.Event, userData interface{}) {

		vlcPlayer, ok := userData.(*VlcPlayer)
		if !ok {
			log.Println("!! [positionChangedCallback] could not vlc user instance")
			return
		}

		log.Println("PositionChange Event")
		player, err := vlcPlayer.player.Player()
		if err != nil {
			log.Println("!! [positionChangedCallback] could not fetch player")
			return
		}

		currPos, err := player.MediaTime()
		if err != nil {
			log.Println("!! [positionChangedCallback] could not media curr time")
			return
		}
		totPos, err := player.MediaLength()
		if err != nil {
			log.Println("!! [positionChangedCallback] could not media total length")
			return
		}

		vlcPlayer.audioState.currentPos = currPos / 1000
		vlcPlayer.audioState.totalLength = totPos / 1000

		log.Println(vlcPlayer.audioState.String())
	}

	mediaListEndedCallback := func(event vlc.Event, userData interface{}) {
		log.Println("Media List Ended Event")

		vlcPlayer, ok := userData.(*VlcPlayer)
		if !ok {
			log.Println("!! [mediaListEndedCallback] no vlc data")
			return
		}

		close(vlcPlayer.quit)
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

func (vlcPlayer *VlcPlayer) close() {
	log.Println("VLC Player closing...")
	vlcPlayer.player.Stop()
	vlcPlayer.mediaList.Release()

	player, err := vlcPlayer.player.Player()
	if err == nil {
		// Retrieve player event manager.
		manager, err := player.EventManager()
		if err == nil {
			log.Println("player events detached")
			manager.Detach(vlcPlayer.eventIDs.player...)
		}
	}
	log.Println("Reached here")

	manager, err := vlcPlayer.player.EventManager()
	if err == nil {
		log.Println("List player event detached")
		manager.Detach(vlcPlayer.eventIDs.listPlayer...)
	} else {
		log.Println(err)
	}

	err = vlcPlayer.player.Release()
	if err != nil {
		log.Println(err)
	}
	log.Println("VLC Player closed")
}

// Playback Control
func (vlcPlayer *VlcPlayer) startPlayback() error {
	return vlcPlayer.player.Play()
}

func (vlcPlayer *VlcPlayer) stopPlayback() error {

	err := vlcPlayer.player.Stop()
	if err != nil {
		return err
	}

	vlcPlayer.audioState.currentTrackIndex = -1
	return nil
}

func (vlcPlayer *VlcPlayer) pauseResume() error {
	return vlcPlayer.player.TogglePause()
}

// info functions
func (vlcPlayer *VlcPlayer) isPlaying() bool {
	return vlcPlayer.player.IsPlaying()
}

func (vlcPlayer *VlcPlayer) getAudioState() AudioState {
	return vlcPlayer.audioState
}

// media control
func (vlcPlayer *VlcPlayer) appendSong(audio *AudioDetails) error {
	mediaState, err := vlcPlayer.getPlayerState()
	if err != nil {
		return err
	}

	if *mediaState == vlc.MediaEnded {
		vlcPlayer.resetMediaList()
		vlcPlayer.addSongToQueue(audio)
		vlcPlayer.startPlayback()
	} else {
		vlcPlayer.addSongToQueue(audio)
	}

	return nil
}

// Internal Functions
func (vlcPlayer *VlcPlayer) addSongToQueue(audio *AudioDetails) error {
	vlcPlayer.audioQueue = append(vlcPlayer.audioQueue, *audio)
	err := vlcPlayer.mediaList.AddMediaFromURL(audio.audioStreamUrl)

	return err
}

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
	log.Println("Getting player state")
	mediaState, err := vlcPlayer.player.MediaState()
	if err != nil {
		return nil, err
	}

	log.Printf("%v", mediaState)

	return &mediaState, nil
}
