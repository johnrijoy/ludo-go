package tui

import "github.com/charmbracelet/lipgloss"

var Green = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render
var Blue = lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff")).Render
var Red = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render
var GreenH = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff73")).Render
var GreenD = lipgloss.NewStyle().Foreground(lipgloss.Color("#00bf30")).Render
var Magenta = lipgloss.NewStyle().Foreground(lipgloss.Color("#8c0175")).Render
var Gray = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Render
var Yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#f7f705")).Render
var Cyan = lipgloss.NewStyle().Foreground(lipgloss.Color("#05f3f7")).Render

var baseStyle = lipgloss.NewStyle().Padding(0, 2, 1).Border(lipgloss.RoundedBorder())

func getBaseHorizontalWidth(m *mainModel) int {
	return m.width - baseStyle.GetHorizontalBorderSize() - baseStyle.GetHorizontalPadding() - baseStyle.GetHorizontalMargins()
}
