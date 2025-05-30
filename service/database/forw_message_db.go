package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (db *appdbimpl) CopyMessageChangingTimeAndSender(messageId int, userPassed string, chat_Passed int) (int, error) {
	userPassed_converted, _, err := db.IDFromNICK(userPassed)
	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("error converting nickname to ID: %w", err)
	}

	var originalMessage struct {
		Text  *string
		Photo *int
	}

	querySelect :=
		`SELECT text, photo FROM message WHERE id = ?;`

	err = db.c.QueryRow(querySelect, messageId).Scan(
		&originalMessage.Text,
		&originalMessage.Photo,
	)

	if !errors.Is(err, nil) {
		if errors.Is(err, sql.ErrNoRows) {
			return 404, fmt.Errorf("no message found with ID %d", messageId)
		}
		return 500, fmt.Errorf("error retrieving original message with ID %d: %w", messageId, err)
	}

	insertQuery := `INSERT INTO message (text, photo, chat, author, idauthor, time, forw, ans) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	result, err := db.c.Exec(insertQuery, originalMessage.Text, originalMessage.Photo, chat_Passed, userPassed, userPassed_converted, time.Now(), true, nil)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error inserting modified message: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error retrieving last inserted message ID: %w", err)
	}

	isGruppo, errorCode, err := db.SearchGroup(chat_Passed)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error checking chat type: %w", err)
	}
	if isGruppo > 0 {
		err = db.CreateGroupMessageStatus(int(lastInsertID))
		if !errors.Is(err, nil) {
			return 500, fmt.Errorf("error creating group message status: %w", err)
		}
		return 0, nil
	} else {
		err = db.CreatePrivateMessageStatus(int(lastInsertID))
		if !errors.Is(err, nil) {
			return 500, fmt.Errorf("error creating private message status: %w", err)
		}
		return 0, nil
	}
}

func (db *appdbimpl) ForwardMessage(userPassed string, newChatId int, messageId int) (int, error) {
	userPassed_converted, errorCode, err := db.IDFromNICK(userPassed)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error converting nickname to ID: %w", err)
	}
	isGruppo, errorCode, err := db.SearchGroup(newChatId)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error checking chat type: %w", err)
	}
	if isGruppo > 0 {
		involved, errorCode, err := db.UserInGroup(userPassed, newChatId)
		if !errors.Is(err, nil) {
			return errorCode, fmt.Errorf("error checking user's membership in group: %w", err)
		}
		if involved == 0 {
			return 401, fmt.Errorf("user is not a member of the group")
		}
	}

	ex, err := db.VerifyMessageExistence(messageId)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error verifying existence of message ID %d: %w", messageId, err)
	}
	if !ex {
		return 404, fmt.Errorf("message with ID %d does not exist", messageId)
	}
	chatID, err := db.GetChatIdByMessage(messageId)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error retrieving chat ID: %w", err)
	}
	isGruppo, errorCode, err = db.SearchGroup(chatID)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error checking chat type: %w", err)
	}
	if isGruppo > 0 {
		involved, errorCode, err := db.UserInGroup(userPassed, chatID)
		if !errors.Is(err, nil) {
			return errorCode, fmt.Errorf("error checking user's membership in group: %w", err)
		}
		if involved == 0 {
			return 401, fmt.Errorf("user is not a member of the group")
		}
	} else {
		PrivateChatId, errorCode, err := db.SearchPrivateConversation(chatID, userPassed_converted)
		if !errors.Is(err, nil) {
			return errorCode, fmt.Errorf("error checking private chat: %w", err)
		}
		if PrivateChatId == 0 {
			return errorCode, fmt.Errorf("user is not involved in private chat")
		}
	}
	errorCode, err = db.CopyMessageChangingTimeAndSender(messageId, userPassed, newChatId)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error copying message: %w", err)
	}

	return 0, nil
}

func (db *appdbimpl) ForwardMessageToNewChat(userPassed string, user2 string, messageId int) (int, error) {
	chat, errorCode, err := db.chatExistenceTraUtenti(userPassed, user2)
	if chat > 0 {
		return 404, fmt.Errorf("chat already exists between the two users, message ID %d not forwarded", messageId)
	}
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("unexpected error: %w", err)
	}
	userPassed_converted, errorCode, err := db.IDFromNICK(userPassed)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error converting nickname to ID: %w", err)
	}
	user2_converted, errorCode, err := db.IDFromNICK(user2)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error converting nickname to ID: %w", err)
	}
	ex, err := db.UserExistence(userPassed)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error checking user existence for message ID %d: %w", messageId, err)
	}
	if !ex {
		return 404, fmt.Errorf("user with ID %d does not exist: %w", userPassed_converted, err)
	}
	ex, err = db.UserExistence(user2)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error checking user existence for message ID %d: %w", messageId, err)
	}
	if !ex {
		return 404, fmt.Errorf("user with ID %d does not exist: %w", user2_converted, err)
	}
	ex, err = db.VerifyMessageExistence(messageId)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error checking message existence for ID %d: %w", messageId, err)
	}
	if !ex {
		return 404, fmt.Errorf("message with ID %d does not exist: %w", messageId, err)
	}
	chatID, err := db.GetChatIdByMessage(messageId)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error retrieving chat ID: %w", err)
	}
	isGruppo, errorCode, err := db.SearchGroup(chatID)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error checking chat type: %w", err)
	}
	if isGruppo > 0 {
		involved, errorCode, err := db.UserInGroup(userPassed, chatID)
		if !errors.Is(err, nil) {
			return errorCode, fmt.Errorf("error checking user's membership in group: %w", err)
		}
		if involved == 0 {
			return 401, fmt.Errorf("user is not a member of the group")
		}
	} else {
		PrivateChatId, errorCode, err := db.SearchPrivateConversation(chatID, userPassed_converted)
		if !errors.Is(err, nil) {
			return errorCode, fmt.Errorf("error checking private chat: %w", err)
		}
		if PrivateChatId == 0 {
			return errorCode, fmt.Errorf("user is not involved in private chat")
		}
	}
	newChat, errorCode, err := db.CreatePrivateChatDB(userPassed, user2)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error creating new private chat: %w", err)
	}
	errorCode, err = db.CopyMessageChangingTimeAndSender(messageId, userPassed, newChat)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error copying message: %w", err)
	}

	return 0, nil
}
