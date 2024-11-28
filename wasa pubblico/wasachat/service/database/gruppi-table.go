package database

import (
	"database/sql"
	"fmt"
)

// Funzione che crea la tabella conversazioni se questa non esiste
func CreateTableGruppi(db *sql.DB) error {

	// La query per creare la tabella 'conversazioni'
	query := `
	CREATE TABLE IF NOT EXISTS gruppi (
		id 										INTEGER PRIMARY KEY AUTOINCREMENT, 
		nome_gruppo								TEXT NOT NULL,    
		foto									BLOB,
		amministratore							INTEGER NOT NULL,         
		ultimo_messaggio_snippet 				INTEGER,      
		FOREIGN KEY (persona_id) REFERENCES persone (id)
		FOREIGN KEY (amministratore) REFERENCES persone (id)
		FOREIGN KEY (ultimo_messaggio_snippet) REFERENCES messaggio (ID)
	);`

	// Esegui la query per creare la tabella
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella conversazioni: %w", err)
	}
	return nil
}
