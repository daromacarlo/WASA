package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// Funzione che crea la tabella conversazione se questa non esiste
func CreaTabellaConversazione(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS conversazione(
			id INTEGER PRIMARY KEY AUTOINCREMENT
	 	);`

	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("errore durante la creazione della tabella conversazione: %w", err)
	}
	return nil
}

// Funzione che crea la tabella gruppo se questa non esiste
func CreaTabellaGruppo(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS gruppo(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nome STRING NOT NULL,
			conversazione UNIQUE NOT NULL,
			foto BLOB NOT NULL,
			FOREIGN KEY (conversazione) REFERENCES conversazione(id)
		);`

	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("errore durante la creazione della tabella gruppo: %w", err)
	}
	return nil
}

// Funzione che crea la tabella utenteInGruppo se questa non esiste
func CreaTabellaUtenteingruppo(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS utenteingruppo(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			gruppo INTEGER NOT NULL,
			utente INTEGER NOT NULL,
			FOREIGN KEY (gruppo) REFERENCES gruppo(id),
			FOREIGN KEY (utente) REFERENCES utente(id)
		);`

	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("errore durante la creazione della tabella utenteingruppo: %w", err)
	}
	return nil
}

// Funzione che crea la tabella chatprivata se questa non esiste
func CreaTabellaConversazionePrivata(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS conversazioneprivata(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			conversazione INTEGER NOT NULL,
			utente1 INTEGER NOT NULL,
			utente2 INTEGER NOT NULL,
			FOREIGN KEY (conversazione) REFERENCES conversazione(id),
			FOREIGN KEY (utente1) REFERENCES utente(id),
			FOREIGN KEY (utente2) REFERENCES utente(id)
		);`

	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("errore durante la creazione della tabella chatprivata: %w", err)
	}
	return nil
}
