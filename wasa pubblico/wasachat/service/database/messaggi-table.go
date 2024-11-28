package database

import (
	"database/sql"
	"fmt"
)

// Funzione che crea la tabella messaggi se questa non esiste
func CreateTableMessaggi(db *sql.DB) error {

	// Query per creare la tabella 'messaggi'
	query := `
	CREATE TABLE IF NOT EXISTS messaggi (
		id INTEGER PRIMARY KEY AUTOINCREMENT,         
		chat INTEGER NOT NULL,
		testo TEXT,
		foto BLOB,
		autore INTEGER NOT NULL,
		dataora DATETIME,
		FOREIGN KEY (autore) REFERENCES persone (id) 
	);`

	// Esegui la query per creare la tabella
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella messaggi: %w", err)
	}
	return nil
}
