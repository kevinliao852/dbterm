package pages

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

	if m.currentModel == nil {
		m.currentModel = m.connectionPage
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.currentModel, cmd = m.currentModel.Update(msg)

	case Navigator:
		log.Println("msg.To", msg.To)

		switch msg.To {

		case ConnectionPageType:
			m.connectionPage = resizeConnectionPage(m.connectionPage, m.windowSizeMsg())
			m.currentModel, cmd = m.connectionPage.Update(msg)

		case ConfirmPageType:
			log.Println("confirm page")
			m.confirmPage = resizeConfirmPage(m.confirmPage, m.windowSizeMsg())
			if msg.Options != nil {
				if db, ok := (*msg.Options)["db"].(*sql.DB); ok {
					m.confirmPage.DB = db
				}
				if driverType, ok := (*msg.Options)["driverType"].(string); ok {
					m.confirmPage.driverType = driverType
				}
			}

			m.currentModel, cmd = m.confirmPage.Update(msg)

		case QueryPageType:
			log.Println(msg, "options")
			m.queryPage = resizeQueryPage(m.queryPage, m.windowSizeMsg())
			if msg.Options != nil {
				if db, ok := (*msg.Options)["db"].(*sql.DB); ok {
					m.queryPage.DB = db
				}
				if driverType, ok := (*msg.Options)["driverType"].(string); ok {
					m.queryPage.driverType = driverType
				}
			}

			m.currentModel, cmd = m.queryPage.Update(msg)
		}
		cmd = tea.Batch(cmd, tea.ClearScreen)
	default:
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

func (m Term) windowSizeMsg() tea.WindowSizeMsg {
	return tea.WindowSizeMsg{Width: m.windowWidth, Height: m.windowHeight}
}

func resizeConnectionPage(page ConnectionPage, msg tea.WindowSizeMsg) ConnectionPage {
	model, _ := page.Update(msg)
	return model.(ConnectionPage)
}

func resizeConfirmPage(page ConfirmPage, msg tea.WindowSizeMsg) ConfirmPage {
	model, _ := page.Update(msg)
	return model.(ConfirmPage)
}

func resizeQueryPage(page QueryPage, msg tea.WindowSizeMsg) QueryPage {
	model, _ := page.Update(msg)
	return *model.(*QueryPage)
}
