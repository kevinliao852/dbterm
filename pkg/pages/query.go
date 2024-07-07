package pages

import (
	"database/sql"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/sirupsen/logrus"
)

type QueryPage struct {
	DbInput    textinput.Model
	DataTable  table.Model
	selectData string
	DB         *sql.DB
	queryStr   string
	selectedDB string
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color("240"))

var _ Pager = &QueryPage{}

var _ tea.Model = &QueryPage{}

func (q *QueryPage) Init() tea.Cmd {
	return q.DbInput.Focus()
}

func (q *QueryPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return q, tea.Quit
		case tea.KeyEnter:
			log.Println("Enter pressed")

			log.Println("DB:", q.DB, &q.DB)

			if q.DB == nil {
				log.Println("here")
				break
			}

			if q.DbInput.Value() == "" {
				q.selectData = "Please enter the SQL code"
				break
			}

			if q.DB.Ping() != nil {
				q.selectData = "DB is not connected"
				break
			}

			q.queryStr = q.DbInput.Value()
			q.selectData = "Querying the database"
			q.readAndQuery()

		}
	}

	q.DbInput, cmd = q.DbInput.Update(msg)

	return q, cmd
}

func (q *QueryPage) View() string {

	view := q.DbInput.View()

	return fmt.Sprintf("Select the database\n\n%s\n\n%s\n%s",
		view,
		q.selectData,
		baseStyle.Render(q.DataTable.View()),
	)
}

func (q *QueryPage) getPageName() string {
	return "queryPage"
}

func (q *QueryPage) readAndQuery() {
	if q.queryStr == "" {
		q.selectData = "Please enter the query"
		return
	}

	rows, err := q.DB.Query(q.queryStr)

	if err != nil {
		q.selectData = "Error executing the query\n" + err.Error()
		return
	}

	tableColumn := []table.Column{}
	tableRowList := []table.Row{}

	types, _ := rows.ColumnTypes()

	for _, col := range types {
		width := len(col.Name())
		tableColumn = append(tableColumn, table.Column{
			Title: col.Name(),
			Width: width,
		})
	}

	for rows.Next() {

		row := make([]interface{}, 0)

		for range types {
			row = append(row, new(interface{}))
		}

		log.Println("debug", types, len(types), row)

		err := rows.Scan(row...)

		if err != nil {
			q.selectData = "Error scanning the row\n" + err.Error()
			return
		}

		var tableRow table.Row

		for _, fields := range row {
			pField := fields.(*interface{})
			strField := fmt.Sprintf("%s", *pField)

			tableRow = append(tableRow, strField)
		}

		log.Println(tableRow)

		tableRowList = append(tableRowList, tableRow)
	}

	log.Println(tableColumn)
	log.Println(tableRowList)

	// make sure to set culumn first!
	q.DataTable.SetColumns(tableColumn)

	if len(tableColumn) > 0 {
		q.DataTable.SetRows(tableRowList)
	}
}
