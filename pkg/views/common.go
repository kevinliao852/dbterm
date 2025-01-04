package views

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	marginRight  = lipgloss.Right
	marginLeft   = lipgloss.Left
	marginTop    = lipgloss.Top
	marginBottom = lipgloss.Bottom
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
	termstyle.Margin(int(marginTop), int(marginRight), int(marginBottom), int(marginTop))
	termstyle.Width(m.width - 5)
	termstyle.Height(m.height - 5)
	text := lipgloss.JoinVertical(0.5, termstyle.Render(innerStr), "press ESC or CRL+C to quit")

	return text
}
