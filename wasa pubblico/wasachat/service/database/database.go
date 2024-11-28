package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase è l'interfaccia di alto livello per il DB
type AppDatabase interface {
	//preesistenti
	GetName() (string, error)
	SetName(name string) error

	//creazione tabelle
	CreateUser(nick string) error
	StartConversazionePrivata(id_utente1 int, id_utente2 int) (interoNullabile, error)
	CreateMessaggioPrivato(nickmittente string, nickdestinatario string, messaggio string) error

	//ricerca
	SearchUser(nick string) (bool, error)
	Get_IdFromNick_Persona(nick string) (interoNullabile, error)
	Get_NickFromId_Persona(id int) (stringaNullabile, error)
	SearchConversazionePrivata(id_utente1 int, id_utente2 int) (interoNullabile, error)
	GetChats(nickname string) ([]Chat, error)

	//modifica
	ChangeNickname(newnickname string, nickname string) error
	ChangePhoto(photo []byte, nickname string) error

	//test
	GetUsers() ([]User, error)
	GetMessaggi() ([]Message, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New restituisce una nuova istanza di AppDatabase basata sulla connessione SQLite `db`.
// `db` è richiesto - verrà restituito un errore se `db` è `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("il database è richiesto per costruire un AppDatabase")
	}

	// Crea la tabella persone se non esiste
	err1 := CreateTablePersone(db)
	if err1 != nil {
		return nil, fmt.Errorf("errore durante la creazione della tabella persone: %w", err1)
	}

	err2 := CreateTableMessaggi(db)
	if err2 != nil {
		return nil, fmt.Errorf("errore durante la creazione della tabella messaggi: %w", err2)
	}

	err3 := CreateTableConversazioni(db)
	if err3 != nil {
		return nil, fmt.Errorf("errore durante la creazione della tabella ConversazioniPrivate: %w", err3)
	}
	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
