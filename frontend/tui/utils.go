package tui

import (
	"strings"

	"github.com/johnrijoy/ludo-go/app"
)

const (
	commandPrompt = ">> "
	commandWidth  = 50
)

type respStatus struct {
	audio       app.AudioBasic
	mediaStatus mediaStat
	pos         int
	total       int
}

// common constants
const defaultForwardRewind = 10

// imode
type imode uint8

const (
	commandMode imode = iota
	listMode
	interactiveListMode
)

// media stat
type mediaStat uint8

const (
	nothing mediaStat = iota
	opening
	buffering
	playing
	paused
	stopped
	ended
	mediaErr
	invalid = 99
)

var mediaStatMap = map[mediaStat]string{
	nothing: "○ Nothing", opening: "opening", buffering: "buffering",
	playing: "▶ Now Playing", paused: "▌▌Paused", stopped: "■ Stopped",
	ended: "Ended", mediaErr: "⚠ Media Error",
}

func (m mediaStat) String() string {
	val, ok := mediaStatMap[m]
	if !ok {
		panic("Invalid constant")
	}
	return val
}

// Warn error

type ErrWarn struct {
	msg string
}

func Warn(msg string) ErrWarn {
	return ErrWarn{msg: msg}
}

func (w ErrWarn) Error() string {
	return w.msg
}

// post Interactive List func type
type postIntList func(index int, m *mainModel)

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

func setInteractiveListMode(m *mainModel, prompt string) {
	m.mode = interactiveListMode
	m.cmdInput.Prompt = prompt
	m.cmdInput.Width = 5
}

func setListMode(m *mainModel) {
	m.mode = listMode
}

func setCommandMode(m *mainModel) {
	m.mode = commandMode
	m.cmdInput.Prompt = commandPrompt
	m.cmdInput.Width = commandWidth
}

func handleErr(err error, m *mainModel) bool {
	if err != nil {
		m.err = err
		return true
	}
	return false
}

func safeTruncString(label string, max int) string {
	var result string

	if len(label) <= max {
		result = label
	} else {
		result = label[0:(max-3)] + "..."
	}

	return result
}

func safeTrimHeight(display string, termHeight int) string {
	textHeight := strings.Count(display, "\n")
	if textHeight > termHeight {
		display = strings.Join(strings.Split(display, "\n")[0:termHeight], "\n")
	}
	return display
}

func safeTrimWidth(display string, termWidth int) string {
	dispList := strings.Split(display, "\n")
	newDispList := make([]string, len(dispList))
	for i, line := range dispList {

		if len(line) > termWidth {
			line = safeTruncString(line, termWidth)
		}
		newDispList[i] = line
	}
	return strings.Join(newDispList, "\n")
}

func safeTrimView(display string, termWidth, termHeight int) string {
	display = safeTrimHeight(display, termHeight)
	display = safeTrimWidth(display, termWidth)
	return display
}
