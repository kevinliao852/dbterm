package views

import (
	"github.com/charmbracelet/lipgloss"
)

type Terminal struct {
	width  int
	height int
}

func NewTerminal(width, height int) Terminal {
	return Terminal{
		width:  width,
		height: height,
	}
}

func TerminalFrame(innerStr string, m Terminal) string {
	if m.width <= 0 || m.height <= 0 {
		return innerStr
	}

	width := max(1, m.width-4)
	height := max(1, m.height-2)
	termstyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(0, 1).
		Margin(0, 1).
		Width(width).
		Height(height)

	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		AppTitleStyle.Render("DBTerm"),
		MutedStyle.Render("  database workspace"),
	)

	return termstyle.Render(header + "\n\n" + innerStr)
}
