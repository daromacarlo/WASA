package database

import (
	"fmt"
)

func (db *appdbimpl) StartConversazionePrivata(id_utente1 int, id_utente2 int) (interoNullabile, error) {
	query := `INSERT INTO conversazioni (tipo_chat, id_persona1, id_persona2) VALUES (?, ?, ?);`

	result, err := db.c.Exec(query, "privata", id_utente1, id_utente2)
	if err != nil {
		return interoNullabile{Valido: false}, fmt.Errorf("errore durante la creazione della conversazione: %w", err)
	}

	idConversazione, err := result.LastInsertId() //lastinsertid() ci consente di recuperare l'ultimo id inserito nel database
	if err != nil {
		return interoNullabile{Valido: false}, fmt.Errorf("errore nel recupero dell'ID della conversazione: %w", err)
	}

	return interoNullabile{Valore: int(idConversazione), Valido: true}, nil
	//bisogna castare il risultato di LastInsertId perché è troppo grande il valore che tira fuori
}
