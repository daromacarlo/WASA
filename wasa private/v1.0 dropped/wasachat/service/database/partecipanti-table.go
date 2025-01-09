package database

import (
	"database/sql"
	"fmt"
)

// Funzione che crea la tabella conversazioni se questa non esiste
func CreateTablePartecipanti(db *sql.DB) error {

	// La query per creare la tabella 'conversazioni'
	query := `
	CREATE TABLE IF NOT EXISTS gruppi (
		id INTEGER 				PRIMARY KEY AUTOINCREMENT,
		id_conversazione 		INTEGER NOT NULL,
		id_persona 				INTEGER NOT NULL,
		FOREIGN KEY (id_conversazione) REFERENCES conversazioni (id),
		FOREIGN KEY (id_persona) REFERENCES persone (id)
	);`

	// Esegui la query per creare la tabella
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella conversazioni: %w", err)
	}
	return nil
}
