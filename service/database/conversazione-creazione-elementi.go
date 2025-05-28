package database

import (
	"errors"
	"fmt"
)

// funzione che crea una conversazione e ritorna il suo id
func (db *appdbimpl) CreaConversazioneDB() (int, error) {
	queryDiInserimento := `INSERT INTO conversazione DEFAULT VALUES;`
	result, err := db.c.Exec(queryDiInserimento)
	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("errore durante la creazione della conversazione: %s", err.Error())
	}
	LastInsertIdid, err := result.LastInsertId()
	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("errore durante il recupero dell'ID della conversazione: %s", err.Error())
	}
	return int(LastInsertIdid), nil
}

// funzione che crea un gruppo e agginge al gruppo l'utente che lo ha creato
func (db *appdbimpl) CreaGruppoDB(UtenteChiamante string, nomeGruppo_Passato string, idfoto_Passata int) (int, error) {
	esisteUtenteChiamante, err := db.EsistenzaUtente(UtenteChiamante)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il controllo dell'esistenza dell'utente chiamante %s: %w", UtenteChiamante, err)
	}
	if !esisteUtenteChiamante {
		return 404, fmt.Errorf("l'utente chiamante %s non esiste", UtenteChiamante)
	}
	utenteconvertito, codiceErrore, err := db.IdUtenteDaNickname(UtenteChiamante)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore nella conversione del nickname in ID: %s", err.Error())
	}
	idConversazione, err := db.CreaConversazioneDB()
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante l'inserimento della nuova conversazione: %s", err.Error())
	}
	queryDiInserimento := `INSERT INTO gruppo (nome, conversazione, foto) VALUES (?,?,?);`
	result, err := db.c.Exec(queryDiInserimento, nomeGruppo_Passato, idConversazione, idfoto_Passata)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante la creazione del gruppo: %s", err.Error())
	}
	LastInsertIdid, err := result.LastInsertId()
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il recupero dell'ID del gruppo: %s", err.Error())
	}

	queryDiInserimentoUtente := `INSERT INTO utenteingruppo (utente, gruppo) VALUES (?, ?);`
	_, err = db.c.Exec(queryDiInserimentoUtente, utenteconvertito, LastInsertIdid)
	if !errors.Is(err, nil) {
		return 304, fmt.Errorf("errore durante l'aggiunta dell'utente al gruppo: %w", err)
	}

	return 0, nil
}

// funzione che crea una conversazione privata, ritorna l'id della conversazione
func (db *appdbimpl) CreaConversazionePrivataDB(utente1_Passato string, utente2_Passato string) (int, int, error) {
	esistenza, err := db.EsistenzaUtente(utente1_Passato)
	if !errors.Is(err, nil) {
		return 0, 500, fmt.Errorf("errore durante il controllo dell'esistenza dell'utente (chiamante): %s", err.Error())
	}
	if !esistenza {
		return 0, 404, fmt.Errorf("l'utente %s non esiste (chiamante)", utente1_Passato)
	}
	esistenza, err = db.EsistenzaUtente(utente2_Passato)
	if !errors.Is(err, nil) {
		return 0, 500, fmt.Errorf("errore durante il controllo dell'esistenza dell'utente (chiamato): %s", err.Error())
	}
	if !esistenza {
		return 0, 404, fmt.Errorf("l'utente %s non esiste (chiamato) o comunque non è stato trovato nel database", utente2_Passato)
	}
	idConversazione, codiceErrore, err := db.EsisteConversazioneTraUtenti(utente1_Passato, utente2_Passato)
	if !errors.Is(err, nil) {
		return 0, codiceErrore, fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %s", err.Error())
	}
	if idConversazione > 0 {
		return 0, 400, fmt.Errorf("l'utente %s ha afjsh", utente2_Passato)
	}
	idConversazione, err = db.CreaConversazioneDB()
	if !errors.Is(err, nil) {
		return 0, 500, fmt.Errorf("errore durante la creazione della nuova conversazione: %s", err.Error())
	}
	utente1_ID, codiceErrore, err := db.IdUtenteDaNickname(utente1_Passato)
	if !errors.Is(err, nil) {
		return 0, codiceErrore, fmt.Errorf("errore nella conversione del nickname %s in ID per utente1: %s", utente1_Passato, err.Error())
	}
	utente2_ID, codiceErrore, err := db.IdUtenteDaNickname(utente2_Passato)
	if !errors.Is(err, nil) {
		return 0, codiceErrore, fmt.Errorf("errore nella conversione del nickname %s in ID per utente2: %s", utente2_Passato, err.Error())
	}
	queryDiInserimento := `INSERT INTO conversazioneprivata (conversazione, utente1, utente2) VALUES (?, ?, ?);`
	_, err = db.c.Exec(queryDiInserimento, idConversazione, utente1_ID, utente2_ID)
	if !errors.Is(err, nil) {
		return 0, 500, fmt.Errorf("errore durante l'inserimento della conversazione nel database: %s", err.Error())
	}
	return idConversazione, 0, nil
}
