package models

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func DBSelectTable() table.Model {
	tr := []table.Row{}
	tc := []table.Column{}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t := table.New(
		table.WithColumns(tc),
		table.WithRows(tr),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	t.SetStyles(s)
	t.SetWidth(1000)

	return t
}

func DBSQLQueryInput() textinput.Model {
	dbi := textinput.New()
	dbi.Placeholder = "SQL Query"
	dbi.Focus()
	dbi.CharLimit = 156
	dbi.Width = 40

	return dbi
}
