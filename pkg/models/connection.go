package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/kevinliao852/dbterm/pkg/views"
)

func ConnectionURIInput() textinput.Model {
	sti := textinput.New()
	sti.Placeholder = "DB_URI"
	sti.Prompt = "› "
	sti.Focus()
	sti.CharLimit = 156
	sti.Width = 70
	sti.PromptStyle = views.InputPromptStyle
	sti.TextStyle = views.InputTextStyle
	sti.PlaceholderStyle = views.InputPlaceholderStyle
	sti.Cursor.Style = views.CursorStyle

	return sti
}
