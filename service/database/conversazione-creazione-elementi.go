package database

import (
	"fmt"
)

// funzione che crea una conversazione e ritorna il suo id
func (db *appdbimpl) CreaConversazioneDB() (int, error) {
	queryDiInserimento := `INSERT INTO conversazione DEFAULT VALUES;`
	result, err := db.c.Exec(queryDiInserimento)
	if err != nil {
		return 0, fmt.Errorf("errore durante la creazione della conversazione: %s", err.Error())
	}
	LastInsertIdid, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("errore durante il recupero dell'ID della conversazione: %s", err.Error())
	}
	fmt.Println("ID conversazione creato:", LastInsertIdid)
	return int(LastInsertIdid), nil
}

// funzione che crea un gruppo e agginge al gruppo l'utente che lo ha creatos
func (db *appdbimpl) CreaGruppoDB(UtenteChiamante string, nomeGruppo_Passato string, idfoto_Passata int) error {
	esisteUtenteChiamante, err := db.EsistenzaUtente(UtenteChiamante)
	if err != nil {
		return fmt.Errorf("errore durante il controllo dell'esistenza dell'utente chiamante %s: %w", UtenteChiamante, err)
	}
	if !esisteUtenteChiamante {
		return fmt.Errorf("l'utente chiamante %s non esiste", UtenteChiamante)
	}
	utenteconvertito, err := db.IdUtenteDaNickname(UtenteChiamante)
	if err != nil {
		return fmt.Errorf("errore nella conversione del nickname in ID: %s", err.Error())
	}
	idConversazione, err := db.CreaConversazioneDB()
	if err != nil {
		return fmt.Errorf("errore durante l'inserimento della nuova conversazione: %s", err.Error())
	}
	queryDiInserimento := `INSERT INTO gruppo (nome, conversazione, foto) VALUES (?,?,?);`
	result, err := db.c.Exec(queryDiInserimento, nomeGruppo_Passato, idConversazione, idfoto_Passata)
	if err != nil {
		return fmt.Errorf("errore durante la creazione del gruppo: %s", err.Error())
	}
	LastInsertIdid, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID del gruppo: %s", err.Error())
	}

	fmt.Println("ID gruppo creato:", LastInsertIdid)
	queryDiInserimentoUtente := `INSERT INTO utenteingruppo (utente, gruppo) VALUES (?, ?);`
	_, err = db.c.Exec(queryDiInserimentoUtente, utenteconvertito, LastInsertIdid)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiunta dell'utente al gruppo: %w", err)
	}

	return nil
}

// funzione che crea una conversazione privata, ritorna l'id della conversazione
func (db *appdbimpl) CreaConversazionePrivataDB(utente1_Passato string, utente2_Passato string) (int, error) {
	//controlliamo che gli utenti coinvolti esistano
	esistenza, err := db.EsistenzaUtente(utente1_Passato)
	if err != nil {
		return 0, fmt.Errorf("errore durante il controllo dell'esistenza dell'utente (chiamante): %s", err.Error())
	}
	if !esistenza {
		return 0, fmt.Errorf("l'utente %s non esiste (chiamante)", utente1_Passato)
	}
	esistenza, err = db.EsistenzaUtente(utente2_Passato)
	if err != nil {
		return 0, fmt.Errorf("errore durante il controllo dell'esistenza dell'utente (chiamato): %s", err.Error())
	}
	if !esistenza {
		return 0, fmt.Errorf("l'utente %s non esiste (chiamato)", utente2_Passato)
	}

	//controlliamo che già non esista una conversazione tra i due, se non esiste la creiamo
	idConversazione, err := db.EsisteConversazioneTraUtenti(utente1_Passato, utente2_Passato)
	if err != nil {
		return 0, fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %s", err.Error())
	}
	if idConversazione > 0 {
		return 0, fmt.Errorf("la conversazione tra %s e %s esiste già", utente1_Passato, utente2_Passato)
	}
	idConversazione, err = db.CreaConversazioneDB()
	if err != nil {
		return 0, fmt.Errorf("errore durante la creazione della nuova conversazione: %s", err.Error())
	}
	utente1_ID, err := db.IdUtenteDaNickname(utente1_Passato)
	if err != nil {
		return 0, fmt.Errorf("errore nella conversione del nickname %s in ID per utente1: %s", utente1_Passato, err.Error())
	}
	utente2_ID, err := db.IdUtenteDaNickname(utente2_Passato)
	if err != nil {
		return 0, fmt.Errorf("errore nella conversione del nickname %s in ID per utente2: %s", utente2_Passato, err.Error())
	}
	queryDiInserimento := `INSERT INTO conversazioneprivata (conversazione, utente1, utente2) VALUES (?, ?, ?);`
	result, err := db.c.Exec(queryDiInserimento, idConversazione, utente1_ID, utente2_ID)
	if err != nil {
		return 0, fmt.Errorf("errore durante l'inserimento della conversazione nel database: %s", err.Error())
	}

	LastInsertIdid, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("errore durante il recupero dell'ID della conversazione appena creata: %s", err.Error())
	}

	fmt.Println("ID della conversazione privata creata:", LastInsertIdid)

	return int(idConversazione), nil
}
