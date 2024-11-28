package database

import "fmt"

func (db *appdbimpl) ChangePhoto(photo []byte, nickname string) error {

	query := "UPDATE persone SET photo = ? WHERE nickname = ?"

	_, err := db.c.Exec(query, photo, nickname)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento della foto: %w", err)
	}
	return nil
}
