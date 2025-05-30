package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// Creates the photo table if it doesn't exist
func CreaTabellaFoto(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS photo(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			photo BLOB NOT NULL
		);`
	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error creating photo table: %w", err)
	}
	return nil
}

// Inserts a new photo and returns its ID
func (db *appdbimpl) CreaFoto(photo string) (int, error) {
	insertQuery := `INSERT INTO photo (photo) VALUES (?);`
	result, err := db.c.Exec(insertQuery, photo)
	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("error inserting photo: %w", err)
	}
	lastIdInserito, err := result.LastInsertId()
	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("error retrieving inserted photo ID: %w", err)
	}
	return int(lastIdInserito), nil
}
