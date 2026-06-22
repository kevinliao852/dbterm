package pages

import (
	"context"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevinliao852/dbterm/pkg/llm"
)

func TestConnectionPageResizesURIInput(t *testing.T) {
	page := NewConnectionPage()
	page.selectDriver()

	narrowModel, _ := page.Update(tea.WindowSizeMsg{Width: 40, Height: 12})
	narrow := narrowModel.(ConnectionPage)
	wideModel, _ := narrow.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	wide := wideModel.(ConnectionPage)

	if narrow.secondTextInput.Width >= wide.secondTextInput.Width {
		t.Fatalf(
			"expected URI input to grow with the viewport: narrow=%d wide=%d",
			narrow.secondTextInput.Width,
			wide.secondTextInput.Width,
		)
	}
}

func TestQueryPageResizesTable(t *testing.T) {
	page := NewQueryPage()
	page.columnNames = []string{"id", "name", "email"}

	page.Update(tea.WindowSizeMsg{Width: 48, Height: 18})
	narrowWidth := page.DataTable.Width()
	narrowHeight := page.DataTable.Height()
	narrowColumnWidth := page.columnWidth()

	page.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	if page.DataTable.Width() <= narrowWidth {
		t.Errorf("expected table width to grow: narrow=%d wide=%d", narrowWidth, page.DataTable.Width())
	}
	if page.DataTable.Height() <= narrowHeight {
		t.Errorf("expected table height to grow: narrow=%d wide=%d", narrowHeight, page.DataTable.Height())
	}
	if page.columnWidth() <= narrowColumnWidth {
		t.Errorf(
			"expected columns to grow: narrow=%d wide=%d",
			narrowColumnWidth,
			page.columnWidth(),
		)
	}
}

func TestTermForwardsPaneResize(t *testing.T) {
	term := NewTermModel()
	model, _ := term.Update(tea.WindowSizeMsg{Width: 72, Height: 20})
	resized := model.(Term)

	if resized.windowWidth != 72 || resized.windowHeight != 20 {
		t.Fatalf("expected root size 72x20, got %dx%d", resized.windowWidth, resized.windowHeight)
	}

	connection, ok := resized.currentModel.(ConnectionPage)
	if !ok {
		t.Fatalf("expected connection page, got %T", resized.currentModel)
	}
	if connection.width != 72 || connection.height != 20 {
		t.Fatalf("expected active page size 72x20, got %dx%d", connection.width, connection.height)
	}
}

func TestQueryComposerCtrlJAddsNewline(t *testing.T) {
	page := NewQueryPage()
	page.DbInput.SetValue("SELECT")

	page.Update(tea.KeyMsg{Type: tea.KeyCtrlJ})

	if !strings.Contains(page.DbInput.Value(), "\n") {
		t.Fatalf("expected ctrl+j to insert a newline, got %q", page.DbInput.Value())
	}
}

func TestQueryComposerEnterSubmitsWithoutNewline(t *testing.T) {
	page := NewQueryPage()
	page.DbInput.SetValue("SELECT 1")

	page.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if page.DbInput.Value() != "SELECT 1" {
		t.Fatalf("expected enter to submit without editing input, got %q", page.DbInput.Value())
	}
	if page.selectData != "DB is not connected" {
		t.Fatalf("expected disconnected status, got %q", page.selectData)
	}
}

func TestQueryPageSwitchesToAITab(t *testing.T) {
	page := NewQueryPage()

	page.Update(tea.KeyMsg{Type: tea.KeyTab})

	if page.activeTab != aiTab {
		t.Fatalf("expected AI tab, got %d", page.activeTab)
	}
}

func TestAIQueryUsesSchemaAndQuestion(t *testing.T) {
	db, err := connectDB("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)"); err != nil {
		t.Fatal(err)
	}

	page := NewQueryPage()
	page.DB = db
	page.driverType = "sqlite"
	page.activeTab = aiTab
	page.AIInput.SetValue("How many users are there?")
	page.generateSQL = func(_ context.Context, dialect, schema, question string) (llm.SQLSuggestion, error) {
		if dialect != "sqlite" {
			t.Fatalf("unexpected dialect %q", dialect)
		}
		if !strings.Contains(schema, "CREATE TABLE users") {
			t.Fatalf("schema was not supplied: %q", schema)
		}
		if question != "How many users are there?" {
			t.Fatalf("unexpected question %q", question)
		}
		return llm.SQLSuggestion{
			SQL:         "SELECT COUNT(*) FROM users",
			Explanation: "Counts users.",
		}, nil
	}

	cmd := page.requestAISuggestion()
	message := cmd()
	page.Update(message)

	if page.generatedSQL != "SELECT COUNT(*) FROM users" {
		t.Fatalf("unexpected generated SQL %q", page.generatedSQL)
	}
}

func TestGeneratedSQLMustBeReadOnly(t *testing.T) {
	for _, query := range []string{
		"INSERT INTO users VALUES (1)",
		"DROP TABLE users",
		"SELECT 1; DELETE FROM users",
		"WITH removed AS (DELETE FROM users RETURNING *) SELECT * FROM removed",
		"EXPLAIN ANALYZE UPDATE users SET name = 'x'",
	} {
		if isReadOnlySQL(query) {
			t.Errorf("expected query to be blocked: %q", query)
		}
	}

	for _, query := range []string{
		"SELECT * FROM users",
		"WITH active AS (SELECT * FROM users) SELECT * FROM active",
		"EXPLAIN SELECT * FROM users",
	} {
		if !isReadOnlySQL(query) {
			t.Errorf("expected query to be allowed: %q", query)
		}
	}
}
