package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"

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

type (
	errMsg error
)
type model struct {
	textInput       textinput.Model
	secondTextInput textinput.Model
	dbInput         textinput.Model
	dataTable       table.Model
	err             error
	num             int
	db              *sql.DB
	selectData      string
	tableRow        []table.Row
	tableColumn     []table.Column
	conntectionPage pages.ConntectionPage
	queryPage       pages.QueryPage
	currentModel    tea.Model
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "TYPE"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

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

	connectionPage := pages.ConntectionPage{
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

	var db *sql.DB

	queryPage := pages.QueryPage{
		DbInput:   dbi,
		DataTable: t,
		DB:        db,
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
		conntectionPage: connectionPage,
		queryPage:       queryPage,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case pages.Navigator:
		log.Println("msg.To", msg.To)

		switch msg.To {

		case "connectionPage":
			m.currentModel, cmd = m.conntectionPage.Update(msg)

		case "queryPage":
			log.Println(msg, "options")
			if msg.Options != nil {
				if db, ok := (*msg.Options)["db"].(*sql.DB); ok {
					m.db = db
				}
			}

			m.currentModel, cmd = m.queryPage.Update(msg)
		}
	}

	log.Println("cmd", cmd)

	if m.currentModel == nil {
		m.currentModel = m.conntectionPage
	}

	m.currentModel, cmd = m.currentModel.Update(msg)

	return m, cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlA:
			fmt.Println("ctrl a")

		case tea.KeyEnter:
			log.Println("Enter pressed")
			m.num++
			if m.db != nil {

				// check if the db is connected
				if m.db.Ping() != nil {
					m.selectData = "DB is not connected"
				}

				m.selectData = ""
				m.tableColumn = []table.Column{}
				m.tableRow = []table.Row{}
				m.dataTable.SetRows(m.tableRow)
				m.dataTable.SetColumns(m.tableColumn)

				rows, err := m.db.Query(m.dbInput.Value())
				if err != nil {

					m.selectData = err.Error()
				} else {
					types, _ := rows.ColumnTypes()
					for rows.Next() {
						row := make([]interface{}, len(types))

						// SELECT * FROM table
						for i := range types {
							row[i] = new(interface{})
						}
						rows.Scan(row...)

						if err != nil {
							log.Fatal(err)
						}

						tableColumns := []table.Column{}

						for _, col := range types {
							width := len(col.Name())
							tableColumns = append(tableColumns, table.Column{
								Title: col.Name(),
								Width: width,
							})
						}
						m.tableColumn = tableColumns

						var tableRow table.Row

						for _, fields := range row {
							pField := fields.(*interface{})
							strField := fmt.Sprintf("%s", *pField)
							tableRow = append(tableRow, strField)
						}

						m.tableRow = append(m.tableRow, tableRow)

					}
				}

			}

		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.num == 1 {
		m.secondTextInput, cmd = m.secondTextInput.Update(msg)
	} else if m.num == 0 {
		m.textInput, cmd = m.textInput.Update(msg)
	} else if m.num == 2 && m.db == nil {
		var err error

		if err != nil {
			m.selectData = err.Error()
			m.num--
		} else {
			m.selectData = "DB is connected"
			m.num++
			m.err = nil
		}

	} else if m.num >= 3 {
		m.dbInput, cmd = m.dbInput.Update(msg)
		m.dataTable.SetColumns(m.tableColumn)
		m.dataTable.SetRows(m.tableRow)
	}
	return m, cmd
}

func (m model) View() string {
	if m.currentModel == nil {
		return ""
	}

	return fmt.Sprintf(
		"Enter the input:\n\n%s\n\n%s",
		m.currentModel.View(),
		"(esc to quit)",
	)
}
