package database

import (
	"fmt"
)

// ImpostaRisposta imposta l'attributo 'risposta' nel messaggio con l'ID del vecchio messaggio
func (db *appdbimpl) ImpostaRisposta(IdMessaggio int, IdNuovoMessaggio int, chat int) error {
	// Esegui un'operazione di aggiornamento sulla tabella del messaggio
	query := `
		UPDATE messaggio
		SET risposta = ?
		WHERE id = ? AND conversazione = ?
	`
	_, err := db.c.Exec(query, IdMessaggio, IdNuovoMessaggio, chat)
	if err != nil {
		return fmt.Errorf("errore inaspettato durante l'impostazione della risposta per il messaggio con ID %d: %w", IdNuovoMessaggio, err)
	}
	return nil
}

// RispondiMessaggioTesto gestisce la risposta a un messaggio testuale
func (db *appdbimpl) RispondiMessaggioTesto(utentePassato string, idConversazionePassato int, idMessaggio int, testoPassato string) error {
	// Verifica se il messaggio con l'ID passato esiste
	esistenza, err := db.EsisteMessaggio(idMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante la verifica dell'esistenza del messaggio con ID %d: %w", idMessaggio, err)
	}
	if !esistenza {
		return fmt.Errorf("errore: il messaggio con ID %d non esiste", idMessaggio)
	}

	// Recupera l'ID della conversazione a partire dal messaggio
	conversazioneID, err := db.GetConversazioneIdByMessaggio(idMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID della conversazione: %w", err)
	}

	// Controlla che il messaggio sia della conversazione corrente
	if conversazioneID != idConversazionePassato {
		return fmt.Errorf("il messaggio non fa parte della chat indicata")
	}

	// Recupera l'ID dell'utente a partire dal nickname
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione del nickname a ID utente: %w", err)
	}

	// Verifica se la conversazione è un gruppo o privata
	isGruppo, err := db.CercaConversazioneGruppo(conversazioneID)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}

	if isGruppo > 0 {
		// Se è un gruppo, verifica che l'utente sia coinvolto
		coinvolto, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazioneID)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return fmt.Errorf("l'utente non è membro del gruppo")
		}
	} else {
		// Se è privata, verifica che l'utente sia coinvolto nella conversazione privata
		idPrivata, err := db.CercaConversazionePrivata(conversazioneID, utente_Passato_convertito)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
		}
		if idPrivata == 0 {
			return fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}
	}

	// Crea il nuovo messaggio testuale nel database
	err = db.CreaMessaggioTestualeDB(utentePassato, idConversazionePassato, testoPassato)
	if err != nil {
		return fmt.Errorf("errore durante la creazione del messaggio: %w", err)
	}

	// Recupera l'ID del nuovo messaggio appena creato
	var idNuovoMessaggio int
	queryID := `SELECT last_insert_rowid()` // Recupera l'ID dell'ultimo messaggio inserito
	err = db.c.QueryRow(queryID).Scan(&idNuovoMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID del nuovo messaggio: %w", err)
	}

	// Imposta il nuovo messaggio come risposta al messaggio precedente
	err = db.ImpostaRisposta(idMessaggio, idNuovoMessaggio, idConversazionePassato)
	if err != nil {
		return fmt.Errorf("errore durante l'impostazione della risposta per il messaggio con ID %d: %w", idNuovoMessaggio, err)
	}

	return nil
}

// RispondiMessaggioFoto gestisce la risposta a un messaggio contenente una foto
func (db *appdbimpl) RispondiMessaggioFoto(utentePassato string, idConversazionePassato int, IdMessaggio int, foto_Passato int) error {
	// Verifica se il messaggio con l'ID passato esiste
	esistenza, err := db.EsisteMessaggio(IdMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante la verifica dell'esistenza del messaggio con ID %d: %w", IdMessaggio, err)
	}
	if !esistenza {
		return fmt.Errorf("errore: il messaggio con ID %d non esiste", IdMessaggio)
	}

	// Recupera l'ID della conversazione a partire dal messaggio
	conversazioneID, err := db.GetConversazioneIdByMessaggio(IdMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID della conversazione: %w", err)
	}

	// Controlla che il messaggio sia della conversazione corrente
	if conversazioneID != idConversazionePassato {
		return fmt.Errorf("il messaggio non fa parte della chat indicata")
	}

	// Recupera l'ID dell'utente a partire dal nickname
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione del nickname a ID utente: %w", err)
	}

	// Verifica se la conversazione è un gruppo o privata
	isGruppo, err := db.CercaConversazioneGruppo(conversazioneID)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}

	if isGruppo > 0 {
		// Se è un gruppo, verifica che l'utente sia coinvolto
		coinvolto, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazioneID)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return fmt.Errorf("l'utente non è membro del gruppo")
		}
	} else {
		// Se è privata, verifica che l'utente sia coinvolto nella conversazione privata
		idPrivata, err := db.CercaConversazionePrivata(conversazioneID, utente_Passato_convertito)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
		}
		if idPrivata == 0 {
			return fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}
	}

	// Crea il nuovo messaggio con foto nel database
	err = db.CreaMessaggioFotoDB(utentePassato, idConversazionePassato, foto_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la creazione del messaggio: %w", err)
	}

	// Recupera l'ID del nuovo messaggio appena creato
	var idNuovoMessaggio int
	queryID := `SELECT last_insert_rowid()` // Recupera l'ID dell'ultimo messaggio inserito
	err = db.c.QueryRow(queryID).Scan(&idNuovoMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID del nuovo messaggio: %w", err)
	}

	// Imposta il nuovo messaggio come risposta al messaggio precedente
	err = db.ImpostaRisposta(IdMessaggio, idNuovoMessaggio, idConversazionePassato)
	if err != nil {
		return fmt.Errorf("errore durante l'impostazione della risposta per il messaggio con ID %d: %w", idNuovoMessaggio, err)
	}

	return nil
}
