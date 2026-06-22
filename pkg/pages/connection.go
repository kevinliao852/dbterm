package pages

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kevinliao852/dbterm/pkg/models"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type ConnectionPage struct {
	secondTextInput textinput.Model
	errorStr        string
	driverType      string
	driverIndex     int
	db              *sql.DB
	isConnected     bool
}

type driverOption struct {
	name     string
	label    string
	template string
}

var driverOptions = []driverOption{
	{
		name:     "mysql",
		label:    "MySQL",
		template: "root:password@tcp(localhost:3306)/database",
	},
	{
		name:     "postgres",
		label:    "PostgreSQL",
		template: "postgres://user:password@localhost:5432/database?sslmode=disable",
	},
	{
		name:     "sqlite",
		label:    "SQLite",
		template: "./database.db",
	},
}

var driverMap = map[string]string{
	"mysql":    "mysql",
	"postgres": "pgx",
	"sqlite":   "sqlite3",
	"sqlite3":  "sqlite3",
}

func NewConnectionPage() ConnectionPage {
	return ConnectionPage{
		secondTextInput: models.ConnectionURIInput(),
	}
}

var _ Pager = &ConnectionPage{}
var _ tea.Model = &ConnectionPage{}

func (q ConnectionPage) Init() tea.Cmd {
	return textinput.Blink
}

func (q ConnectionPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return q, tea.Quit
		}

		if q.driverType == "" {
			switch keyMsg.String() {
			case "up", "k":
				q.driverIndex = (q.driverIndex - 1 + len(driverOptions)) % len(driverOptions)
			case "down", "j":
				q.driverIndex = (q.driverIndex + 1) % len(driverOptions)
			case "1", "2", "3":
				q.driverIndex = int(keyMsg.Runes[0] - '1')
				q.selectDriver()
			case "enter":
				q.selectDriver()
			}
			return q, nil
		}

		if keyMsg.Type == tea.KeyEnter {
			dbURI := q.secondTextInput.Value()
			if dbURI == "" {
				q.errorStr = "\nPlease enter a database URI\n"
				return q, nil
			}

			q.errorStr = ""
			var err error
			q.db, err = connectDB(driverMap[q.driverType], dbURI)
			if err != nil {
				q.errorStr = "\nError connecting to the database\n" + err.Error() + "\n"
				return q, nil
			}
			q.isConnected = true
		}
	}

	if q.isConnected {
		options := &map[string]interface{}{"db": q.db}
		n := Navigator{
			To:      ConfirmPageType,
			Options: options,
		}
		return q, n.NavigateTo
	}

	var cmd tea.Cmd
	q.secondTextInput, cmd = q.secondTextInput.Update(msg)
	return q, cmd
}

func (q ConnectionPage) View() string {
	baseBorder := lipgloss.NewStyle().BorderStyle(lipgloss.ThickBorder())

	if q.driverType == "" {
		options := "Choose a database:\n\n"
		for index, option := range driverOptions {
			cursor := "  "
			if index == q.driverIndex {
				cursor = "> "
			}
			options += cursor + option.label + "\n"
		}
		options += "\n↑/↓ or j/k: select  •  1/2/3: quick select  •  enter: confirm"
		return baseBorder.Render(options)
	}

	selectedDriver := driverOptions[q.driverIndex]
	return baseBorder.Render(
		"Connect to " + selectedDriver.label + "\n\n" +
			"Edit the connection URI:\n" +
			q.secondTextInput.View() +
			"\n\nA template is filled in; replace its example values and press enter." +
			q.errorStr,
	)
}

func (q ConnectionPage) getPageName() string {
	return "connectionPage"
}

func (q ConnectionPage) isValidDriverType(s string, driverMap map[string]string) bool {
	_, exists := driverMap[s]
	return exists
}

func (q *ConnectionPage) selectDriver() {
	option := driverOptions[q.driverIndex]
	q.driverType = option.name
	q.errorStr = ""
	q.secondTextInput.SetValue(option.template)
	q.secondTextInput.CursorEnd()
	q.secondTextInput.Focus()
}

func connectDB(driverName, dbURI string) (*sql.DB, error) {
	log.Println("Connecting to the database...")

	db, err := sql.Open(driverName, dbURI)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
