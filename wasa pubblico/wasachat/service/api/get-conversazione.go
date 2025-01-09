package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetConversazionePrivata(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Ottieni il nickname dall'URL
	chiamante := ps.ByName("utente")
	destinatario := ps.ByName("destinatario")

	// Recupera la lista dei messaggi dal database
	lista, err := rt.db.GetConversazionePrivata(chiamante, destinatario)

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
func (rt *_router) GetConversazioneGruppo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Ottieni il nickname dall'URL
	chiamante := ps.ByName("utente")
	gruppoStr := ps.ByName("gruppo")

	// Converte il parametro 'gruppo' da stringa a intero
	gruppo, err := strconv.Atoi(gruppoStr)
	if err != nil {
		http.Error(w, "Il parametro 'gruppo' deve essere un numero intero valido", http.StatusBadRequest)
		return
	}

	// Recupera la lista dei messaggi dal database
	lista, err := rt.db.GetConversazioneGruppo(chiamante, gruppo)

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
