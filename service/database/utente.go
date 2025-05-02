package database

import (
	"database/sql"
	"fmt"
)

// Funzione che crea la tabella utente se questa non esiste
// l'utente avrà un id, un nickname, una foto.
func CreaTabellaUtente(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS utente (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nickname TEXT UNIQUE NOT NULL,
			foto,
			FOREIGN KEY (foto) REFERENCES foto(id)
		);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("errore durante la creazione della tabella utente: %w", err)
	}
	return nil
}

// CreaUtente crea un nuovo utente nel database con nickname "nicknamePassato" e id della foto "idfotoPassata",
// l'utente avrà un nuovo id diverso dagli altri (autoincrementante)
func (db *appdbimpl) CreaUtente(nicknamePassato string, idfotoPassata int) error {
	esistenza, err := db.EsistenzaUtente(nicknamePassato)
	if esistenza {
		return fmt.Errorf("il nickname è già in uso: %w", err)
	}

	// Eseguiamo la query da inserimento
	queryDiInserimento := `INSERT INTO utente (nickname, foto) VALUES (?,?);`

	// Eseguiamo la query e otteniamo il risultato
	_, err = db.c.Exec(queryDiInserimento, nicknamePassato, idfotoPassata)
	if err != nil {
		return fmt.Errorf("errore inaspettato durante la creazione dell'utente: %w", err)
	}
	return nil
}

func (db *appdbimpl) Login(nicknamePassato string) (string, error) {
	esistenza, err := db.EsistenzaUtente(nicknamePassato)
	if !esistenza {
		err = db.CreaUtente(nicknamePassato, 0)
		if err != nil {
			return "", fmt.Errorf("errore durante la creazione del nuovo utente a seguito della sua non esistenza nel database")
		}
	}
	if err != nil {
		return "", fmt.Errorf("errore durante la verifica dell'esistenza dell'utente")
	}

	return nicknamePassato, nil
}

// funzione che controlla se l'utente con nome "nickname_Passato" esiste già
func (db *appdbimpl) EsistenzaUtente(nicknamePassato string) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM utente WHERE nickname = ?;`
	// Eseguiamo la query per verificare l'esistenza dell'utente
	err := db.c.QueryRow(query, nicknamePassato).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			// Se non ci sono righe, significa che l'utente non esiste
			return false, nil
		}
		// Se c'è un altro errore, lo solleviamo
		return false, fmt.Errorf("errore inaspettato durante la verifica dell'esistenza: %w", err)
	}
	// Se count > 0, significa che l'utente "nickname_Passato" esiste, sollevo un errore
	return count > 0, nil
}

// IdUtenteDaNickname restituisce l'ID dell'utente dato il nickname
func (db *appdbimpl) IdUtenteDaNickname(nicknamePassato string) (int, int, error) {
	var id int
	query := `SELECT id FROM utente WHERE nickname = ?;`
	err := db.c.QueryRow(query, nicknamePassato).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 404, fmt.Errorf("utente con nickname '%s' non trovato", nicknamePassato)
		}
		return 0, 500, fmt.Errorf("errore durante il recupero dell'ID utente: %w", err)
	}
	return id, 0, nil
}

// NicknameUtenteDaId restituisce il nickname dell'utente dato il suo ID
func (db *appdbimpl) NicknameUtenteDaId(idPassato int) (string, int, error) {
	var nickname string
	query := `SELECT nickname FROM utente WHERE id = ?;`
	err := db.c.QueryRow(query, idPassato).Scan(&nickname)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 404, fmt.Errorf("utente con ID '%d' non trovato", idPassato)
		}
		return "", 500, fmt.Errorf("errore durante il recupero del nickname: %w", err)
	}
	return nickname, 0, nil
}

// impostazioni utente

// ImpostaFotoProfilo aggiorna la foto del profilo dell'utente
func (db *appdbimpl) ImpostaFotoProfilo(nicknamePassato string, idfotoPassata int) error {
	// Verifica che l'utente esista
	esiste, err := db.EsistenzaUtente(nicknamePassato)
	if err != nil {
		return fmt.Errorf("errore durante il controllo dell'esistenza dell'utente: %w", err)
	}
	if !esiste {
		return fmt.Errorf("l'utente %s non esiste", nicknamePassato)
	}

	// Aggiorna la foto del profilo
	queryUpdateFoto := `UPDATE utente SET foto = ? WHERE nickname = ?`
	_, err = db.c.Exec(queryUpdateFoto, idfotoPassata, nicknamePassato)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento della foto del profilo: %w", err)
	}

	return nil
}

// ImpostaNome aggiorna il nickname dell'utente
func (db *appdbimpl) ImpostaNome(nicknamePassato string, nuovoNickPassato string) (int, error) {
	// Verifica che l'utente esista
	esiste, err := db.EsistenzaUtente(nicknamePassato)
	if err != nil {
		return 500, fmt.Errorf("errore durante il controllo dell'esistenza dell'utente: %w", err)
	}
	if !esiste {
		return 404, fmt.Errorf("l'utente %s non esiste", nicknamePassato)
	}

	// Verifica che il nuovo nickname non sia già in uso
	esisteNuovoNick, err := db.EsistenzaUtente(nuovoNickPassato)
	if err != nil {
		return 500, fmt.Errorf("errore durante il controllo del nuovo nickname: %w", err)
	}
	if esisteNuovoNick {
		return 304, fmt.Errorf("il nickname %s è già in uso", nuovoNickPassato)
	}

	// Verifica che il nuovo nickname non sia uguale al vecchio
	if nicknamePassato == nuovoNickPassato {
		return 304, fmt.Errorf("il nuovo nickname è uguale a quello vecchio")
	}

	// Aggiorna il nickname dell'utente
	queryUpdateNome := `UPDATE utente SET nickname = ? WHERE nickname = ?`
	_, err = db.c.Exec(queryUpdateNome, nuovoNickPassato, nicknamePassato)
	if err != nil {
		return 500, fmt.Errorf("errore durante l'aggiornamento del nickname: %w", err)
	}

	return 0, nil
}

// Struttura per memorizzare il nome e la foto profilo dell'utente
type Profilo struct {
	Nickname string
}

// Funzione di test.
// Funzione per ottenere tutti i nomi e foto profilo degli utenti
func (db *appdbimpl) UsersInGroup(chiamante string, chat int) ([]Profilo, int, error) {
	esiste, err := db.EsisteConversazione(chat)
	if err != nil {
		return nil, 500, fmt.Errorf("errore durante la verifica dell'esistenza della conversazione: %w", err)
	}
	if !esiste {
		return nil, 404, fmt.Errorf("la conversazione con ID %d non esiste", chat)
	}
	esisteUtenteChiamante, err := db.EsistenzaUtente(chiamante)
	if err != nil {
		return nil, 500, fmt.Errorf("errore durante il controllo dell'esistenza dell'utente chiamante %s: %w", chiamante, err)
	}
	if !esisteUtenteChiamante {
		return nil, 404, fmt.Errorf("l'utente chiamante %s non esiste", chiamante)
	}
	chiamantePresente, codiceErrore, err := db.UtenteCoinvoltoGruppo(chiamante, chat)
	if err != nil {
		return nil, codiceErrore, fmt.Errorf("errore durante il controllo della presenza dell'utente nel gruppo: %w", err)
	}
	if chiamantePresente == 0 {
		return nil, 401, fmt.Errorf("l'utente %s non fa parte del gruppo", chiamante)
	}

	query := `SELECT u.nickname
			  FROM utente as u
			  JOIN utenteingruppo as uig ON u.id = uig.utente
			  JOIN gruppo as g ON g.id = uig.gruppo
			  WHERE g.id = ?
			  `

	rows, err := db.c.Query(query, chat)
	if err != nil {
		return nil, 500, fmt.Errorf("errore durante il recupero dei profili utente: %w", err)
	}

	var lista []Profilo

	for rows.Next() {
		var nickname string

		if err := rows.Scan(&nickname); err != nil {
			return nil, 500, fmt.Errorf("errore durante la lettura dei dati: %w", err)
		}

		lista = append(lista, Profilo{
			Nickname: nickname,
		})
	}

	return lista, 0, nil
}
