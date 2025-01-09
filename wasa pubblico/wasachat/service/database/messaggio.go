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
			-- Uno dei due
			testo TEXT,
			foto INTEGER,
			--
			tempo TIME,
			FOREIGN KEY (foto) REFERENCES foto(id),
			FOREIGN KEY (autore) REFERENCES utente(id),
			FOREIGN KEY (conversazione) REFERENCES conversazione(id)
		);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella messaggio: %w", err)
	}
	return nil
}

// Funzione che crea la tabella commento se questa non esiste
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
func (db *appdbimpl) CreaMessaggioTestualeDBGruppo(utente_Passato string, conversazione_Passata int, testo_Passato string) error {
	esistenza, err := db.EsisteConversazione(conversazione_Passata)
	if err != nil {
		return fmt.Errorf("errore durante la verifica di esistenza: %w", err)
	}
	if !esistenza {
		return fmt.Errorf("errore, la chat non esiste")
	}
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utente_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a id: %w", err)
	}
	utenteCoinvolto, err := db.UtenteCoinvoltoGruppo(utente_Passato, conversazione_Passata)
	if err != nil {
		return fmt.Errorf("errore durante la verifica della partecipazione dell'utente: %w", err)
	}
	if utenteCoinvolto == 0 {
		return fmt.Errorf("l'utente %s non è coinvolto nella conversazione", utente_Passato)
	}
	queryDiInserimento := `INSERT INTO messaggio (testo, conversazione, autore, tempo) VALUES (?, ?, ?, ?);`
	result, err := db.c.Exec(queryDiInserimento, testo_Passato, conversazione_Passata, utente_Passato_convertito, time.Now())
	if err != nil {
		return fmt.Errorf("errore durante la creazione del messaggio: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID dell'ultimo messaggio: %w", err)
	}
	err = db.CreaStatoMessaggioGruppo(int(lastInsertID))
	if err != nil {
		return fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
	}
	return nil
}

// Funzione che aggiunge un nuovo messaggio con foto nel database
func (db *appdbimpl) CreaMessaggioFotoDBGruppo(utente_Passato string, conversazione_Passata int, foto_Passata int) error {
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utente_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a id: %w", err)
	}
	utenteCoinvolto, err := db.UtenteCoinvoltoGruppo(utente_Passato, conversazione_Passata)
	if err != nil {
		return fmt.Errorf("errore durante la verifica della partecipazione dell'utente: %w", err)
	}
	if utenteCoinvolto == 0 {
		return fmt.Errorf("l'utente %s non è coinvolto nella conversazione", utente_Passato)
	}
	queryDiInserimento := `INSERT INTO messaggio (autore, conversazione, foto, tempo) VALUES (?, ?, ?, ?);`
	result, err := db.c.Exec(queryDiInserimento, utente_Passato_convertito, conversazione_Passata, foto_Passata, time.Now())
	if err != nil {
		return fmt.Errorf("errore durante la creazione del messaggio: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID dell'ultimo messaggio: %w", err)
	}
	err = db.CreaStatoMessaggioGruppo(int(lastInsertID))
	if err != nil {
		return fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
	}
	return nil

}

// Funzione che aggiunge un nuovo messaggio testuale nel database
func (db *appdbimpl) CreaMessaggioTestualeDBPrivato(utente_Passato string, destinatario_Passato string, testo_Passato string) error {
	// Verifica che l'utente e il destinatario siano coinvolti in una conversazione privata
	chat, err := db.UtenteCoinvoltoPrivato(utente_Passato, destinatario_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
	}

	utente_Passato_convertito, err := db.IdUtenteDaNickname(utente_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a id: %w", err)
	}
	queryDiInserimento := `INSERT INTO messaggio (testo, conversazione, autore, tempo) VALUES (?, ?, ?, ?);`
	result, err := db.c.Exec(queryDiInserimento, testo_Passato, chat, utente_Passato_convertito, time.Now())
	if err != nil {
		return fmt.Errorf("errore durante la creazione del messaggio: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID dell'ultimo messaggio: %w", err)
	}
	err = db.CreaStatoMessaggioPrivato(int(lastInsertID))
	if err != nil {
		return fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
	}
	return nil
}

// Funzione che aggiunge un nuovo messaggio con foto nel database
func (db *appdbimpl) CreaMessaggioFotoDBPrivato(utente_Passato string, destinatario_Passato string, foto_Passata int) error {
	// Verifica che l'utente e il destinatario siano coinvolti in una conversazione privata
	chat, err := db.UtenteCoinvoltoPrivato(utente_Passato, destinatario_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
	}
	if chat == 0 {
		return fmt.Errorf("utente non coinvolto nella chat")
	}
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utente_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a id: %w", err)
	}
	queryDiInserimento := `INSERT INTO messaggio (autore, conversazione, foto, tempo) VALUES (?, ?, ?, ?);`
	result, err := db.c.Exec(queryDiInserimento, utente_Passato_convertito, chat, foto_Passata, time.Now())
	if err != nil {
		return fmt.Errorf("errore durante la creazione del messaggio: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID dell'ultimo messaggio: %w", err)
	}
	err = db.CreaStatoMessaggioPrivato(int(lastInsertID))
	if err != nil {
		return fmt.Errorf("errore durante la creazione dello stato del messaggio: %w", err)
	}

	return nil
}

func (db *appdbimpl) EliminaMessaggio(utente_Passato string, id_messaggio int, id_chat int) error {
	// Converti il nickname dell'utente in un ID
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utente_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a ID: %w", err)
	}

	// Query per eliminare il messaggio
	queryDiEliminazione := `DELETE FROM messaggio WHERE autore = ? AND id = ? AND conversazione = ?;`

	// Esegui la query e ottieni il numero di righe interessate
	result, err := db.c.Exec(queryDiEliminazione, utente_Passato_convertito, id_messaggio, id_chat)
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
func (db *appdbimpl) AggiungiCommento(utente_Passato string, messaggio_Passato int, reazione_Passata string) error {
	// Converti il nickname dell'utente in un ID
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utente_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a ID: %w", err)
	}

	// Verifica che il messaggio esista e recupera la conversazione associata
	var conversazioneId int
	queryVerificaMessaggio := `SELECT conversazione FROM messaggio WHERE id = ?;`
	err = db.c.QueryRow(queryVerificaMessaggio, messaggio_Passato).Scan(&conversazioneId)
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
		coinvolto, err := db.UtenteCoinvoltoGruppo(utente_Passato, conversazioneId)
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
	err = db.c.QueryRow(queryVerificaCommento, utente_Passato_convertito, messaggio_Passato).Scan(&commentoId)
	if err == nil {
		// Se il commento esiste già, restituisci un errore
		return fmt.Errorf("l'utente ha già commentato questo messaggio")
	} else if err != sql.ErrNoRows {
		// Se si verifica un altro errore, lo restituiamo
		return fmt.Errorf("errore durante la verifica del commento: %w", err)
	}

	// Aggiungi il commento nel database
	queryDiInserimento := `INSERT INTO commento (autore, messaggio, reazione) VALUES (?, ?, ?);`
	_, err = db.c.Exec(queryDiInserimento, utente_Passato_convertito, messaggio_Passato, reazione_Passata)
	if err != nil {
		return fmt.Errorf("errore durante l'inserimento del commento: %w", err)
	}

	return nil
}

// Funzione che elimina un commento dato l'ID del commento e il nome dell'utente
func (db *appdbimpl) EliminaCommento(utente_Passato string, id_commento int) error {
	// Converti il nickname dell'utente in un ID
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utente_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione da nickname a ID: %w", err)
	}

	// Verifica che il commento esista e prendi l'autore del commento
	var autoreCommento int
	queryVerificaCommento := `SELECT autore FROM commento WHERE id = ?;`
	err = db.c.QueryRow(queryVerificaCommento, id_commento).Scan(&autoreCommento)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("commento non trovato")
		}
		return fmt.Errorf("errore durante la verifica del commento: %w", err)
	}

	// Controlla che l'utente passato sia l'autore del commento
	if autoreCommento != utente_Passato_convertito {
		return fmt.Errorf("l'utente non è l'autore di questo commento, non può eliminarlo")
	}

	// Elimina il commento dal database
	queryDiEliminazione := `DELETE FROM commento WHERE id = ?;`
	_, err = db.c.Exec(queryDiEliminazione, id_commento)
	if err != nil {
		return fmt.Errorf("errore durante l'eliminazione del commento: %w", err)
	}

	return nil
}
