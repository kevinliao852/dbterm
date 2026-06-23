package do

import (
	"database/sql"
	"fmt"
	"strings"
)

func Schema(db *sql.DB, driverType string) (string, error) {
	switch driverType {
	case "sqlite", "sqlite3":
		return sqliteSchema(db)
	case "postgres":
		return informationSchema(db, `
			SELECT table_schema, table_name, column_name, data_type
			FROM information_schema.columns
			WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
			ORDER BY table_schema, table_name, ordinal_position`)
	case "mysql":
		return informationSchema(db, `
			SELECT table_schema, table_name, column_name, data_type
			FROM information_schema.columns
			WHERE table_schema = DATABASE()
			ORDER BY table_name, ordinal_position`)
	default:
		return "", fmt.Errorf("unsupported database driver %q", driverType)
	}
}

func sqliteSchema(db *sql.DB) (string, error) {
	rows, err := db.Query(`
		SELECT sql
		FROM sqlite_master
		WHERE type IN ('table', 'view') AND sql IS NOT NULL
		ORDER BY name`)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var statements []string
	for rows.Next() {
		var statement string
		if err := rows.Scan(&statement); err != nil {
			return "", err
		}
		statements = append(statements, statement+";")
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	if len(statements) == 0 {
		return "(database has no tables)", nil
	}
	return strings.Join(statements, "\n"), nil
}

func informationSchema(db *sql.DB, query string) (string, error) {
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var lines []string
	for rows.Next() {
		var schemaName, tableName, columnName, dataType string
		if err := rows.Scan(&schemaName, &tableName, &columnName, &dataType); err != nil {
			return "", err
		}
		lines = append(lines, fmt.Sprintf("%s.%s.%s %s", schemaName, tableName, columnName, dataType))
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	if len(lines) == 0 {
		return "(database has no tables)", nil
	}
	return strings.Join(lines, "\n"), nil
}
