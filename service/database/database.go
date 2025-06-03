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

type AppDatabase interface {

	// settings
	SetPhotoDB(nickname string, idPhoto int) error
	SetNameDB(nickname string, newNickname string) (int, error)
	SetGroupPhotoDB(callerUser string, idPhoto int, groupId int) (int, error)
	SetGroupNameDB(callerUser string, newGroupName string, groupId int) (int, error)

	// user
	CreateUser(nickname string, idPhoto int) error
	IDFromNICK(nickname string) (int, int, error)
	NICKFromID(idPassed int) (string, int, error)
	UsersInGroup(nickname string, chat int) ([]Profile, int, error)
	Login(nickname string) (string, error)

	// photo
	CreaFoto(passedPhoto string) (int, error)

	// messages
	CreatePhotoMessageDB(userPassed string, chatId int, passedPhoto int) (int, int, error)
	CreateTextMessageDB(userPassed string, chatId int, passedText string) (int, int, error)
	CreatePrivateMessageStatus(messageId int) error
	ReadPrivateMessage(user2 string, chatID int) error
	CreateGroupMessageStatus(messageId int) error
	ReadGroupMessages(user1 string, chatID int) error
	CheckGroupMessageReadStatus(chatID int) error
	DeleteMessage(userPassed string, messageId int, idchat int) error
	CopyMessageChangingTimeAndSender(messageId int, newAuthor string, chat int) (int, error)
	ForwardMessage(userPassed string, newChatId int, messageId int) (int, error)
	setAns(messageId int, newMessageId int) error
	AnsPhotoMessage(userPassed string, groupId int, messageId int, passedPhoto int) (int, error)
	AnsTextMessage(userPassed string, groupId int, messageId int, passedText string) (int, error)
	ForwardMessageToNewChat(userPassed string, user2 string, messageId int) (int, error)

	// comments
	DeleteComment(userPassed string, idcommento int) error
	AddComment(userPassed string, messageId int, reaction string) (int, error)

	// chat
	CreateConversationDB() (int, error)
	CreateGroupDB(callerUser string, nomeGruppoPassato string, idPhoto int) (int, error)
	CreatePrivateChatDB(user1 string, user2 string) (int, int, error)
	AddToGroupDB(chatId int, callerUser string, targetUser string) (int, error)
	LeaveGroup(chatId int, callerUser string) (int, error)
	GetConversationsDB(userPassed string) ([]Conversation, int, error)
	GetChat(userPassed string, chatId int) ([]MessageData, int, error)
	GetPrivateChat(user1 string, user2 string) ([]MessageData, int, error)
	GetGroupChat(user1 string, chatId int) ([]MessageData, int, error)

	// check
	chatExistence(chatId int) (bool, error)
	chatExistenceTraUtenti(user1 string, user2 string) (int, int, error)
	UserInvolved(userPassed string, targetUser string) (int, int, error)
	SearchGroup(chatId int) (int, int, error)
	SearchPrivateConversation(chatID int, userPassed_converted int) (int, int, error)
	UserInGroup(userPassed string, chatId int) (int, int, error)
	UserExistence(nickname string) (bool, error)
	UserExistenceId(userId int) (bool, error)

	// test
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building an AppDatabase")
	}

	if err := CreateTableUser(db); err != nil {
		return nil, err
	}
	if err := CreaTabellaFoto(db); err != nil {
		return nil, err
	}
	if err := CreateTableChat(db); err != nil {
		return nil, err
	}
	if err := CreateTableGroup(db); err != nil {
		return nil, err
	}
	if err := CreateTablePrivateChat(db); err != nil {
		return nil, err
	}
	if err := CreateTableMessages(db); err != nil {
		return nil, err
	}
	if err := CreateTableComment(db); err != nil {
		return nil, err
	}
	if err := CreateTableUseringruppo(db); err != nil {
		return nil, err
	}
	if err := CreateTablePrivateMessageStatus(db); err != nil {
		return nil, err
	}
	if err := CreateTableMessagesStatusGroup(db); err != nil {
		return nil, err
	}
	if err := CreateTableMessagesStatusGroupPersona(db); err != nil {
		return nil, err
	}
	if err := CreateGroupMessageReceiptStatusTable(db); err != nil {
		return nil, err
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
