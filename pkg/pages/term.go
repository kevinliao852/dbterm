package pages

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kevinliao852/dbterm/pkg/views"
)

type Term struct {
	connectionPage ConnectionPage
	queryPage      QueryPage
	currentModel   tea.Model
	windowWidth    int
	windowHeight   int
}

func NewTermModel() Term {
	return Term{
		connectionPage: NewConnectionPage(),
		queryPage:      NewQueryPage(),
		currentModel:   nil,
		windowWidth:    0,
		windowHeight:   0,
	}
}

func (m Term) Init() tea.Cmd {
	return textinput.Blink
}

func (m Term) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height

	case Navigator:
		switch msg.On {
		case ConnectionPageType:
			if msg.Options != nil {
				if db, ok := (*msg.Options)["db"].(*sql.DB); ok {
					m.queryPage.DB = db
				}
			}

			m.currentModel, cmd = m.queryPage.Update(msg)
		}
	}

	if m.currentModel == nil {
		m.currentModel = m.connectionPage
	}

	// If this is not a Navigator message, update the current model
	if _, ok := msg.(Navigator); !ok {
		m.currentModel, cmd = m.currentModel.Update(msg)
	}

	return m, tea.Batch(cmd)
}

func (m Term) View() string {
	var view string

	if m.currentModel != nil {
		view = m.currentModel.View()
	}

	return views.TerminalFrame(view, views.NewTerminal(m.windowWidth, m.windowHeight))
}
