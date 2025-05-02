package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) usersInGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	chiamante := ps.ByName("utente")
	gruppostr := ps.ByName("gruppo")
	gruppo, err := strconv.Atoi(gruppostr)
	if err != nil {
		CreaErroreJson(w, "Errore durante la conversione del nome: "+err.Error(), http.StatusInternalServerError)
		return
	}

	lista, codiceErrore, err := rt.db.UsersInGroup(chiamante, gruppo) // <-- ordine corretto
	if err != nil {
		CreaErroreJson(w, "Errore durante il recupero degli utenti: "+err.Error(), codiceErrore)
		return
	}
	if len(lista) == 0 {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lista); err != nil {
		CreaErroreJson(w, "Errore durante la codifica della risposta in JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	var input struct {
		Foto string `json:"foto"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		CreaErroreJson(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	if len(input.Foto) == 0 {
		CreaErroreJson(w, "La foto è obbligatoria", http.StatusBadRequest)
		return
	}

	idFoto, err := rt.db.CreaFoto(input.Foto)
	if err != nil {
		CreaErroreJson(w, "Errore durante l'inserimento della foto profilo nella funzione setMyPhoto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = rt.db.ImpostaFotoProfilo(UtenteChiamante, idFoto)
	if err != nil {
		CreaErroreJson(w, "Errore durante l'aggiornamento della foto profilo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	CreaRispostaJson(w, "Foto aggiornata con successo", http.StatusOK)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	var input struct {
		Nome string `json:"nome"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		CreaErroreJson(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	existe, err := rt.db.EsistenzaUtente(input.Nome)
	if err != nil {
		CreaErroreJson(w, "Errore durante il controllo dell'esistenza del nome utente: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if existe {
		CreaErroreJson(w, "Il nickname scritto è già in uso", http.StatusConflict)
		return
	}

	codiceErrore, err := rt.db.ImpostaNome(UtenteChiamante, input.Nome)
	if err != nil {
		CreaErroreJson(w, "Errore durante l'aggiornamento del nome: "+err.Error(), codiceErrore)
		return
	}

	CreaRispostaJson(w, "Nome aggiornato con successo", http.StatusOK)
}
