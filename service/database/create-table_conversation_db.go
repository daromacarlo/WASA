package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func CreateTableChat(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS chat(
			id INTEGER PRIMARY KEY AUTOINCREMENT
	 	);`

	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error while creating the chat table: %w", err)
	}
	return nil
}

func CreateTableGroup(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS groups(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name STRING NOT NULL,
			chat UNIQUE NOT NULL,
			photo BLOB NOT NULL,
			FOREIGN KEY (chat) REFERENCES chat(id)
		);`

	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error while creating the groups table: %w", err)
	}
	return nil
}

func CreateTableUseringruppo(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS user_in_group(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			groups INTEGER NOT NULL,
			user INTEGER NOT NULL,
			FOREIGN KEY (groups) REFERENCES groups(id),
			FOREIGN KEY (user) REFERENCES user(id)
		);`

	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error while creating the user_in_group table: %w", err)
	}
	return nil
}

func CreateTablePrivateChat(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS privconv(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat INTEGER NOT NULL,
			user1 INTEGER NOT NULL,
			user2 INTEGER NOT NULL,
			FOREIGN KEY (chat) REFERENCES chat(id),
			FOREIGN KEY (user1) REFERENCES user(id),
			FOREIGN KEY (user2) REFERENCES user(id)
		);`

	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error while creating the private chat table: %w", err)
	}
	return nil
}
