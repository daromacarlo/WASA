package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) setAns(messageId int, newMessageId int) error {
	query := `
		UPDATE message
		SET ans = ?
		WHERE id = ?
	`
	_, err := db.c.Exec(query, messageId, newMessageId)
	if !errors.Is(err, nil) {
		return fmt.Errorf("unexpected error while setting answer for message with ID %d: %w", newMessageId, err)
	}
	return nil
}

func (db *appdbimpl) verifyConversation(userPassed string, chatId int, messageIdAns int) (int, int, error) {
	ex, err := db.VerifyMessageExistence(messageIdAns)
	if !errors.Is(err, nil) {
		return 0, 500, fmt.Errorf("error while verifying existence of message with ID %d: %w", messageIdAns, err)
	}
	if !ex {
		return 0, 404, fmt.Errorf("error: message with ID %d does not exist", messageIdAns)
	}
	chatID, err := db.GetChatIdByMessage(messageIdAns)
	if !errors.Is(err, nil) {
		return 0, 500, fmt.Errorf("error retrieving chat ID for message ID %d: %w", messageIdAns, err)
	}
	if chatID != chatId {
		return 0, 401, fmt.Errorf("message does not belong to the specified chat (ID: %d)", chatId)
	}
	userID, errorCode, err := db.IDFromNICK(userPassed)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error converting nickname '%s' to user ID: %w", userPassed, err)
	}
	isGroup, errorCode, err := db.SearchGroup(chatID)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error checking chat type for ID %d: %w", chatID, err)
	}
	if isGroup > 0 {
		involved, errorCode, err := db.UserInGroup(userPassed, chatID)
		if !errors.Is(err, nil) {
			return 0, errorCode, fmt.Errorf("error verifying participation of user '%s' in group: %w", userPassed, err)
		}
		if involved == 0 {
			return 0, 401, fmt.Errorf("user '%s' is not a member of the group with ID %d", userPassed, chatID)
		}
	} else {
		idPrivateChat, errorCode, err := db.SearchPrivateConversation(chatID, userID)
		if !errors.Is(err, nil) {
			return 0, errorCode, fmt.Errorf("error checking private chat for user '%s' (ID: %d): %w", userPassed, userID, err)
		}
		if idPrivateChat == 0 {
			return 0, 401, fmt.Errorf("user '%s' is not involved in the private chat with ID %d", userPassed, chatID)
		}
	}

	return chatID, 0, nil
}

func (db *appdbimpl) AnsTextMessage(userPassed string, chatId int, messageId int, passedText string) (int, error) {
	_, errorCode, err := db.verifyConversation(userPassed, chatId, messageId)
	if !errors.Is(err, nil) {
		return errorCode, err
	}
	newMessageId, errorCode, err := db.CreateTextMessageDB(userPassed, chatId, passedText)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error inserting text message for user '%s': %w", userPassed, err)
	}
	err = db.setAns(messageId, newMessageId)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error setting answer for message with ID %d (reply to %d): %w", newMessageId, messageId, err)
	}
	return 0, nil
}

func (db *appdbimpl) AnsPhotoMessage(userPassed string, chatId int, messageId int, passedPhoto int) (int, error) {
	_, errorCode, err := db.verifyConversation(userPassed, chatId, messageId)
	if !errors.Is(err, nil) {
		return errorCode, err
	}
	newMessageId, errorCode, err := db.CreatePhotoMessageDB(userPassed, chatId, passedPhoto)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error inserting photo message for user '%s': %w", userPassed, err)
	}
	err = db.setAns(messageId, newMessageId)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error setting answer for message with ID %d (reply to %d): %w", newMessageId, messageId, err)
	}
	return 0, nil
}

func (db *appdbimpl) AnsPhotoTextMessage(userPassed string, chatId int, messageId int, passedPhoto int, passedString string) (int, error) {
	_, errorCode, err := db.verifyConversation(userPassed, chatId, messageId)
	if !errors.Is(err, nil) {
		return errorCode, err
	}
	newMessageId, errorCode, err := db.CreatePhotoTextMessageDB(userPassed, chatId, passedPhoto, passedString)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error inserting photo and text message for user '%s': %w", userPassed, err)
	}
	err = db.setAns(messageId, newMessageId)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error setting answer for message with ID %d (reply to %d): %w", newMessageId, messageId, err)
	}
	return 0, nil
}
