package models

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevinliao852/dbterm/pkg/views"
)

func DBSelectTable() table.Model {
	tr := []table.Row{}
	tc := []table.Column{}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(views.BorderColor).
		BorderBottom(true).
		Bold(true).
		Foreground(views.PrimaryColor)
	s.Cell = s.Cell.
		Foreground(views.TextColor).
		Padding(0, 1)
	s.Selected = s.Selected.
		Foreground(views.AccentColor).
		Bold(true)

	t := table.New(
		table.WithColumns(tc),
		table.WithRows(tr),
		table.WithFocused(true),
		table.WithHeight(5),
	)
	t.SetStyles(s)

	return t
}

func DBSQLQueryInput() textarea.Model {
	return composer("Write a SQL query…")
}

func DBNaturalLanguageInput() textarea.Model {
	return composer("Ask a question about your data…")
}

func composer(placeholder string) textarea.Model {
	dbi := textarea.New()
	dbi.Placeholder = placeholder
	dbi.Prompt = "  "
	dbi.ShowLineNumbers = false
	dbi.EndOfBufferCharacter = ' '
	dbi.Focus()
	dbi.CharLimit = 4000
	dbi.SetWidth(60)
	dbi.SetHeight(4)
	dbi.FocusedStyle.Base = views.BodyStyle
	dbi.FocusedStyle.Prompt = views.InputPromptStyle
	dbi.FocusedStyle.Text = views.InputTextStyle
	dbi.FocusedStyle.Placeholder = views.InputPlaceholderStyle
	dbi.FocusedStyle.CursorLine = views.BodyStyle
	dbi.BlurredStyle = dbi.FocusedStyle
	dbi.Cursor.Style = views.CursorStyle

	return dbi
}
