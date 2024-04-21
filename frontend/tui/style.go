package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var Green = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render
var Blue = lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff")).Render
var Red = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render

var GreenH = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff73")).Render
var GreenD = lipgloss.NewStyle().Foreground(lipgloss.Color("#00bf30")).Render
var Magenta = lipgloss.NewStyle().Foreground(lipgloss.Color("#9D44C0")).Render
var Pink = lipgloss.NewStyle().Foreground(lipgloss.Color("#EC53B0")).Render
var Gray = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Render
var Yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#f7f705")).Render

var BlueD = lipgloss.NewStyle().Foreground(lipgloss.Color("#4D2DB7"))
var Cyan = lipgloss.NewStyle().Foreground(lipgloss.Color("#05f3f7"))
var AquaD = lipgloss.NewStyle().Foreground(lipgloss.Color("#5FBDFF"))
var Aqua = lipgloss.NewStyle().Foreground(lipgloss.Color("#96EFFF"))
var Purple = lipgloss.NewStyle().Foreground(lipgloss.Color("#7B66FF"))
var NoStyle = lipgloss.NewStyle()

var borderType = lipgloss.RoundedBorder()
var baseStyle = lipgloss.NewStyle().Margin(0).Padding(0, 2, 1).Border(borderType, false, true, true)
var titleStyle = lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("#EC53B0")).Italic(true).Bold(true).Render

func getBaseHorizontalWidth(m *mainModel) int {
	return m.width - baseStyle.GetHorizontalBorderSize() - baseStyle.GetHorizontalPadding() - baseStyle.GetHorizontalMargins()
}

func (m *mainModel) getAppTitle() string {
	s := ""
	termWidth := m.width
	titleWidth := lipgloss.Width(appTitle)
	cornerWidth := lipgloss.Width(borderType.TopLeft)
	fillerWidth := termWidth - (2 * cornerWidth) - titleWidth

	if fillerWidth < 0 {
		s += borderType.TopLeft + appTitle + borderType.TopRight
		return s
	}

	leftFillerWidth := fillerWidth / 2
	rightFillerWidth := leftFillerWidth
	if fillerWidth%2 != 0 {
		rightFillerWidth += 1
	}
	s += borderType.TopLeft + strings.Repeat(borderType.Top, leftFillerWidth)
	s += titleStyle(appTitle)
	s += strings.Repeat(borderType.Top, rightFillerWidth) + borderType.TopRight
	return s
}
