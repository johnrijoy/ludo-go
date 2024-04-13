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
var Magenta = lipgloss.NewStyle().Foreground(lipgloss.Color("#8c0175")).Render
var Gray = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Render
var Yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#f7f705")).Render
var Cyan = lipgloss.NewStyle().Foreground(lipgloss.Color("#05f3f7")).Render

var borderType = lipgloss.RoundedBorder()
var baseStyle = lipgloss.NewStyle().Margin(0).Padding(0, 2, 1).Border(borderType, false, true, true)

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
	s += appTitle
	s += strings.Repeat(borderType.Top, rightFillerWidth) + borderType.TopRight
	return s
}
