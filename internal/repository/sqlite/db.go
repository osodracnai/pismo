package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const initSchema = `CREATE TABLE IF NOT EXISTS accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_number TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    account_id INTEGER NOT NULL,
    operation_type_id INTEGER NOT NULL,
    amount REAL NOT NULL,
    event_date TEXT NOT NULL,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);`

func Open(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := migrate(context.Background(), db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, initSchema); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}
	return nil
}
