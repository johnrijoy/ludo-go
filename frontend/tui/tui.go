package tui

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func Run() {
	p := tea.NewProgram(newMainModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

type mainModel struct {
	cmdInput      textinput.Model
	currentStatus respStatus
	statusChan    chan respStatus
	resultMsg     string
	searchList    []string
	mode          imode
	err           error
}

func newMainModel() mainModel {
	m := mainModel{}

	m.cmdInput = textinput.New()
	m.cmdInput.Prompt = ">> "
	m.cmdInput.Focus()
	m.cmdInput.CharLimit = 200
	m.cmdInput.Width = 50

	m.statusChan = make(chan respStatus)

	m.currentStatus = respStatus{mediaStatus: nothing, total: 100}

	m.mode = commandMode
	return m
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(startActivity(m.statusChan), listenActivity(m.statusChan))
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.cmdInput, cmd = m.cmdInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			cmd := m.cmdInput.Value()
			if m.mode == searchMode {
				parseSearch(cmd, &m)
			} else {
				parseCommand(cmd, &m)
			}
			m.cmdInput.Reset()
			return m, nil
		case "runes":
			m.err = nil
			return m, nil
		}
	case respStatus:
		m.currentStatus = msg
		return m, listenActivity(m.statusChan)
	}
	return m, cmd
}

func (m mainModel) View() string {
	s := "LudoGo\n"
	s += fmt.Sprintf("%s | %d / %d\n", m.currentStatus.mediaStatus, m.currentStatus.pos, m.currentStatus.total)
	s += fmt.Sprintf("\n%s\n", m.cmdInput.View())

	if m.mode == commandMode {
		if m.resultMsg != "" {
			s += fmt.Sprintf("\n%s\n", m.resultMsg)
		}
	}

	if m.mode == listMode {
		s += fmt.Sprintln()
		for i, item := range m.searchList {
			s += fmt.Sprintf("%d - %s\n", i+1, item)
		}
	}

	if m.mode == searchMode {
		s += fmt.Sprintln()
		for i, item := range m.searchList {
			s += fmt.Sprintf("%d - %s\n", i+1, item)
		}
	}

	if m.err != nil {
		s += fmt.Sprintf("\nError: %s\n", m.err.Error())
	}

	return s
}

func parseCommand(cmd string, m *mainModel) {
	m.mode = commandMode
	switch cmd {
	case "play":
		m.resultMsg = "playing song 1"
	case "nothing":
		m.err = errors.New("not supported")
	case "search":
		m.mode = searchMode
		m.cmdInput.Prompt = "> Please select a song: "
		m.searchList = []string{"item1", "item2", "item3"}
	case "show":
		m.mode = listMode
		m.searchList = []string{"item4", "item5", "item6"}
	default:
		m.err = errors.New("invalid command")
	}
}

func parseSearch(ind string, m *mainModel) {
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

func startActivity(status chan respStatus) tea.Cmd {
	return func() tea.Msg {
		i := 0
		for {
			stat := playing
			time.Sleep(time.Second)
			if i%2 == 0 {
				stat = paused
			}
			if i%5 == 0 {
				stat = stopped
			}

			if i > 10 {
				stat = mediaErr
			}
			status <- respStatus{pos: i, total: 100, mediaStatus: stat}
			i++
		}
	}
}

func listenActivity(status chan respStatus) tea.Cmd {
	return func() tea.Msg {
		return <-status
	}
}
