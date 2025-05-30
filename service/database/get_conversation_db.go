package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) GetChat(requestingUser string, chatId int) ([]MessageData, int, error) {
	exists, err := db.chatExistence(chatId)
	if !errors.Is(err, nil) {
		return nil, 500, fmt.Errorf("error checking chat existence: %w", err)
	}
	if !exists {
		return nil, 404, fmt.Errorf("error, chat does not exist or was not found in database")
	}

	conversationType, errorCode, err := db.SearchGroup(chatId)
	if !errors.Is(err, nil) {
		return nil, errorCode, fmt.Errorf("error checking conversation type: %w", err)
	}
	if conversationType > 0 {
		return db.GetGroupChat(requestingUser, chatId)
	} else {
		otherUser, err := db.GetNameUserInvolved(chatId, requestingUser)
		if !errors.Is(err, nil) {
			return nil, 500, fmt.Errorf("error retrieving name of involved user: %w", err)
		}
		return db.GetPrivateChat(requestingUser, otherUser)
	}
}

type MessageData struct {
	MessageID int       `json:"message_id"`
	Author    string    `json:"author"`
	AuthorID  int       `json:"idauthor"`
	Text      *string   `json:"text"`
	Photo     *string   `json:"photo"`
	Time      string    `json:"time"`
	Received  bool      `json:"rec"`
	Read      bool      `json:"read"`
	Comments  []Comment `json:"comments"`
	Forwarded bool      `json:"forw"`
	ReplyTo   *int      `json:"ans"`
}

type Comment struct {
	CommentID int    `json:"comment_id"`
	Author    string `json:"author"`
	AuthorID  int    `json:"idauthor"`
	Reaction  string `json:"reaction"`
}

func (db *appdbimpl) GetPrivateChat(user1 string, user2 string) ([]MessageData, int, error) {
	chatID, errorCode, err := db.chatExistenceTraUtenti(user1, user2)
	if !errors.Is(err, nil) {
		return nil, errorCode, fmt.Errorf("error retrieving chat between users: %w", err)
	}

	err = db.ReadPrivateMessage(user2, chatID)
	if !errors.Is(err, nil) {
		return nil, 500, fmt.Errorf("error updating message state: %w", err)
	}

	query := `
		SELECT m.id, m.author, m.idauthor, m.text, f.photo, m.time, sm.rec, sm.read, m.forw, m.ans
		FROM message m
		JOIN private_state as sm on m.id = sm.message 
		LEFT JOIN photo as f on m.photo = f.id
		WHERE m.chat = ?
		ORDER BY m.time ASC;`

	rows, err := db.c.Query(query, chatID)
	if !errors.Is(err, nil) {
		return nil, 500, fmt.Errorf("error retrieving messages: %w", err)
	}

	var messages []MessageData

	for rows.Next() {
		var m MessageData
		if err := rows.Scan(&m.MessageID, &m.Author, &m.AuthorID, &m.Text, &m.Photo, &m.Time, &m.Received, &m.Read, &m.Forwarded, &m.ReplyTo); err != nil {
			return nil, 500, fmt.Errorf("error reading message data: %w", err)
		}

		m.Comments, errorCode, err = db.GetCommentFromMessage(m.MessageID)
		if !errors.Is(err, nil) {
			return nil, errorCode, fmt.Errorf("error retrieving comments for message %d: %w", m.MessageID, err)
		}

		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		return nil, 500, fmt.Errorf("error iterating private chat messages")
	}

	return messages, 0, nil
}

func (db *appdbimpl) GetGroupChat(user string, chatID int) ([]MessageData, int, error) {
	inGroup, errorCode, err := db.UserInGroup(user, chatID)
	if !errors.Is(err, nil) {
		return nil, errorCode, fmt.Errorf("error checking group membership: %w", err)
	}
	if inGroup == 0 {
		return nil, 401, fmt.Errorf("user %s is not a member of chat %d", user, chatID)
	}

	err = db.ReadGroupMessages(user, chatID)
	if !errors.Is(err, nil) {
		return nil, 418, fmt.Errorf("error updating group messages: %w", err)
	}
	err = db.CheckGroupMessageReadStatus(chatID)
	if !errors.Is(err, nil) {
		return nil, 500, fmt.Errorf("error checking message read status: %w", err)
	}

	query := `
		SELECT m.id, m.author ,m.idauthor, m.text, f.photo, m.time, smg.read, smg.rec, m.forw, m.ans
		FROM message m
		JOIN group_state as smg on smg.message = m.id
		LEFT JOIN photo as f on m.photo = f.id
		WHERE m.chat = ?
		ORDER BY m.time ASC;`

	rows, err := db.c.Query(query, chatID)
	if !errors.Is(err, nil) {
		return nil, 500, fmt.Errorf("error retrieving group messages: %w", err)
	}

	var messages []MessageData

	for rows.Next() {
		var m MessageData
		if err := rows.Scan(&m.MessageID, &m.Author, &m.AuthorID, &m.Text, &m.Photo, &m.Time, &m.Read, &m.Received, &m.Forwarded, &m.ReplyTo); err != nil {
			return nil, 500, fmt.Errorf("error reading message data: %w", err)
		}

		m.Comments, errorCode, err = db.GetCommentFromMessage(m.MessageID)
		if !errors.Is(err, nil) {
			return nil, errorCode, fmt.Errorf("error retrieving comments for message %d: %w", m.MessageID, err)
		}

		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		return nil, 500, fmt.Errorf("error iterating group chat messages")
	}

	return messages, 0, nil
}

func (db *appdbimpl) GetCommentFromMessage(messageID int) ([]Comment, int, error) {
	query := `
		SELECT c.id, c.author, c.idauthor, c.reaction
		FROM comment c
		WHERE c.message = ?;`

	rows, err := db.c.Query(query, messageID)
	if !errors.Is(err, nil) {
		return nil, 500, fmt.Errorf("error retrieving comments: %w", err)
	}

	var comments []Comment

	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.CommentID, &c.Author, &c.AuthorID, &c.Reaction); err != nil {
			return nil, 500, fmt.Errorf("error reading comment data: %w", err)
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, 500, fmt.Errorf("error iterating comments")
	}

	return comments, 0, nil
}
