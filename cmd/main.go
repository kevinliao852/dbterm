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
	log "github.com/sirupsen/logrus"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color("240"))

type LoggerOption struct {
	log    *log.Logger
	prefix string
}

// SetOutput (io.Writer)
// SetPrefix (string)
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

	return model{
		textInput:       ti,
		secondTextInput: sti,
		dbInput:         dbi,
		err:             nil,
		num:             0,
		dataTable:       t,
		tableRow:        tr,
		tableColumn:     tc,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

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
		m.db, err = conntectDB(m.secondTextInput.Value())

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
	var view string

	if m.num > 0 {
		view = m.secondTextInput.View()
	} else {
		view = m.textInput.View()
	}

	if m.num >= 2 {
		view = m.dbInput.View()

		return fmt.Sprintf("Select the DB\n\n%s\n\n%s\n%s",
			view,
			m.selectData,
			baseStyle.Render(m.dataTable.View()),
		)
	}

	return fmt.Sprintf(
		"Enter the intput\n\n%s\n\n%s\n\n%s",
		view,
		m.selectData,
		"(esc to quit)",
	)
}

func conntectDB(dbURI string) (*sql.DB, error) {
	log.Println("Connecting to the database...")

	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, err
	}

	// Check if the connection to the database is successful
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
