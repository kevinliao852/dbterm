package pages

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Pager interface {
	tea.Model
	getPageName() string
}
