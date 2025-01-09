package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetConversazioni(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Ottieni il nickname dall'URL
	chiamante := ps.ByName("utente")

	// Recupera la lista dei messaggi dal database
	lista, err := rt.db.GetConversazioni(chiamante)

	// Gestisci gli errori
	if err != nil {
		http.Error(w, "Errore durante il recupero della conversazione: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Controlla se la lista è vuota
	if len(lista) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Nessuna conversazione aperta"}`))
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
