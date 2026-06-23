package main

import (
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kevinliao852/dbterm/pkg/logger"
	"github.com/kevinliao852/dbterm/pkg/pages"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func main() {

	// add logger
	if len(os.Getenv("DEBUG")) > 0 {
		tlog := logger.NewLoggerOption(log.New())
		f, err := tea.LogToFileWith("debug.log", "DEBUG", tlog)
		if err != nil {
			fmt.Println("fatal:", err)
		} else {
			defer func() {
				if err := f.Close(); err != nil {
					log.Error(err)
				}
			}()
		}
	} else {
		log.SetOutput(io.Discard)
	}

	// start the program
	p := tea.NewProgram(pages.NewTermModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
