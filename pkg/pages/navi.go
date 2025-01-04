package pages

import tea "github.com/charmbracelet/bubbletea"

type Navigator struct {
	On      int
	Options *map[string]interface{}
}

func (n *Navigator) NavigateOn() tea.Msg {
	return Navigator{On: n.On, Options: n.Options}
}
