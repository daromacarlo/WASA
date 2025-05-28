package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (db *appdbimpl) CopiaMessaggioCambiandoOraEMitente(idMessaggio int, utente_Passato string, conversazione_Passata int) (int, error) {

	utente_Passato_convertito, _, err := db.IdUtenteDaNickname(utente_Passato)
	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("errore durante la conversione da nickname a id: %w", err)
	}

	var messaggioOriginale struct {
		Testo *string
		Foto  *int
	}

	querySelect := `
        SELECT testo, foto 
        FROM messaggio 
        WHERE id = ?;
    `
	err = db.c.QueryRow(querySelect, idMessaggio).Scan(
		&messaggioOriginale.Testo,
		&messaggioOriginale.Foto,
	)

	if !errors.Is(err, nil) {
		if errors.Is(err, sql.ErrNoRows) {
			return 404, fmt.Errorf("nessun messaggio trovato con ID %d", idMessaggio)
		}
		return 500, fmt.Errorf("errore nel recupero del messaggio originale con ID %d: %w", idMessaggio, err)
	}

	queryDiInserimento := `INSERT INTO messaggio (testo, foto, conversazione, autore, idautore, tempo, inoltrato, risposta) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	result, err := db.c.Exec(queryDiInserimento, messaggioOriginale.Testo, messaggioOriginale.Foto, conversazione_Passata, utente_Passato, utente_Passato_convertito, time.Now(), true, nil)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante l'inserimento del messaggio modificato: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il recupero dell'ID dell'ultimo messaggio: %w", err)
	}

	isGruppo, codiceErrore, err := db.CercaConversazioneGruppo(conversazione_Passata)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	if isGruppo > 0 {
		err = db.CreaStatoMessaggioGruppo(int(lastInsertID))
		if !errors.Is(err, nil) {
			return 500, fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
		}
		return 0, nil
	} else {
		err = db.CreaStatoMessaggioPrivato(int(lastInsertID))
		if !errors.Is(err, nil) {
			return 500, fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
		}
		return 0, nil
	}
}

func (db *appdbimpl) InoltraMessaggio(utente_Passato string, idChatNuova int, IdMessaggio int) (int, error) {
	utente_Passato_convertito, codiceErrore, err := db.IdUtenteDaNickname(utente_Passato)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante la conversione da nickname a ID: %w", err)
	}
	isGruppo, codiceErrore, err := db.CercaConversazioneGruppo(idChatNuova)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	if isGruppo > 0 {
		coinvolto, codiceErrore, err := db.UtenteCoinvoltoGruppo(utente_Passato, idChatNuova)
		if !errors.Is(err, nil) {
			return codiceErrore, fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return 401, fmt.Errorf("l'utente non è membro del gruppo")
		}
	}

	esistenza, err := db.EsisteMessaggio(IdMessaggio)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il recupero dell'ID della conversazione per il messaggio con ID %d: %w", IdMessaggio, err)
	}
	if !esistenza {
		return 404, fmt.Errorf("errore: il messaggio con ID %d non esiste", IdMessaggio)
	}
	conversazioneID, err := db.GetConversazioneIdByMessaggio(IdMessaggio)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il recupero dell'ID conversazione: %w", err)
	}
	isGruppo, codiceErrore, err = db.CercaConversazioneGruppo(conversazioneID)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	if isGruppo > 0 {
		coinvolto, codiceErrore, err := db.UtenteCoinvoltoGruppo(utente_Passato, conversazioneID)
		if !errors.Is(err, nil) {
			return codiceErrore, fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return 401, fmt.Errorf("l'utente non è membro del gruppo")
		}
	} else {
		idPrivata, codiceErrore, err := db.CercaConversazionePrivata(conversazioneID, utente_Passato_convertito)
		if !errors.Is(err, nil) {
			return codiceErrore, fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
		}
		if idPrivata == 0 {
			return codiceErrore, fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}
	}
	codiceErrore, err = db.CopiaMessaggioCambiandoOraEMitente(IdMessaggio, utente_Passato, idChatNuova)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante la copia del messaggio: %w", err)
	}

	return 0, nil
}

func (db *appdbimpl) InoltraMessaggioANuovaChat(utente_Passato string, utente2_Passato string, IdMessaggio int) (int, error) {
	if utente2_Passato == "" {
		return 400, fmt.Errorf("richiesta mal formata, l'utente deve avere un nome")
	}

	chat, codiceErrore, err := db.EsisteConversazioneTraUtenti(utente_Passato, utente2_Passato)
	if chat > 0 {
		return 404, fmt.Errorf("chat già esistente tra i due utenti, non è stato inviato il messaggio con ID %d", IdMessaggio)
	}
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore inaspettato %w", err)
	}
	utente_Passato_convertito, codiceErrore, err := db.IdUtenteDaNickname(utente_Passato)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante la conversione da nickname a ID: %w", err)
	}
	utente2_Passato_convertito, codiceErrore, err := db.IdUtenteDaNickname(utente2_Passato)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante la conversione da nickname a ID: %w", err)
	}
	esistenza, err := db.EsistenzaUtente(utente_Passato)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il recupero dell'ID della conversazione per il messaggio con ID %d: %w", IdMessaggio, err)
	}
	if !esistenza {
		return 404, fmt.Errorf("l'utente con ID %d non esiste %w", utente_Passato_convertito, err)
	}
	esistenza, err = db.EsistenzaUtente(utente2_Passato)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il recupero dell'ID della conversazione per il messaggio con ID %d: %w", IdMessaggio, err)
	}
	if !esistenza {
		return 404, fmt.Errorf("l'utente con ID %d non esiste %w", utente2_Passato_convertito, err)
	}
	esistenza, err = db.EsisteMessaggio(IdMessaggio)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il recupero dell'ID della conversazione per il messaggio con ID %d: %w", IdMessaggio, err)
	}
	if !esistenza {
		return 404, fmt.Errorf("messaggio con ID %d non esiste %w", IdMessaggio, err)
	}
	esistenza, err = db.EsisteMessaggio(IdMessaggio)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il recupero dell'ID della conversazione per il messaggio con ID %d: %w", IdMessaggio, err)
	}
	if !esistenza {
		return 404, fmt.Errorf("errore: il messaggio con ID %d non esiste", IdMessaggio)
	}
	conversazioneID, err := db.GetConversazioneIdByMessaggio(IdMessaggio)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il recupero dell'ID conversazione: %w", err)
	}
	isGruppo, codiceErrore, err := db.CercaConversazioneGruppo(conversazioneID)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	if isGruppo > 0 {
		coinvolto, codiceErrore, err := db.UtenteCoinvoltoGruppo(utente_Passato, conversazioneID)
		if !errors.Is(err, nil) {
			return codiceErrore, fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return 401, fmt.Errorf("l'utente non è membro del gruppo")
		}
	} else {
		idPrivata, codiceErrore, err := db.CercaConversazionePrivata(conversazioneID, utente_Passato_convertito)
		if !errors.Is(err, nil) {
			return codiceErrore, fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
		}
		if idPrivata == 0 {
			return codiceErrore, fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}
	}
	nuovaConversazione, codiceErrore, err := db.CreaConversazionePrivataDB(utente_Passato, utente2_Passato)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante la creazione della nuova conversazione privata: %w", err)
	}
	codiceErrore, err = db.CopiaMessaggioCambiandoOraEMitente(IdMessaggio, utente_Passato, nuovaConversazione)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante la copia del messaggio: %w", err)
	}

	return 0, nil
}
