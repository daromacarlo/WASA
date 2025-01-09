package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) Get_IdFromNick_Persona(nick string) (interoNullabile, error) {
	var idTrovato int
	err := db.c.QueryRow(`SELECT id FROM persone WHERE nickname = ?`, nick).Scan(&idTrovato)

	// Se l'errore è sql.ErrNoRows, significa che non è stato trovato nessun utente con quel nickname
	if err == sql.ErrNoRows {
		// Se non esistono utenti con il nickname
		return interoNullabile{Valido: false}, fmt.Errorf("persona non trovata nel database: %w", err)
	}

	// Se si verifica un altro errore, restituire l'errore del database
	if err != nil {
		return interoNullabile{Valido: false}, fmt.Errorf("errore database: %w", err)
	}

	// Se la query ha trovato un risultato, significa che il nickname esiste già
	return interoNullabile{Valore: idTrovato, Valido: true}, nil
}

func (db *appdbimpl) Get_NickFromId_Persona(id int) (stringaNullabile, error) {
	var nickTrovato string
	err := db.c.QueryRow(`SELECT nickname FROM persone WHERE id = ?`, id).Scan(&nickTrovato)

	// Se l'errore è sql.ErrNoRows, significa che non è stato trovato nessun utente con quel nickname
	if err == sql.ErrNoRows {
		// Se non esistono utenti con il nickname
		return stringaNullabile{Valido: false}, fmt.Errorf("persona non trovata nel database: %w", err)
	}

	// Se si verifica un altro errore, restituire l'errore del database
	if err != nil {
		return stringaNullabile{Valido: false}, fmt.Errorf("errore database: %w", err)
	}

	// Se la query ha trovato un risultato, significa che il nickname esiste già
	return stringaNullabile{Valore: nickTrovato, Valido: true}, nil
}
