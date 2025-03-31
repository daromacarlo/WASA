package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) registrare(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Nickname string `json:"nickname"`
		Foto     string `json:"foto"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	if len(input.Nickname) == 0 {
		http.Error(w, "Il nickname è obbligatorio", http.StatusBadRequest)
		return
	}

	if len(input.Foto) == 0 {
		http.Error(w, "La foto è obbligatoria", http.StatusBadRequest)
		return
	}

	idFoto, err := rt.db.CreaFoto(input.Foto)
	if err != nil {
		http.Error(w, "Errore durante l'inserimento della foto profilo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = rt.db.CreaUtente(input.Nickname, idFoto)
	if err != nil {
		http.Error(w, "Errore durante la creazione dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Nickname string `json:"nickname"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	if len(input.Nickname) == 0 {
		http.Error(w, "Il nickname è obbligatorio", http.StatusBadRequest)
		return
	}

	utente, err := rt.db.Login(input.Nickname)
	if err != nil {
		http.Error(w, "Errore durante la verifica dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if utente == "" {
		http.Error(w, "Credenziali non valide", http.StatusUnauthorized)
		return
	}
}
