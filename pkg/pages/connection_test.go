package pages

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSQLiteDriverTypes(t *testing.T) {
	for _, driverType := range []string{"sqlite", "sqlite3"} {
		if !(ConnectionPage{}).isValidDriverType(driverType, driverMap) {
			t.Errorf("expected %q to be a valid driver type", driverType)
		}

		if driverMap[driverType] != "sqlite3" {
			t.Errorf("expected %q to use the sqlite3 driver", driverType)
		}
	}
}

func TestConnectSQLite(t *testing.T) {
	db, err := connectDB(driverMap["sqlite"], ":memory:")
	if err != nil {
		t.Fatalf("connect to in-memory SQLite database: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE example (id INTEGER PRIMARY KEY)"); err != nil {
		t.Fatalf("create SQLite table: %v", err)
	}
}

func TestSelectDriverPrefillsURI(t *testing.T) {
	page := NewConnectionPage()
	page.driverIndex = 2

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyEnter})
	selectedPage := model.(ConnectionPage)

	if selectedPage.driverType != "sqlite" {
		t.Fatalf("expected sqlite driver, got %q", selectedPage.driverType)
	}

	if selectedPage.secondTextInput.Value() != "./database.db" {
		t.Fatalf("expected SQLite URI template, got %q", selectedPage.secondTextInput.Value())
	}
}
