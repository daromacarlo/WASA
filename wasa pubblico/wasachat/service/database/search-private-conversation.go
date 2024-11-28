package database

import (
	"database/sql"
	"fmt"
)

// SearchConversazionePrivata cerca una conversazione privata tra due utenti nel database.
func (db *appdbimpl) SearchConversazionePrivata(id_utente1 int, id_utente2 int) (interoNullabile, error) {
	query := `
	SELECT id_chat 
	FROM conversazioni
	WHERE (id_persona1 = ? AND id_persona2 = ?) 
	   OR (id_persona1 = ? AND id_persona2 = ?);`

	var idConversazione int
	err := db.c.QueryRow(query, id_utente1, id_utente2, id_utente2, id_utente1).Scan(&idConversazione)

	if err != nil {
		if err == sql.ErrNoRows {
			// Nessuna conversazione trovata
			return interoNullabile{Valido: false}, nil
		}
		// Errore generico nella query
		return interoNullabile{Valido: false}, fmt.Errorf("errore durante la ricerca della conversazione: %w", err)
	}

	// Conversazione trovata
	return interoNullabile{Valore: idConversazione, Valido: true}, nil
}
