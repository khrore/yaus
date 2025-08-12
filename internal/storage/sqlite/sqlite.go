package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "yaus.sqlite.New"

	if storagePath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("failed to get user path: %s", err)
		}
		storagePath = homeDir + storagePath[1:]
	}
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL,
			data TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
		`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != err {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db}, nil
}
