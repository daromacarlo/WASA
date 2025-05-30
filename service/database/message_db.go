package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func CreateTableMessages(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS message(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author TEXT NOT NULL,
			idauthor INTEGER NOT NULL,
			chat INTEGER NOT NULL,
			forw BOOL NOT NULL default FALSE,
			ans INTEGER default NULL,
			-- One of the two
			text TEXT,
			photo INTEGER,
			--
			time TIME,
			FOREIGN KEY (photo) REFERENCES photo(id),
			FOREIGN KEY (author) REFERENCES user(nickname),
			FOREIGN KEY (idauthor) REFERENCES user(id),
			FOREIGN KEY (chat) REFERENCES chat(id),
			FOREIGN KEY (ans) REFERENCES message(id)
		);`
	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error during message table creation: %w", err)
	}
	return nil
}

func CreateTableComment(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS comment(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author TEXT NOT NULL,
			idauthor INTEGER NOT NULL,
			message INTEGER NOT NULL,
			reaction TEXT NOT NULL,	
			FOREIGN KEY (author) REFERENCES user(nickname),
			FOREIGN KEY (idauthor) REFERENCES user(id),
			FOREIGN KEY (message) REFERENCES message(id)
		);`
	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error during comment table creation: %w", err)
	}
	return nil
}

func (db *appdbimpl) CreateTextMessageDB(userPassed string, chatId int, passedText string) (int, int, error) {
	ex, err := db.chatExistence(chatId)
	if !errors.Is(err, nil) {
		return 0, 500, fmt.Errorf("error during existence check: %w", err)
	}
	if !ex {
		return 0, 404, fmt.Errorf("error, chat does not exist")
	}
	userPassed_converted, errorCode, err := db.IDFromNICK(userPassed)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error converting nickname to id: %w", err)
	}
	isGroup, errorCode, err := db.SearchGroup(chatId)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error during chat type check: %w", err)
	}
	var messageId int
	if isGroup > 0 {
		involved, errorCode, err := db.UserInGroup(userPassed, chatId)
		if !errors.Is(err, nil) {
			return 0, errorCode, fmt.Errorf("error checking user participation in group: %w", err)
		}
		if involved == 0 {
			return 0, 401, fmt.Errorf("user is not a member of the group")
		}
		messageId, err = db.insertMessage(chatId, userPassed, passedText, true)
		if !errors.Is(err, nil) {
			return 0, 0, err
		}

	} else {
		PrivateChatId, errorCode, err := db.SearchPrivateConversation(chatId, userPassed_converted)
		if !errors.Is(err, nil) {
			return 0, errorCode, fmt.Errorf("error during chat verification: %w", err)
		}
		if PrivateChatId == 0 {
			return 0, 401, fmt.Errorf("user is not involved in the private chat")
		}
		messageId, err = db.insertMessage(chatId, userPassed, passedText, false)
		if !errors.Is(err, nil) {
			return 0, 500, err
		}
	}

	return messageId, 0, nil
}

func (db *appdbimpl) insertMessage(chatId int, userPassed string, passedText string, isGroup bool) (int, error) {
	userPassed_converted, _, err := db.IDFromNICK(userPassed)
	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("error converting nickname to id: %w", err)
	}
	insertQuery := `INSERT INTO message (author, idauthor, chat, text, time) VALUES (?, ?, ?, ?, ?);`
	result, err := db.c.Exec(insertQuery, userPassed, userPassed_converted, chatId, passedText, time.Now())

	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("error creating new message: %w", err)
	}

	lastInsertID, err := result.LastInsertId()
	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("error retrieving last message ID: %w", err)
	}

	if isGroup {
		err = db.CreateGroupMessageStatus(int(lastInsertID))
	} else {
		err = db.CreatePrivateMessageStatus(int(lastInsertID))
	}

	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("error creating message status: %w", err)
	}

	return int(lastInsertID), nil
}

func (db *appdbimpl) CreatePhotoMessageDB(userPassed string, chatId int, passedPhoto int) (int, int, error) {
	ex, err := db.chatExistence(chatId)
	if !errors.Is(err, nil) {
		return 0, 500, fmt.Errorf("error during existence check: %w", err)
	}
	if !ex {
		return 0, 404, fmt.Errorf("error, chat does not exist")
	}
	userPassed_converted, errorCode, err := db.IDFromNICK(userPassed)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error converting nickname to id: %w", err)
	}
	isGroup, errorCode, err := db.SearchGroup(chatId)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error during chat type check: %w", err)
	}

	var messageId int
	if isGroup > 0 {
		involved, errorCode, err := db.UserInGroup(userPassed, chatId)
		if !errors.Is(err, nil) {
			return 0, errorCode, fmt.Errorf("error checking user participation in group: %w", err)
		}
		if involved == 0 {
			return 0, 401, fmt.Errorf("user is not a member of the group")
		}
		messageId, err = db.insertMessagePhoto(chatId, userPassed, passedPhoto, true)
		if !errors.Is(err, nil) {
			return 0, 500, fmt.Errorf("error inserting message: %w", err)
		}

	} else {
		PrivateChatId, errorCode, err := db.SearchPrivateConversation(chatId, userPassed_converted)
		if !errors.Is(err, nil) {
			return 0, errorCode, fmt.Errorf("error during chat verification: %w", err)
		}
		if PrivateChatId == 0 {
			return 0, 401, fmt.Errorf("user is not involved in the private chat")
		}
		messageId, err = db.insertMessagePhoto(chatId, userPassed, passedPhoto, false)
		if !errors.Is(err, nil) {
			return 0, 500, fmt.Errorf("error inserting message: %w", err)
		}
	}

	return messageId, 0, nil
}

func (db *appdbimpl) insertMessagePhoto(chatId int, userPassed string, passedPhoto int, isGroup bool) (int, error) {
	userPassed_converted, _, err := db.IDFromNICK(userPassed)
	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("error converting nickname to id: %w", err)
	}
	insertQuery := `INSERT INTO message (author, idauthor, chat, photo, time) VALUES (?, ?, ?, ?, ?);`
	result, err := db.c.Exec(insertQuery, userPassed, userPassed_converted, chatId, passedPhoto, time.Now())

	if !errors.Is(err, nil) {
		return 0, fmt.Errorf(": %w", err)
	}

	lastInsertID, err := result.LastInsertId()
	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("error retrieving last message ID: %w", err)
	}

	if isGroup {
		err = db.CreateGroupMessageStatus(int(lastInsertID))
	} else {
		err = db.CreatePrivateMessageStatus(int(lastInsertID))
	}

	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("error creating message status: %w", err)
	}

	return int(lastInsertID), nil
}

func (db *appdbimpl) DeleteMessage(userPassed string, idMessage int, idChat int) error {
	userPassed_converted, _, err := db.IDFromNICK(userPassed)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error converting nickname to id: %w", err)
	}
	DeleteQuery := `DELETE FROM message WHERE idauthor = ? AND id = ? AND chat = ?;`
	result, err := db.c.Exec(DeleteQuery, userPassed_converted, idMessage, idChat)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error deleting message: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if !errors.Is(err, nil) {
		return fmt.Errorf("error checking affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no message found matching the specified criteria")
	}

	return nil
}

func (db *appdbimpl) AddComment(userPassed string, messageId int, reaction string) (int, error) {
	userPassed_converted, errorCode, err := db.IDFromNICK(userPassed)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error converting nickname to ID: %w", err)
	}

	var chatId int
	checkQuery := `SELECT chat FROM message WHERE id = ?;`
	err = db.c.QueryRow(checkQuery, messageId).Scan(&chatId)
	if !errors.Is(err, nil) {
		if errors.Is(err, sql.ErrNoRows) {
			return 404, fmt.Errorf("message not found")
		}
		return 500, fmt.Errorf("error checking message: %w", err)
	}

	isGroup, errorCode, err := db.SearchGroup(chatId)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error during chat type check: %w", err)
	}

	if isGroup > 0 {
		involved, errorCode, err := db.UserInGroup(userPassed, chatId)
		if !errors.Is(err, nil) {
			return errorCode, fmt.Errorf("error checking user participation in group: %w", err)
		}
		if involved == 0 {
			return 401, fmt.Errorf("user is not a member of the group")
		}
	} else {
		PrivateChatId, errorCode, err := db.SearchPrivateConversation(chatId, userPassed_converted)
		if !errors.Is(err, nil) {
			return errorCode, fmt.Errorf("error checking user participation: %w", err)
		}
		if PrivateChatId == 0 {
			return 401, fmt.Errorf("user is not involved in the private chat")
		}
	}

	verifyCommentQuery := `SELECT id FROM comment WHERE idauthor = ? AND message = ?;`
	var commentId int
	err = db.c.QueryRow(verifyCommentQuery, userPassed_converted, messageId).Scan(&commentId)

	switch {
	case err == nil:
		updateQuery := `UPDATE comment SET reaction = ? WHERE id = ?;`
		_, err = db.c.Exec(updateQuery, reaction, commentId)
		if err != nil {
			return 500, fmt.Errorf("error updating comment: %w", err)
		}
	case errors.Is(err, sql.ErrNoRows):
		insertQuery := `INSERT INTO comment (author, idauthor, message, reaction) VALUES (?, ?, ?, ?);`
		_, err = db.c.Exec(insertQuery, userPassed, userPassed_converted, messageId, reaction)
		if err != nil {
			return 500, fmt.Errorf("error inserting comment: %w", err)
		}
	default:
		return 500, fmt.Errorf("error checking comment: %w", err)
	}

	return 0, nil
}

// DeleteComment deletes a specific comment of a user given the message ID
func (db *appdbimpl) DeleteComment(user string, messageId int) error {
	// Check if a comment exists from the specific user for this message
	userPassed_converted, _, err := db.IDFromNICK(user)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error converting nickname to id: %w", err)
	}

	var commentId int

	verifyCommentQuery := `
        SELECT id 
        FROM comment
        WHERE message = ? AND idauthor = ?;
    `

	err = db.c.QueryRow(verifyCommentQuery, messageId, userPassed_converted).Scan(&commentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no comment found for this user and message")
		}
		return fmt.Errorf("error searching for comment: %w", err)
	}
	deleteQuery := `
        DELETE FROM comment
        WHERE id = ?;
    `

	_, err = db.c.Exec(deleteQuery, commentId)
	if err != nil {
		return fmt.Errorf("error deleting comment: %w", err)
	}

	return nil
}
