package database

import (
	"database/sql"
	"fmt"
)

// Funzione che crea la tabella statomessaggio se questa non esiste
func CreaTabellaStatoMessaggioPrivato(db *sql.DB) error {
	queryTabella := `
		CREATE TABLE IF NOT EXISTS statomessaggioprivato (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			messaggio INTEGER NOT NULL,
			ricevuto BOOL DEFAULT FALSE NOT NULL,
			letto BOOL DEFAULT FALSE NOT NULL,
			FOREIGN KEY (messaggio) REFERENCES messaggio(id)
		);`

	_, err := db.Exec(queryTabella)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella statomessaggio: %w", err)
	}
	return nil
}

// Funzione per creare un nuovo stato del messaggio nella tabella statomessaggio
func (db *appdbimpl) CreaStatoMessaggioPrivato(idmessaggioPassato int) error {
	queryDiInserimento := `INSERT INTO statomessaggioprivato (messaggio, ricevuto, letto) VALUES (?, ?, ?);`
	_, err := db.c.Exec(queryDiInserimento, idmessaggioPassato, false, false)
	if err != nil {
		return fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
	}
	return nil
}

func (db *appdbimpl) LeggiMessaggiPrivati(utente2Passato string, conversazioneID int) error {

	queryUpdate := `
	UPDATE statomessaggioprivato
	SET letto = 1, ricevuto = 1
	WHERE messaggio IN (
		SELECT m.id
		FROM messaggio m
		WHERE m.conversazione = ? 
		AND m.autore == ? 
	)`

	_, err := db.c.Exec(queryUpdate, conversazioneID, utente2Passato)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento dei messaggi come letti e ricevuti: %w", err)
	}

	return nil
}

// Funzione che crea la tabella statomessaggio se questa non esiste e il trigger associato
func CreaTabellaStatoMessaggioGruppo(db *sql.DB) error {
	queryTabella := `
		CREATE TABLE IF NOT EXISTS statomessaggiogruppo (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			messaggio INTEGER NOT NULL,
			ricevuto BOOL DEFAULT FALSE NOT NULL,
			letto BOOL DEFAULT FALSE NOT NULL,
			FOREIGN KEY (messaggio) REFERENCES messaggio(id)
		);`

	_, err := db.Exec(queryTabella)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella statomessaggio: %w", err)
	}
	return nil
}

// Funzione che crea la tabella statomessaggiogruppopersona se questa non esiste e il trigger associato
func CreaTabellaStatoMessaggioGruppoPersona(db *sql.DB) error {
	queryTabella := `
		CREATE TABLE IF NOT EXISTS statomessaggiogruppopersona (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			messaggio INTEGER NOT NULL,
			utente TEXT NOT NULL,
			UNIQUE (messaggio, utente), -- Restrizione di unicità sulla coppia messaggio e utente
			FOREIGN KEY (messaggio) REFERENCES messaggio(id),
			FOREIGN KEY (utente) REFERENCES utente(nickname)
		);`

	_, err := db.Exec(queryTabella)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella statomessaggiogruppopersona: %w", err)
	}
	return nil
}

// Funzione che crea la tabella statomessaggiogruppopersona se questa non esiste e il trigger associato
func CreaTabellaStatoMessaggioGruppoPersonaRicevimento(db *sql.DB) error {
	queryTabella := `
		CREATE TABLE IF NOT EXISTS statomessaggiogruppopersonaricevimento (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			messaggio INTEGER NOT NULL,
			utente TEXT NOT NULL,
			UNIQUE (messaggio, utente), -- Restrizione di unicità sulla coppia messaggio e utente
			FOREIGN KEY (messaggio) REFERENCES messaggio(id),
			FOREIGN KEY (utente) REFERENCES utente(nickname)
		);`

	_, err := db.Exec(queryTabella)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella statomessaggiogruppopersonaricevimento: %w", err)
	}
	return nil
}

// Funzione per creare un nuovo stato del messaggio nella tabella statomessaggio
func (db *appdbimpl) CreaStatoMessaggioGruppo(idmessaggioPassato int) error {
	queryDiInserimento := `INSERT INTO statomessaggiogruppo (messaggio, ricevuto, letto) VALUES (?, ?, ?);`
	_, err := db.c.Exec(queryDiInserimento, idmessaggioPassato, false, false)
	if err != nil {
		return fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
	}
	return nil
}

func (db *appdbimpl) LeggiMessaggiGruppo(utentePassato string, conversazioneID int) error {
	// Converti il nickname in ID utente
	// Avvia una transazione
	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("errore nell'avvio della transazione: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Ripristina il panic
		} else if err != nil {
			tx.Rollback() // Annulla la transazione in caso di errore
		} else {
			err = tx.Commit() // Conferma la transazione
		}
	}()

	// Query per ottenere i messaggi
	queryMessaggi := `
		SELECT m.id
		FROM messaggio m
		JOIN statomessaggiogruppo smg on smg.messaggio = m.id
		JOIN utente as u on m.autore = u.nickname
		WHERE m.conversazione = ? AND smg.letto = false AND u.nickname != ?;`

	rows, err := tx.Query(queryMessaggi, conversazioneID, utentePassato)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dei messaggi: %w", err)
	}
	defer rows.Close()

	// Query di inserimento
	queryDiInserimento := `
		INSERT INTO statomessaggiogruppopersona (messaggio, utente)
		VALUES (?, ?)
		ON CONFLICT DO NOTHING;`

	// Ciclo per inserire i dati
	for rows.Next() {
		var messaggioID int
		if err := rows.Scan(&messaggioID); err != nil {
			return fmt.Errorf("errore durante la lettura dei messaggi: %w", err)
		}

		_, err := tx.Exec(queryDiInserimento, messaggioID, utentePassato)
		if err != nil {
			return fmt.Errorf("errore durante la creazione dello stato del messaggio (messaggio ID %d): %w", messaggioID, err)
		}
	}

	// Controlla errori durante l'iterazione
	if err := rows.Err(); err != nil {
		return fmt.Errorf("errore durante l'iterazione dei messaggi: %w", err)
	}

	return nil
}

func (db *appdbimpl) CheckLetturaMessaggiGruppo(conversazioneID int) error {
	// Aggiorna lo stato dei messaggi come letti se soddisfano la condizione
	queryUpdate := `
		UPDATE statomessaggiogruppo
		SET letto = 1, ricevuto = 1
		WHERE messaggio IN (
			SELECT m.id
			FROM messaggio m
			WHERE m.conversazione = ? AND (
				SELECT COUNT(*)
				FROM statomessaggiogruppopersona smgp
				WHERE smgp.messaggio = m.id
			) >= (
				SELECT COUNT(*)
				FROM utenteingruppo uig
				JOIN gruppo as g on uig.gruppo = g.id
				WHERE uig.gruppo = g.id and g.conversazione = ?
			) - 1
		);`

	// Esegui la query
	_, err := db.c.Exec(queryUpdate, conversazioneID, conversazioneID)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento dei messaggi come letti e ricevuti: %w", err)
	}

	return nil
}

func (db *appdbimpl) CheckRicevimentoMessaggiGruppo(conversazioneID int) error {
	// Aggiorna lo stato dei messaggi come letti se soddisfano la condizione
	queryUpdate := `
		UPDATE statomessaggiogruppo
		SET ricevuto = 1
		WHERE messaggio IN (
			SELECT m.id
			FROM messaggio m
			WHERE m.conversazione = ? AND (
				SELECT COUNT(*)
				FROM statomessaggiogruppopersonaricevimento smgp
				WHERE smgp.messaggio = m.id
			) >= (
				SELECT COUNT(*)
				FROM utenteingruppo uig
				JOIN gruppo as g on uig.gruppo = g.id
				WHERE uig.gruppo = g.id and g.conversazione = ?
			) - 1
		);`

	// Esegui la query
	_, err := db.c.Exec(queryUpdate, conversazioneID, conversazioneID)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento dei messaggi come letti e ricevuti: %w", err)
	}

	return nil
}

// Segna i messaggi privati come ricevuti per l'utente specificato
func (db *appdbimpl) segnaMessaggiPrivatiRicevuti(utenteID string, conversazioneID int) error {
	queryUpdate := `
    UPDATE statomessaggioprivato
    SET ricevuto = 1
    WHERE messaggio IN (
        SELECT m.id
        FROM messaggio m
        WHERE m.conversazione = ? 
        AND m.autore != ?
    )`

	_, err := db.c.Exec(queryUpdate, conversazioneID, utenteID)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento dei messaggi privati come ricevuti: %w", err)
	}

	return nil
}

// Segna i messaggi di gruppo come ricevuti per l'utente specificato
func (db *appdbimpl) segnaMessaggiGruppoRicevuti(utenteID string, conversazioneID int) error {
	// Converti il nickname in ID utente
	// Avvia una transazione
	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("errore nell'avvio della transazione: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Ripristina il panic
		} else if err != nil {
			tx.Rollback() // Annulla la transazione in caso di errore
		} else {
			err = tx.Commit() // Conferma la transazione
		}
	}()

	// Query per ottenere i messaggi
	queryMessaggi := `
		SELECT m.id
		FROM messaggio m
		JOIN statomessaggiogruppo smg on smg.messaggio = m.id
		JOIN utente as u on m.autore = u.nickname
		WHERE m.conversazione = ? AND smg.ricevuto = false AND u.nickname != ?;`

	rows, err := tx.Query(queryMessaggi, conversazioneID, utenteID)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dei messaggi: %w", err)
	}
	defer rows.Close()

	// Query di inserimento
	queryDiInserimento := `
		INSERT INTO statomessaggiogruppopersonaricevimento (messaggio, utente)
		VALUES (?, ?)
		ON CONFLICT DO NOTHING;`

	// Ciclo per inserire i dati
	for rows.Next() {
		var messaggioID int
		if err := rows.Scan(&messaggioID); err != nil {
			return fmt.Errorf("errore durante la lettura dei messaggi: %w", err)
		}

		_, err := tx.Exec(queryDiInserimento, messaggioID, utenteID)
		if err != nil {
			return fmt.Errorf("errore durante la creazione dello stato del messaggio (messaggio ID %d): %w", messaggioID, err)
		}
	}

	// Controlla errori durante l'iterazione
	if err := rows.Err(); err != nil {
		return fmt.Errorf("errore durante l'iterazione dei messaggi: %w", err)
	}

	return nil
}
