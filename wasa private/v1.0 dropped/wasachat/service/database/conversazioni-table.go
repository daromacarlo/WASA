package database

import (
	"database/sql"
	"fmt"
)

// CreateTableConversazioniPrivate crea la tabella 'conversazioni_private' se non esiste.
func CreateTableConversazioni(db *sql.DB) error {
	// Query per creare la tabella
	query := `
	CREATE TABLE IF NOT EXISTS conversazioni (
		id_chat     INTEGER PRIMARY KEY AUTOINCREMENT, 
		tipo_chat  	STRING NOT NULL,
		id_persona1 INT,
		id_persona2 INT,
		FOREIGN KEY (id_persona1) REFERENCES persone (id),
		FOREIGN KEY (id_persona2) REFERENCES persone (id)
	);`

	// Esegui la query per creare la tabella
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella conversazioni_private: %w", err)
	}

	return nil
}
