package database

import (
	"fmt"
)

// CreateUser crea un nuovo utente nel database con il nickname passato ('nick')
func (db *appdbimpl) CreateUser(nick string) error {

	// La query per inserire un nuovo utente
	query := `INSERT INTO persone (nickname) VALUES (?);`

	// Controlliamo che non esista già un utente con lo stesso nome e che non ci siano errori
	_, err := db.c.Exec(query, nick)
	if err != nil {
		return fmt.Errorf("errore durante la creazione dell'utente: %w", err)
	}
	return nil
}
