package database

import (
	"database/sql"
	"fmt"
)

// CreaTabellaStatoMessaggioPrivato crea la tabella per lo stato dei messaggi privati
func CreaTabellaStatoMessaggioPrivato(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS statomessaggioprivato (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		messaggio INTEGER NOT NULL,
		ricevuto BOOL DEFAULT FALSE NOT NULL,
		letto BOOL DEFAULT FALSE NOT NULL,
		FOREIGN KEY (messaggio) REFERENCES messaggio(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore creazione tabella statomessaggioprivato: %w", err)
	}
	return nil
}

// CreaStatoMessaggioPrivato inserisce un nuovo stato messaggio privato
func (db *appdbimpl) CreaStatoMessaggioPrivato(idmessaggio int) error {
	_, err := db.c.Exec(
		"INSERT INTO statomessaggioprivato (messaggio, ricevuto, letto) VALUES (?, ?, ?)",
		idmessaggio, false, false,
	)
	if err != nil {
		return fmt.Errorf("errore creazione stato messaggio privato: %w", err)
	}
	return nil
}

// LeggiMessaggiPrivati marca i messaggi come letti e ricevuti
func (db *appdbimpl) LeggiMessaggiPrivati(destinatario string, conversazioneID int) error {
	_, err := db.c.Exec(`
		UPDATE statomessaggioprivato 
		SET letto = 1, ricevuto = 1 
		WHERE messaggio IN (
			SELECT id FROM messaggio 
			WHERE conversazione = ? AND autore = ?
		)`,
		conversazioneID, destinatario,
	)
	if err != nil {
		return fmt.Errorf("errore aggiornamento stato messaggi privati: %w", err)
	}
	return nil
}

// CreaTabellaStatoMessaggioGruppo crea la tabella per lo stato dei messaggi di gruppo
func CreaTabellaStatoMessaggioGruppo(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS statomessaggiogruppo (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		messaggio INTEGER NOT NULL,
		ricevuto BOOL DEFAULT FALSE NOT NULL,
		letto BOOL DEFAULT FALSE NOT NULL,
		FOREIGN KEY (messaggio) REFERENCES messaggio(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore creazione tabella statomessaggiogruppo: %w", err)
	}
	return nil
}

// CreaTabellaStatoMessaggioGruppoPersona crea la tabella per lo stato lettura messaggi gruppo
func CreaTabellaStatoMessaggioGruppoPersona(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS statomessaggiogruppopersona (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		messaggio INTEGER NOT NULL,
		utente TEXT NOT NULL,
		UNIQUE(messaggio, utente),
		FOREIGN KEY (messaggio) REFERENCES messaggio(id),
		FOREIGN KEY (utente) REFERENCES utente(nickname)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore creazione tabella statomessaggiogruppopersona: %w", err)
	}
	return nil
}

// CreaTabellaStatoMessaggioGruppoPersonaRicevimento crea la tabella per lo stato ricezione messaggi gruppo
func CreaTabellaStatoMessaggioGruppoPersonaRicevimento(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS statomessaggiogruppopersonaricevimento (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		messaggio INTEGER NOT NULL,
		utente TEXT NOT NULL,
		UNIQUE(messaggio, utente),
		FOREIGN KEY (messaggio) REFERENCES messaggio(id),
		FOREIGN KEY (utente) REFERENCES utente(nickname)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore creazione tabella statomessaggiogruppopersonaricevimento: %w", err)
	}
	return nil
}

// CreaStatoMessaggioGruppo inserisce un nuovo stato messaggio gruppo
func (db *appdbimpl) CreaStatoMessaggioGruppo(idmessaggio int) error {
	_, err := db.c.Exec(
		"INSERT INTO statomessaggiogruppo (messaggio, ricevuto, letto) VALUES (?, ?, ?)",
		idmessaggio, false, false,
	)
	if err != nil {
		return fmt.Errorf("errore creazione stato messaggio gruppo: %w", err)
	}
	return nil
}

// LeggiMessaggiGruppo marca i messaggi di gruppo come letti dall'utente
func (db *appdbimpl) LeggiMessaggiGruppo(utente string, conversazioneID int) error {
	// Inserisce nella tabella statomessaggiogruppopersona i messaggi non ancora letti
	_, err := db.c.Exec(`
		INSERT INTO statomessaggiogruppopersona (messaggio, utente)
		SELECT m.id, ?
		FROM messaggio m
		JOIN statomessaggiogruppo smg ON smg.messaggio = m.id
		JOIN utente u ON m.autore = u.nickname
		WHERE m.conversazione = ? AND smg.letto = false AND u.nickname != ?
		ON CONFLICT DO NOTHING`,
		utente, conversazioneID, utente,
	)
	if err != nil {
		return fmt.Errorf("errore aggiornamento lettura messaggi gruppo: %w", err)
	}
	return nil
}

// CheckLetturaMessaggiGruppo verifica e aggiorna lo stato lettura messaggi gruppo
func (db *appdbimpl) CheckLetturaMessaggiGruppo(conversazioneID int) error {
	_, err := db.c.Exec(`
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
				JOIN gruppo g ON uig.gruppo = g.id
				WHERE g.conversazione = ?
			) - 1
		)`,
		conversazioneID, conversazioneID,
	)
	if err != nil {
		return fmt.Errorf("errore verifica lettura messaggi gruppo: %w", err)
	}
	return nil
}

// CheckRicevimentoMessaggiGruppo verifica e aggiorna lo stato ricezione messaggi gruppo
func (db *appdbimpl) CheckRicevimentoMessaggiGruppo(conversazioneID int) error {
	_, err := db.c.Exec(`
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
				JOIN gruppo g ON uig.gruppo = g.id
				WHERE g.conversazione = ?
			) - 1
		)`,
		conversazioneID, conversazioneID,
	)
	if err != nil {
		return fmt.Errorf("errore verifica ricezione messaggi gruppo: %w", err)
	}
	return nil
}

// SegnaMessaggiPrivatiRicevuti marca i messaggi privati come ricevuti
func (db *appdbimpl) SegnaMessaggiPrivatiRicevuti(utente string, conversazioneID int) error {
	_, err := db.c.Exec(`
		UPDATE statomessaggioprivato
		SET ricevuto = 1
		WHERE messaggio IN (
			SELECT id FROM messaggio
			WHERE conversazione = ? AND autore != ?
		)`,
		conversazioneID, utente,
	)
	if err != nil {
		return fmt.Errorf("errore aggiornamento ricezione messaggi privati: %w", err)
	}
	return nil
}

// SegnaMessaggiGruppoRicevuti marca i messaggi di gruppo come ricevuti dall'utente
func (db *appdbimpl) SegnaMessaggiGruppoRicevuti(utente string, conversazioneID int) error {
	// Inserisce nella tabella statomessaggiogruppopersonaricevimento i messaggi non ancora ricevuti
	_, err := db.c.Exec(`
		INSERT INTO statomessaggiogruppopersonaricevimento (messaggio, utente)
		SELECT m.id, ?
		FROM messaggio m
		JOIN statomessaggiogruppo smg ON smg.messaggio = m.id
		JOIN utente u ON m.autore = u.nickname
		WHERE m.conversazione = ? AND smg.ricevuto = false AND u.nickname != ?
		ON CONFLICT DO NOTHING`,
		utente, conversazioneID, utente,
	)
	if err != nil {
		return fmt.Errorf("errore aggiornamento ricezione messaggi gruppo: %w", err)
	}
	return nil
}
