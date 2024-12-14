package db

import (
	"database/sql"
	"log"
)

// DB represents a database connection
type DB struct {
	driverName string
	uri        string
	logger     *log.Logger
}

func (db *DB) connectDB() (*sql.DB, error) {
	db.logger.Println("Connecting to the database...")

	// Open a connection to the target database
	dbInstance, err := sql.Open(db.driverName, db.uri)
	if err != nil {
		return nil, err
	}

	// Check if the connection to the database is successful
	err = dbInstance.Ping()
	if err != nil {
		return nil, err
	}

	return dbInstance, err
}
