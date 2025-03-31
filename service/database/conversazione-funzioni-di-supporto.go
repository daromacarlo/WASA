package database

import (
	"database/sql"
	"fmt"
)

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
			return 0, fmt.Errorf("nessun messaggio trovato con l'ID %d", idMessaggio)
		}
		return 0, fmt.Errorf("errore durante il recupero dell'ID della conversazione: %w", err)
	}

	return conversazioneID, nil
}

func (db *appdbimpl) EsisteMessaggio(idMessaggio int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM messaggio WHERE id = ?;"
	err := db.c.QueryRow(query, idMessaggio).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("errore durante la verifica dell'esistenza del messaggio con ID %d: %w", idMessaggio, err)
	}
	return count > 0, nil
}

// Passando un id di una conversazione la funzione restituisce true se esiste una conversazione con questo id, false se non esiste
func (db *appdbimpl) EsisteConversazione(idConversazione int) (bool, error) {
	var count int
	// Query per verificare se esiste una conversazione con l'id specificato
	query := `SELECT COUNT(*) FROM conversazione WHERE id = ?`
	err := db.c.QueryRow(query, idConversazione).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %w", err)
	}

	// Se count è maggiore di 0, la conversazione esiste
	if count > 0 {
		return true, nil
	}

	// Se count è 0, la conversazione non esiste
	return false, nil
}

// Passando i nickname degli utenti alla seguente funzione essa ritorna l'id della conversazine se questa esiste altrimenti ritona 0 e un errore
func (db *appdbimpl) EsisteConversazioneTraUtenti(utente1Passato string, utente2Passato string) (int, error) {
	var idConversazione int

	utente1ID, err := db.IdUtenteDaNickname(utente1Passato)
	if err != nil {
		return 0, fmt.Errorf("errore nella conversione del nickname %s in ID per utente1: %s", utente1Passato, err.Error())
	}
	utente2ID, err := db.IdUtenteDaNickname(utente2Passato)
	if err != nil {
		return 0, fmt.Errorf("errore nella conversione del nickname %s in ID per utente2: %s", utente2Passato, err.Error())
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
func (db *appdbimpl) UtenteCoinvoltoPrivato(utentePassato string, destinatarioPassato string) (int, error) {
	esiste, err := db.EsisteConversazioneTraUtenti(utentePassato, destinatarioPassato)
	if err != nil {
		return 0, fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %s", err.Error())
	}
	// caso in cui non esiste una conversazione privata tra i due utenti
	if esiste == 0 {
		// creo la conversazione
		id, err := db.CreaConversazionePrivataDB(utentePassato, destinatarioPassato)
		if err != nil {
			return 0, fmt.Errorf("errore durante la creazione della conversazione privata: %s", err.Error())
		}
		return id, nil
	} else {
		// caso in cui esiste
		return esiste, nil
	}
}
func (db *appdbimpl) GetNomeUtenteCoinvolto(conversazioneID int, nomeUtentePassato string) (string, error) {
	// Prima otteniamo l'ID dell'utente passato tramite il suo nome
	queryUtenteID := `
		SELECT id
		FROM utente
		WHERE nickname = ?;`

	var utenteID int
	err := db.c.QueryRow(queryUtenteID, nomeUtentePassato).Scan(&utenteID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Se l'utente non esiste
			return "", fmt.Errorf("utente con nome %s non trovato", nomeUtentePassato)
		}
		return "", fmt.Errorf("errore durante il recupero dell'ID dell'utente: %w", err)
	}

	// Esegui la query per ottenere l'altro utente coinvolto nella conversazione privata
	queryVerificaPartecipazione := `
		SELECT 
			CASE 
				WHEN cp.utente1 = ? THEN u2.nickname
				WHEN cp.utente2 = ? THEN u1.nickname
			END as utente_coinvolto
		FROM conversazione as c
		JOIN conversazioneprivata as cp ON cp.conversazione = c.id
		JOIN utente u1 ON cp.utente1 = u1.id
		JOIN utente u2 ON cp.utente2 = u2.id
		WHERE c.id = ? AND (cp.utente1 = ? OR cp.utente2 = ?);`

	var nomeUtenteCoinvolto string
	err = db.c.QueryRow(queryVerificaPartecipazione, utenteID, utenteID, conversazioneID, utenteID, utenteID).Scan(&nomeUtenteCoinvolto)
	if err != nil {
		if err == sql.ErrNoRows {
			// Se non ci sono risultati, la conversazione non esiste o l'utente non è coinvolto
			return "", fmt.Errorf("l'utente non è coinvolto in questa conversazione")
		}
		return "", fmt.Errorf("errore durante la ricerca dell'utente coinvolto: %w", err)
	}

	// Restituisci il nome dell'altro utente coinvolto
	return nomeUtenteCoinvolto, nil
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
			return 0, fmt.Errorf("la conversazione non esiste o l'utente non ne è coinvolto")
		}
		return 0, fmt.Errorf("errore durante la verifica della partecipazione dell'utente: %w", err)
	}

	// Se la conversazione privata esiste, ritorna il suo ID
	return conversazionePrivataID, nil
}

// La funzione controlla se l'utente è coinvolto nel gruppo passato (viene passato in input l'id di una conversazione non l'id di un gruppo)
// La funzione ritorna l'id del gruppo se l'utente è presente alrimenti 0
func (db *appdbimpl) UtenteCoinvoltoGruppo(utentePassato string, conversazionePassata int) (int, error) {
	// cerco l'id del gruppo
	idGruppo, err := db.CercaConversazioneGruppo(conversazionePassata)
	if err != nil {
		return 0, fmt.Errorf("errore durante la ricerca dell'ID del gruppo: %w", err)
	}
	// caso in cui il gruppo non esiste
	if idGruppo == 0 {
		return 0, fmt.Errorf("non esiste un gruppo con queste caratteristiche: %w", err)
	}

	idUtente, err := db.IdUtenteDaNickname(utentePassato)
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
func (db *appdbimpl) CercaConversazioneGruppo(conversazionePassata int) (int, error) {
	query := `SELECT id FROM gruppo WHERE conversazione = ?`
	var idGruppo int

	err := db.c.QueryRow(query, conversazionePassata).Scan(&idGruppo)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %w", err)
	}

	return idGruppo, nil
}
