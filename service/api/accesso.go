package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Funzione per registrare un nuovo utente
func (rt *_router) register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Nickname string `json:"nickname"`
		Foto     string `json:"foto"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if !errors.Is(err, nil) {
		CreaErroreJson(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}
	if len(input.Nickname) == 0 {
		CreaErroreJson(w, "Il nickname è obbligatorio", http.StatusBadRequest)
		return
	}
	if len(input.Foto) == 0 {
		CreaErroreJson(w, "La foto è obbligatoria", http.StatusBadRequest)
		return
	}
	idFoto, err := rt.db.CreaFoto(input.Foto)
	if !errors.Is(err, nil) {
		CreaErroreJson(w, "Errore durante l'inserimento della foto profilo durante la registrazione: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = rt.db.CreaUtente(input.Nickname, idFoto)
	if !errors.Is(err, nil) {
		CreaErroreJson(w, "Errore durante la creazione dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}
	CreaRispostaJson(w, "Registrazione avvenuta con successo", http.StatusOK)
}

// Funzione per effettuare il login
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Nickname string `json:"nickname"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if !errors.Is(err, nil) {
		CreaErroreJson(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}
	if len(input.Nickname) == 0 {
		CreaErroreJson(w, "Il nickname è obbligatorio", http.StatusBadRequest)
		return
	}
	utente, err := rt.db.Login(input.Nickname)
	if !errors.Is(err, nil) {
		CreaErroreJson(w, "Errore durante la verifica dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if utente == "" {
		CreaErroreJson(w, "Credenziali non valide", http.StatusBadRequest)
		return
	}
	CreaRispostaJson(w, "Login effettuato con successo", http.StatusOK)
}
