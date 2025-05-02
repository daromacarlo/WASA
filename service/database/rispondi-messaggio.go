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
func (db *appdbimpl) verificaConversazione(utentePassato string, idConversazionePassato int, idMessaggioRisposta int) (int, int, error) {
	esistenza, err := db.EsisteMessaggio(idMessaggioRisposta)
	if err != nil {
		return 0, 500, fmt.Errorf("errore durante la verifica dell'esistenza del messaggio con ID %d: %w", idMessaggioRisposta, err)
	}
	if !esistenza {
		return 0, 404, fmt.Errorf("errore: il messaggio con ID %d non esiste", idMessaggioRisposta)
	}
	conversazioneID, err := db.GetConversazioneIdByMessaggio(idMessaggioRisposta)
	if err != nil {
		return 0, 500, fmt.Errorf("errore durante il recupero dell'ID della conversazione per il messaggio ID %d: %w", idMessaggioRisposta, err)
	}
	if conversazioneID != idConversazionePassato {
		return 0, 401, fmt.Errorf("il messaggio non fa parte della conversazione indicata (ID: %d)", idConversazionePassato)
	}
	utenteID, codiceErrore, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return 0, codiceErrore, fmt.Errorf("errore durante la conversione del nickname '%s' a ID utente: %w", utentePassato, err)
	}
	isGruppo, codiceErrore, err := db.CercaConversazioneGruppo(conversazioneID)
	if err != nil {
		return 0, codiceErrore, fmt.Errorf("errore durante la verifica del tipo di conversazione per ID %d: %w", conversazioneID, err)
	}
	if isGruppo > 0 {
		coinvolto, codiceErrore, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazioneID)
		if err != nil {
			return 0, codiceErrore, fmt.Errorf("errore durante la verifica della partecipazione dell'utente '%s' al gruppo: %w", utentePassato, err)
		}
		if coinvolto == 0 {
			return 0, 401, fmt.Errorf("l'utente '%s' non è membro del gruppo con ID %d", utentePassato, conversazioneID)
		}
	} else {
		idPrivata, codiceErrore, err := db.CercaConversazionePrivata(conversazioneID, utenteID)
		if err != nil {
			return 0, codiceErrore, fmt.Errorf("errore durante la verifica della conversazione privata per l'utente '%s' (ID: %d): %w", utentePassato, utenteID, err)
		}
		if idPrivata == 0 {
			return 0, 401, fmt.Errorf("l'utente '%s' non è coinvolto nella conversazione privata con ID %d", utentePassato, conversazioneID)
		}
	}

	return conversazioneID, 0, nil
}

// RispondiMessaggioTesto gestisce la risposta a un messaggio testuale
func (db *appdbimpl) RispondiMessaggioTesto(utentePassato string, idConversazionePassato int, idMessaggio int, testoPassato string) (int, error) {
	_, codiceErrore, err := db.verificaConversazione(utentePassato, idConversazionePassato, idMessaggio)
	if err != nil {
		return codiceErrore, err
	}
	idNuovoMessaggio, codiceErrore, err := db.CreaMessaggioTestualeDB(utentePassato, idConversazionePassato, testoPassato)
	if err != nil {
		return codiceErrore, fmt.Errorf("errore durante l'immissione del messaggio testuale per l'utente '%s': %w", utentePassato, err)
	}
	err = db.ImpostaRisposta(idMessaggio, idNuovoMessaggio)
	if err != nil {
		return 500, fmt.Errorf("errore durante l'impostazione della risposta per il messaggio con ID %d (risposta a %d): %w", idNuovoMessaggio, idMessaggio, err)
	}
	return 0, nil
}

func (db *appdbimpl) RispondiMessaggioFoto(utentePassato string, idConversazionePassato int, idMessaggio int, fotoPassato int) (int, error) {
	_, codiceErrore, err := db.verificaConversazione(utentePassato, idConversazionePassato, idMessaggio)
	if err != nil {
		return codiceErrore, err
	}
	idNuovoMessaggio, codiceErrore, err := db.CreaMessaggioFotoDB(utentePassato, idConversazionePassato, fotoPassato)
	if err != nil {
		return codiceErrore, fmt.Errorf("errore durante l'immissione dell foto per l'utente '%s': %w", utentePassato, err)
	}
	err = db.ImpostaRisposta(idMessaggio, idNuovoMessaggio)
	if err != nil {
		return 500, fmt.Errorf("errore durante l'impostazione della risposta per il messaggio con ID %d (risposta a %d): %w", idNuovoMessaggio, idMessaggio, err)
	}
	return 0, nil
}
