package postgresql

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "yaus.postgres.New"

	dbPath, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("%s: cannot get user home direcory: %w", op, err)
	}
	dbPath = filepath.Join(dbPath, ".local/share/yaus")
	os.Mkdir(dbPath, 0755)

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != err {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db}, nil
}
