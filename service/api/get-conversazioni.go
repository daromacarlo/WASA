package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	chiamante := ps.ByName("utente")
	lista, err := rt.db.GetConversazioni(chiamante)

	if err != nil {
		http.Error(w, "Errore durante il recupero della conversazione: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(lista) == 0 {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lista); err != nil {
		http.Error(w, "Errore durante la codifica della risposta in JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
