package database

import (
	"database/sql"
	"fmt"
	"sort"
)

type Conversazione struct {
	Id       int     `json:"chat_id"`
	Nome     string  `json:"nome"`
	Foto     *string `json:"foto"`
	Time     *string `json:"time"`
	Ultimo   *string `json:"ultimosnip"`
	IsGruppo bool    `json:"isgruppo"`
}

func (db *appdbimpl) GetConversazioni(utente_Passato_string string) ([]Conversazione, error) {
	utente_Passato, err := db.IdUtenteDaNickname(utente_Passato_string)
	if err != nil {
		return nil, err
	}

	queryConversazioniPrivate := `
    SELECT c.id, 
           CASE 
               WHEN cp.utente1 = ? THEN u2.nickname 
               ELSE u1.nickname 
           END AS nome,
           MAX(m.tempo) AS time, 
           m.testo AS ultimosnip,
           CASE 
               WHEN cp.utente1 = ? THEN f2.foto
               ELSE f1.foto
           END AS foto,
           false AS isgruppo
    FROM conversazione AS c
    JOIN conversazioneprivata AS cp ON c.id = cp.conversazione
    JOIN utente AS u1 ON u1.id = cp.utente1
    JOIN utente AS u2 ON u2.id = cp.utente2
    LEFT JOIN foto AS f1 ON f1.id = u1.foto
    LEFT JOIN foto AS f2 ON f2.id = u2.foto
    LEFT JOIN messaggio AS m ON m.conversazione = c.id
    WHERE cp.utente1 = ? OR cp.utente2 = ?
    GROUP BY c.id, nome, f1.foto, f2.foto
    `

	queryConversazioniGruppo := `
    SELECT c.id, g.nome, MAX(m.tempo) AS time, m.testo AS ultimosnip, f.foto AS foto, true AS isgruppo
    FROM conversazione AS c
    JOIN gruppo AS g ON g.conversazione = c.id
    JOIN utenteingruppo AS ug ON g.id = ug.gruppo
    JOIN foto AS f ON f.id = g.foto
    LEFT JOIN messaggio AS m ON m.conversazione = c.id
    WHERE ug.utente = ?
    GROUP BY c.id, g.nome, f.foto
    `

	var conversazioni []Conversazione
	rowsPrivate, err := db.c.Query(queryConversazioniPrivate, utente_Passato, utente_Passato, utente_Passato, utente_Passato)
	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero delle conversazioni private: %w", err)
	}
	defer rowsPrivate.Close()

	var idsConversazioniPrivate []int
	for rowsPrivate.Next() {
		var conv Conversazione
		var foto sql.NullString
		if err := rowsPrivate.Scan(&conv.Id, &conv.Nome, &conv.Time, &conv.Ultimo, &foto, &conv.IsGruppo); err != nil {
			return nil, fmt.Errorf("errore durante la lettura delle conversazioni private: %w", err)
		}
		if foto.Valid && foto.String != "" {
			conv.Foto = &foto.String
		} else {
			conv.Foto = nil
		}
		if conv.Ultimo != nil && len(*conv.Ultimo) > 15 {
			*conv.Ultimo = (*conv.Ultimo)[:15]
		}
		conversazioni = append(conversazioni, conv)
		idsConversazioniPrivate = append(idsConversazioniPrivate, conv.Id)
	}

	rowsGruppo, err := db.c.Query(queryConversazioniGruppo, utente_Passato)
	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero delle conversazioni di gruppo: %w", err)
	}
	defer rowsGruppo.Close()

	var idsConversazioniGruppo []int
	for rowsGruppo.Next() {
		var conv Conversazione
		var foto sql.NullString
		if err := rowsGruppo.Scan(&conv.Id, &conv.Nome, &conv.Time, &conv.Ultimo, &foto, &conv.IsGruppo); err != nil {
			return nil, fmt.Errorf("errore durante la lettura delle conversazioni di gruppo: %w", err)
		}
		if foto.Valid && foto.String != "" {
			conv.Foto = &foto.String
		} else {
			conv.Foto = nil
		}
		if conv.Ultimo != nil && len(*conv.Ultimo) > 15 {
			*conv.Ultimo = (*conv.Ultimo)[:15]
		}
		conversazioni = append(conversazioni, conv)
		idsConversazioniGruppo = append(idsConversazioniGruppo, conv.Id)
	}

	for _, idConv := range idsConversazioniPrivate {
		if err := db.segnaMessaggiPrivatiRicevuti(utente_Passato_string, idConv); err != nil {
			return nil, fmt.Errorf("errore durante la segnalazione dei messaggi privati come ricevuti: %w", err)
		}
	}

	for _, idConv := range idsConversazioniGruppo {
		if err := db.segnaMessaggiGruppoRicevuti(utente_Passato_string, idConv); err != nil {
			return nil, fmt.Errorf("errore durante la segnalazione dei messaggi di gruppo come ricevuti: %w", err)
		}
		err = db.CheckRicevimentoMessaggiGruppo(idConv)
		if err != nil {
			return nil, fmt.Errorf("errore durante il check di lettura dei messaggi: %w", err)
		}
	}

	sort.Slice(conversazioni, func(i, j int) bool {
		if conversazioni[i].Time == nil {
			return false
		}
		if conversazioni[j].Time == nil {
			return true
		}
		return *conversazioni[i].Time > *conversazioni[j].Time
	})

	return conversazioni, nil
}
