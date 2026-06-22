package pages

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	dbschema "github.com/kevinliao852/dbterm/pkg/db"
	"github.com/kevinliao852/dbterm/pkg/llm"
	"github.com/kevinliao852/dbterm/pkg/models"
	"github.com/kevinliao852/dbterm/pkg/views"
	log "github.com/sirupsen/logrus"
)

const (
	sqlTab = iota
	aiTab
)

var writeStatementPattern = regexp.MustCompile(
	`(?i)\b(INSERT|UPDATE|DELETE|DROP|ALTER|TRUNCATE|CREATE|REPLACE|MERGE|GRANT|REVOKE|VACUUM|ATTACH|DETACH)\b`,
)

type generateSQLFunc func(context.Context, string, string, string) (llm.SQLSuggestion, error)

type aiSuggestionMsg struct {
	suggestion llm.SQLSuggestion
	err        error
}

type QueryPage struct {
	DbInput      textarea.Model
	AIInput      textarea.Model
	DataTable    table.Model
	selectData   string
	DB           *sql.DB
	driverType   string
	queryStr     string
	width        int
	height       int
	columnNames  []string
	activeTab    int
	aiStatus     string
	generatedSQL string
	explanation  string
	generateSQL  generateSQLFunc
}

func NewQueryPage() QueryPage {
	client := llm.NewClientFromEnv()
	return QueryPage{
		DbInput:     models.DBSQLQueryInput(),
		AIInput:     models.DBNaturalLanguageInput(),
		DataTable:   models.DBSelectTable(),
		selectData:  "",
		DB:          nil,
		queryStr:    "",
		generateSQL: client.GenerateSQL,
	}
}

var _ Pager = &QueryPage{}

var _ tea.Model = &QueryPage{}

func (q *QueryPage) Init() tea.Cmd {
	return textarea.Blink
}

func (q *QueryPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		q.width = msg.Width
		q.height = msg.Height
		q.resize()
		return q, nil
	case aiSuggestionMsg:
		if msg.err != nil {
			q.aiStatus = msg.err.Error()
			return q, nil
		}
		q.generatedSQL = msg.suggestion.SQL
		q.explanation = msg.suggestion.Explanation
		q.aiStatus = "SQL generated. Review it, then press ctrl+e to execute."
		q.resize()
		return q, nil
	case tea.KeyMsg:
		if msg.Type == tea.KeyTab {
			q.activeTab = (q.activeTab + 1) % 2
			return q, nil
		}

		if msg.String() == "ctrl+j" {
			var cmd tea.Cmd
			if q.activeTab == aiTab {
				q.AIInput, cmd = q.AIInput.Update(tea.KeyMsg{Type: tea.KeyEnter})
			} else {
				q.DbInput, cmd = q.DbInput.Update(tea.KeyMsg{Type: tea.KeyEnter})
			}
			return q, cmd
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return q, tea.Quit
		case tea.KeyCtrlE:
			if q.activeTab == aiTab {
				q.executeGeneratedSQL()
				return q, nil
			}
		case tea.KeyEnter:
			log.Println("Enter pressed")
			if q.activeTab == aiTab {
				return q, q.requestAISuggestion()
			}
			q.executeQuery(q.DbInput.Value())
			return q, nil
		}
	}

	if q.activeTab == aiTab {
		q.AIInput, cmd = q.AIInput.Update(msg)
	} else {
		q.DbInput, cmd = q.DbInput.Update(msg)
	}

	return q, cmd
}

func (q *QueryPage) executeQuery(query string) {
	if q.DB == nil {
		q.selectData = "DB is not connected"
		return
	}

	if strings.TrimSpace(query) == "" {
		q.selectData = "Please enter the SQL code"
		return
	}

	if q.DB.Ping() != nil {
		q.selectData = "DB is not connected"
		return
	}

	q.queryStr = query
	q.selectData = "Querying the database"
	q.readAndQuery()
}

func (q *QueryPage) View() string {
	tabs := q.renderTabs()
	if q.activeTab == aiTab {
		return tabs + "\n\n" + q.renderAIWorkspace()
	}
	return tabs + "\n\n" + q.renderSQLWorkspace()
}

func (q *QueryPage) renderSQLWorkspace() string {
	composerHelp := views.KeyStyle("enter") + views.HelpStyle.Render(" run  ") +
		views.KeyStyle("ctrl+j") + views.HelpStyle.Render(" newline")
	composer := q.DbInput.View() + "\n" + composerHelp
	queryInput := views.ComposerStyle(q.tableWidth()).Render(composer)
	resultTable := views.PanelStyle(q.tableWidth()).Render(q.DataTable.View())

	status := ""
	if q.selectData != "" {
		status = q.renderStatus() + "\n\n"
	}

	return views.LabelStyle.Render("RESULTS") + "\n" +
		resultTable + "\n\n" +
		status +
		views.LabelStyle.Render("COMPOSER") + "\n" +
		queryInput
}

func (q *QueryPage) renderAIWorkspace() string {
	composerHelp := views.KeyStyle("enter") + views.HelpStyle.Render(" generate  ") +
		views.KeyStyle("ctrl+j") + views.HelpStyle.Render(" newline")
	if q.generatedSQL != "" {
		composerHelp += views.HelpStyle.Render("  ") +
			views.KeyStyle("ctrl+e") + views.HelpStyle.Render(" execute SQL")
	}
	composer := q.AIInput.View() + "\n" + composerHelp
	queryInput := views.ComposerStyle(q.tableWidth()).Render(composer)
	resultTable := views.PanelStyle(q.tableWidth()).Render(q.DataTable.View())

	preview := ""
	if q.generatedSQL != "" {
		preview = views.LabelStyle.Render("GENERATED SQL") + "\n" +
			views.PanelStyle(q.tableWidth()).Render(
				views.BodyStyle.Render(q.generatedSQL)+"\n\n"+
					views.MutedStyle.Render(q.explanation),
			) + "\n\n"
	}

	status := ""
	if q.aiStatus != "" {
		status = q.renderAIStatus() + "\n\n"
	}

	return preview +
		views.LabelStyle.Render("RESULTS") + "\n" +
		resultTable + "\n\n" +
		status +
		views.LabelStyle.Render("QUESTION") + "\n" +
		queryInput
}

func (q QueryPage) renderTabs() string {
	sqlLabel := " SQL "
	aiLabel := " Ask AI "
	if q.activeTab == sqlTab {
		sqlLabel = views.ActiveTabStyle.Render(sqlLabel)
		aiLabel = views.InactiveTabStyle.Render(aiLabel)
	} else {
		sqlLabel = views.InactiveTabStyle.Render(sqlLabel)
		aiLabel = views.ActiveTabStyle.Render(aiLabel)
	}
	return sqlLabel + " " + aiLabel + views.HelpStyle.Render("   tab switch")
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
	defer rows.Close()

	tableColumn := []table.Column{}
	tableRowList := []table.Row{}

	types, err := rows.ColumnTypes()
	if err != nil {
		q.selectData = "Error reading columns\n" + err.Error()
		return
	}
	q.columnNames = make([]string, len(types))

	for rows.Next() {

		row := make([]interface{}, 0)

		for range types {
			row = append(row, new(interface{}))
		}

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

		tableRowList = append(tableRowList, tableRow)
	}

	if err := rows.Err(); err != nil {
		q.selectData = "Error reading rows\n" + err.Error()
		return
	}

	for index, col := range types {
		q.columnNames[index] = col.Name()
		tableColumn = append(tableColumn, table.Column{
			Title: col.Name(),
		})
	}

	// make sure to set column first!
	if len(tableColumn) == 0 {
		var c []table.Column
		c = append(c, table.Column{Title: "Message", Width: 16})
		var r []table.Row
		r = append(r, table.Row{"No Rows Returned"})
		q.DataTable.SetColumns(c)
		q.DataTable.SetRows(r)
		return
	}

	q.DataTable.SetColumns(tableColumn)
	q.DataTable.SetRows(tableRowList)
	q.resizeColumns()
	q.selectData = fmt.Sprintf("%d row(s)", len(tableRowList))
}

func (q *QueryPage) resize() {
	q.DbInput.SetWidth(max(12, q.tableWidth()-4))
	q.DbInput.SetHeight(q.composerHeight())
	q.AIInput.SetWidth(max(12, q.tableWidth()-4))
	q.AIInput.SetHeight(q.composerHeight())
	q.DataTable.SetWidth(max(12, q.tableWidth()-4))
	q.DataTable.SetHeight(max(1, q.height-q.composerHeight()-q.verticalOverhead()))
	q.resizeColumns()
}

func (q *QueryPage) resizeColumns() {
	if len(q.columnNames) == 0 {
		return
	}

	columnWidth := q.columnWidth()
	columns := make([]table.Column, len(q.columnNames))
	for index, name := range q.columnNames {
		columns[index] = table.Column{Title: name, Width: columnWidth}
	}
	q.DataTable.SetColumns(columns)
}

func (q QueryPage) columnWidth() int {
	if len(q.columnNames) == 0 {
		return 0
	}
	tableWidth := max(12, q.tableWidth()-4)
	return max(3, (tableWidth-len(q.columnNames)-1)/len(q.columnNames))
}

func (q QueryPage) tableWidth() int {
	if q.width <= 0 {
		return 76
	}
	return max(16, q.width-10)
}

func (q QueryPage) composerHeight() int {
	if q.height > 32 {
		return 5
	}
	return 3
}

func (q QueryPage) verticalOverhead() int {
	overhead := 23
	if q.activeTab == aiTab && q.generatedSQL != "" {
		overhead += 7
	}
	return overhead
}

func (q QueryPage) renderStatus() string {
	status := q.selectData
	switch {
	case strings.HasPrefix(status, "Error"),
		strings.HasPrefix(status, "Please"),
		strings.HasPrefix(status, "DB is not"):
		return views.ErrorStyle.Render("! " + status)
	case status == "Querying the database":
		return views.MutedStyle.Render("… " + status)
	default:
		return views.SuccessStyle.Render("✓ " + status)
	}
}

func (q QueryPage) renderAIStatus() string {
	if strings.HasPrefix(q.aiStatus, "Generating") {
		return views.MutedStyle.Render("… " + q.aiStatus)
	}
	if strings.HasPrefix(q.aiStatus, "SQL generated") {
		return views.SuccessStyle.Render("✓ " + q.aiStatus)
	}
	return views.ErrorStyle.Render("! " + q.aiStatus)
}

func (q *QueryPage) requestAISuggestion() tea.Cmd {
	question := strings.TrimSpace(q.AIInput.Value())
	if question == "" {
		q.aiStatus = "Please enter a question."
		return nil
	}
	if q.DB == nil {
		q.aiStatus = "DB is not connected"
		return nil
	}

	q.aiStatus = "Generating SQL…"
	q.generatedSQL = ""
	q.explanation = ""

	db := q.DB
	driverType := q.driverType
	generate := q.generateSQL
	return func() tea.Msg {
		schema, err := dbschema.Schema(db, driverType)
		if err != nil {
			return aiSuggestionMsg{err: fmt.Errorf("read database schema: %w", err)}
		}
		suggestion, err := generate(context.Background(), driverType, schema, question)
		return aiSuggestionMsg{suggestion: suggestion, err: err}
	}
}

func (q *QueryPage) executeGeneratedSQL() {
	if q.generatedSQL == "" {
		q.aiStatus = "Generate SQL before executing it."
		return
	}
	if !isReadOnlySQL(q.generatedSQL) {
		q.aiStatus = "Generated SQL was blocked because it is not read-only."
		return
	}

	q.executeQuery(q.generatedSQL)
	q.aiStatus = q.selectData
}

func isReadOnlySQL(query string) bool {
	query = strings.TrimSpace(strings.TrimSuffix(query, ";"))
	if query == "" || strings.Contains(query, ";") || writeStatementPattern.MatchString(query) {
		return false
	}

	firstField := strings.ToUpper(strings.Fields(query)[0])
	switch firstField {
	case "SELECT", "WITH", "EXPLAIN":
		return true
	default:
		return false
	}
}
