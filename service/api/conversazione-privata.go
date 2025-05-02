package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createPrivateConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Utente string `json:"utente"`
	}

	UtenteChiamante := ps.ByName("utente")

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		CreaErroreJson(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	if len(input.Utente) == 0 {
		CreaErroreJson(w, "Il nome dell'utente è obbligatorio", http.StatusBadRequest)
		return
	}

	if len(UtenteChiamante) == 0 {
		CreaErroreJson(w, "Utente chiamante non esplicitato correttamente", http.StatusBadRequest)
		return
	}

	_, codiceErrore, err := rt.db.CreaConversazionePrivataDB(UtenteChiamante, input.Utente)
	if err != nil {
		CreaErroreJson(w, "Errore durante la creazione della conversazione: "+err.Error(), codiceErrore)
		return
	}
	CreaRispostaJson(w, "Conversazione privata creata con successo", http.StatusOK)
}
