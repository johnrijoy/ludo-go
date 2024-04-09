package tui

import "strings"

type respStatus struct {
	mediaStatus mediaStat
	pos         int
	total       int
}

type imode uint8

const (
	commandMode imode = iota
	listMode
	interactiveListMode
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

// helpers
func parseCommand(command string) (string, string) {
	cmdList := strings.Split(command, " ")
	command = cmdList[0]
	arg := ""
	if len(cmdList) > 1 {
		arg = strings.Join(cmdList[1:], " ")
	}

	return command, arg
}
