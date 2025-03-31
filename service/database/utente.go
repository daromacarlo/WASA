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
		db.CreaUtente(nicknamePassato, 0)
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
func (db *appdbimpl) IdUtenteDaNickname(nicknamePassato string) (int, error) {
	var id int
	query := `SELECT id FROM utente WHERE nickname = ?;`
	err := db.c.QueryRow(query, nicknamePassato).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("utente con nickname '%s' non trovato", nicknamePassato)
		}
		return 0, fmt.Errorf("errore durante il recupero dell'ID utente: %w", err)
	}
	return id, nil
}

// NicknameUtenteDaId restituisce il nickname dell'utente dato il suo ID
func (db *appdbimpl) NicknameUtenteDaId(idPassato int) (string, error) {
	var nickname string
	query := `SELECT nickname FROM utente WHERE id = ?;`
	err := db.c.QueryRow(query, idPassato).Scan(&nickname)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("utente con ID '%d' non trovato", idPassato)
		}
		return "", fmt.Errorf("errore durante il recupero del nickname: %w", err)
	}
	return nickname, nil
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
func (db *appdbimpl) ImpostaNome(nicknamePassato string, nuovoNickPassato string) error {
	// Verifica che l'utente esista
	esiste, err := db.EsistenzaUtente(nicknamePassato)
	if err != nil {
		return fmt.Errorf("errore durante il controllo dell'esistenza dell'utente: %w", err)
	}
	if !esiste {
		return fmt.Errorf("l'utente %s non esiste", nicknamePassato)
	}

	// Verifica che il nuovo nickname non sia già in uso
	esisteNuovoNick, err := db.EsistenzaUtente(nuovoNickPassato)
	if err != nil {
		return fmt.Errorf("errore durante il controllo del nuovo nickname: %w", err)
	}
	if esisteNuovoNick {
		return fmt.Errorf("il nickname %s è già in uso", nuovoNickPassato)
	}

	// Verifica che il nuovo nickname non sia uguale al vecchio
	if nicknamePassato == nuovoNickPassato {
		return fmt.Errorf("il nuovo nickname è uguale a quello vecchio")
	}

	// Aggiorna il nickname dell'utente
	queryUpdateNome := `UPDATE utente SET nickname = ? WHERE nickname = ?`
	_, err = db.c.Exec(queryUpdateNome, nuovoNickPassato, nicknamePassato)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento del nickname: %w", err)
	}

	return nil
}

// Struttura per memorizzare il nome e la foto profilo dell'utente
type Profilo struct {
	Nickname string
	Foto     *string
}

// Funzione di test.
// Funzione per ottenere tutti i nomi e foto profilo degli utenti
func (db *appdbimpl) VediProfili(chiamante string) ([]Profilo, error) {
	// Query per ottenere il nickname e la foto profilo di tutti gli utenti
	esiste, err := db.EsistenzaUtente(chiamante)
	if err != nil {
		return nil, fmt.Errorf("errore durante il controllo dell'esistenza dell'utente: %w", err)
	}
	if !esiste {
		return nil, fmt.Errorf("l'utente %s non esiste", chiamante)
	}

	query := `SELECT u.nickname, f.foto
			  FROM utente as u
			  JOIN foto as f ON f.id = u.foto`

	rows, err := db.c.Query(query)
	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero dei profili utente: %w", err)
	}
	defer rows.Close()

	// Array per memorizzare i profili utente
	var userProfiles []Profilo

	// Itera sui risultati della query
	for rows.Next() {
		var nickname string
		var foto *string

		// Scansiona i dati della riga nella struttura
		if err := rows.Scan(&nickname, &foto); err != nil {
			return nil, fmt.Errorf("errore durante la lettura dei dati: %w", err)
		}

		// Aggiungi il profilo alla lista
		userProfiles = append(userProfiles, Profilo{
			Nickname: nickname,
			Foto:     foto,
		})
	}

	// Controlla se ci sono errori durante l'iterazione
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("errore durante l'iterazione del risultato: %w", err)
	}

	return userProfiles, nil
}
