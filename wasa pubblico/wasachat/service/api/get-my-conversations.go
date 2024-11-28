package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetMyConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Ottieni il nickname dall'URL
	nickname := ps.ByName("nickname")

	// Recupera la lista delle chat dal database
	lista, err := rt.db.GetChats(nickname)

	// Gestisci gli errori
	if err != nil {
		http.Error(w, "Errore durante il recupero delle conversazioni: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Controlla se la lista è vuota
	if len(lista) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Nessuna conversazione trovata per l'utente"}`))
		return
	}

	// Converte la lista in JSON e restituisce la risposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(lista); err != nil {
		http.Error(w, "Errore durante la codifica della risposta in JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
