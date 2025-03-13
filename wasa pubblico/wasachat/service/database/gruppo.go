package database

import (
	"fmt"
)

// Funzione che consente di aggiungere persone ad un gruppo
func (db *appdbimpl) AggiungiAGruppoDB(idConversazione int, UtenteChiamante string, UtenteDaAggiungere string) error {
	//controlliamo che l'utente chiamante esista
	esisteUtenteChiamante, err := db.EsistenzaUtente(UtenteChiamante)
	if err != nil {
		return fmt.Errorf("errore durante il controllo dell'esistenza dell'utente chiamante %s: %w", UtenteChiamante, err)
	}
	if !esisteUtenteChiamante {
		return fmt.Errorf("l'utente chiamante %s non esiste", UtenteChiamante)
	}
	//controlliamo che questo sia presente nel gruppo
	chiamantePresente, err := db.UtenteCoinvoltoGruppo(UtenteChiamante, idConversazione)
	if err != nil {
		return fmt.Errorf("errore durante il controllo della presenza dell'utente nel gruppo: %w", err)
	}
	if chiamantePresente == 0 {
		return fmt.Errorf("l'utente %s non fa parte del gruppo", UtenteChiamante)
	}
	// controlliamo che l'utente da aggiungere esista
	esisteUtentedaAggiungere, err := db.EsistenzaUtente(UtenteDaAggiungere)
	if err != nil {
		return fmt.Errorf("errore durante il controllo dell'esistenza dell'utente da aggiungere %s: %w", UtenteDaAggiungere, err)
	}
	if !esisteUtentedaAggiungere {
		return fmt.Errorf("l'utente da aggiungere %s non esiste", UtenteDaAggiungere)
	}
	//controlliamo che non sia già presente nel gruppo
	utenteGiaPresente, err := db.UtenteCoinvoltoGruppo(UtenteDaAggiungere, idConversazione)
	if err != nil {
		return fmt.Errorf("errore durante il controllo della presenza dell'utente nel gruppo: %w", err)
	}
	if utenteGiaPresente > 0 {
		return fmt.Errorf("l'utente %s è già presente nel gruppo", UtenteDaAggiungere)
	}
	//troviamo il suo id
	utenteDaAggiungere_convertito, err := db.IdUtenteDaNickname(UtenteDaAggiungere)
	if err != nil {
		return fmt.Errorf("errore nella conversione del nickname %s in ID: %w", UtenteDaAggiungere, err)
	}
	//controlliamo che esista la conversazione
	esiste, err := db.EsisteConversazione(idConversazione)
	if err != nil {
		return fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %w", err)
	}
	if !esiste {
		return fmt.Errorf("la conversazione con ID %d non esiste", idConversazione)
	}
	//controlliamo che questa sia un gruppo e troviamo l'id del gruppo
	idGruppo, err := db.CercaConversazioneGruppo(idConversazione)
	if err != nil {
		return fmt.Errorf("errore durante la ricerca del gruppo associato alla conversazione: %w", err)
	}
	//inseriamo l'utente nel gruppo con questa query di inserimento
	queryDiInserimento := `INSERT INTO utenteingruppo (utente, gruppo) VALUES (?, ?);`
	_, err = db.c.Exec(queryDiInserimento, utenteDaAggiungere_convertito, idGruppo)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiunta dell'utente al gruppo: %w", err)
	}
	return nil
}

// Funzione che consente di lasciare un gruppo
func (db *appdbimpl) LasciaGruppo(idConversazione int, UtenteChiamante string) error {
	//controlliamo che l'utente chiamante esista
	esisteUtenteChiamante, err := db.EsistenzaUtente(UtenteChiamante)
	if err != nil {
		return fmt.Errorf("errore durante il controllo dell'esistenza dell'utente chiamante %s: %w", UtenteChiamante, err)
	}
	if !esisteUtenteChiamante {
		return fmt.Errorf("l'utente chiamante %s non esiste", UtenteChiamante)
	}
	//controlliamo che questo sia presente nel gruppo
	chiamantePresente, err := db.UtenteCoinvoltoGruppo(UtenteChiamante, idConversazione)
	if err != nil {
		return fmt.Errorf("errore durante il controllo della presenza dell'utente nel gruppo: %w", err)
	}
	if chiamantePresente == 0 {
		return fmt.Errorf("l'utente %s non fa parte del gruppo", UtenteChiamante)
	}
	//troviamo il suo id
	utenteChiamante_convertito, err := db.IdUtenteDaNickname(UtenteChiamante)
	if err != nil {
		return fmt.Errorf("errore nella conversione del nickname %s in ID: %w", UtenteChiamante, err)
	}
	//controlliamo che esista la conversazione
	esiste, err := db.EsisteConversazione(idConversazione)
	if err != nil {
		return fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %w", err)
	}
	if !esiste {
		return fmt.Errorf("la conversazione con ID %d non esiste", idConversazione)
	}
	//controlliamo che questa sia un gruppo e troviamo l'id del gruppo
	idGruppo, err := db.CercaConversazioneGruppo(idConversazione)
	if err != nil {
		return fmt.Errorf("errore durante la ricerca del gruppo associato alla conversazione: %w", err)
	}
	//rimuovo l'utente dal gruppo con questa query
	queryDiInserimento := `DELETE from utenteingruppo WHERE utente = ? AND gruppo = ?;`
	_, err = db.c.Exec(queryDiInserimento, utenteChiamante_convertito, idGruppo)
	if err != nil {
		return fmt.Errorf("errore durante la rimozione dell'utente al gruppo: %w", err)
	}
	return nil
}

// ImpostaFotoGruppo aggiorna la foto del gruppo
func (db *appdbimpl) ImpostaFotoGruppo(UtenteChiamante string, id_foto_Passata int, id_gruppo_Passato int) error {
	// Verifica che la conversazione esista
	esiste, err := db.EsisteConversazione(id_gruppo_Passato)
	if err != nil {
		return fmt.Errorf("errore durante il controllo dell'esistenza della conversazione: %w", err)
	}
	if !esiste {
		return fmt.Errorf("la conversazione specificata non esiste")
	}

	// Verifica che la conversazione sia un gruppo
	tipoGruppo, err := db.CercaConversazioneGruppo(id_gruppo_Passato)
	if err != nil {
		return fmt.Errorf("errore durante il controllo del tipo di conversazione: %w", err)
	}
	if tipoGruppo == 0 {
		return fmt.Errorf("la conversazione specificata non è un gruppo")
	}

	// Verifica che l'utente sia coinvolto nel gruppo
	coinvolto, err := db.UtenteCoinvoltoGruppo(UtenteChiamante, id_gruppo_Passato)
	if err != nil {
		return fmt.Errorf("errore durante il controllo della partecipazione dell'utente: %w", err)
	}
	if coinvolto == 0 {
		return fmt.Errorf("l'utente %s non fa parte del gruppo %d", UtenteChiamante, id_gruppo_Passato)
	}

	// Aggiorna la foto del gruppo
	query := `UPDATE gruppo SET foto = ? WHERE id = ?`
	_, err = db.c.Exec(query, id_foto_Passata, id_gruppo_Passato)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento della foto del gruppo: %w", err)
	}
	return nil
}

// Funzione che consente di impostare il nome di un gruppo
func (db *appdbimpl) ImpostaNomeGruppo(UtenteChiamante string, nomeGruppo_Passato string, id_gruppo_Passato int) error {
	// Verifica che la conversazione esista
	esiste, err := db.EsisteConversazione(id_gruppo_Passato)
	if err != nil {
		return fmt.Errorf("errore durante il controllo dell'esistenza della conversazione: %w", err)
	}
	if !esiste {
		return fmt.Errorf("la conversazione specificata non esiste")
	}

	// Verifica che la conversazione sia un gruppo
	tipoGruppo, err := db.CercaConversazioneGruppo(id_gruppo_Passato)
	if err != nil {
		return fmt.Errorf("errore durante il controllo del tipo di conversazione: %w", err)
	}
	if tipoGruppo == 0 {
		return fmt.Errorf("la conversazione specificata non è un gruppo")
	}

	// Verifica che l'utente sia coinvolto nel gruppo
	coinvolto, err := db.UtenteCoinvoltoGruppo(UtenteChiamante, id_gruppo_Passato)
	if err != nil {
		return fmt.Errorf("errore durante il controllo della partecipazione dell'utente: %w", err)
	}
	if coinvolto == 0 {
		return fmt.Errorf("l'utente %s non fa parte del gruppo %d", UtenteChiamante, id_gruppo_Passato)
	}

	// Recupera il vecchio nome del gruppo
	var vecchioNome string
	queryGetNome := `SELECT nome FROM gruppo WHERE id = ?`
	err = db.c.QueryRow(queryGetNome, id_gruppo_Passato).Scan(&vecchioNome)
	if err != nil {
		return fmt.Errorf("errore durante il recupero del nome del gruppo: %w", err)
	}

	// Controlla se il nuovo nome è uguale al vecchio
	if vecchioNome == nomeGruppo_Passato {
		return fmt.Errorf("il nuovo nome del gruppo è uguale al vecchio")
	}

	// Aggiorna il nome del gruppo
	queryUpdate := `UPDATE gruppo SET nome = ? WHERE id = ?`
	_, err = db.c.Exec(queryUpdate, nomeGruppo_Passato, id_gruppo_Passato)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento del nome del gruppo: %w", err)
	}
	return nil
}
