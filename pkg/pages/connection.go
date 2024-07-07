package pages

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
)

type ConntectionPage struct {
	TextInput       textinput.Model
	SecondTextInput textinput.Model
	errorStr        string
	driverType      string
	db              *sql.DB
	isConnected     bool
}

var _ Pager = &ConntectionPage{}

var _ tea.Model = &ConntectionPage{}

func (q ConntectionPage) Init() tea.Cmd {
	return textinput.Blink
}

func (q ConntectionPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	log.Println("ConntectionPage Update")

	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				return q, tea.Quit
			case tea.KeyEnter:
				driverType := q.TextInput.Value()
				dbUri := q.SecondTextInput.Value()

				if !q.isValidDriverType(driverType) {
					q.errorStr = "\nInvalid driver type.\nonly 'mysql' and 'postgres' are supported\n"
					break
				}
				q.errorStr = ""
				q.driverType = q.TextInput.Value()

				if driverType != "" && dbUri != "" {
					// connect to the database
					var err error
					q.db, err = conntectDB(q.SecondTextInput.Value())

					if err != nil {
						q.errorStr = "\nError connecting to the database\n" + err.Error() + "\n"
						break
					}

					q.isConnected = true
				}

			}
		}

	}

	var cmd tea.Cmd

	if q.isConnected {
		options := &map[string]interface{}{}
		(*options)["db"] = q.db

		n := Navigator{
			To:      QueryPageType,
			Options: options,
		}

		return q, n.NavigateTo
	}

	if q.driverType == "" {
		q.TextInput, cmd = q.TextInput.Update(msg)
	} else {
		q.SecondTextInput, cmd = q.SecondTextInput.Update(msg)
	}

	return q, cmd
}

func (q ConntectionPage) View() string {
	if q.driverType == "" {
		q.TextInput.Placeholder = "Enter the driver type (mysql, postgres)"
		q.TextInput.Focus()
		return q.TextInput.View() + q.errorStr
	}

	q.SecondTextInput.Placeholder = "Enter the database uri"
	q.SecondTextInput.Focus()

	return q.SecondTextInput.View() + q.errorStr

}

func (q ConntectionPage) getPageName() string {
	return "conntectionPage"
}

func (q ConntectionPage) isValidDriverType(s string) bool {
	if s == "mysql" || s == "postgres" {
		return true
	}

	return false
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
