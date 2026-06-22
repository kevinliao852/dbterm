package models

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func ConnectionURIInput() textinput.Model {
	sti := textinput.New()
	sti.Placeholder = "DB_URI"
	sti.Focus()
	sti.CharLimit = 156
	sti.Width = 70

	return sti
}
