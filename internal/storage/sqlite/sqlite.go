package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"yaus/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	if storagePath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("failed to get user path: %s", err)
		}
		storagePath = homeDir + storagePath[1:]
	}
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("falied to open database: %w", err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
		`)
	if err != nil {
		return nil, fmt.Errorf("falied to create statement for database: %w", err)
	}

	_, err = stmt.Exec()
	if err != err {
		return nil, fmt.Errorf("failed to insert statement to database: %w", err)
	}

	return &Storage{db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("falied to create statement for database: %w", err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("error: cannot add two same url aliases: %w", storage.ErrURLExists)
		}
		return fmt.Errorf("failed to insert statement to database: %w", err)
	}

	_, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert: %w", err)
	}

	return nil
}
