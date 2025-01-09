package database

import (
	"fmt"
)

// Struttura per rappresentare una chat
type Chat struct {
	ID       int `json:"chatid"`
	Persona1 int `json:"persona1"`
	Persona2 int `json:"persona2"`
}

func (db *appdbimpl) GetChats(nickname string) ([]Chat, error) {
	// Ottieni l'ID dell'utente basato sul nickname
	id, err := db.Get_IdFromNick_Persona(nickname)
	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero dell'ID utente: %w", err)
	}

	// Query per ottenere tutte le conversazioni dell'utente
	query := `
		SELECT id_chat, id_persona1, id_persona2
		FROM conversazioni
		WHERE id_persona1 = ? OR id_persona2 = ?;`

	// Esegui la query
	rows, err := db.c.Query(query, id.Valore, id.Valore)
	if err != nil {
		return nil, fmt.Errorf("errore durante la ricerca delle conversazioni: %w", err)
	}
	defer rows.Close()

	// Raccogli gli ID delle chat
	var chats []Chat
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ID, &chat.Persona1, &chat.Persona2); err != nil {
			return nil, fmt.Errorf("errore durante la scansione del risultato: %w", err)
		}
		chats = append(chats, chat)
	}

	// Controlla eventuali errori al termine della lettura delle righe
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("errore durante l'iterazione delle righe: %w", err)
	}

	return chats, nil
}
