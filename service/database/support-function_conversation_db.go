package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) GetChatIdByMessage(messageId int) (int, error) {
	query := `
		SELECT chat
		FROM message
		WHERE id = ?;
	`
	var chatID int
	err := db.c.QueryRow(query, messageId).Scan(&chatID)
	if !errors.Is(err, nil) {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("no message found with ID %d", messageId)
		}
		return 0, fmt.Errorf("error retrieving chat ID: %w", err)
	}
	return chatID, nil
}

func (db *appdbimpl) VerifyMessageExistence(messageId int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM message WHERE id = ?;"
	err := db.c.QueryRow(query, messageId).Scan(&count)
	if !errors.Is(err, nil) {
		return false, fmt.Errorf("error verifying existence of message with ID %d: %w", messageId, err)
	}
	return count > 0, nil
}

// Returns true if a chat with the given ID exists, false otherwise
func (db *appdbimpl) chatExistence(chatId int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM chat WHERE id = ?`
	err := db.c.QueryRow(query, chatId).Scan(&count)
	if !errors.Is(err, nil) {
		return false, fmt.Errorf("error verifying chat existence: %w", err)
	}
	return count > 0, nil
}

// Returns the ID of the private chat between two users if it exists, otherwise 0
func (db *appdbimpl) chatExistenceTraUtenti(user1 string, user2 string) (int, int, error) {
	var chatId int

	user1ID, errorCode, err := db.IDFromNICK(user1)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error converting nickname %s to ID for user1: %s", user1, err.Error())
	}
	user2ID, errorCode, err := db.IDFromNICK(user2)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error converting nickname %s to ID for user2: %s", user2, err.Error())
	}

	query := `
		SELECT c.id
		FROM chat c
		LEFT JOIN privconv cp ON cp.chat = c.id
		WHERE (cp.user1 = ? AND cp.user2 = ?) OR (cp.user1 = ? AND cp.user2 = ?);`

	err = db.c.QueryRow(query, user1ID, user2ID, user2ID, user1ID).Scan(&chatId)
	if !errors.Is(err, nil) {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, nil
		}
		return 0, 500, fmt.Errorf("error verifying chat existence: %s", err.Error())
	}
	return chatId, 0, nil
}

// Ensures a private chat exists between two users and returns its ID (creates it if not present)
func (db *appdbimpl) UserInvolved(userPassed string, targetUser string) (int, int, error) {
	ex, errorCode, err := db.chatExistenceTraUtenti(userPassed, targetUser)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error verifying chat existence: %s", err.Error())
	}
	if ex == 0 {
		id, errorCode, err := db.CreatePrivateChatDB(userPassed, targetUser)
		if !errors.Is(err, nil) {
			return 0, errorCode, fmt.Errorf("error creating private chat: %s", err.Error())
		}
		return id, 0, nil
	} else {
		return ex, 0, nil
	}
}

func (db *appdbimpl) GetNameUserInvolved(chatID int, userPassed string) (string, error) {
	queryUserID := `
		SELECT id
		FROM user
		WHERE nickname = ?;`

	var userID int
	err := db.c.QueryRow(queryUserID, userPassed).Scan(&userID)
	if !errors.Is(err, nil) {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user with name %s not found", userPassed)
		}
		return "", fmt.Errorf("error retrieving user ID: %w", err)
	}

	query := `
		SELECT 
			CASE 
				WHEN cp.user1 = ? THEN u2.nickname
				WHEN cp.user2 = ? THEN u1.nickname
			END as user_involved
		FROM chat as c
		JOIN privconv as cp ON cp.chat = c.id
		JOIN user u1 ON cp.user1 = u1.id
		JOIN user u2 ON cp.user2 = u2.id
		WHERE c.id = ? AND (cp.user1 = ? OR cp.user2 = ?);`

	var userInvolved string
	err = db.c.QueryRow(query, userID, userID, chatID, userID, userID).Scan(&userInvolved)
	if !errors.Is(err, nil) {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user is not involved in this chat")
		}
		return "", fmt.Errorf("error retrieving involved user: %w", err)
	}
	return userInvolved, nil
}

func (db *appdbimpl) SearchPrivateConversation(chatID int, userPassed_converted int) (int, int, error) {
	query := `
		SELECT cp.id
		FROM chat as c
		JOIN privconv as cp on cp.chat = c.id
		WHERE c.id = ? AND (cp.user1 = ? OR cp.user2 = ?);`

	var chatPrivateID int
	err := db.c.QueryRow(query, chatID, userPassed_converted, userPassed_converted).Scan(&chatPrivateID)
	if !errors.Is(err, nil) {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 404, fmt.Errorf("chat does not exist or user is not involved")
		}
		return 0, 500, fmt.Errorf("error verifying user involvement: %w", err)
	}
	return chatPrivateID, 0, nil
}

// Checks if a user is in the group related to the given chat ID, returns group ID or 0
func (db *appdbimpl) UserInGroup(userPassed string, chatId int) (int, int, error) {
	idGroup, errorCode, err := db.SearchGroup(chatId)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error searching for group ID: %w", err)
	}
	if idGroup == 0 {
		return 0, 404, fmt.Errorf("no group with these characteristics: %w", err)
	}

	idUser, errorCode, err := db.IDFromNICK(userPassed)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error retrieving user ID: %w", err)
	}

	query := `SELECT 1 FROM user_in_group WHERE groups = ? AND user = ?`
	var result bool
	err = db.c.QueryRow(query, idGroup, idUser).Scan(&result)
	if !errors.Is(err, nil) {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 404, nil
		}
		return 0, 500, fmt.Errorf("error verifying user in group: %w", err)
	}
	if result {
		return idGroup, 0, nil
	}
	return 0, 401, fmt.Errorf("user is not in the group: %w", err)
}

// Checks if there is a group linked to the given chat ID, returns group ID or 0
func (db *appdbimpl) SearchGroup(chatId int) (int, int, error) {
	query := `SELECT id FROM groups WHERE chat = ?`
	var idGroup int

	err := db.c.QueryRow(query, chatId).Scan(&idGroup)
	if !errors.Is(err, nil) {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 404, nil
		}
		return 0, 500, fmt.Errorf("error checking chat existence: %w", err)
	}
	return idGroup, 0, nil
}
