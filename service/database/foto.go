package database

import (
	"database/sql"
	"fmt"
)

// Funzione che crea la tabella foto se questa non esiste
func CreaTabellaFoto(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS foto(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			foto BLOB NOT NULL
		);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella foto: %w", err)
	}
	return nil
}

// CreaFoto crea una nuova foto nel database con percorso "percorso_Passato" e foto "foto_Passata", la foto avra un nuovo id diverso dagli altri (autoincrementante)
func (db *appdbimpl) CreaFoto(foto_Passata string) (int, error) {
	queryDiInserimento := `INSERT INTO foto (foto) VALUES (?);`
	result, err := db.c.Exec(queryDiInserimento, foto_Passata)
	if err != nil {
		return 0, fmt.Errorf("errore durante la creazione della foto: %w", err)
	}
	ultimoIdInserito, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("errore durante il recupero dell'ID dell'ultimo elemento inserito: %w", err)
	}
	return int(ultimoIdInserito), nil

}
