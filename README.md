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
# mysql

root:my-secret-pw@tcp(localhost:3306)/my_db

# postgres
postgres://postgres:postgres@localhost:5432/demo
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

### tmux

DBTerm responds to tmux pane resizing automatically. For consistent 256-color
and true-color rendering, use this tmux configuration:

```tmux
set -g default-terminal "tmux-256color"
set -as terminal-features ",xterm-256color:RGB"
```

After changing the configuration, reload it with:

```bash
tmux source-file ~/.tmux.conf
```

### Ask AI

After connecting, press `Tab` to switch between the SQL and Ask AI workspaces.
Ask AI sends the database schema and your question to the OpenAI Responses API;
it does not send table rows.

Set your API key before starting DBTerm:

```bash
export OPENAI_API_KEY="your-api-key"
```

Optional configuration:

```bash
export OPENAI_MODEL="gpt-5.5"
export OPENAI_BASE_URL="https://api.openai.com/v1"
```

In the Ask AI tab:

- `Enter` generates SQL.
- `Ctrl+J` inserts a newline.
- `Ctrl+E` executes generated SQL after review.

Generated SQL is restricted to a single read-only `SELECT`, `WITH`, or
`EXPLAIN` statement.

### Clean Up

```
make clean
```
