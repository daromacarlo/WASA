package database

import (
	"fmt"
)

// ImpostaRisposta imposta l'attributo 'risposta' nel messaggio con l'ID del vecchio messaggio
func (db *appdbimpl) ImpostaRisposta(IdMessaggio int, IdNuovoMessaggio int, chat int) error {
	// Esegui un'operazione di aggiornamento sulla tabella del messaggio
	// Imposto la colonna 'risposta' del messaggio con l'ID del messaggio a cui si sta rispondendo
	query := `
		UPDATE messaggio
		SET risposta = ?
		WHERE id = ? AND conversazione = ?
	`
	_, err := db.c.Exec(query, IdMessaggio, IdNuovoMessaggio, chat)
	if err != nil {
		return fmt.Errorf("errore inaspettato durante l'impostazione della risposta per il messaggio con ID %d: %w", IdNuovoMessaggio, err)
	}
	return nil
}

// RispondiMessaggioPrivatoTesto gestisce la risposta a un messaggio privato
func (db *appdbimpl) RispondiMessaggioPrivatoTesto(utentePassato string, destinatarioPassato string, IdMessaggio int, testoPassato string) error {
	// Verifica se il messaggio con l'ID passato esiste
	esistenza, err := db.EsisteMessaggio(IdMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante la verifica dell'esistenza del messaggio con ID %d: %w", IdMessaggio, err)
	}
	if !esistenza {
		return fmt.Errorf("errore: il messaggio con ID %d non esiste", IdMessaggio)
	}

	// Verifica se esiste una conversazione privata tra gli utenti
	chat, err := db.UtenteCoinvoltoPrivato(utentePassato, destinatarioPassato)
	if err != nil {
		return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
	}

	// Recupera l'ID della conversazione a partire dal messaggio
	conversazioneID, err := db.GetConversazioneIdByMessaggio(IdMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID della conversazione: %w", err)
	}

	// Recupera l'ID dell'utente a partire dal nickname
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione del nickname a ID utente: %w", err)
	}

	// Verifica se la conversazione è un gruppo o privata
	isGruppo, err := db.CercaConversazioneGruppo(conversazioneID)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}

	if isGruppo > 0 {
		// Se è un gruppo, verifica che l'utente sia coinvolto
		coinvolto, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazioneID)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return fmt.Errorf("l'utente non è membro del gruppo")
		}
	} else {
		// Se è privata, verifica che l'utente sia coinvolto nella conversazione privata
		idPrivata, err := db.CercaConversazionePrivata(conversazioneID, utente_Passato_convertito)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
		}
		if idPrivata == 0 {
			return fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}
	}

	// Crea il nuovo messaggio testuale nel database
	err = db.CreaMessaggioTestualeDBPrivato(utentePassato, destinatarioPassato, testoPassato)
	if err != nil {
		return fmt.Errorf("errore durante la creazione del messaggio: %w", err)
	}

	// Recupera l'ID del nuovo messaggio appena creato
	var idNuovoMessaggio int
	queryID := `SELECT last_insert_rowid()` // Questa query recupera l'ID dell'ultimo messaggio inserito
	err = db.c.QueryRow(queryID).Scan(&idNuovoMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID del nuovo messaggio: %w", err)
	}

	// Debug: Stampa l'ID del nuovo messaggio per il debug
	fmt.Printf("ID del nuovo messaggio: %d\n", idNuovoMessaggio)

	// Imposta il nuovo messaggio come risposta al messaggio precedente
	err = db.ImpostaRisposta(IdMessaggio, idNuovoMessaggio, chat)
	if err != nil {
		return fmt.Errorf("errore durante l'impostazione della risposta per il messaggio con ID %d: %w", idNuovoMessaggio, err)
	}

	return nil
}

func (db *appdbimpl) RispondiMessaggioPrivatoFoto(utentePassato string, destinatarioPassato string, idMessaggio int, fotoPassato int) error {
	esistenza, err := db.EsisteMessaggio(idMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante la verifica dell'esistenza del messaggio con ID %d: %w", idMessaggio, err)
	}
	if !esistenza {
		return fmt.Errorf("errore: il messaggio con ID %d non esiste", idMessaggio)
	}
	chat, err := db.UtenteCoinvoltoPrivato(utentePassato, destinatarioPassato)
	if err != nil {
		return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
	}
	conversazioneID, err := db.GetConversazioneIdByMessaggio(idMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID della conversazione: %w", err)
	}
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione del nickname a ID utente: %w", err)
	}
	isGruppo, err := db.CercaConversazioneGruppo(conversazioneID)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	if isGruppo > 0 {
		coinvolto, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazioneID)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return fmt.Errorf("l'utente non è membro del gruppo")
		}
	} else {
		idPrivata, err := db.CercaConversazionePrivata(conversazioneID, utente_Passato_convertito)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
		}
		if idPrivata == 0 {
			return fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}
	}
	err = db.CreaMessaggioFotoDBPrivato(utentePassato, destinatarioPassato, fotoPassato)
	if err != nil {
		return fmt.Errorf("errore durante la creazione del messaggio: %w", err)
	}
	var idNuovoMessaggio int
	queryID := `SELECT last_insert_rowid()`
	err = db.c.QueryRow(queryID).Scan(&idNuovoMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID del nuovo messaggio: %w", err)
	}
	fmt.Printf("ID del nuovo messaggio: %d\n", idNuovoMessaggio)
	err = db.ImpostaRisposta(idMessaggio, idNuovoMessaggio, chat)
	if err != nil {
		return fmt.Errorf("errore durante l'impostazione della risposta per il messaggio con ID %d: %w", idNuovoMessaggio, err)
	}
	return nil
}

// RispondiMessaggioPrivatoTesto gestisce la risposta a un messaggio privato
func (db *appdbimpl) RispondiMessaggioGruppoTesto(utentePassato string, idGruppoPassato int, idMessaggio int, testoPassato string) error {
	// Verifica se il messaggio con l'ID passato esiste
	esistenza, err := db.EsisteMessaggio(idMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante la verifica dell'esistenza del messaggio con ID %d: %w", idMessaggio, err)
	}
	if !esistenza {
		return fmt.Errorf("errore: il messaggio con ID %d non esiste", idMessaggio)
	}

	// Verifica se esiste una conversazione privata tra gli utenti
	chat, err := db.UtenteCoinvoltoGruppo(utentePassato, idGruppoPassato)
	if err != nil {
		return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
	}

	// Recupera l'ID della conversazione a partire dal messaggio
	conversazioneID, err := db.GetConversazioneIdByMessaggio(idMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID della conversazione: %w", err)
	}

	//controlla che il messaggio sia della conversazione corrente
	if conversazioneID != idGruppoPassato {
		return fmt.Errorf("il messaggio non fa parte della chat indicata: %w", err)
	}
	// Recupera l'ID dell'utente a partire dal nickname
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione del nickname a ID utente: %w", err)
	}

	// Verifica se la conversazione è un gruppo o privata
	isGruppo, err := db.CercaConversazioneGruppo(conversazioneID)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}

	if isGruppo > 0 {
		// Se è un gruppo, verifica che l'utente sia coinvolto
		coinvolto, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazioneID)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return fmt.Errorf("l'utente non è membro del gruppo")
		}
	} else {
		// Se è privata, verifica che l'utente sia coinvolto nella conversazione privata
		idPrivata, err := db.CercaConversazionePrivata(conversazioneID, utente_Passato_convertito)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
		}
		if idPrivata == 0 {
			return fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}
	}

	// Crea il nuovo messaggio testuale nel database
	err = db.CreaMessaggioTestualeDBGruppo(utentePassato, idGruppoPassato, testoPassato)
	if err != nil {
		return fmt.Errorf("errore durante la creazione del messaggio: %w", err)
	}

	// Recupera l'ID del nuovo messaggio appena creato
	var idNuovoMessaggio int
	queryID := `SELECT last_insert_rowid()` // Questa query recupera l'ID dell'ultimo messaggio inserito
	err = db.c.QueryRow(queryID).Scan(&idNuovoMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID del nuovo messaggio: %w", err)
	}

	// Debug: Stampa l'ID del nuovo messaggio per il debug
	fmt.Printf("ID del nuovo messaggio: %d\n", idNuovoMessaggio)

	// Imposta il nuovo messaggio come risposta al messaggio precedente
	err = db.ImpostaRisposta(idMessaggio, idNuovoMessaggio, chat)
	if err != nil {
		return fmt.Errorf("errore durante l'impostazione della risposta per il messaggio con ID %d: %w", idNuovoMessaggio, err)
	}

	return nil
}

func (db *appdbimpl) RispondiMessaggioGruppoFoto(utentePassato string, idGruppoPassato int, IdMessaggio int, foto_Passato int) error {
	esistenza, err := db.EsisteMessaggio(IdMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante la verifica dell'esistenza del messaggio con ID %d: %w", IdMessaggio, err)
	}
	if !esistenza {
		return fmt.Errorf("errore: il messaggio con ID %d non esiste", IdMessaggio)
	}
	conversazioneID, err := db.GetConversazioneIdByMessaggio(IdMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID della conversazione: %w", err)
	}

	if conversazioneID != idGruppoPassato {
		return fmt.Errorf("il messaggio non fa parte della chat indicata: %w", err)
	}
	utente_Passato_convertito, err := db.IdUtenteDaNickname(utentePassato)
	if err != nil {
		return fmt.Errorf("errore durante la conversione del nickname a ID utente: %w", err)
	}
	chat, err := db.UtenteCoinvoltoGruppo(utentePassato, idGruppoPassato)
	if err != nil {
		return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
	}
	isGruppo, err := db.CercaConversazioneGruppo(conversazioneID)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del tipo di conversazione: %w", err)
	}
	if isGruppo > 0 {
		coinvolto, err := db.UtenteCoinvoltoGruppo(utentePassato, conversazioneID)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della partecipazione dell'utente al gruppo: %w", err)
		}
		if coinvolto == 0 {
			return fmt.Errorf("l'utente non è membro del gruppo")
		}
	} else {
		idPrivata, err := db.CercaConversazionePrivata(conversazioneID, utente_Passato_convertito)
		if err != nil {
			return fmt.Errorf("errore durante la verifica della conversazione privata: %w", err)
		}
		if idPrivata == 0 {
			return fmt.Errorf("l'utente non è coinvolto nella conversazione privata")
		}
	}
	err = db.CreaMessaggioFotoDBGruppo(utentePassato, idGruppoPassato, foto_Passato)
	if err != nil {
		return fmt.Errorf("errore durante la creazione del messaggio: %w", err)
	}
	var idNuovoMessaggio int
	queryID := `SELECT last_insert_rowid()`
	err = db.c.QueryRow(queryID).Scan(&idNuovoMessaggio)
	if err != nil {
		return fmt.Errorf("errore durante il recupero dell'ID del nuovo messaggio: %w", err)
	}
	fmt.Printf("ID del nuovo messaggio: %d\n", idNuovoMessaggio)
	err = db.ImpostaRisposta(IdMessaggio, idNuovoMessaggio, chat)
	if err != nil {
		return fmt.Errorf("errore durante l'impostazione della risposta per il messaggio con ID %d: %w", idNuovoMessaggio, err)
	}
	return nil
}
