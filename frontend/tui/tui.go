package tui

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/johnrijoy/ludo-go/app"
	"golang.org/x/term"
)

func Run() {
	if err := app.Init(); err != nil {
		panic(err)
	}
	defer app.Close()

	isPiped = app.IsSourcePiped()

	p := tea.NewProgram(newMainModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

type mainModel struct {
	currentStatus    respStatus
	statusChan       chan respStatus
	cmdInput         textinput.Model
	cmdHist          []string
	cmdHistIndex     int
	resultMsg        string
	listTitle        string
	searchList       []string
	highlightIndices []int
	postSearchFunc   postIntList
	mode             imode
	help             viewport.Model
	err              error
	width, height    int
	quit             bool
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

	m.help = viewport.New(50, 50)
	m.help.SetContent(showHelp())
	m.help.KeyMap = viewport.DefaultKeyMap()
	m.help.Init()
	//m.help.HighPerformanceRendering = true

	return m
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("Ludo Go"), tea.EnterAltScreen, startActivity(m.statusChan), listenActivity(m.statusChan), resizeTicker)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.quit {
		return m, tea.Quit
	}

	// General events
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, nil
	case resizeTickMsg:
		w, h, _ := term.GetSize(int(os.Stdout.Fd()))
		// if w != m.width || h != m.height {
		// 	m.updateSize(w, h)
		// }
		return m, tea.Batch(resizeTicker, func() tea.Msg { return tea.WindowSizeMsg{Width: w, Height: h} })
	case respStatus:
		m.currentStatus = msg
		return m, listenActivity(m.statusChan)
	}

	// help mode
	if m.mode == helpMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type.String() {
			case "esc":
				setCommandMode(&m)
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.help, cmd = m.help.Update(msg)
		return m, cmd
	}

	m.cmdInput, cmd = m.cmdInput.Update(msg)
	// Player events
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+h":
			setHelpMode(&m)
			m.help, cmd = m.help.Update(msg)
			return m, cmd
		case "enter":
			inp := m.cmdInput.Value()
			if m.mode == interactiveListMode {
				doInterativeList(inp, &m)
			} else {
				doCommand(inp, &m)
			}
			m.cmdHist = append(m.cmdHist, inp)
			m.cmdHistIndex = len(m.cmdHist)
			m.cmdInput.Reset()

			return m, m.cmdInput.Focus()
		case "up":
			m.cmdHistIndex--
			if m.cmdHistIndex < 0 {
				m.cmdHistIndex = 0
			}
			if len(m.cmdHist) > 0 {
				m.cmdInput.SetValue(m.cmdHist[m.cmdHistIndex])
			}
		case "down":
			m.cmdHistIndex++
			if m.cmdHistIndex >= len(m.cmdHist) {
				m.cmdHistIndex = len(m.cmdHist)
			}
			if m.cmdHistIndex == len(m.cmdHist) {
				m.cmdInput.Reset()
			} else {
				val := m.cmdHist[m.cmdHistIndex]
				m.cmdInput.SetValue(val)
			}
		case "runes":
			m.err = nil
			m.resultMsg = ""
			return m, nil
		}
	}

	return m, cmd
}

func (m mainModel) View() string {
	s := ""

	if m.mode == helpMode {
		m.help.Height = m.height - 4
		m.help.Width = getBaseHorizontalWidth(&m)
		return lipgloss.JoinVertical(lipgloss.Left, m.getAppTitle(), baseStyle.Height(m.height-2).Width(m.width-2).Render(m.help.View()))
	}

	if m.mode == commandMode {
		if m.resultMsg != "" {
			s += fmt.Sprintf("\n%s\n", m.resultMsg)
		}
	}

	if m.mode == listMode || m.mode == interactiveListMode {
		s += fmt.Sprintln()

		dispWidth := getBaseHorizontalWidth(&m)
		for i, item := range m.searchList {
			itemView := safeTruncString(item, dispWidth)
			if sliceContains(m.highlightIndices, i) {
				itemView = Magenta(itemView)
			}
			s += itemView
			s += "\n"
		}

		// s += lipgloss.JoinHorizontal(lipgloss.Left, m.searchList...)
		// s += "\n"
	}

	if m.err != nil {
		if _, ok := m.err.(ErrWarn); ok {
			s += fmt.Sprintf("\n%s %s\n", Yellow("WARN:"), m.err.Error())
		} else {
			s += fmt.Sprintf("\n%s %s\n", Red("ERROR:"), m.err.Error())
		}
	}

	s = lipgloss.JoinVertical(lipgloss.Left, "\n", m.viewCurrentAudio(), m.cmdInput.View(), s)

	s = safeTrimHeight(s, m.height-4)
	// s = safeTrimWidth(s, m.width)

	//s += fmt.Sprintf("TextHeight: %d WindowHeight: %d WindowWidth: %d", height, m.height, m.width)
	return lipgloss.JoinVertical(lipgloss.Left, m.getAppTitle(), baseStyle.Height(m.height-2).Width(m.width-2).Render(s))
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

// view current song

func (m *mainModel) viewCurrentAudio() string {
	s := ""

	scale := 50
	dispWidth := getBaseHorizontalWidth(m)
	if dispWidth < scale && dispWidth > 0 {
		scale = dispWidth
	}

	currPos, totPos := m.currentStatus.pos, m.currentStatus.total

	navMsg := Gray(strings.Repeat("-", scale))
	if totPos > currPos {
		scaledCurrPos := int(math.Round((float64(currPos) / float64(totPos)) * float64(scale)))
		restPosition := scale - scaledCurrPos

		navMsg = Magenta(strings.Repeat(">", scaledCurrPos)) + Gray(strings.Repeat("-", restPosition))
	}

	audTitle := safeTruncString(m.currentStatus.audio.Title, 30)
	audUploader := safeTruncString(m.currentStatus.audio.Uploader, 20)

	s += fmt.Sprintf("%s%s%s\n",
		NoStyle.Width(scale*3/5).Render(m.currentStatus.mediaStatus.String()),
		NoStyle.Width(scale/5).AlignHorizontal(lipgloss.Right).Render(app.GetFormattedTime(currPos)),
		NoStyle.Width(scale/5).AlignHorizontal(lipgloss.Right).Render(app.GetFormattedTime(totPos)),
	)
	s += fmt.Sprintf("%s\n", navMsg)
	s += fmt.Sprintf("%s%s\n",
		Aqua.Width(scale*3/5).Render(audTitle),
		AquaD.Width(scale*2/5).AlignHorizontal(lipgloss.Right).Render(audUploader))

	// s += fmt.Sprintf("\n%s  | %s / %s\n", m.currentStatus.mediaStatus, app.GetFormattedTime(m.currentStatus.pos), app.GetFormattedTime(m.currentStatus.total))
	// s += fmt.Sprintf("%s\n", m.currentStatus.audio.Title)
	// s += fmt.Sprintf("%-20s %10s\n", m.currentStatus.audio.Uploader, m.currentStatus.audio.GetFormattedDuration())
	return s
}
