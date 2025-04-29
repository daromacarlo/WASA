package database

import (
	"fmt"
)

// ImpostaRisposta imposta l'attributo 'risposta' nel messaggio con l'ID del vecchio messaggio
func (db *appdbimpl) ImpostaRisposta(idMessaggio int, idNuovoMessaggio int) error {
	query := `
		UPDATE messaggio
		SET risposta = ?
		WHERE id = ?
	`
	_, err := db.c.Exec(query, idMessaggio, idNuovoMessaggio)
	if err != nil {
		return fmt.Errorf("errore inaspettato durante l'impostazione della risposta per il messaggio con ID %d: %w", idNuovoMessaggio, err)
	}
	return nil
}

// verificaConversazione verifica tutti i controlli comuni per entrambi i tipi di messaggio
func (db *appdbimpl) verificaConversazione(utentePassato string, idConversazionePassato int, idMessaggioRisposta int) (int, error) {
	// Verifica se il messaggio con l'ID passato esiste
	esistenza, err := db.EsisteMessaggio(idMessaggioRisposta)
	if err != nil {
		return 0, fmt.Errorf("errore durante la verifica dell'esistenza del messaggio con ID %d: %w", idMessaggioRisposta, err)
	}
	if !esistenza {
		return 0, fmt.Errorf("errore: il messaggio con ID %d non esiste", idMessaggioRisposta)
	}

	// Recupera l'ID della conversazione a partire dal messaggio
	conversazioneID, err := db.GetConversazioneIdByMessaggio(idMessaggioRisposta)
	if err != nil {
		return 0, fmt.Errorf("errore durante il recupero dell'ID della conversazione per il messaggio ID %d: %w", idMessaggioRisposta, err)
	}

	// Controlla che il messaggio sia della conversazione corrente
	if conversazioneID != idConversazionePassato {
		return 0, fmt.Errorf("il messaggio non fa parte della conversazione indicata (ID: %d)", idConversazionePassato)
	}

	// Recupera l'ID dell'utente a partire dal nickname
	utenteID, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return 0, fmt.Errorf("errore durante la conversione del nickname '%s' a ID utente: %w", utentePassato, err)
	}

	// Verifica se la conversazione è un gruppo o privata
	isGruppo, err := db.CercaConversazioneGruppo(conversazioneID)
	if err != nil {
		return 0, fmt.Errorf("errore durante la verifica del tipo di conversazione per ID %d: %w", conversazioneID, err)
	}

	if isGruppo > 0 {
		// Se è un gruppo, verifica che l'utente sia coinvolto
		coinvolto, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazioneID)
		if err != nil {
			return 0, fmt.Errorf("errore durante la verifica della partecipazione dell'utente '%s' al gruppo: %w", utentePassato, err)
		}
		if coinvolto == 0 {
			return 0, fmt.Errorf("l'utente '%s' non è membro del gruppo con ID %d", utentePassato, conversazioneID)
		}
	} else {
		// Se è privata, verifica che l'utente sia coinvolto nella conversazione privata
		idPrivata, err := db.CercaConversazionePrivata(conversazioneID, utenteID)
		if err != nil {
			return 0, fmt.Errorf("errore durante la verifica della conversazione privata per l'utente '%s' (ID: %d): %w", utentePassato, utenteID, err)
		}
		if idPrivata == 0 {
			return 0, fmt.Errorf("l'utente '%s' non è coinvolto nella conversazione privata con ID %d", utentePassato, conversazioneID)
		}
	}

	return conversazioneID, nil
}

// RispondiMessaggioTesto gestisce la risposta a un messaggio testuale
func (db *appdbimpl) RispondiMessaggioTesto(utentePassato string, idConversazionePassato int, idMessaggio int, testoPassato string) error {
	// Verifica la conversazione
	_, err := db.verificaConversazione(utentePassato, idConversazionePassato, idMessaggio)
	if err != nil {
		return err
	}

	// Crea il nuovo messaggio testuale nel database e ottieni direttamente l'ID
	idNuovoMessaggio, err := db.CreaMessaggioTestualeDB(utentePassato, idConversazionePassato, testoPassato)
	if err != nil {
		return fmt.Errorf(" testuale per l'utente '%s': %w", utentePassato, err)
	}

	// Imposta il nuovo messaggio come risposta al messaggio precedente
	err = db.ImpostaRisposta(idMessaggio, idNuovoMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante l'impostazione della risposta per il messaggio con ID %d (risposta a %d): %w", idNuovoMessaggio, idMessaggio, err)
	}

	return nil
}

// RispondiMessaggioFoto gestisce la risposta a un messaggio contenente una foto
func (db *appdbimpl) RispondiMessaggioFoto(utentePassato string, idConversazionePassato int, idMessaggio int, fotoPassato int) error {
	// Verifica la conversazione
	_, err := db.verificaConversazione(utentePassato, idConversazionePassato, idMessaggio)
	if err != nil {
		return err
	}

	// Crea il nuovo messaggio con foto nel database e ottieni direttamente l'ID
	idNuovoMessaggio, err := db.CreaMessaggioFotoDB(utentePassato, idConversazionePassato, fotoPassato)
	if err != nil {
		return fmt.Errorf(" foto per l'utente '%s': %w", utentePassato, err)
	}

	// Imposta il nuovo messaggio come risposta al messaggio precedente
	err = db.ImpostaRisposta(idMessaggio, idNuovoMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante l'impostazione della risposta per il messaggio con ID %d (risposta a %d): %w", idNuovoMessaggio, idMessaggio, err)
	}

	return nil
}
