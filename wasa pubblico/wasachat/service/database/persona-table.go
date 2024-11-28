package database

import (
	"database/sql"
	"fmt"
)

// Funzione che crea la tabella persona se questa non esiste
func CreateTablePersone(db *sql.DB) error {

	// La query per creare la tabella 'persone'
	query := `
	CREATE TABLE IF NOT EXISTS persone (
		id 				INTEGER PRIMARY KEY AUTOINCREMENT,
		nickname 		TEXT UNIQUE NOT NULL,
		foto 			BLOB
	);`

	// Esegui la query per creare la tabella
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella persone: %w", err)
	}
	return nil
}
