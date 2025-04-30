package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) usersInGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Ottieni il nickname dall'URL
	chiamante := ps.ByName("utente")
	gruppostr := ps.ByName("gruppo")
	gruppo, err := strconv.Atoi(gruppostr)
	if err != nil {
		http.Error(w, "Errore durante la conversione del nome: "+err.Error(), http.StatusInternalServerError)
	}

	lista, err := rt.db.UsersInGroup(chiamante, gruppo)
	if err != nil {
		http.Error(w, "Errore durante il recupero degli utenti: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(lista) == 0 {
		return
	}
	if err := json.NewEncoder(w).Encode(lista); err != nil {
		http.Error(w, "Errore durante la codifica della risposta in JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	UtenteChiamante := ps.ByName("utente")

	var input struct {
		Foto string `json:"foto"`
	}

	// Decodifica il corpo della richiesta
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	// Validazione: la foto è obbligatoria
	if len(input.Foto) == 0 {
		http.Error(w, "La foto è obbligatoria", http.StatusBadRequest)
		return
	}

	// Chiamata al database per creare la foto profilo
	idFoto, err := rt.db.CreaFoto(input.Foto)
	if err != nil {
		http.Error(w, "Errore durante l'inserimento della foto profilo nella funzione setMyPhoto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = rt.db.ImpostaFotoProfilo(UtenteChiamante, idFoto)
	if err != nil {
		http.Error(w, "Errore durante l'aggiornamento della foto profilo: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")

	var input struct {
		Nome string `json:"nome"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	existe, err := rt.db.EsistenzaUtente(input.Nome)
	if err != nil {
		http.Error(w, "Errore durante il controllo dell'esistenza del nome utente: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if existe {
		http.Error(w, "Il nickname scritto è già in uso", http.StatusConflict)
		return
	}

	err = rt.db.ImpostaNome(UtenteChiamante, input.Nome)
	if err != nil {
		http.Error(w, "Errore durante l'aggiornamento del nome: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
