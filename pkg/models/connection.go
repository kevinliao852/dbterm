package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func ConnectionTypeInput() textinput.Model {
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = "the type of database(mysql, postgres)"
	ti.CharLimit = 156
	ti.Width = 100
	ti.TextStyle.Background(lipgloss.Color("63"))
	ti.TextStyle.Foreground(lipgloss.Color("63"))
	ti.SetSuggestions([]string{"mysql", "postgres"})
	ti.ShowSuggestions = true
	ti.PromptStyle.Border(lipgloss.NormalBorder())
	ti.TextStyle.BorderStyle(lipgloss.NormalBorder())

	return ti
}

func ConnectionURIInput() textinput.Model {
	sti := textinput.New()
	sti.Placeholder = "DB_URI"
	sti.Focus()
	sti.CharLimit = 156
	sti.Width = 50

	return sti
}
