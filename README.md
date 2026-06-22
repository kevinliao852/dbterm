# DBTerm

DBTerm is a terminal-based user interface (TUI) for interacting with relational databases using Go.

![DBTerm Demo](https://imgur.com/Fvojw3V.gif)

## Features

- **Simple Interface:** Clean and intuitive TUI design for easy navigation.
- **SQL Querying:** Execute SQL queries directly from the terminal.
- **Table Visualization:** Display query results in a tabular format for better readability.
- **Multi-Database Support:** Connect to and manage multiple databases seamlessly.

## Getting Started

### Prerequisites

- Go (version 1.18 or newer)
- [Your Database System] (e.g., MySQL, PostgreSQL) installed and accessible

### Installation

```bash
go get -u github.com/kevinliao852/dbterm
```

### Build

```
make build
```

### Development

Run the application from source with debug logging enabled:

```
make dev
```

Debug output is written to `debug.log`.

### Connection Examples

Select a database with the arrow keys or the `1`, `2`, and `3` shortcuts.
DBTerm then fills in an editable connection URI template.

MySQL:

```
root:my-secret-pw@tcp(localhost:3306)/my_db
```

PostgreSQL:

```
postgres://user:password@localhost:5432/my_db
```

SQLite:

```
./database.db
```

For an in-memory SQLite database:

```
:memory:
```

### Clean Up

```
make clean
```
