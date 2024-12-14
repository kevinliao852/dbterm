package pages

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kevinliao852/dbterm/pkg/views"
	log "github.com/sirupsen/logrus"
)

type Term struct {
	connectionPage ConnectionPage
	confirmPage    ConfirmPage
	queryPage      QueryPage
	currentModel   tea.Model
	windowWidth    int
	windowHeight   int
}

func NewTermModel() Term {
	return Term{
		connectionPage: NewConnectionPage(),
		confirmPage:    NewConfirmPage(),
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
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height

	case Navigator:
		log.Println("msg.To", msg.To)

		switch msg.To {

		case ConnectionPageType:
			m.currentModel, cmd = m.connectionPage.Update(msg)

		case ConfirmPageType:
			log.Println("confirm page")
			if msg.Options != nil {
				if db, ok := (*msg.Options)["db"].(*sql.DB); ok {
					m.confirmPage.DB = db
				}
			}

			m.currentModel, cmd = m.confirmPage.Update(msg)

		case QueryPageType:
			log.Println(msg, "options")
			if msg.Options != nil {
				if db, ok := (*msg.Options)["db"].(*sql.DB); ok {
					m.queryPage.DB = db
				}
			}

			m.currentModel, cmd = m.queryPage.Update(msg)
		}
	default:
		if m.currentModel == nil {
			m.currentModel = m.connectionPage
		}

		if _, ok := msg.(Navigator); !ok {
			m.currentModel, cmd = m.currentModel.Update(msg)
		}
	}

	return m, cmd
}

func (m Term) View() string {
	var view string

	if m.currentModel != nil {
		view = m.currentModel.View()
	}

	return views.TerminalFrame(view, views.NewTerminal(m.windowWidth, m.windowHeight))
}
