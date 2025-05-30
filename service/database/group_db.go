package database

import (
	"errors"
	"fmt"
)

// Adds a user to a group chat if conditions are met
func (db *appdbimpl) AddToGroupDB(chatID int, callingUser string, targetUser string) (int, error) {
	callingUserExists, err := db.UserExistence(callingUser)
	if err != nil {
		return 500, fmt.Errorf("error checking existence of callingUser %s: %w", callingUser, err)
	}
	if !callingUserExists {
		return 404, fmt.Errorf("callingUser %s does not exist", callingUser)
	}

	isCallingUserInGroup, errorCode, err := db.UserInGroup(callingUser, chatID)
	if err != nil {
		return errorCode, fmt.Errorf("error checking if callingUser is in group: %w", err)
	}
	if isCallingUserInGroup == 0 {
		return 401, fmt.Errorf("user %s is not part of the group", callingUser)
	}

	targetUserExists, err := db.UserExistence(targetUser)
	if err != nil {
		return 500, fmt.Errorf("error checking existence of targetUser %s: %w", targetUser, err)
	}
	if !targetUserExists {
		return 404, fmt.Errorf("targetUser %s does not exist", targetUser)
	}

	isTargetUserInGroup, errorCode, err := db.UserInGroup(targetUser, chatID)
	if err != nil {
		return errorCode, fmt.Errorf("error checking if targetUser is in group: %w", err)
	}
	if isTargetUserInGroup > 0 {
		return 400, fmt.Errorf("user %s is already in the group", targetUser)
	}

	targetUserID, errorCode, err := db.IDFromNICK(targetUser)
	if err != nil {
		return errorCode, fmt.Errorf("error converting nickname %s to ID: %w", targetUser, err)
	}

	chatExists, err := db.chatExistence(chatID)
	if err != nil {
		return 500, fmt.Errorf("error verifying chat existence: %w", err)
	}
	if !chatExists {
		return 404, fmt.Errorf("chat with ID %d does not exist", chatID)
	}

	groupID, errorCode, err := db.SearchGroup(chatID)
	if err != nil {
		return errorCode, fmt.Errorf("error finding group linked to chat: %w", err)
	}

	_, err = db.c.Exec(`INSERT INTO user_in_group (user, groups) VALUES (?, ?);`, targetUserID, groupID)
	if err != nil {
		return 500, fmt.Errorf("error adding user to group: %w", err)
	}

	return 0, nil
}

// Removes a user from a group
func (db *appdbimpl) LeaveGroup(chatId int, callingUser string) (int, error) {
	excallingUser, err := db.UserExistence(callingUser)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error checking existence of callingUser %s: %w", callingUser, err)
	}
	if !excallingUser {
		return 404, fmt.Errorf("callingUser %s does not exist", callingUser)
	}
	callingUserPresente, errorCode, err := db.UserInGroup(callingUser, chatId)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error checking if callingUser is in group: %w", err)
	}
	if callingUserPresente == 0 {
		return 401, fmt.Errorf("user %s is not part of the group", callingUser)
	}
	callerUser_converted, errorCode, err := db.IDFromNICK(callingUser)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error converting nickname %s to ID: %w", callingUser, err)
	}
	ex, err := db.chatExistence(chatId)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error verifying chat existence: %w", err)
	}
	if !ex {
		return 404, fmt.Errorf("chat with ID %d does not exist", chatId)
	}
	idGruppo, errorCode, err := db.SearchGroup(chatId)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error finding group linked to chat: %w", err)
	}
	insertQuery := `DELETE from user_in_group WHERE user = ? AND groups = ?;`
	_, err = db.c.Exec(insertQuery, callerUser_converted, idGruppo)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error removing user from group: %w", err)
	}
	return 0, nil
}

// Sets a new group profile picture
func (db *appdbimpl) SetGroupPhotoDB(callingUser string, idPhoto int, idGroup int) (int, error) {
	ex, err := db.chatExistence(idGroup)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error checking chat existence: %w", err)
	}
	if !ex {
		return 404, fmt.Errorf("specified chat does not exist")
	}
	groupType, errorCode, err := db.SearchGroup(idGroup)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error checking chat type: %w", err)
	}
	if groupType == 0 {
		return 400, fmt.Errorf("specified chat is not a group")
	}
	involved, errorCode, err := db.UserInGroup(callingUser, idGroup)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error checking user membership: %w", err)
	}
	if involved == 0 {
		return 401, fmt.Errorf("user %s is not part of group %d", callingUser, idGroup)
	}
	query := `UPDATE groups SET photo = ? WHERE chat = ?`
	_, err = db.c.Exec(query, idPhoto, idGroup)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error updating group photo: %w", err)
	}
	return 0, nil
}

// Updates the group name if different from the current one
func (db *appdbimpl) SetGroupNameDB(callingUser string, nomeGruppo_passeed string, idGroup int) (int, error) {
	ex, err := db.chatExistence(idGroup)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error checking chat existence: %w", err)
	}
	if !ex {
		return 404, fmt.Errorf("specified chat does not exist")
	}
	groupType, errorCode, err := db.SearchGroup(idGroup)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error checking chat type: %w", err)
	}
	if groupType == 0 {
		return 400, fmt.Errorf("specified chat is not a group")
	}
	involved, errorCode, err := db.UserInGroup(callingUser, idGroup)
	if !errors.Is(err, nil) {
		return errorCode, fmt.Errorf("error checking user membership: %w", err)
	}
	if involved == 0 {
		return 401, fmt.Errorf("user %s is not part of group %d", callingUser, idGroup)
	}
	var oldName string
	queryGetNome := `SELECT name FROM groups WHERE chat = ?`
	err = db.c.QueryRow(queryGetNome, idGroup).Scan(&oldName)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error retrieving group name: %w", err)
	}
	if oldName == nomeGruppo_passeed {
		return 400, fmt.Errorf("new name is the same as the old name")
	}
	queryUpdate := `UPDATE groups SET name = ? WHERE chat = ?`
	_, err = db.c.Exec(queryUpdate, nomeGruppo_passeed, idGroup)
	if !errors.Is(err, nil) {
		return 500, fmt.Errorf("error updating group name: %w", err)
	}
	return 0, nil
}
