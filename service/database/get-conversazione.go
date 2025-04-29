package database

import (
	"fmt"
)

func (db *appdbimpl) GetConversazione(utentePassato string, conversazionePassata int) ([]MessageData, error) {
	// Controllo se la chat esiste
	esistenza, err := db.EsisteConversazione(conversazionePassata)
	if err != nil {
		return nil, fmt.Errorf("errore durante la verifica di esistenza: %w", err)
	}
	if !esistenza {
		return nil, fmt.Errorf("errore, la chat non esiste")
	}

	tipoConversazione, err := db.CercaConversazioneGruppo(conversazionePassata)
	if err != nil {
		return nil, fmt.Errorf("errore durante il controllo del tipo di conversazione: %w", err)
	}
	if tipoConversazione > 0 {
		// La conversazione è un gruppo, chiama GetConversazioneGruppo
		return db.GetConversazioneGruppo(utentePassato, conversazionePassata)
	} else {
		nome, err := db.GetNomeUtenteCoinvolto(conversazionePassata, utentePassato)
		if err != nil {
			return nil, fmt.Errorf("errore durante il controllo del tipo di conversazione: %w", err)
		}
		// La conversazione è privata, chiama GetConversazionePrivata
		return db.GetConversazionePrivata(utentePassato, nome)
	}
}

// Creo uno struct per contenere i messaggi
type MessageData struct {
	MessageID int        `json:"message_id"`
	Autore    string     `json:"autore"`
	Text      *string    `json:"text"`
	Foto      *string    `json:"foto"`
	Time      string     `json:"time"`
	Ricevuto  bool       `json:"ricevuto"`
	Letto     bool       `json:"letto"`
	Commenti  []Commento `json:"commenti"`
	Inoltrato bool       `json:"inoltrato"`
	Risposta  *int       `json:"risposta"`
}

// Creo uno struct per contenere i commenti
type Commento struct {
	CommentID int    `json:"comment_id"`
	Autore    string `json:"autore"`
	Reazione  string `json:"reazione"`
}

func (db *appdbimpl) GetConversazionePrivata(utente1_Passato string, utente2_Passato string) ([]MessageData, error) {
	conversazioneID, err := db.EsisteConversazioneTraUtenti(utente1_Passato, utente2_Passato)
	if err != nil {
		return nil, fmt.Errorf("errore durante la ricerca della conversazione: %w", err)
	}

	err = db.LeggiMessaggiPrivati(utente2_Passato, conversazioneID)
	if err != nil {
		return nil, fmt.Errorf("errore durante la modifica dello stato dei messaggi: %w", err)
	}

	querySelect := `
		SELECT m.id, m.autore, m.testo, f.foto, m.tempo, sm.ricevuto, sm.letto, m.inoltrato, m.risposta
		FROM messaggio m
		JOIN statomessaggioprivato as sm on m.id = sm.messaggio 
		LEFT JOIN foto as f on m.foto = f.id
		WHERE m.conversazione = ?
		ORDER BY m.tempo ASC;`

	rows, err := db.c.Query(querySelect, conversazioneID)
	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero dei messaggi: %w", err)
	}

	var messageData []MessageData

	for rows.Next() {
		var messageID int
		var autore string
		var text *string
		var foto *string
		var time string
		var ricevuto bool
		var letto bool
		var inoltrato bool
		var risposta *int

		if err := rows.Scan(&messageID, &autore, &text, &foto, &time, &ricevuto, &letto, &inoltrato, &risposta); err != nil {
			return nil, fmt.Errorf("errore durante la lettura dei dati: %w", err)
		}

		commenti, err := db.GetCommentiMessaggio(messageID)
		if err != nil {
			return nil, fmt.Errorf("errore durante il recupero dei commenti per il messaggio %d: %w", messageID, err)
		}

		messageData = append(messageData, MessageData{
			MessageID: messageID,
			Autore:    autore,
			Text:      text,
			Foto:      foto,
			Time:      time,
			Ricevuto:  ricevuto,
			Letto:     letto,
			Commenti:  commenti,
			Inoltrato: inoltrato,
			Risposta:  risposta,
		})
	}

	return messageData, nil
}

func (db *appdbimpl) GetConversazioneGruppo(utente1_Passato string, id_conversazione int) ([]MessageData, error) {
	coinvolto, err := db.UtenteCoinvoltoGruppo(utente1_Passato, id_conversazione)
	if err != nil {
		return nil, fmt.Errorf("errore durante il controllo della partecipazione dell'utente: %w", err)
	}
	if coinvolto == 0 {
		return nil, fmt.Errorf("l'utente %s non fa parte della conversazione %d", utente1_Passato, id_conversazione)
	}

	err = db.LeggiMessaggiGruppo(utente1_Passato, id_conversazione)
	if err != nil {
		return nil, fmt.Errorf("errore durante la modifica dei messaggi: %w", err)
	}
	err = db.CheckLetturaMessaggiGruppo(id_conversazione)
	if err != nil {
		return nil, fmt.Errorf("errore durante il check di lettura dei messaggi: %w", err)
	}

	querySelect := `
		SELECT m.id, m.autore, m.testo, f.foto, m.tempo, smg.letto, smg.ricevuto, m.inoltrato, m.risposta
		FROM messaggio m
		JOIN statomessaggiogruppo as smg on smg.messaggio = m.id
		LEFT JOIN foto as f on m.foto = f.id
		WHERE m.conversazione = ?
		ORDER BY m.tempo ASC;`

	rows, err := db.c.Query(querySelect, id_conversazione)
	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero dei messaggi: %w", err)
	}

	var messageData []MessageData

	for rows.Next() {
		var messageID int
		var autore string
		var text *string
		var foto *string
		var time string
		var ricevuto bool
		var letto bool
		var inoltrato bool
		var risposta *int

		if err := rows.Scan(&messageID, &autore, &text, &foto, &time, &letto, &ricevuto, &inoltrato, &risposta); err != nil {
			return nil, fmt.Errorf("errore durante la lettura dei dati: %w", err)
		}

		commenti, err := db.GetCommentiMessaggio(messageID)
		if err != nil {
			return nil, fmt.Errorf("errore durante il recupero dei commenti per il messaggio %d: %w", messageID, err)
		}

		messageData = append(messageData, MessageData{
			MessageID: messageID,
			Autore:    autore,
			Text:      text,
			Foto:      foto,
			Time:      time,
			Ricevuto:  ricevuto,
			Letto:     letto,
			Commenti:  commenti,
			Inoltrato: inoltrato,
			Risposta:  risposta,
		})
	}

	return messageData, nil
}

func (db *appdbimpl) GetCommentiMessaggio(messageID int) ([]Commento, error) {
	querySelect := `
		SELECT c.id, c.autore, c.reazione
		FROM commento c
		WHERE c.messaggio = ?;`

	rows, err := db.c.Query(querySelect, messageID)
	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero dei commenti: %w", err)
	}

	var commenti []Commento

	for rows.Next() {
		var commentID int
		var autore string
		var reazione string

		if err := rows.Scan(&commentID, &autore, &reazione); err != nil {
			return nil, fmt.Errorf("errore durante la lettura dei dati del commento: %w", err)
		}

		commenti = append(commenti, Commento{
			CommentID: commentID,
			Autore:    autore,
			Reazione:  reazione,
		})
	}

	return commenti, nil
}
