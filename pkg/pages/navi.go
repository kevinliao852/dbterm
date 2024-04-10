package pages

import tea "github.com/charmbracelet/bubbletea"

type Navigator struct {
	To      string
	Options *map[string]interface{}
}

func (n *Navigator) NavigateTo() tea.Msg {
	return Navigator{To: n.To, Options: n.Options}
}
