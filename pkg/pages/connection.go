package pages

import (
	"database/sql"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kevinliao852/dbterm/pkg/models"
	"github.com/kevinliao852/dbterm/pkg/views"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type ConnectionPage struct {
	secondTextInput textinput.Model
	errorStr        string
	driverType      string
	driverIndex     int
	width           int
	height          int
	db              *sql.DB
	isConnected     bool
}

type driverOption struct {
	name     string
	label    string
	detail   string
	template string
}

var driverOptions = []driverOption{
	{
		name:     "mysql",
		label:    "MySQL",
		detail:   "Network database",
		template: "root:password@tcp(localhost:3306)/database",
	},
	{
		name:     "postgres",
		label:    "PostgreSQL",
		detail:   "Network database",
		template: "postgres://user:password@localhost:5432/database?sslmode=disable",
	},
	{
		name:     "sqlite",
		label:    "SQLite",
		detail:   "Local database file",
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
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		q.width = sizeMsg.Width
		q.height = sizeMsg.Height
		q.resize()
		return q, nil
	}

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
				q.errorStr = "Please enter a database URI."
				return q, nil
			}

			q.errorStr = ""
			var err error
			q.db, err = connectDB(driverMap[q.driverType], dbURI)
			if err != nil {
				q.errorStr = "Could not connect: " + err.Error()
				return q, nil
			}
			q.isConnected = true
		}
	}

	if q.isConnected {
		options := &map[string]interface{}{
			"db":         q.db,
			"driverType": q.driverType,
		}
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
	if q.driverType == "" {
		rows := make([]string, 0, len(driverOptions))
		for index, option := range driverOptions {
			if index == q.driverIndex {
				rows = append(rows, views.ActiveItemStyle.Render("● "+option.label)+
					views.MutedStyle.Render("  "+option.detail))
				continue
			}
			rows = append(rows, views.InactiveItemStyle.Render("○ "+option.label)+
				views.MutedStyle.Render("  "+option.detail))
		}

		help := views.KeyStyle("↑/↓") + views.HelpStyle.Render(" select  ") +
			views.KeyStyle("enter") + views.HelpStyle.Render(" confirm")
		if q.contentWidth() < 55 {
			help = views.KeyStyle("↑/↓") + views.HelpStyle.Render(" select  ") +
				views.KeyStyle("enter") + views.HelpStyle.Render(" confirm")
		} else {
			help = views.KeyStyle("↑/↓") + views.HelpStyle.Render(" select  ") +
				views.KeyStyle("1–3") + views.HelpStyle.Render(" quick select  ") +
				views.KeyStyle("enter") + views.HelpStyle.Render(" confirm")
		}

		content := views.PageTitleStyle.Render("Choose a database") + "\n" +
			views.MutedStyle.Render("Select the engine for this connection.") + "\n\n" +
			strings.Join(rows, "\n\n") + "\n\n" + help

		return views.CardStyle(q.contentWidth()).Render(content)
	}

	selectedDriver := driverOptions[q.driverIndex]
	content := views.PageTitleStyle.Render("Connect to "+selectedDriver.label) + "\n" +
		views.MutedStyle.Render("Edit the generated connection URI.") + "\n\n" +
		views.LabelStyle.Render("CONNECTION URI") + "\n" +
		q.secondTextInput.View() + "\n\n" +
		views.HelpStyle.Render("Replace the example values, then press ") +
		views.KeyStyle("enter") + views.HelpStyle.Render(".")

	if q.errorStr != "" {
		content += "\n\n" + views.ErrorStyle.Render("! "+q.errorStr)
	}

	return views.CardStyle(q.contentWidth()).Render(content)
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
	q.resize()
}

func (q *ConnectionPage) resize() {
	q.secondTextInput.Width = max(12, min(70, q.contentWidth()-4))
}

func (q ConnectionPage) contentWidth() int {
	if q.width <= 0 {
		return 76
	}
	return max(20, q.width-8)
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
