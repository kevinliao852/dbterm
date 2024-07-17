package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
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
	textInput       textinput.Model
	secondTextInput textinput.Model
	dbInput         textinput.Model
	dataTable       table.Model
	err             error
	num             int
	tableRow        []table.Row
	tableColumn     []table.Column
	connectionPage  pages.ConnectionPage
	queryPage       pages.QueryPage
	currentModel    tea.Model
	windowWidth     int
	windowHeight    int
}

func initialModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = "the type of database(mysql, postgres)"
	ti.CharLimit = 156
	ti.Width = 50
	ti.TextStyle.Background(lipgloss.Color("63"))
	ti.TextStyle.Foreground(lipgloss.Color("63"))
	ti.SetSuggestions([]string{"mysql", "postgres"})
	ti.ShowSuggestions = true
	ti.PromptStyle.Border(lipgloss.NormalBorder())
	ti.TextStyle.BorderStyle(lipgloss.NormalBorder())

	sti := textinput.New()
	sti.Placeholder = "DB_URI"
	sti.Focus()
	sti.CharLimit = 156
	sti.Width = 50

	dbi := textinput.New()
	dbi.Placeholder = "SQL Query"
	dbi.Focus()
	dbi.CharLimit = 156
	dbi.Width = 40

	connectionPage := pages.ConnectionPage{
		TextInput:       ti,
		SecondTextInput: sti,
	}

	var tr []table.Row = []table.Row{}
	var tc []table.Column = []table.Column{}

	t := table.New(
		table.WithColumns(tc),
		table.WithRows(tr),
		table.WithFocused(true),
		table.WithHeight(7),
	)

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
	t.SetStyles(s)
	t.SetWidth(1000)

	queryPage := pages.QueryPage{
		DbInput:   dbi,
		DataTable: t,
	}

	return model{
		textInput:       ti,
		secondTextInput: sti,
		dbInput:         dbi,
		err:             nil,
		num:             0,
		dataTable:       t,
		tableRow:        tr,
		tableColumn:     tc,
		connectionPage:  connectionPage,
		queryPage:       queryPage,
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
	termstyle := lipgloss.NewStyle().Align(lipgloss.Center).BorderStyle(lipgloss.ThickBorder())

	termstyle.Margin(1, 1, 1, 1)
	termstyle.Width(m.windowWidth - 5)
	termstyle.Height(m.windowHeight - 5)

	if m.currentModel == nil {
		return ""
	}

	return termstyle.Render(fmt.Sprintf(
		"Enter the input:\n\nwidth:%s height:%s\n\n%s\n\n%s",
		strconv.Itoa(m.windowWidth),
		strconv.Itoa(m.windowHeight),
		m.currentModel.View(),
		"(press ESC or CRL+C to quit)",
	))
}
