package views

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestTerminalFrameFitsViewport(t *testing.T) {
	for _, size := range []struct {
		width  int
		height int
	}{
		{width: 40, height: 12},
		{width: 80, height: 24},
		{width: 160, height: 50},
	} {
		rendered := TerminalFrame("content", NewTerminal(size.width, size.height))

		if width := lipgloss.Width(rendered); width > size.width {
			t.Errorf("rendered width %d exceeds viewport width %d", width, size.width)
		}
		if height := lipgloss.Height(rendered); height > size.height {
			t.Errorf("rendered height %d exceeds viewport height %d", height, size.height)
		}
	}
}
