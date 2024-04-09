package tui

import (
	"fmt"
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
			if m.mode == interactiveListMode {
				doInterativeList(cmd, &m)
			} else {
				doCommand(cmd, &m)
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

	if m.mode == interactiveListMode {
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
