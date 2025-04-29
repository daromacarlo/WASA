package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func CreaTabellaMessaggio(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS messaggio(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			autore TEXT NOT NULL,
			conversazione INTEGER NOT NULL,
			inoltrato BOOL NOT NULL default FALSE,
			risposta INTEGER default NULL,
			-- Uno dei due
			testo TEXT,
			foto INTEGER,
			--
			tempo TIME,
			FOREIGN KEY (foto) REFERENCES foto(id),
			FOREIGN KEY (autore) REFERENCES utente(nickname),
			FOREIGN KEY (conversazione) REFERENCES conversazione(id),
			FOREIGN KEY (risposta) REFERENCES messaggio(id)
		);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella messaggio: %w", err)
	}
	return nil
}

func CreaTabellaCommento(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS commento(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			autore TEXT NOT NULL,
			messaggio INTEGER NOT NULL,
			reazione TEXT NOT NULL,	
			FOREIGN KEY (autore) REFERENCES utente(nickname),
			FOREIGN KEY (messaggio) REFERENCES messaggio(id)
		);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella commento: %w", err)
	}
	return nil
}

func (db *appdbimpl) CreaMessaggioTestualeDB(utentePassato string, conversazionePassata int, testoPassato string) (int, error) {
	esistenza, err := db.EsisteConversazione(conversazionePassata)
	if err != nil {
		return 0, fmt.Errorf("errore durante la verifica di esistenza: %w", err)
	}
	if !esistenza {
		return 0, fmt.Errorf("errore, la chat non esiste")
	}
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return 0, fmt.Errorf("errore durante la conversione da nickname a id: %w", err)
	}
	isGruppo, err := db.CercaConversazioneGruppo(conversazionePassata)
	if err != nil {
		return 0, fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	var idMessaggio int
	if isGruppo > 0 {
		coinvolto, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazionePassata)
		if err != nil {
			return 0, fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return 0, fmt.Errorf("l'utente non è membro del gruppo")
		}
		idMessaggio, err = db.inserisciMessaggio(conversazionePassata, utentePassato, testoPassato, true)
		if err != nil {
			return 0, err
		}

	} else {
		idPrivata, err := db.CercaConversazionePrivata(conversazionePassata, utente_Passato_convertito)
		if err != nil {
			return 0, fmt.Errorf("errore durante la verifica della conversazione: %w", err)
		}
		if idPrivata == 0 {
			return 0, fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}
		idMessaggio, err = db.inserisciMessaggio(conversazionePassata, utentePassato, testoPassato, false)
		if err != nil {
			return 0, err
		}
	}

	return idMessaggio, nil
}

func (db *appdbimpl) inserisciMessaggio(conversazionePassata int, utente_Passato string, testoPassato string, isGruppo bool) (int, error) {
	var queryDiInserimento string
	var result sql.Result
	var err error

	queryDiInserimento = `INSERT INTO messaggio (autore, conversazione, testo, tempo) VALUES (?, ?, ?, ?);`
	result, err = db.c.Exec(queryDiInserimento, utente_Passato, conversazionePassata, testoPassato, time.Now())

	if err != nil {
		return 0, fmt.Errorf("errore durante la creazione del nuovo messaggio : %w", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("errore durante il recupero dell'ID dell'ultimo messaggio: %w", err)
	}

	if isGruppo {
		err = db.CreaStatoMessaggioGruppo(int(lastInsertID))
	} else {
		err = db.CreaStatoMessaggioPrivato(int(lastInsertID))
	}

	if err != nil {
		return 0, fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
	}

	return int(lastInsertID), nil
}

func (db *appdbimpl) CreaMessaggioFotoDB(utentePassato string, conversazionePassata int, fotoPassata int) (int, error) {
	esistenza, err := db.EsisteConversazione(conversazionePassata)
	if err != nil {
		return 0, fmt.Errorf("errore durante la verifica di esistenza: %w", err)
	}
	if !esistenza {
		return 0, fmt.Errorf("errore, la chat non esiste")
	}
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return 0, fmt.Errorf("errore durante la conversione da nickname a id: %w", err)
	}
	isGruppo, err := db.CercaConversazioneGruppo(conversazionePassata)
	if err != nil {
		return 0, fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}

	var idMessaggio int
	if isGruppo > 0 {
		coinvolto, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazionePassata)
		if err != nil {
			return 0, fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return 0, fmt.Errorf("l'utente non è membro del gruppo")
		}
		idMessaggio, err = db.inserisciMessaggioFoto(conversazionePassata, utentePassato, fotoPassata, true)
		if err != nil {
			return 0, fmt.Errorf("errore durante l'inserimento del messaggio: %w", err)
		}

	} else {
		idPrivata, err := db.CercaConversazionePrivata(conversazionePassata, utente_Passato_convertito)
		if err != nil {
			return 0, fmt.Errorf("errore durante la verifica della conversazione: %w", err)
		}
		if idPrivata == 0 {
			return 0, fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}
		idMessaggio, err = db.inserisciMessaggioFoto(conversazionePassata, utentePassato, fotoPassata, false)
		if err != nil {
			return 0, fmt.Errorf("errore durante l'inserimento del messaggio: %w", err)
		}
	}

	return idMessaggio, nil
}

func (db *appdbimpl) inserisciMessaggioFoto(conversazionePassata int, utente_Passato string, fotoPassata int, isGruppo bool) (int, error) {
	queryDiInserimento := `INSERT INTO messaggio (autore, conversazione, foto, tempo) VALUES (?, ?, ?, ?);`
	result, err := db.c.Exec(queryDiInserimento, utente_Passato, conversazionePassata, fotoPassata, time.Now())

	if err != nil {
		return 0, fmt.Errorf(": %w", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("errore durante il recupero dell'ID dell'ultimo messaggio: %w", err)
	}

	if isGruppo {
		err = db.CreaStatoMessaggioGruppo(int(lastInsertID))
	} else {
		err = db.CreaStatoMessaggioPrivato(int(lastInsertID))
	}

	if err != nil {
		return 0, fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
	}

	return int(lastInsertID), nil
}

func (db *appdbimpl) EliminaMessaggio(utentePassato string, idmessaggio int, idchat int) error {
	queryDiEliminazione := `DELETE FROM messaggio WHERE autore = ? AND id = ? AND conversazione = ?;`
	result, err := db.c.Exec(queryDiEliminazione, utentePassato, idmessaggio, idchat)
	if err != nil {
		return fmt.Errorf("errore durante l'eliminazione del messaggio: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("errore durante il controllo delle righe interessate: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nessun messaggio trovato con i criteri specificati")
	}

	return nil
}

func (db *appdbimpl) AggiungiCommento(utentePassato string, messaggioPassato int, reazionePassata string) error {
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a ID: %w", err)
	}

	var conversazioneId int
	queryVerificaMessaggio := `SELECT conversazione FROM messaggio WHERE id = ?;`
	err = db.c.QueryRow(queryVerificaMessaggio, messaggioPassato).Scan(&conversazioneId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("messaggio non trovato")
		}
		return fmt.Errorf("errore durante la verifica del messaggio: %w", err)
	}

	isGruppo, err := db.CercaConversazioneGruppo(conversazioneId)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}

	if isGruppo > 0 {
		coinvolto, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazioneId)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return fmt.Errorf("l'utente non è membro del gruppo")
		}
	} else {
		idPrivata, err := db.CercaConversazionePrivata(conversazioneId, utente_Passato_convertito)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della partecipazione dell'utente: %w", err)
		}
		if idPrivata == 0 {
			return fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}
	}

	queryVerificaCommento := `SELECT id FROM commento WHERE autore = ? AND messaggio = ?;`
	var commentoId int
	err = db.c.QueryRow(queryVerificaCommento, utentePassato, messaggioPassato).Scan(&commentoId)

	switch {
	case err == nil:
		queryDiAggiornamento := `UPDATE commento SET reazione = ?, data_commento = CURRENT_TIMESTAMP WHERE id = ?;`
		_, err = db.c.Exec(queryDiAggiornamento, reazionePassata, commentoId)
		if err != nil {
			return fmt.Errorf("errore durante l'aggiornamento del commento: %w", err)
		}
	case errors.Is(err, sql.ErrNoRows):
		queryDiInserimento := `INSERT INTO commento (autore, messaggio, reazione) VALUES (?, ?, ?);`
		_, err = db.c.Exec(queryDiInserimento, utentePassato, messaggioPassato, reazionePassata)
		if err != nil {
			return fmt.Errorf("errore durante l'inserimento del commento: %w", err)
		}
	default:
		return fmt.Errorf("errore durante la verifica del commento: %w", err)
	}

	return nil
}

// EliminaCommento elimina un commento specifico di un utente dato l'ID del messaggio
func (db *appdbimpl) EliminaCommento(utente string, idmessaggio int) error {
	// Verifica che esista un commento dell'utente specifico per questo messaggio
	var idcommento int

	queryVerificaCommento := `
        SELECT id 
        FROM commento 
        WHERE messaggio = ? AND autore = ?;
    `

	err := db.c.QueryRow(queryVerificaCommento, idmessaggio, utente).Scan(&idcommento)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("nessun commento trovato per questo utente e messaggio")
		}
		return fmt.Errorf("errore durante la ricerca del commento: %w", err)
	}

	// Elimina il commento dal database
	queryEliminazione := `
        DELETE FROM commento 
        WHERE id = ?;
    `

	_, err = db.c.Exec(queryEliminazione, idcommento)
	if err != nil {
		return fmt.Errorf("errore durante l'eliminazione del commento: %w", err)
	}

	return nil
}
