package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"yaus/internal/config"
	"yaus/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fnName = "yaus.sqlite.New"
	if storagePath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("%s: failed to get user path: %s", fnName, err)
		}
		_, err = os.Stat(storagePath)
		if storagePath == config.DefaultDBPath && os.IsNotExist(err) {
			os.Mkdir(homeDir+"/.local/share/yaus", 0755)
		}
		storagePath = homeDir + storagePath[1:]
	}

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open database: %w", fnName, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
		`)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert statement to database: %w", fnName, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create database file: %w", fnName, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const fnName = "yaus.Storage.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %s: %w", fnName, storage.FailedToCreateStatement, err)
	}

	response, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: cannot add two same url aliases: %w", fnName, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: failed to insert statement to database: %w", fnName, err)
	}

	id, err := response.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", fnName, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const fnName = "yaus.Storage.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %s: %w", fnName, storage.FailedToCreateStatement, err)
	}

	var resultURL string
	err = stmt.QueryRow(alias).Scan(&resultURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: failed to select url \"%s\" from database: %w", fnName, alias, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: failed to select url \"%s\" from database: %w", fnName, alias, err)
	}

	return resultURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const fnName = "yaus.Storage.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %s: %w", fnName, storage.FailedToCreateStatement, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: failed to delete column, using %s alias: %w", fnName, alias, storage.ErrURLExists)
		}
		return fmt.Errorf("%s: failed to delete column, using %s alias: %w", fnName, alias, err)
	}

	return nil
}
