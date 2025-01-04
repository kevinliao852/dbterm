package models

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func ConnectionTypeInput() textinput.Model {
	ti := textinput.New()
	ti.CharLimit = 156
	ti.Focus()
	ti.Placeholder = "the type of database(mysql, postgres)"

	return ti
}

func ConnectionURIInput() textinput.Model {
	sti := textinput.New()
	sti.Placeholder = "DB_URI"
	sti.Focus()
	sti.Width = 100

	return sti
}
