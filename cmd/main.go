package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kevinliao852/dbterm/pkg/pages"
	log "github.com/sirupsen/logrus"
)

type LoggerOption struct {
	log    *log.Logger
	prefix string
}

func (l *LoggerOption) SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func (l *LoggerOption) SetPrefix(s string) {
	l.prefix = s
}

func main() {

	if len(os.Getenv("DEBUG")) > 0 {
		log := log.New()

		lo := &LoggerOption{
			log: log,
		}

		f, err := tea.LogToFileWith("debug.log", "DEBUG", lo)

		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	} else {
		fmt.Println("need to export DEBUG=true")
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	err            error
	num            int
	connectionPage pages.ConnectionPage
	queryPage      pages.QueryPage
	currentModel   tea.Model
	windowWidth    int
	windowHeight   int
}

func initialModel() model {
	return model{
		err:            nil,
		num:            0,
		connectionPage: pages.NewConnectionPage(),
		queryPage:      pages.NewQueryPage(),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height

	case pages.Navigator:
		log.Println("msg.To", msg.To)

		switch msg.To {

		case pages.ConnectionPageType:
			m.currentModel, cmd = m.connectionPage.Update(msg)

		case pages.QueryPageType:
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

		if _, ok := msg.(pages.Navigator); !ok {
			m.currentModel, cmd = m.currentModel.Update(msg)
		}
	}

	return m, cmd
}

func (m model) View() string {

	if m.currentModel == nil {
		return ""
	}

	return m.frame(m.currentModel.View())

}

func (m model) frame(inner string) string {
	termstyle := lipgloss.NewStyle().Align(lipgloss.Center).BorderStyle(lipgloss.ThickBorder())
	termstyle.Margin(1, 1, 1, 1)
	termstyle.Width(m.windowWidth - 5)
	termstyle.Height(m.windowHeight - 5)

	return termstyle.Render(fmt.Sprintf(
		"Enter the input:\n\nwidth:%s height:%s\n\n%s\n\n%s",
		strconv.Itoa(m.windowWidth),
		strconv.Itoa(m.windowHeight),
		inner,
		"(press ESC or CRL+C to quit)",
	))
}
