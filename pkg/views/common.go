package views

import (
	"fmt"
	"strconv"

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
	termstyle := lipgloss.NewStyle().Align(lipgloss.Center).BorderStyle(lipgloss.ThickBorder())
	termstyle.Margin(1, 1, 1, 1)
	termstyle.Width(m.width - 5)
	termstyle.Height(m.height - 5)

	return termstyle.Render(fmt.Sprintf(
		"Enter the input:\n\nwidth:%s height:%s\n\n%s\n\n%s",
		strconv.Itoa(m.width),
		strconv.Itoa(m.height),
		innerStr,
		"(press ESC or CRL+C to quit)",
	))
}
