package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	chiamante := ps.ByName("utente")
	chatStr := ps.ByName("chat")

	chat, err := strconv.Atoi(chatStr)
	if err != nil {
		CreaErroreJson(w, "Il parametro 'gruppo' deve essere un numero intero valido", http.StatusBadRequest)
		return
	}

	lista, codiceErrore, err := rt.db.GetConversazione(chiamante, chat)

	if err != nil {
		CreaErroreJson(w, "Errore durante il recupero della conversazione: "+err.Error(), codiceErrore)
		return
	}

	if len(lista) == 0 {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(lista); err != nil {
		CreaErroreJson(w, "Errore durante la codifica della risposta in JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
