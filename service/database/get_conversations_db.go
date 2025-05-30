package database

import (
	"database/sql"
	"fmt"
	"sort"
)

type Conversation struct {
	Id        int     `json:"chat_id"`
	Name      string  `json:"name"`
	Photo     *string `json:"photo"`
	Time      *string `json:"time"`
	LastSnip  *string `json:"lastsnip"`
	LastPhoto *string `json:"lastphoto"`
	IsGroup   bool    `json:"isgroup"`
}

func (db *appdbimpl) GetConversationsDB(userNickname string) ([]Conversation, int, error) {
	userId, errorCode, err := db.IDFromNICK(userNickname)
	if err != nil {
		return nil, errorCode, err
	}

	privateConversationsQuery := `
    SELECT c.id, 
           CASE 
               WHEN cp.user1 = ? THEN u2.nickname 
               ELSE u1.nickname 
           END AS name,
           MAX(m.time) AS time, 
           m.text AS last, photo.photo AS lastPhoto,
           CASE 
               WHEN cp.user1 = ? THEN f2.photo
               ELSE f1.photo
           END AS photo,
           false AS isgroup
    FROM chat AS c
    JOIN privconv AS cp ON c.id = cp.chat
    JOIN user AS u1 ON u1.id = cp.user1
    JOIN user AS u2 ON u2.id = cp.user2
    LEFT JOIN photo AS f1 ON f1.id = u1.photo
    LEFT JOIN photo AS f2 ON f2.id = u2.photo
    LEFT JOIN message AS m ON m.chat = c.id
	LEFT JOIN photo ON m.photo = photo.id
    WHERE cp.user1 = ? OR cp.user2 = ?
    GROUP BY c.id, name, f1.photo, f2.photo
    `

	groupConversationsQuery := `
    SELECT c.id, g.name, MAX(m.time) AS time, m.text AS last, photo.photo AS lastPhoto, f.photo AS photo, true AS isgroup
    FROM chat AS c
    JOIN groups AS g ON g.chat = c.id
    JOIN user_in_group AS ug ON g.id = ug.groups
    JOIN photo AS f ON f.id = g.photo
    LEFT JOIN message AS m ON m.chat = c.id
	LEFT JOIN photo ON m.photo = photo.id
    WHERE ug.user = ?
    GROUP BY c.id, g.name, f.photo
    `

	var conversations []Conversation

	privateRows, err := db.c.Query(privateConversationsQuery, userId, userId, userId, userId)
	if err != nil {
		return nil, 500, fmt.Errorf("error retrieving private conversations: %w", err)
	}

	var privateConversationIDs []int
	for privateRows.Next() {
		var conv Conversation
		var photo sql.NullString
		if err := privateRows.Scan(&conv.Id, &conv.Name, &conv.Time, &conv.LastSnip, &conv.LastPhoto, &photo, &conv.IsGroup); err != nil {
			return nil, 500, fmt.Errorf("error reading private conversations: %w", err)
		}
		if photo.Valid && photo.String != "" {
			conv.Photo = &photo.String
		} else {
			conv.Photo = nil
		}
		if conv.LastSnip != nil && len(*conv.LastSnip) > 15 {
			*conv.LastSnip = (*conv.LastSnip)[:15] + "..."
		}
		conversations = append(conversations, conv)
		privateConversationIDs = append(privateConversationIDs, conv.Id)
	}

	if err := privateRows.Err(); err != nil {
		return nil, 500, fmt.Errorf("error iterating private conversations")
	}

	groupRows, err := db.c.Query(groupConversationsQuery, userId)
	if err != nil {
		return nil, 500, fmt.Errorf("error retrieving group conversations: %w", err)
	}

	var groupConversationIDs []int
	for groupRows.Next() {
		var conv Conversation
		var photo sql.NullString
		if err := groupRows.Scan(&conv.Id, &conv.Name, &conv.Time, &conv.LastSnip, &conv.LastPhoto, &photo, &conv.IsGroup); err != nil {
			return nil, 500, fmt.Errorf("error reading group conversations: %w", err)
		}
		if photo.Valid && photo.String != "" {
			conv.Photo = &photo.String
		} else {
			conv.Photo = nil
		}
		if conv.LastSnip != nil && len(*conv.LastSnip) > 15 {
			*conv.LastSnip = (*conv.LastSnip)[:15] + "..."
		}
		conversations = append(conversations, conv)
		groupConversationIDs = append(groupConversationIDs, conv.Id)
	}

	if err := groupRows.Err(); err != nil {
		return nil, 500, fmt.Errorf("error iterating group conversations")
	}

	for _, convId := range privateConversationIDs {
		if err := db.MarkPrivateMessagesAsReceived(userNickname, convId); err != nil {
			return nil, 500, fmt.Errorf("error marking private messages as received: %w", err)
		}
	}

	for _, convId := range groupConversationIDs {
		if err := db.MarkGroupMessagesAsReceived(userNickname, convId); err != nil {
			return nil, 500, fmt.Errorf("error marking group messages as received: %w", err)
		}
		if err := db.CheckGroupMessageReceipt(convId); err != nil {
			return nil, 500, fmt.Errorf("error checking read status of group messages: %w", err)
		}
	}

	sort.Slice(conversations, func(i, j int) bool {
		if conversations[i].Time == nil {
			return false
		}
		if conversations[j].Time == nil {
			return true
		}
		return *conversations[i].Time > *conversations[j].Time
	})

	return conversations, 0, nil
}
