package database

import (
	"database/sql"
	"fmt"
	"time"
)

func (db *appdbimpl) CopiaMessaggioCambiandoOraEMitente(idMessaggio int, utente_Passato string, conversazione_Passata int) error {

	var messaggioOriginale struct {
		Testo *string
		Foto  *int
	}

	querySelect := `
        SELECT testo, foto 
        FROM messaggio 
        WHERE id = ?;
    `
	err := db.c.QueryRow(querySelect, idMessaggio).Scan(
		&messaggioOriginale.Testo,
		&messaggioOriginale.Foto,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("nessun messaggio trovato con ID %d", idMessaggio)
		}
		return fmt.Errorf("errore nel recupero del messaggio originale con ID %d: %w", idMessaggio, err)
	}

	queryDiInserimento := `INSERT INTO messaggio (testo, foto, conversazione, autore, tempo, inoltrato, risposta) VALUES (?, ?, ?, ?, ?, ?, ?);`
	result, err := db.c.Exec(queryDiInserimento, messaggioOriginale.Testo, messaggioOriginale.Foto, conversazione_Passata, utente_Passato, time.Now(), true, nil)
	if err != nil {
		return fmt.Errorf("errore durante la creazione del messaggio: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID dell'ultimo messaggio: %w", err)
	}

	isGruppo, err := db.CercaConversazioneGruppo(conversazione_Passata)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	if isGruppo > 0 {
		err = db.CreaStatoMessaggioGruppo(int(lastInsertID))
		if err != nil {
			return fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
		}
		return nil

	} else {
		err = db.CreaStatoMessaggioPrivato(int(lastInsertID))
		if err != nil {
			return fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
		}
		return nil
	}

}

func (db *appdbimpl) InoltraMessaggio(utente_Passato string, idChatNuova int, IdMessaggio int) error {
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utente_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a ID: %w", err)
	}
	isGruppo, err := db.CercaConversazioneGruppo(idChatNuova)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	if isGruppo > 0 {
		coinvolto, err := db.UtenteCoinvoltoGruppo(utente_Passato, idChatNuova)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return fmt.Errorf("l'utente non è membro del gruppo")
		}
	} else {
		idPrivata, err := db.CercaConversazionePrivata(idChatNuova, utente_Passato_convertito)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
		}
		if idPrivata == 0 {
			return fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}
	}

	esistenza, err := db.EsisteMessaggio(IdMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID della conversazione per il messaggio con ID %d: %w", IdMessaggio, err)
	}
	if !esistenza {
		return fmt.Errorf("errore: il messaggio con ID %d non esiste", IdMessaggio)
	}
	conversazioneID, err := db.GetConversazioneIdByMessaggio(IdMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID conversazione: %w", err)
	}
	isGruppo, err = db.CercaConversazioneGruppo(conversazioneID)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	if isGruppo > 0 {
		coinvolto, err := db.UtenteCoinvoltoGruppo(utente_Passato, conversazioneID)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return fmt.Errorf("l'utente non è membro del gruppo")
		}
	} else {
		idPrivata, err := db.CercaConversazionePrivata(conversazioneID, utente_Passato_convertito)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
		}
		if idPrivata == 0 {
			return fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}
	}
	err = db.CopiaMessaggioCambiandoOraEMitente(IdMessaggio, utente_Passato, idChatNuova)
	if err != nil {
		return fmt.Errorf("errore durante la copia del messaggio: %w", err)
	}

	return nil
}
