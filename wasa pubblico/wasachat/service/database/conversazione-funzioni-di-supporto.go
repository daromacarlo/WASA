package database

import (
	"database/sql"
	"fmt"
)

// Passando un id di una conversazione la funzione restituisce true se esiste una conversazione con questo id, false se non esiste
func (db *appdbimpl) EsisteConversazione(idConversazione int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM conversazioneprivata WHERE id = ?`
	err := db.c.QueryRow(query, idConversazione).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %w", err)
	}
	return true, nil
}

// Passando i nickname degli utenti alla seguente funzione essa ritorna l'id della conversazine se questa esiste altrimenti ritona 0 e un errore
func (db *appdbimpl) EsisteConversazioneTraUtenti(utente1_Passato string, utente2_Passato string) (int, error) {
	var idConversazione int

	utente1ID, err := db.IdUtenteDaNickname(utente1_Passato)
	if err != nil {
		return 0, fmt.Errorf("errore nella conversione del nickname %s in ID per utente1: %s", utente1_Passato, err.Error())
	}
	utente2ID, err := db.IdUtenteDaNickname(utente2_Passato)
	if err != nil {
		return 0, fmt.Errorf("errore nella conversione del nickname %s in ID per utente2: %s", utente2_Passato, err.Error())
	}
	// Cerchiamo l'id della conversazione
	query := `
		SELECT c.id
		FROM conversazione c
		LEFT JOIN conversazioneprivata cp ON cp.conversazione = c.id
		WHERE (cp.utente1 = ? AND cp.utente2 = ?) OR (cp.utente1 = ? AND cp.utente2 = ?);`

	err = db.c.QueryRow(query, utente1ID, utente2ID, utente2ID, utente1ID).Scan(&idConversazione)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %s", err.Error())
	}
	return idConversazione, nil
}

// Passando i nickname degli utenti di una conversazione privata questa funzione ritorna sempre l'id di una conversazione privata, se non esiste viene creata
// e il suo id viene ritornato
func (db *appdbimpl) UtenteCoinvoltoPrivato(utente_Passato string, destinatario_Passato string) (int, error) {
	esiste, err := db.EsisteConversazioneTraUtenti(utente_Passato, destinatario_Passato)
	if err != nil {
		return 0, fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %s", err.Error())
	}
	// caso in cui non esiste una conversazione privata tra i due utenti
	if esiste == 0 {
		// creo la conversazione
		id, err := db.CreaConversazionePrivataDB(utente_Passato, destinatario_Passato)
		if err != nil {
			return 0, fmt.Errorf("errore durante la creazione della conversazione privata: %s", err.Error())
		}
		return id, nil
	} else {
		// caso in cui esiste
		return esiste, nil
	}
}

func (db *appdbimpl) CercaConversazionePrivata(conversazioneID int, utente_Passato_convertito int) (int, error) {
	// Se è una conversazione privata, verifica che l'utente sia coinvolto
	queryVerificaPartecipazione := `
		SELECT cp.id
		FROM conversazione as c
		JOIN conversazioneprivata as cp on cp.conversazione = c.id
		WHERE c.id = ? AND (cp.utente1 = ? OR cp.utente2 = ?);`

	var conversazionePrivataID int
	err := db.c.QueryRow(queryVerificaPartecipazione, conversazioneID, utente_Passato_convertito, utente_Passato_convertito).Scan(&conversazionePrivataID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Se non ci sono risultati, la conversazione non esiste
			return 0, fmt.Errorf("la conversazione privata non esiste o l'utente non è coinvolto")
		}
		return 0, fmt.Errorf("errore durante la verifica della partecipazione dell'utente: %w", err)
	}

	// Se la conversazione privata esiste, ritorna il suo ID
	return conversazionePrivataID, nil
}

// La funzione controlla se l'utente è coinvolto nel gruppo passato (viene passato in input l'id di una conversazione non l'id di un gruppo)
// La funzione ritorna l'id del gruppo se l'utente è presente alrimenti 0
func (db *appdbimpl) UtenteCoinvoltoGruppo(utente_Passato string, conversazione_Passata int) (int, error) {
	// cerco l'id del gruppo
	idGruppo, err := db.CercaConversazioneGruppo(conversazione_Passata)
	if err != nil {
		return 0, fmt.Errorf("errore durante la ricerca dell'ID del gruppo: %w", err)
	}
	// caso in cui il gruppo non esiste
	if idGruppo == 0 {
		return 0, fmt.Errorf("non esiste un gruppo con queste caratteristiche: %w", err)
	}

	idUtente, err := db.IdUtenteDaNickname(utente_Passato)
	if err != nil {
		return 0, fmt.Errorf("errore durante la ricerca dell'ID utente: %w", err)
	}

	query := `SELECT 1 FROM utenteingruppo WHERE gruppo = ? AND utente = ?`
	var result bool
	err = db.c.QueryRow(query, idGruppo, idUtente).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("errore durante la verifica dell'utente nel gruppo: %w", err)
	}
	if result {
		return idGruppo, nil
	}
	return 0, fmt.Errorf("errore, l'utente non è presente nel gruppo: %w", err)
}

// La funzione cerca tra i gruppi se ne esiste uno che abbia l'attributo "conversazione" uguale al parametro passato in input
// se sì ritorna l'id del gruppo, se no ritorna 0
func (db *appdbimpl) CercaConversazioneGruppo(conversazione_Passata int) (int, error) {
	query := `SELECT id FROM gruppo WHERE conversazione = ?`
	var idGruppo int

	err := db.c.QueryRow(query, conversazione_Passata).Scan(&idGruppo)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %w", err)
	}

	return idGruppo, nil
}
