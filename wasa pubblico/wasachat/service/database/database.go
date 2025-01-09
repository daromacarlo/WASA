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

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	//impostazioni
	ImpostaFotoProfilo(nicknamePassato string, idfotoPassata int) error
	ImpostaNome(nicknamePassato string, nuovoNickPassato string) error
	ImpostaFotoGruppo(UtenteChiamante string, id_foto_Passata int, id_gruppo_passato int) error
	ImpostaNomeGruppo(UtenteChiamante string, nomeGruppo_Passato string, id_gruppo_passato int) error

	//utente
	CreaUtente(nicknamePassato string, idfotoPassata int) error
	IdUtenteDaNickname(nickname_Passato string) (int, error)
	NicknameUtenteDaId(id_passato int) (string, error)
	VediProfili(nickname_Passato string) ([]Profilo, error)

	//foto
	CreaFoto(percorso_Passato string, foto_Passata []byte) (int, error)

	//messaggi
	CreaMessaggioFotoDBPrivato(utente_Passato string, destinatario_Passato string, foto_Passata int) error
	CreaMessaggioTestualeDBPrivato(utente_Passato string, destinatario_Passato string, testo_Passato string) error
	CreaMessaggioFotoDBGruppo(utente_Passato string, conversazione_Passata int, foto_Passata int) error
	CreaMessaggioTestualeDBGruppo(utente_Passato string, conversazione_Passata int, testo_Passato string) error
	CreaStatoMessaggioPrivato(id_messaggio_Passato int) error
	LeggiMessaggiPrivati(utente1_Passato string, utente2_Passato string, conversazioneID int) error
	CreaStatoMessaggioGruppo(id_messaggio_Passato int) error
	LeggiMessaggiGruppo(utente1_Passato string, conversazioneID int) error
	CheckLetturaMessaggiGruppo(conversazioneID int) error
	EliminaMessaggio(utente_Passato string, id_messaggio int, id_chat int) error

	//commenti
	EliminaCommento(utente_Passato string, id_commento int) error
	AggiungiCommento(utente_Passato string, messaggio_Passato int, reazione_Passata string) error

	//conversazione
	CreaConversazioneDB() (int, error)
	CreaGruppoDB(UtenteChiamante string, nomeGruppo_Passato string, idfoto_Passata int) error
	CreaConversazionePrivataDB(utente1_Passato string, utente2_Passato string) (int, error)
	AggiungiAGruppoDB(idConversazione int, UtenteChiamante string, UtenteDaAggiungere string) error
	GetConversazionePrivata(utente1_Passato string, utente2_Passato string) ([]MessageData, error)
	GetConversazioneGruppo(utente1_Passato string, id_conversazione int) ([]MessageData, error)
	LasciaGruppo(idConversazione int, UtenteChiamante string) error
	GetConversazioni(utente_Passato string) ([]Conversazione, error)

	//check
	EsisteConversazione(idConversazione int) (bool, error)
	EsisteConversazioneTraUtenti(utente1_Passato string, utente2_Passato string) (int, error)
	UtenteCoinvoltoPrivato(utente_Passato string, destinatario_Passato string) (int, error)
	CercaConversazioneGruppo(conversazione_Passata int) (int, error)
	CercaConversazionePrivata(conversazioneID int, utente_Passato_convertito int) (int, error)
	UtenteCoinvoltoGruppo(utente_Passato string, conversazione_Passata int) (int, error)
	EsistenzaUtente(nickname_Passato string) (bool, error)

	//cosa manca?
	/*
		inoltra messaggio
	*/

	//test
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Crea la tabella persone se non esiste (guarda file utente.go per maggiori informazioni)
	err := CreaTabellaUtente(db)
	if err != nil {
		return nil, err
	}

	err = CreaTabellaFoto(db)
	if err != nil {
		return nil, err
	}

	err = CreaTabellaConversazione(db)
	if err != nil {
		return nil, err
	}

	err = CreaTabellaGruppo(db)
	if err != nil {
		return nil, err
	}

	err = CreaTabellaConversazionePrivata(db)
	if err != nil {
		return nil, err
	}

	err = CreaTabellaMessaggio(db)
	if err != nil {
		return nil, err
	}

	err = CreaTabellaCommento(db)
	if err != nil {
		return nil, err
	}

	err = CreaTabellaUtenteingruppo(db)
	if err != nil {
		return nil, err
	}

	err = CreaTabellaStatoMessaggioPrivato(db)
	if err != nil {
		return nil, err
	}

	err = CreaTabellaStatoMessaggioGruppo(db)
	if err != nil {
		return nil, err
	}

	err = CreaTabellaStatoMessaggioGruppoPersona(db)
	if err != nil {
		return nil, err
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
