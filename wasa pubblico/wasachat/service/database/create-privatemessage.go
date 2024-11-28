package database

import (
	"errors"
	"time"
)

// CreateMessaggioPrivato gestisce la creazione di un messaggio privato tra due utenti.
// Se non esiste una conversazione privata tra i due, viene creata.
func (db *appdbimpl) CreateMessaggioPrivato(nickmittente string, nickdestinatario string, testo string) error {
	// Ottieni l'ID del mittente
	idmittente, err := db.Get_IdFromNick_Persona(nickmittente)
	if err != nil {
		return errors.New("errore durante il recupero dell'ID del mittente: " + err.Error())
	}

	// Ottieni l'ID del destinatario
	iddestinatario, err := db.Get_IdFromNick_Persona(nickdestinatario)
	if err != nil {
		return errors.New("errore durante il recupero dell'ID del destinatario: " + err.Error())
	}

	// Cerca o crea una conversazione privata tra i due utenti
	idchat, err := db.SearchConversazionePrivata(idmittente.Valore, iddestinatario.Valore)
	if err != nil {
		return errors.New("errore durante la ricerca della conversazione privata: " + err.Error())
	}

	// Se la conversazione non esiste, creala
	if !idchat.Valido {
		idchat, err = db.StartConversazionePrivata(idmittente.Valore, iddestinatario.Valore)
		if err != nil {
			return errors.New("errore durante la creazione della conversazione privata: " + err.Error())
		}
	}

	// Query per inserire il messaggio nella tabella messaggi
	query := `INSERT INTO messaggi (chat, testo, foto, dataora, autore) VALUES (?, ?, ?, ?, ?);`

	// Inserisci il messaggio
	_, err = db.c.Exec(query, idchat.Valore, testo, nil, time.Now(), idmittente.Valore)
	if err != nil {
		return errors.New("errore durante l'inserimento del messaggio: " + err.Error())
	}

	// Ritorna nil se tutto è andato a buon fine
	return nil
}
