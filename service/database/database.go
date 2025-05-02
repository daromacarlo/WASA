/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
)

// AppDatabase è una interfaccia ad alto livello del database.
type AppDatabase interface {
	// impostazioni
	ImpostaFotoProfilo(nicknamePassato string, idfotoPassata int) error
	ImpostaNome(nicknamePassato string, nuovoNickPassato string) (int, error)
	ImpostaFotoGruppo(utenteChiamante string, idfotoPassata int, idgruppoPassato int) (int, error)
	ImpostaNomeGruppo(utenteChiamante string, nomeGruppoPassato string, idgruppoPassato int) (int, error)

	// utente
	CreaUtente(nicknamePassato string, idfotoPassata int) error
	IdUtenteDaNickname(nicknamePassato string) (int, int, error)
	NicknameUtenteDaId(idPassato int) (string, int, error)
	UsersInGroup(nicknamePassato string, chat int) ([]Profilo, int, error)
	Login(nicknamePassato string) (string, error)

	// foto
	CreaFoto(fotoPassata string) (int, error)

	// messaggi
	CreaMessaggioFotoDB(utentePassato string, conversazionePassata int, fotoPassata int) (int, int, error)
	CreaMessaggioTestualeDB(utentePassato string, conversazionePassata int, testoPassato string) (int, int, error)
	CreaStatoMessaggioPrivato(idmessaggioPassato int) error
	LeggiMessaggiPrivati(utente2Passato string, conversazioneID int) error
	CreaStatoMessaggioGruppo(idmessaggioPassato int) error
	LeggiMessaggiGruppo(utente1Passato string, conversazioneID int) error
	CheckLetturaMessaggiGruppo(conversazioneID int) error
	EliminaMessaggio(utentePassato string, id_messaggio int, idchat int) error
	CopiaMessaggioCambiandoOraEMitente(idMessaggio int, nuovoAutore string, chat int) (int, error)
	InoltraMessaggio(utentePassato string, idChatNuova int, IdMessaggio int) (int, error)
	ImpostaRisposta(IdMessaggio int, IdNuovoMessaggio int) error
	RispondiMessaggioFoto(utentePassato string, idGruppoPassato int, IdMessaggio int, fotoPassato int) (int, error)
	RispondiMessaggioTesto(utentePassato string, idGruppoPassato int, IdMessaggio int, testoPassato string) (int, error)
	InoltraMessaggioANuovaChat(utente_Passato string, utente2_Passato string, IdMessaggio int) (int, error)

	// commenti
	EliminaCommento(utentePassato string, idcommento int) error
	AggiungiCommento(utentePassato string, messaggioPassato int, reazionePassata string) (int, error)

	// conversazione
	CreaConversazioneDB() (int, error)
	CreaGruppoDB(utenteChiamante string, nomeGruppoPassato string, idfotoPassata int) (int, error)
	CreaConversazionePrivataDB(utente1_Passato string, utente2_Passato string) (int, int, error)
	AggiungiAGruppoDB(idConversazione int, utenteChiamante string, utenteDaAggiungere string) (int, error)
	LasciaGruppo(idConversazione int, utenteChiamante string) (int, error)
	GetConversazioni(utentePassato string) ([]Conversazione, int, error)
	GetConversazione(utentePassato string, conversazionePassata int) ([]MessageData, int, error)
	GetConversazionePrivata(utente1Passato string, utente2Passato string) ([]MessageData, int, error)
	GetConversazioneGruppo(utente1Passato string, idConversazione int) ([]MessageData, int, error)

	// check
	EsisteConversazione(idConversazione int) (bool, error)
	EsisteConversazioneTraUtenti(utente1Passato string, utente2Passato string) (int, int, error)
	UtenteCoinvoltoPrivato(utentePassato string, destinatarioPassato string) (int, int, error)
	CercaConversazioneGruppo(conversazionePassata int) (int, int, error)
	CercaConversazionePrivata(conversazioneID int, utentePassatoConvertito int) (int, int, error)
	UtenteCoinvoltoGruppo(utentePassato string, conversazionePassata int) (int, int, error)
	EsistenzaUtente(nicknamePassato string) (bool, error)

	// test
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building an AppDatabase")
	}

	if err := CreaTabellaUtente(db); err != nil {
		return nil, err
	}
	if err := CreaTabellaFoto(db); err != nil {
		return nil, err
	}
	if err := CreaTabellaConversazione(db); err != nil {
		return nil, err
	}
	if err := CreaTabellaGruppo(db); err != nil {
		return nil, err
	}
	if err := CreaTabellaConversazionePrivata(db); err != nil {
		return nil, err
	}
	if err := CreaTabellaMessaggio(db); err != nil {
		return nil, err
	}
	if err := CreaTabellaCommento(db); err != nil {
		return nil, err
	}
	if err := CreaTabellaUtenteingruppo(db); err != nil {
		return nil, err
	}
	if err := CreaTabellaStatoMessaggioPrivato(db); err != nil {
		return nil, err
	}
	if err := CreaTabellaStatoMessaggioGruppo(db); err != nil {
		return nil, err
	}
	if err := CreaTabellaStatoMessaggioGruppoPersona(db); err != nil {
		return nil, err
	}
	if err := CreaTabellaStatoMessaggioGruppoPersonaRicevimento(db); err != nil {
		return nil, err
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
