package database

import (
	"database/sql"
	"fmt"
	"time"
)

// Funzione che crea la tabella messaggio se questa non esiste
func CreaTabellaMessaggio(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS messaggio(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			autore INTEGER NOT NULL,
			conversazione INTEGER NOT NULL,
			inoltrato BOOL NOT NULL default FALSE,
			risposta INTEGER default NULL,
			-- Uno dei due
			testo TEXT,
			foto INTEGER,
			--
			tempo TIME,
			FOREIGN KEY (foto) REFERENCES foto(id),
			FOREIGN KEY (autore) REFERENCES utente(id),
			FOREIGN KEY (conversazione) REFERENCES conversazione(id),
			FOREIGN KEY (risposta) REFERENCES messaggio(id)
		);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella messaggio: %w", err)
	}
	return nil
}

// Funzione che crea la tabella commento se questa non esiste
// I commenti non sono messaggi, sono come delle emoji associate ai messaggi
func CreaTabellaCommento(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS commento(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			autore INTEGER NOT NULL,
			messaggio INTEGER NOT NULL,
			reazione TEXT NOT NULL,	
			UNIQUE (messaggio, autore),
			FOREIGN KEY (autore) REFERENCES utente(id),
			FOREIGN KEY (messaggio) REFERENCES messaggio(id)
		);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella commento: %w", err)
	}
	return nil
}

// Funzione che aggiunge un nuovo messaggio testuale nel database
func (db *appdbimpl) CreaMessaggioTestualeDB(utentePassato string, conversazionePassata int, testoPassato string) error {
	// Controllo se la chat esiste
	esistenza, err := db.EsisteConversazione(conversazionePassata)
	if err != nil {
		return fmt.Errorf("errore durante la verifica di esistenza: %w", err)
	}
	if !esistenza {
		return fmt.Errorf("errore, la chat non esiste")
	}

	// Trovo l'id utente associato al nickname
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a id: %w", err)
	}

	// Verifica se la conversazione è di tipo gruppo
	isGruppo, err := db.CercaConversazioneGruppo(conversazionePassata)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}

	if isGruppo > 0 {
		// Se è un gruppo, controllo se l'utente è coinvolto
		coinvolto, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazionePassata)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return fmt.Errorf("l'utente non è membro del gruppo")
		}

		// Esegui l'inserimento del messaggio nel gruppo
		if err := db.inserisciMessaggio(conversazionePassata, utente_Passato_convertito, testoPassato, true); err != nil {
			return err
		}

	} else {
		// Se è una conversazione privata, controllo la partecipazione dell'utente
		idPrivata, err := db.CercaConversazionePrivata(conversazionePassata, utente_Passato_convertito)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della conversazione: %w", err)
		}
		if idPrivata == 0 {
			return fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}

		// Esegui l'inserimento del messaggio nella conversazione privata
		if err := db.inserisciMessaggio(conversazionePassata, utente_Passato_convertito, testoPassato, false); err != nil {
			return err
		}
	}
	return nil
}

// Funzione che esegue l'inserimento del messaggio
func (db *appdbimpl) inserisciMessaggio(conversazionePassata, utente_Passato_convertito int, testoPassato string, isGruppo bool) error {
	var queryDiInserimento string
	var result sql.Result
	var err error

	// Determina la query a seconda del tipo di conversazione
	if isGruppo {
		queryDiInserimento = `INSERT INTO messaggio (testo, conversazione, autore, tempo) VALUES (?, ?, ?, ?);`
		result, err = db.c.Exec(queryDiInserimento, testoPassato, conversazionePassata, utente_Passato_convertito, time.Now())
	} else {
		queryDiInserimento = `INSERT INTO messaggio (autore, conversazione, testo, tempo) VALUES (?, ?, ?, ?);`
		result, err = db.c.Exec(queryDiInserimento, utente_Passato_convertito, conversazionePassata, testoPassato, time.Now())
	}

	if err != nil {
		return fmt.Errorf("errore durante la creazione del messaggio: %w", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID dell'ultimo messaggio: %w", err)
	}

	// Crea lo stato del messaggio, letto e ricevuto saranno impostati a false
	if isGruppo {
		err = db.CreaStatoMessaggioGruppo(int(lastInsertID))
	} else {
		err = db.CreaStatoMessaggioPrivato(int(lastInsertID))
	}

	if err != nil {
		return fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
	}

	return nil
}

// Funzione che aggiunge un nuovo messaggio con foto nel database
func (db *appdbimpl) CreaMessaggioFotoDB(utentePassato string, conversazionePassata int, fotoPassata int) error {
	// Controllo se la chat esiste
	esistenza, err := db.EsisteConversazione(conversazionePassata)
	if err != nil {
		return fmt.Errorf("errore durante la verifica di esistenza: %w", err)
	}
	if !esistenza {
		return fmt.Errorf("errore, la chat non esiste")
	}

	// Trovo l'id utente associato al nickname
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a id: %w", err)
	}

	// Verifica se l'utente è coinvolto nella conversazione (sia gruppo che privata)
	isGruppo, err := db.CercaConversazioneGruppo(conversazionePassata)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}

	if isGruppo > 0 {
		// Se è un gruppo, controllo se l'utente è coinvolto
		coinvolto, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazionePassata)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return fmt.Errorf("l'utente non è membro del gruppo")
		}

		// Esegui l'inserimento del messaggio con foto nel gruppo
		if err := db.inserisciMessaggioFoto(conversazionePassata, utente_Passato_convertito, fotoPassata, true); err != nil {
			return err
		}

	} else {
		// Se è una conversazione privata, controllo la partecipazione dell'utente
		idPrivata, err := db.CercaConversazionePrivata(conversazionePassata, utente_Passato_convertito)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della conversazione: %w", err)
		}
		if idPrivata == 0 {
			return fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}

		// Esegui l'inserimento del messaggio con foto nella conversazione privata
		if err := db.inserisciMessaggioFoto(conversazionePassata, utente_Passato_convertito, fotoPassata, false); err != nil {
			return err
		}
	}

	return nil
}

// Funzione che esegue l'inserimento del messaggio
func (db *appdbimpl) inserisciMessaggioFoto(conversazionePassata, utente_Passato_convertito int, fotoPassata int, isGruppo bool) error {
	var queryDiInserimento string
	var result sql.Result
	var err error

	// Determina la query a seconda del tipo di conversazione (gruppo o privata)
	if isGruppo {
		queryDiInserimento = `INSERT INTO messaggio (autore, conversazione, foto, tempo) VALUES (?, ?, ?, ?);`
		result, err = db.c.Exec(queryDiInserimento, utente_Passato_convertito, conversazionePassata, fotoPassata, time.Now())
	} else {
		queryDiInserimento = `INSERT INTO messaggio (autore, conversazione, foto, tempo) VALUES (?, ?, ?, ?);`
		result, err = db.c.Exec(queryDiInserimento, utente_Passato_convertito, conversazionePassata, fotoPassata, time.Now())
	}

	// Gestione dell'errore durante l'inserimento del messaggio
	if err != nil {
		return fmt.Errorf("errore durante la creazione del messaggio: %w", err)
	}

	// Recupero l'ID dell'ultimo messaggio inserito
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID dell'ultimo messaggio: %w", err)
	}

	// Crea lo stato del messaggio, che verrà impostato come "non letto" per il gruppo
	if isGruppo {
		err = db.CreaStatoMessaggioGruppo(int(lastInsertID))
	} else {
		// Puoi aggiungere una logica per lo stato del messaggio per le conversazioni private
		err = db.CreaStatoMessaggioPrivato(int(lastInsertID))
	}

	if err != nil {
		return fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
	}

	return nil
}

// funzione che elimina il messaggio passato tramite id
func (db *appdbimpl) EliminaMessaggio(utentePassato string, idmessaggio int, idchat int) error {
	// Converti il nickname dell'utente in un ID
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a ID: %w", err)
	}

	// Query per eliminare il messaggio
	queryDiEliminazione := `DELETE FROM messaggio WHERE autore = ? AND id = ? AND conversazione = ?;`

	// Esegui la query e ottieni il numero di righe interessate
	result, err := db.c.Exec(queryDiEliminazione, utente_Passato_convertito, idmessaggio, idchat)
	if err != nil {
		return fmt.Errorf("errore durante l'eliminazione del messaggio: %w", err)
	}

	// Verifica se è stato eliminato almeno un messaggio
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("errore durante il controllo delle righe interessate: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nessun messaggio trovato con i criteri specificati")
	}

	return nil
}

// Funzione che aggiunge un commento a un messaggio, con controllo se l'utente ha già commentato
// e verifica della partecipazione alla conversazione (privata o di gruppo)
func (db *appdbimpl) AggiungiCommento(utentePassato string, messaggioPassato int, reazionePassata string) error {
	// Converti il nickname dell'utente in un ID
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a ID: %w", err)
	}

	// Verifica che il messaggio esista e recupera la conversazione associata
	var conversazioneId int
	queryVerificaMessaggio := `SELECT conversazione FROM messaggio WHERE id = ?;`
	err = db.c.QueryRow(queryVerificaMessaggio, messaggioPassato).Scan(&conversazioneId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("messaggio non trovato")
		}
		return fmt.Errorf("errore durante la verifica del messaggio: %w", err)
	}

	// Verifica se la conversazione è un gruppo o privata
	isGruppo, err := db.CercaConversazioneGruppo(conversazioneId)
	if isGruppo > 0 {
		// Se è un gruppo, verifica che l'utente sia coinvolto
		coinvolto, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazioneId)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return fmt.Errorf("l'utente non è membro del gruppo")
		}
	}
	if err != nil {
		idPrivata, err := db.CercaConversazionePrivata(conversazioneId, utente_Passato_convertito)
		if idPrivata > 0 {
			if err != nil {
				return fmt.Errorf("errore durante la verifica della partecipazione dell'utente: %w", err)
			}
			if idPrivata == 0 {
				return fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
			}
		}
	}
	// Verifica che l'utente non abbia già commentato questo messaggio
	queryVerificaCommento := `SELECT id FROM commento WHERE autore = ? AND messaggio = ?;`
	var commentoId int
	err = db.c.QueryRow(queryVerificaCommento, utente_Passato_convertito, messaggioPassato).Scan(&commentoId)
	if err == nil {
		// Se il commento esiste già, restituisci un errore
		return fmt.Errorf("l'utente ha già commentato questo messaggio")
	} else if err != sql.ErrNoRows {
		// Se si verifica un altro errore, lo restituiamo
		return fmt.Errorf("errore durante la verifica del commento: %w", err)
	}

	// Aggiungi il commento nel database
	queryDiInserimento := `INSERT INTO commento (autore, messaggio, reazione) VALUES (?, ?, ?);`
	_, err = db.c.Exec(queryDiInserimento, utente_Passato_convertito, messaggioPassato, reazionePassata)
	if err != nil {
		return fmt.Errorf("errore durante l'inserimento del commento: %w", err)
	}

	return nil
}

// Funzione che elimina un commento dato l'ID del commento e il nome dell'utente
func (db *appdbimpl) EliminaCommento(utentePassato string, idcommento int) error {
	// Converti il nickname dell'utente in un ID
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a ID: %w", err)
	}

	// Verifica che il commento esista e prendi l'autore del commento
	var autoreCommento int
	queryVerificaCommento := `SELECT autore FROM commento WHERE id = ?;`
	err = db.c.QueryRow(queryVerificaCommento, idcommento).Scan(&autoreCommento)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("commento non trovato")
		}
		return fmt.Errorf("errore durante la verifica del commento: %w", err)
	}

	// Controlla che l'utente passato sia l'autore del commento
	// Se non è l'autore solleva un errore
	if autoreCommento != utente_Passato_convertito {
		return fmt.Errorf("l'utente non è l'autore di questo commento, non può eliminarlo")
	}

	// Elimina il commento dal database
	queryDiEliminazione := `DELETE FROM commento WHERE id = ?;`
	_, err = db.c.Exec(queryDiEliminazione, idcommento)
	if err != nil {
		return fmt.Errorf("errore durante l'eliminazione del commento: %w", err)
	}

	return nil
}
