package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/johnrijoy/ludo-go/app"
)

func Run() {
	if err := app.Init(); err != nil {
		panic(err)
	}
	defer app.Close()

	p := tea.NewProgram(newMainModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

type mainModel struct {
	cmdInput       textinput.Model
	currentStatus  respStatus
	statusChan     chan respStatus
	resultMsg      string
	listTitle      string
	searchList     []string
	postSearchFunc postIntList
	mode           imode
	err            error
	width, height  int
	quit           bool
}

func newMainModel() mainModel {
	m := mainModel{}

	m.cmdInput = textinput.New()
	m.cmdInput.Prompt = commandPrompt
	m.cmdInput.Focus()
	m.cmdInput.CharLimit = 200
	m.cmdInput.Width = 50

	m.statusChan = make(chan respStatus)

	m.currentStatus = respStatus{mediaStatus: nothing, total: 0}

	m.mode = commandMode
	m.quit = false
	return m
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, startActivity(m.statusChan), listenActivity(m.statusChan))
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.quit {
		return m, tea.Quit
	}
	m.cmdInput, cmd = m.cmdInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			inp := m.cmdInput.Value()
			if m.mode == interactiveListMode {
				doInterativeList(inp, &m)
			} else {
				doCommand(inp, &m)
			}
			m.cmdInput.Reset()

			return m, m.cmdInput.Focus()
		case "runes":
			m.err = nil
			m.resultMsg = ""
			return m, nil
		}
	case respStatus:
		m.currentStatus = msg
		return m, listenActivity(m.statusChan)
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, nil
	}

	return m, cmd
}

func (m mainModel) View() string {
	s := "LudoGo\n"
	s += fmt.Sprintf("\n%s | %d / %d\n", m.currentStatus.mediaStatus, m.currentStatus.pos, m.currentStatus.total)
	s += fmt.Sprintf("%s\n", &m.currentStatus.audio)
	s += fmt.Sprintf("\n%s\n", m.cmdInput.View())

	if m.mode == commandMode {
		if m.resultMsg != "" {
			s += fmt.Sprintf("\n%s\n", m.resultMsg)
		}
	}

	if m.mode == listMode || m.mode == interactiveListMode {
		s += fmt.Sprintln()
		dispWidth := getBaseHorizontalWidth(&m)
		for i, item := range m.searchList {
			s += safeTruncString(fmt.Sprintf("%-2d - %s", i+1, item), dispWidth)
			s += "\n"
		}
	}

	if m.err != nil {
		s += fmt.Sprintf("\nError: %s\n", m.err.Error())
	}

	s = safeTrimHeight(s, m.height)

	//s += fmt.Sprintf("TextHeight: %d WindowHeight: %d WindowWidth: %d", height, m.height, m.width)
	return baseStyle.Render(s)
}

func startActivity(status chan respStatus) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Second)
			stat := app.MediaPlayer().FetchPlayerState()
			curr, pos := app.MediaPlayer().GetMediaPosition()
			aud := app.MediaPlayer().GetAudioState().AudioBasic
			status <- respStatus{pos: curr, total: pos, mediaStatus: mediaStat(stat), audio: aud}
		}
	}
}

func listenActivity(status chan respStatus) tea.Cmd {
	return func() tea.Msg {
		return <-status
	}
}
