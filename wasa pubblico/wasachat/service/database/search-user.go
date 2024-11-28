package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) SearchUser(nick string) (bool, error) {
	var nickTrovato string
	err := db.c.QueryRow(`SELECT nickname FROM persone WHERE nickname = ?`, nick).Scan(&nickTrovato)

	// Se l'errore è sql.ErrNoRows, significa che non è stato trovato nessun utente con quel nickname
	if err == sql.ErrNoRows {
		// Se non esistono utenti con il nickname
		return false, nil
	}

	// Se si verifica un altro errore, restituire l'errore del database
	if err != nil {
		return false, fmt.Errorf("errore database: %w", err)
	}

	// Se la query ha trovato un risultato, significa che il nickname esiste già
	return true, nil
}
