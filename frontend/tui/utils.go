package tui

type respStatus struct {
	mediaStatus mediaStat
	pos         int
	total       int
}

type imode uint8

const (
	commandMode imode = iota
	listMode
	searchMode
)

type mediaStat uint8

const (
	nothing mediaStat = iota
	playing
	paused
	stopped
	mediaErr
)

var mediaStatMap = map[mediaStat]string{
	nothing: "○ Nothing", playing: "▶ Now Playing", paused: "▌▌Paused", stopped: "■ Stopped", mediaErr: "⚠ Media Error",
}

func (m mediaStat) String() string {
	val, ok := mediaStatMap[m]
	if !ok {
		panic("Invalid constant")
	}
	return val
}
