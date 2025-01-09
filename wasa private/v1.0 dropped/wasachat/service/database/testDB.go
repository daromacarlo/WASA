package database

import (
	"database/sql"
	"fmt"
)

// Struttura per rappresentare un utente
type User struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
}

// Funzione per ottenere tutti gli utenti dal database
func (db *appdbimpl) GetUsers() ([]User, error) {
	query := `
		SELECT id, nickname
		FROM persone
	`
	rows, err := db.c.Query(query)
	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero degli utenti: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Nickname); err != nil {
			return nil, fmt.Errorf("errore durante la scansione del risultato: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("errore durante l'iterazione delle righe: %w", err)
	}

	return users, nil
}

// Struttura per rappresentare un messaggio
type Message struct {
	ID      int          `json:"id"`
	ChatID  int          `json:"chatid"`
	Testo   string       `json:"testo"`
	DataOra sql.NullTime `json:"dataora"`
	Autore  int          `json:"autore"`
}

// Funzione per ottenere tutti i messaggi dal database
func (db *appdbimpl) GetMessaggi() ([]Message, error) {
	query := `
		SELECT id, chat, testo, autore, dataora
		FROM messaggi
	`
	rows, err := db.c.Query(query)
	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero dei messaggi: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.ChatID, &message.Testo, &message.Autore, &message.DataOra); err != nil {
			return nil, fmt.Errorf("errore durante la scansione del risultato: %w", err)
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("errore durante l'iterazione delle righe: %w", err)
	}

	return messages, nil
}
