package database

import (
	"database/sql"
	"fmt"
	"time"
)

func (db *appdbimpl) EsisteMessaggio(idMessaggio int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM messaggio WHERE id = ?;"
	err := db.c.QueryRow(query, idMessaggio).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("errore durante la verifica dell'esistenza del messaggio con ID %d: %w", idMessaggio, err)
	}
	return count > 0, nil
}

func (db *appdbimpl) GetConversazioneIdByMessaggio(idMessaggio int) (int, error) {
	// Query per ottenere l'ID della conversazione dal messaggio
	query := `
		SELECT conversazione
		FROM messaggio
		WHERE id = ?;
	`
	var conversazioneID int
	err := db.c.QueryRow(query, idMessaggio).Scan(&conversazioneID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Maggiore chiarezza nel messaggio d'errore
			return 0, fmt.Errorf("nessun messaggio trovato con l'ID %d", idMessaggio)
		}
		return 0, fmt.Errorf("errore durante il recupero dell'ID della conversazione: %w", err)
	}

	return conversazioneID, nil
}

func (db *appdbimpl) CopiaMessaggioCambiandoOraEMitente(idMessaggio int, utente_Passato int, conversazione_Passata int) error {

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
			// Se il messaggio non esiste, restituisci errore
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
func (db *appdbimpl) InoltraMessaggioPrivato(utente_Passato string, destinatario_Passato string, IdMessaggio int) error {
	// Verifica che l'utente e il destinatario siano coinvolti in una conversazione privata
	// Verifica se il messaggio esiste

	esistenza, err := db.EsisteMessaggio(IdMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID della conversazione per il messaggio con ID %d: %w", IdMessaggio, err)
	}
	if !esistenza {
		return fmt.Errorf("errore: il messaggio con ID %d non esiste", IdMessaggio)
	}

	// Recupera l'ID della conversazione a partire dal messaggio
	conversazioneID, err := db.GetConversazioneIdByMessaggio(IdMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID conversazione: %w", err)
	}

	// Recupera l'ID dell'utente a partire dal nickname
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utente_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a ID: %w", err)
	}

	// Verifica se la conversazione è un gruppo o privata
	isGruppo, err := db.CercaConversazioneGruppo(conversazioneID)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	if isGruppo > 0 {
		// Se è un gruppo, verifica che l'utente sia coinvolto
		coinvolto, err := db.UtenteCoinvoltoGruppo(utente_Passato, conversazioneID)
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

	chat, err := db.UtenteCoinvoltoPrivato(utente_Passato, destinatario_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
	}

	isGruppo, err = db.CercaConversazioneGruppo(chat)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	if isGruppo > 0 {
		return fmt.Errorf("non hai passato una conversazione privata: %w", err)
	}

	err = db.CopiaMessaggioCambiandoOraEMitente(IdMessaggio, utente_Passato_convertito, chat)
	if err != nil {
		return fmt.Errorf("errore durante la copia del messaggio: %w", err)
	}
	return nil

}

func (db *appdbimpl) InoltraMessaggioGruppo(utente_Passato string, idChatNuova int, IdMessaggio int) error {
	// Verifica se il messaggio esiste

	esistenza, err := db.EsisteMessaggio(IdMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID della conversazione per il messaggio con ID %d: %w", IdMessaggio, err)
	}
	if !esistenza {
		return fmt.Errorf("errore: il messaggio con ID %d non esiste", IdMessaggio)
	}

	// Recupera l'ID della conversazione a partire dal messaggio
	conversazioneID, err := db.GetConversazioneIdByMessaggio(IdMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID conversazione: %w", err)
	}

	// Recupera l'ID dell'utente a partire dal nickname
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utente_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a ID: %w", err)
	}

	// Verifica se la conversazione è un gruppo o privata
	isGruppo, err := db.CercaConversazioneGruppo(conversazioneID)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	if isGruppo > 0 {
		// Se è un gruppo, verifica che l'utente sia coinvolto
		coinvolto, err := db.UtenteCoinvoltoGruppo(utente_Passato, conversazioneID)
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
	isGruppo, err = db.CercaConversazioneGruppo(idChatNuova)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	if isGruppo == 0 {
		return fmt.Errorf("non hai passato una conversazione di gruppo: %w", err)
	}

	// Copia il messaggio cambiando l'ora e il mittente
	err = db.CopiaMessaggioCambiandoOraEMitente(IdMessaggio, utente_Passato_convertito, idChatNuova)
	if err != nil {
		return fmt.Errorf("errore durante la copia del messaggio: %w", err)
	}

	return nil
}
