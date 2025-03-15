package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Funzione che serve a gettare una conversazione di gruppo
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Ottieni il nickname dall'URL
	chiamante := ps.ByName("utente")
	chatStr := ps.ByName("chat")

	// Converte il parametro 'gruppo' da stringa a intero
	chat, err := strconv.Atoi(chatStr)
	if err != nil {
		http.Error(w, "Il parametro 'gruppo' deve essere un numero intero valido", http.StatusBadRequest)
		return
	}

	// Recupera la lista dei messaggi dal database
	lista, err := rt.db.GetConversazione(chiamante, chat)

	// Gestisci gli errori
	if err != nil {
		http.Error(w, "Errore durante il recupero della conversazione: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Controlla se la lista è vuota
	if len(lista) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Nessun messaggio trovato in questa conversazione"}`))
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
