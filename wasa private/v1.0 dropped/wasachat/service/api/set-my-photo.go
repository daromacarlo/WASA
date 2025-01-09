package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) SetMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Recuperiamo il nickname corrente dall'URL
	nickname := ps.ByName("nickname")
	if len(nickname) == 0 {
		http.Error(w, `{"error": "Parametro nickname mancante nell'URL"}`, http.StatusBadRequest)
		return
	}

	// Variabile per contenere i dati della nuova foto
	var newPhoto []byte

	// Leggiamo l'input dal corpo della richiesta
	var input struct {
		NewPhoto string `json:"newphoto"` // Utilizziamo una stringa perché JSON non gestisce direttamente []byte
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || len(input.NewPhoto) == 0 {
		http.Error(w, `{"error": "Parametro newPhoto mancante o non valido"}`, http.StatusBadRequest)
		return
	}

	// Convertiamo la stringa in []byte
	newPhoto = []byte(input.NewPhoto)

	// Controlliamo se l'utente esiste
	esistenza, err := rt.db.SearchUser(nickname)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Errore durante la ricerca dell'utente: %s"}`, err), http.StatusInternalServerError)
		return
	}

	// Se l'utente non esiste
	if !esistenza {
		http.Error(w, fmt.Sprintf(`{"error": "L'utente '%s' non esiste"}`, nickname), http.StatusNotFound)
		return
	}

	// Modifica della foto nel database
	err = rt.db.ChangePhoto(newPhoto, nickname)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Errore durante la modifica della foto profilo: %s"}`, err), http.StatusInternalServerError)
		return
	}

	// Risposta di successo
	fmt.Fprintf(w, `{"message": "La foto profilo è stata cambiata con successo per l'utente '%s'"}`, nickname)
}
