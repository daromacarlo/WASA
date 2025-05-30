package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// CreateTablePrivateMessageStatus creates the table for private message status
func CreateTablePrivateMessageStatus(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS private_state (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message INTEGER NOT NULL,
		rec BOOL DEFAULT FALSE NOT NULL,
		read BOOL DEFAULT FALSE NOT NULL,
		FOREIGN KEY (message) REFERENCES message(id)
	)`
	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error creating private_state table: %w", err)
	}
	return nil
}

// CreatePrivateMessageStatus inserts a new private message status entry
func (db *appdbimpl) CreatePrivateMessageStatus(messageID int) error {
	_, err := db.c.Exec(
		"INSERT INTO private_state (message, rec, read) VALUES (?, ?, ?)",
		messageID, false, false,
	)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error creating private message status: %w", err)
	}
	return nil
}

// ReadPrivateMessage marks private messages as read and received
func (db *appdbimpl) ReadPrivateMessage(recipient string, chatID int) error {
	recipientID, _, err := db.IDFromNICK(recipient)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error converting nickname to ID: %w", err)
	}

	_, err = db.c.Exec(`
		UPDATE private_state 
		SET read = 1, rec = 1 
		WHERE message IN (
			SELECT id FROM message 
			WHERE chat = ? AND idauthor = ?
		)`,
		chatID, recipientID,
	)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error updating private message status: %w", err)
	}
	return nil
}

// CreateTableMessagesStatusGroup creates the table for group message status
func CreateTableMessagesStatusGroup(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS group_state (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message INTEGER NOT NULL,
		rec BOOL DEFAULT FALSE NOT NULL,
		read BOOL DEFAULT FALSE NOT NULL,
		FOREIGN KEY (message) REFERENCES message(id)
	)`
	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error creating group_state table: %w", err)
	}
	return nil
}

// CreateTableMessagesStatusGroupPersona creates the table for per-user group message read tracking
func CreateTableMessagesStatusGroupPersona(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS group_state_user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message INTEGER NOT NULL,
		user TEXT NOT NULL,
		iduser INTEGER NOT NULL,
		UNIQUE(message, iduser),
		FOREIGN KEY (message) REFERENCES message(id),
		FOREIGN KEY (user) REFERENCES user(nickname),
		FOREIGN KEY (iduser) REFERENCES user(id)
	)`
	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error creating group_state_user table: %w", err)
	}
	return nil
}

// CreateGroupMessageReceiptStatusTable creates the table for group message receipt tracking
func CreateGroupMessageReceiptStatusTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS group_state_user_rec (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message INTEGER NOT NULL,
		user TEXT NOT NULL,
		iduser INTEGER NOT NULL,
		UNIQUE(message, iduser),
		FOREIGN KEY (message) REFERENCES message(id),
		FOREIGN KEY (user) REFERENCES user(nickname),
		FOREIGN KEY (iduser) REFERENCES user(id)
	)`
	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error creating group_state_user_rec table: %w", err)
	}
	return nil
}

// CreateGroupMessageStatus inserts a new group message status entry
func (db *appdbimpl) CreateGroupMessageStatus(messageID int) error {
	_, err := db.c.Exec(
		"INSERT INTO group_state (message, rec, read) VALUES (?, ?, ?)",
		messageID, false, false,
	)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error creating group message status: %w", err)
	}
	return nil
}

// ReadGroupMessages marks group messages as read by the user
func (db *appdbimpl) ReadGroupMessages(user string, chatID int) error {
	// Insert unread messages into group_state_user for tracking
	userID, _, err := db.IDFromNICK(user)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error converting nickname to ID: %w", err)
	}
	_, err = db.c.Exec(`
		INSERT INTO group_state_user (message, user, iduser)
		SELECT m.id, ?, ? 
		FROM message m
		JOIN group_state smg ON smg.message = m.id
		JOIN user u ON m.idauthor = u.id
		WHERE m.chat = ? AND smg.read = false AND u.id != ?
		ON CONFLICT DO NOTHING`,
		user, userID, chatID, userID,
	)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error updating group message read status: %w", err)
	}
	return nil
}

// CheckGroupMessageReadStatus verifies and updates read status for group messages
func (db *appdbimpl) CheckGroupMessageReadStatus(chatID int) error {
	_, err := db.c.Exec(`
		UPDATE group_state
		SET read = 1, rec = 1
		WHERE message IN (
			SELECT m.id
			FROM message m
			WHERE m.chat = ? AND (
				SELECT COUNT(*) 
				FROM group_state_user smgp
				WHERE smgp.message = m.id
			) >= (
				SELECT COUNT(*)
				FROM user_in_group uig
				JOIN groups g ON uig.groups = g.id
				WHERE g.chat = ?
			) - 1
		)`,
		chatID, chatID,
	)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error verifying group message read status: %w", err)
	}
	return nil
}

// CheckGroupMessageReceipt verifies and updates receipt status for group messages
func (db *appdbimpl) CheckGroupMessageReceipt(chatID int) error {
	_, err := db.c.Exec(`
		UPDATE group_state
		SET rec = 1
		WHERE message IN (
			SELECT m.id
			FROM message m
			WHERE m.chat = ? AND (
				SELECT COUNT(*)
				FROM group_state_user_rec smgp
				WHERE smgp.message = m.id
			) >= (
				SELECT COUNT(*)
				FROM user_in_group uig
				JOIN groups g ON uig.groups = g.id
				WHERE g.chat = ?
			) - 1
		)`,
		chatID, chatID,
	)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error verifying group message receipt status: %w", err)
	}
	return nil
}

// MarkPrivateMessagesAsReceived marks private messages as received
func (db *appdbimpl) MarkPrivateMessagesAsReceived(user string, chatID int) error {
	userID, _, err := db.IDFromNICK(user)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error converting nickname to ID: %w", err)
	}
	_, err = db.c.Exec(`
		UPDATE private_state
		SET rec = 1
		WHERE message IN (
			SELECT id FROM message
			WHERE chat = ? AND idauthor != ?
		)`,
		chatID, userID,
	)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error updating private message receipt status: %w", err)
	}
	return nil
}

// MarkGroupMessagesAsReceived marks group messages as received by the user
func (db *appdbimpl) MarkGroupMessagesAsReceived(user string, chatID int) error {
	userID, _, err := db.IDFromNICK(user)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error converting nickname to ID: %w", err)
	}

	_, err = db.c.Exec(`
		INSERT INTO group_state_user_rec (message, user, iduser)
		SELECT m.id, ?, ?
		FROM message m
		JOIN group_state smg ON smg.message = m.id
		JOIN user u ON m.author = u.nickname
		WHERE m.chat = ? AND smg.rec = false AND u.id != ?
		ON CONFLICT DO NOTHING`,
		user, userID, chatID, userID,
	)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error updating group message receipt status: %w", err)
	}
	return nil
}
