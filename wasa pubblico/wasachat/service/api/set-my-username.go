package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) SetMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Recuperiamo il nickname corrente dall'URL
	nickname := ps.ByName("nickname")
	if len(nickname) == 0 {
		http.Error(w, "Parametro nickname mancante nell'URL", http.StatusBadRequest)
		return
	}
	var newnickname string

	// Leggiamo il nuovo nickname dal body
	var input struct {
		NewNickname string `json:"newnickname"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || len(input.NewNickname) == 0 {
		http.Error(w, "Parametro newnickname mancante o non valido", http.StatusBadRequest)
		return
	}
	newnickname = input.NewNickname

	// Log di debug per verificare i parametri
	fmt.Printf("Nickname corrente: '%s', Nuovo nickname: '%s'\n ", nickname, newnickname)

	// Se il nuovo nickname è uguale a quello corrente
	if nickname == newnickname {
		fmt.Fprintf(w, `Il nuovo nome utente è uguale a quello già in uso, nessuna modifica effettuata `)
		return
	}

	//controllo se esiste l'utetne esiste
	esistenza, err := rt.db.SearchUser(nickname)
	if err != nil {
		http.Error(w, fmt.Sprintf("Errore durante la ricerca dell'utente: %s ", err), http.StatusInternalServerError)
		return
	}

	//se non esiste
	if !esistenza {
		fmt.Fprintf(w, `l'utente '%s' non esiste `, nickname)
		return
	}

	// Controlla se il nuovo nickname è già in uso
	giaInUso, err := rt.db.SearchUser(newnickname)
	if err != nil {
		http.Error(w, fmt.Sprintf("Errore durante la ricerca dell'utente: %s ", err), http.StatusInternalServerError)
		return
	}

	if giaInUso {
		fmt.Fprintf(w, `Il nome utente '%s' è già in uso, prova con uno più originale! `, newnickname)
		return
	}

	// Modifico il nickname nel database
	err = rt.db.ChangeNickname(newnickname, nickname)
	if err != nil {
		http.Error(w, fmt.Sprintf("Errore durante la modifica del nome utente: %s ", err), http.StatusInternalServerError)
		return
	}

	// Risposta di successo
	fmt.Fprintf(w, `Il tuo nome utente è stato aggiornato con successo da '%s' a '%s' `, nickname, newnickname)
}
