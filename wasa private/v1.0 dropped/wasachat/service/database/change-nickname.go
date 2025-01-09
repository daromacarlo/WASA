package database

import "fmt"

func (db *appdbimpl) ChangeNickname(newnickname string, nickname string) error {

	query := "UPDATE persone SET nickname = ? WHERE nickname = ?"

	_, err := db.c.Exec(query, newnickname, nickname)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento del nickname: %w", err)
	}
	return nil
}
