package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kevinliao852/dbterm/pkg/logger"
	"github.com/kevinliao852/dbterm/pkg/pages"
	log "github.com/sirupsen/logrus"
)

func main() {

	// add logger
	if len(os.Getenv("DEBUG")) > 0 {
		tlog := logger.NewLoggerOption(log.New())
		f, err := tea.LogToFileWith("debug.log", "DEBUG", tlog)
		defer (func() {
			err := f.Close()
			if err != nil {
				log.Error(err)
			}
		})()

		if err != nil {
			fmt.Println("fatal:", err)
		}
	}

	// start the program
	p := tea.NewProgram(pages.NewTermModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
