package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type Profile struct {
	Nickname string
}

func CreateTableUser(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS user (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nickname TEXT UNIQUE NOT NULL,
			photo,
			FOREIGN KEY (photo) REFERENCES photo(id)
		);`
	_, err := db.Exec(query)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error while creating user table: %w", err)
	}
	return nil
}

func (db *appdbimpl) CreateUser(nickname string, idphoto int) error {
	ex, err := db.UserExistence(nickname)
	if ex {
		return fmt.Errorf("nickname is already in use: %w", err)
	}

	insertQuery := `INSERT INTO user (nickname, photo) VALUES (?,?);`

	_, err = db.c.Exec(insertQuery, nickname, idphoto)
	if !errors.Is(err, nil) {
		return fmt.Errorf("unexpected error while creating user: %w", err)
	}
	return nil
}

func (db *appdbimpl) Login(nickname string) (string, error) {
	ex, err := db.UserExistence(nickname)
	if !ex {
		err = db.CreateUser(nickname, 0)
		if !errors.Is(err, nil) {
			return "", fmt.Errorf("error creating a new user after not finding it in the database")
		}
	}
	if !errors.Is(err, nil) {
		return "", fmt.Errorf("error checking if user exists")
	}

	return nickname, nil
}

func (db *appdbimpl) UserExistence(nickname string) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM user WHERE nickname = ?;`
	err := db.c.QueryRow(query, nickname).Scan(&count)
	if !errors.Is(err, nil) {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("unexpected error while checking user existence: %w", err)
	}
	return count > 0, nil
}

func (db *appdbimpl) IDFromNICK(nickname string) (int, int, error) {
	var id int
	query := `SELECT id FROM user WHERE nickname = ?;`
	err := db.c.QueryRow(query, nickname).Scan(&id)
	if !errors.Is(err, nil) {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 404, fmt.Errorf("user with nickname '%s' not found", nickname)
		}
		return 0, 500, fmt.Errorf("error retrieving user ID: %w", err)
	}
	return id, 0, nil
}

func (db *appdbimpl) NICKFromID(idPassato int) (string, int, error) {
	var nickname string
	query := `SELECT nickname FROM user WHERE id = ?;`
	err := db.c.QueryRow(query, idPassato).Scan(&nickname)
	if !errors.Is(err, nil) {
		if errors.Is(err, sql.ErrNoRows) {
			return "", 404, fmt.Errorf("user with ID '%d' not found", idPassato)
		}
		return "", 500, fmt.Errorf("error retrieving user nickname: %w", err)
	}
	return nickname, 0, nil
}

func (db *appdbimpl) SetPhotoDB(nickname string, idPhoto int) error {
	ex, err := db.UserExistence(nickname)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if !ex {
		return fmt.Errorf("user %s does not exist", nickname)
	}

	queryUpdateFoto := `UPDATE user SET photo = ? WHERE nickname = ?`
	_, err = db.c.Exec(queryUpdateFoto, idPhoto, nickname)
	if !errors.Is(err, nil) {
		return fmt.Errorf("error updating profile photo: %w", err)
	}

	return nil
}

func (db *appdbimpl) SetNameDB(nickname string, newNickname string) (int, error) {
	if nickname == "" || newNickname == "" {
		return 400, fmt.Errorf("malformed request")
	}

	if nickname == newNickname {
		return 409, fmt.Errorf("new nickname is the same as the old one")
	}

	exists, err := db.UserExistence(nickname)
	if err != nil {
		return 500, fmt.Errorf("error checking if the user exists: %w", err)
	}
	if !exists {
		return 404, fmt.Errorf("user %s does not exist", nickname)
	}

	exNewNick, err := db.UserExistence(newNickname)
	if err != nil {
		return 500, fmt.Errorf("error checking if new nickname is in use: %w", err)
	}
	if exNewNick {
		return 409, fmt.Errorf("nickname %s is already in use", newNickname)
	}

	queryUpdate := `UPDATE user SET nickname = ? WHERE nickname = ?`
	_, err = db.c.Exec(queryUpdate, newNickname, nickname)
	if err != nil {
		return 500, fmt.Errorf("error updating nickname: %w", err)
	}

	return 200, nil
}

func (db *appdbimpl) UsersInGroup(callingUser string, chat int) ([]Profile, int, error) {
	ex, err := db.chatExistence(chat)
	if !errors.Is(err, nil) {
		return nil, 500, fmt.Errorf("error checking if chat exists: %w", err)
	}
	if !ex {
		return nil, 404, fmt.Errorf("chat with ID %d does not exist", chat)
	}
	exc, err := db.UserExistence(callingUser)
	if !errors.Is(err, nil) {
		return nil, 500, fmt.Errorf("error checking if callingUser %s exists: %w", callingUser, err)
	}
	if !exc {
		return nil, 404, fmt.Errorf("user callingUser %s does not exist", callingUser)
	}
	callingUseringroup, errorCode, err := db.UserInGroup(callingUser, chat)
	if !errors.Is(err, nil) {
		return nil, errorCode, fmt.Errorf("error checking if user is in group: %w", err)
	}
	if callingUseringroup == 0 {
		return nil, 401, fmt.Errorf("user %s is not part of the group", callingUser)
	}

	query := `SELECT u.nickname
			  FROM user as u
			  JOIN user_in_group as uig on u.id = uig.user
			  JOIN groups as g on uig.groups = g.id
			  JOIN chat as c on c.id = g.chat
			  WHERE c.id = ?
			  `

	rows, err := db.c.Query(query, chat)
	if !errors.Is(err, nil) {
		return nil, 500, fmt.Errorf("error retrieving user profiles: %w", err)
	}

	var list []Profile

	for rows.Next() {
		var nickname string

		if err := rows.Scan(&nickname); err != nil {
			return nil, 500, fmt.Errorf("error reading data: %w", err)
		}

		list = append(list, Profile{
			Nickname: nickname,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, 500, fmt.Errorf("error iterating over users")
	}

	return list, 0, nil
}
