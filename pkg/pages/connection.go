package pages

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevinliao852/dbterm/pkg/models"
	log "github.com/sirupsen/logrus"
)

type ConnectionPage struct {
	textInput       textinput.Model
	secondTextInput textinput.Model
	errorStr        string
	driverType      string
	db              *sql.DB
	isConnected     bool
}

func NewConnectionPage() ConnectionPage {
	return ConnectionPage{
		textInput:       models.ConnectionTypeInput(),
		secondTextInput: models.ConnectionURIInput(),
	}
}

var _ Pager = &ConnectionPage{}

var _ tea.Model = &ConnectionPage{}

var driverMap = map[string]string{
	"mysql":    "mysql",
	"postgres": "pgx",
	"sqlite3": "sqlite3",
}

func (q ConnectionPage) Init() tea.Cmd {
	return textinput.Blink
}

func (q ConnectionPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				return q, tea.Quit
			case tea.KeyEnter:
				driverType := q.textInput.Value()
				dbUri := q.secondTextInput.Value()

				if !q.isValidDriverType(driverType, driverMap) {
					q.errorStr = "\nInvalid driver type.\nonly 'mysql' and 'postgres' are supported\n"
					break
				}
				q.errorStr = ""
				q.driverType = q.textInput.Value()

				if driverType != "" && dbUri != "" {
					// connect to the database
					driverName := driverMap[driverType]
					var err error
					q.db, err = connectDB(driverName, q.secondTextInput.Value())

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
			To:      ConfirmPageType,
			Options: options,
		}

		return q, n.NavigateTo
	}

	if q.driverType == "" {
		q.textInput, cmd = q.textInput.Update(msg)
	} else {
		q.secondTextInput, cmd = q.secondTextInput.Update(msg)
	}

	return q, cmd
}

func (q ConnectionPage) View() string {
	baseBorder := lipgloss.NewStyle().BorderStyle(lipgloss.ThickBorder())

	if q.driverType == "" {
		q.textInput.Placeholder = "Enter the driver type (mysql, postgres)"
		q.textInput.Focus()
		return baseBorder.Render(q.textInput.View() + q.errorStr)
	}

	q.secondTextInput.Placeholder = "Enter the database uri"
	q.secondTextInput.Focus()

	return baseBorder.Render(q.secondTextInput.View() + q.errorStr)

}

func (q ConnectionPage) getPageName() string {
	return "connectionPage"
}

func (q ConnectionPage) isValidDriverType(s string, driverMap map[string]string) bool {
	_, exists := driverMap[s]
	return exists
}

func connectDB(driverName, dbURI string) (*sql.DB, error) {
	log.Println("Connecting to the database...")

	// Open a connection to the target database
	db, err := sql.Open(driverName, dbURI)
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
