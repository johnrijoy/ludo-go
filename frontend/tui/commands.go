package tui

import (
	"errors"
	"fmt"
	"strconv"
)

func doCommand(cmd string, m *mainModel) {
	m.mode = commandMode
	switch cmd {
	case "play":
		m.resultMsg = "playing song 1"
	case "nothing":
		m.err = errors.New("not supported")
	case "search":
		m.mode = interactiveListMode
		m.cmdInput.Prompt = "> Please select a song: "
		m.searchList = []string{"item1", "item2", "item3"}
	case "show":
		m.mode = listMode
		m.searchList = []string{"item4", "item5", "item6"}
	default:
		m.err = errors.New("invalid command")
	}
}

func doInterativeList(ind string, m *mainModel) {
	defer func() { m.cmdInput.Prompt = ">> " }()

	i, err := strconv.Atoi(ind)
	if err != nil {
		m.err = errors.New("invalid index")
		m.mode = commandMode
		return
	}

	i--

	if i < 0 || i >= len(m.searchList) {
		m.err = errors.New("index out of bounds")
		m.mode = commandMode
		return
	}

	m.resultMsg = fmt.Sprintf("playing song < %s >", m.searchList[i])
	m.mode = commandMode
}
