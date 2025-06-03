package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) CreateConversationDB() (int, error) {
	insertQuery := `INSERT INTO chat DEFAULT VALUES;`
	result, err := db.c.Exec(insertQuery)
	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("error while creating chat: %s", err.Error())
	}
	lastInsertID, err := result.LastInsertId()
	if !errors.Is(err, nil) {
		return 0, fmt.Errorf("error while retrieving chat ID: %s", err.Error())
	}
	return int(lastInsertID), nil
}

func (db *appdbimpl) CreateGroupDB(creatorNick string, groupName string, photoID int) (int, error) {
	exists, err := db.UserExistence(creatorNick)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error while checking existence of user %s: %w", creatorNick, err)
	}
	if !exists {
		return 404, fmt.Errorf("user %s does not exist", creatorNick)
	}

	creatorID, errorCode, err := db.IDFromNICK(creatorNick)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error converting nickname to ID: %s", err.Error())
	}

	chatID, err := db.CreateConversationDB()
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error inserting new chat: %s", err.Error())
	}

	insertQuery := `INSERT INTO groups (name, chat, photo) VALUES (?, ?, ?);`
	result, err := db.c.Exec(insertQuery, groupName, chatID, photoID)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error creating group: %s", err.Error())
	}

	groupID, err := result.LastInsertId()
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error retrieving group ID: %s", err.Error())
	}

	insertUserQuery := `INSERT INTO user_in_group (user, groups) VALUES (?, ?);`
	_, err = db.c.Exec(insertUserQuery, creatorID, groupID)
	if !errors.Is(err, nil) {
		return 304, fmt.Errorf("error adding user to group: %w", err)
	}

	return 0, nil
}

func (db *appdbimpl) CreatePrivateChatDB(nick1 string, nick2 string) (int, int, error) {
	exists1, err := db.UserExistence(nick1)
	if !errors.Is(err, nil) {
		return 0, 500, fmt.Errorf("error checking existence of user (caller): %s", err.Error())
	}
	if !exists1 {
		return 0, 404, fmt.Errorf("user %s does not exist (caller)", nick1)
	}

	exists2, err := db.UserExistence(nick2)
	if !errors.Is(err, nil) {
		return 0, 500, fmt.Errorf("error checking existence of user (target): %s", err.Error())
	}
	if !exists2 {
		return 0, 404, fmt.Errorf("user %s does not exist (target) or is not found in database", nick2)
	}

	existingChatID, errorCode, err := db.chatExistenceTraUtenti(nick1, nick2)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error checking if chat already exists: %s", err.Error())
	}
	if existingChatID > 0 {
		return 0, 400, fmt.Errorf("a chat already exists between users")
	}

	newChatID, err := db.CreateConversationDB()
	if !errors.Is(err, nil) {
		return 0, 500, fmt.Errorf("error creating new chat: %s", err.Error())
	}

	user1ID, errorCode, err := db.IDFromNICK(nick1)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error converting nickname %s to ID (user1): %s", nick1, err.Error())
	}
	user2ID, errorCode, err := db.IDFromNICK(nick2)
	if !errors.Is(err, nil) {
		return 0, errorCode, fmt.Errorf("error converting nickname %s to ID (user2): %s", nick2, err.Error())
	}

	insertQuery := `INSERT INTO privconv (chat, user1, user2) VALUES (?, ?, ?);`
	_, err = db.c.Exec(insertQuery, newChatID, user1ID, user2ID)
	if !errors.Is(err, nil) {
		return 0, 500, fmt.Errorf("error inserting private chat into database: %s", err.Error())
	}

	return newChatID, 0, nil
}
