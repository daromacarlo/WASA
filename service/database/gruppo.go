package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) AggiungiAGruppoDB(idConversazione int, UtenteChiamante string, UtenteDaAggiungere string) (int, error) {
	esisteUtenteChiamante, err := db.EsistenzaUtente(UtenteChiamante)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il controllo dell'esistenza dell'utente chiamante %s: %w", UtenteChiamante, err)
	}
	if !esisteUtenteChiamante {
		return 404, fmt.Errorf("l'utente chiamante %s non esiste", UtenteChiamante)
	}
	chiamantePresente, codiceErrore, err := db.UtenteCoinvoltoGruppo(UtenteChiamante, idConversazione)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante il controllo della presenza dell'utente nel gruppo: %w", err)
	}
	if chiamantePresente == 0 {
		return 401, fmt.Errorf("l'utente %s non fa parte del gruppo", UtenteChiamante)
	}
	esisteUtentedaAggiungere, err := db.EsistenzaUtente(UtenteDaAggiungere)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il controllo dell'esistenza dell'utente da aggiungere %s: %w", UtenteDaAggiungere, err)
	}
	if !esisteUtentedaAggiungere {
		return 404, fmt.Errorf("l'utente da aggiungere %s non esiste", UtenteDaAggiungere)
	}
	utenteGiaPresente, codiceErrore, err := db.UtenteCoinvoltoGruppo(UtenteDaAggiungere, idConversazione)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante il controllo della presenza dell'utente nel gruppo: %w", err)
	}
	if utenteGiaPresente > 0 {
		return 304, fmt.Errorf("l'utente %s è già presente nel gruppo", UtenteDaAggiungere)
	}
	utenteDaAggiungere_convertito, codiceErrore, err := db.IdUtenteDaNickname(UtenteDaAggiungere)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore nella conversione del nickname %s in ID: %w", UtenteDaAggiungere, err)
	}
	esiste, err := db.EsisteConversazione(idConversazione)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %w", err)
	}
	if !esiste {
		return 404, fmt.Errorf("la conversazione con ID %d non esiste", idConversazione)
	}
	idGruppo, codiceErrore, err := db.CercaConversazioneGruppo(idConversazione)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante la ricerca del gruppo associato alla conversazione: %w", err)
	}
	queryDiInserimento := `INSERT INTO utenteingruppo (utente, gruppo) VALUES (?, ?);`
	_, err = db.c.Exec(queryDiInserimento, utenteDaAggiungere_convertito, idGruppo)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante l'aggiunta dell'utente al gruppo: %w", err)
	}
	return 0, nil
}

func (db *appdbimpl) LasciaGruppo(idConversazione int, UtenteChiamante string) (int, error) {
	esisteUtenteChiamante, err := db.EsistenzaUtente(UtenteChiamante)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il controllo dell'esistenza dell'utente chiamante %s: %w", UtenteChiamante, err)
	}
	if !esisteUtenteChiamante {
		return 404, fmt.Errorf("l'utente chiamante %s non esiste", UtenteChiamante)
	}
	chiamantePresente, codiceErrore, err := db.UtenteCoinvoltoGruppo(UtenteChiamante, idConversazione)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante il controllo della presenza dell'utente nel gruppo: %w", err)
	}
	if chiamantePresente == 0 {
		return 401, fmt.Errorf("l'utente %s non fa parte del gruppo", UtenteChiamante)
	}
	utenteChiamante_convertito, codiceErrore, err := db.IdUtenteDaNickname(UtenteChiamante)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore nella conversione del nickname %s in ID: %w", UtenteChiamante, err)
	}
	esiste, err := db.EsisteConversazione(idConversazione)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %w", err)
	}
	if !esiste {
		return 404, fmt.Errorf("la conversazione con ID %d non esiste", idConversazione)
	}
	idGruppo, codiceErrore, err := db.CercaConversazioneGruppo(idConversazione)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante la ricerca del gruppo associato alla conversazione: %w", err)
	}
	queryDiInserimento := `DELETE from utenteingruppo WHERE utente = ? AND gruppo = ?;`
	_, err = db.c.Exec(queryDiInserimento, utenteChiamante_convertito, idGruppo)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante la rimozione dell'utente al gruppo: %w", err)
	}
	return 0, nil
}

func (db *appdbimpl) ImpostaFotoGruppo(UtenteChiamante string, id_foto_Passata int, id_gruppo_Passato int) (int, error) {
	esiste, err := db.EsisteConversazione(id_gruppo_Passato)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il controllo dell'esistenza della conversazione: %w", err)
	}
	if !esiste {
		return 404, fmt.Errorf("la conversazione specificata non esiste")
	}

	tipoGruppo, codiceErrore, err := db.CercaConversazioneGruppo(id_gruppo_Passato)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante il controllo del tipo di conversazione: %w", err)
	}
	if tipoGruppo == 0 {
		return 400, fmt.Errorf("la conversazione specificata non è un gruppo")
	}
	coinvolto, codiceErrore, err := db.UtenteCoinvoltoGruppo(UtenteChiamante, id_gruppo_Passato)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante il controllo della partecipazione dell'utente: %w", err)
	}
	if coinvolto == 0 {
		return 401, fmt.Errorf("l'utente %s non fa parte del gruppo %d", UtenteChiamante, id_gruppo_Passato)
	}
	query := `UPDATE gruppo SET foto = ? WHERE conversazione = ?`
	_, err = db.c.Exec(query, id_foto_Passata, id_gruppo_Passato)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante l'aggiornamento della foto del gruppo: %w", err)
	}
	return 0, nil
}

func (db *appdbimpl) ImpostaNomeGruppo(UtenteChiamante string, nomeGruppo_Passato string, id_gruppo_Passato int) (int, error) {
	esiste, err := db.EsisteConversazione(id_gruppo_Passato)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il controllo dell'esistenza della conversazione: %w", err)
	}
	if !esiste {
		return 404, fmt.Errorf("la conversazione specificata non esiste")
	}
	tipoGruppo, codiceErrore, err := db.CercaConversazioneGruppo(id_gruppo_Passato)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante il controllo del tipo di conversazione: %w", err)
	}
	if tipoGruppo == 0 {
		return 400, fmt.Errorf("la conversazione specificata non è un gruppo")
	}

	coinvolto, codiceErrore, err := db.UtenteCoinvoltoGruppo(UtenteChiamante, id_gruppo_Passato)
	if !errors.Is(err, nil) {
		return codiceErrore, fmt.Errorf("errore durante il controllo della partecipazione dell'utente: %w", err)
	}
	if coinvolto == 0 {
		return 401, fmt.Errorf("l'utente %s non fa parte del gruppo %d", UtenteChiamante, id_gruppo_Passato)
	}

	var vecchioNome string
	queryGetNome := `SELECT nome FROM gruppo WHERE conversazione = ?`
	err = db.c.QueryRow(queryGetNome, id_gruppo_Passato).Scan(&vecchioNome)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante il recupero del nome del gruppo: %w", err)
	}

	if vecchioNome == nomeGruppo_Passato {
		return 304, fmt.Errorf("il nuovo nome del gruppo è uguale al vecchio")
	}

	queryUpdate := `UPDATE gruppo SET nome = ? WHERE conversazione = ?`
	_, err = db.c.Exec(queryUpdate, nomeGruppo_Passato, id_gruppo_Passato)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("errore durante l'aggiornamento del nome del gruppo: %w", err)
	}
	return 0, nil
}
