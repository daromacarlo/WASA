package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) CreaConversazionePrivata(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	var input struct {
		Utente string `json:"utente"`
	}

	UtenteChiamante := ps.ByName("utente")

	// Decodifica il corpo della richiesta
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	// controlliamo che la richiesta sia stata strutturata nella maniera corretta
	if len(input.Utente) == 0 {
		http.Error(w, "Il nome è obbligatorio", http.StatusBadRequest)
		return
	}

	if len(UtenteChiamante) == 0 {
		http.Error(w, "utente chiamante non esplicitato correttamente", http.StatusBadRequest)
		return
	}

	_, err = rt.db.CreaConversazionePrivataDB(UtenteChiamante, input.Utente)
	if err != nil {
		http.Error(w, "Errore durante la creazione della conversazione: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Risposta di successo
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Conversazione creata con successo "))
}

// test
/*
curl -X POST http://localhost:3000/wasachat/test/Mario/conversazioniprivate \
-H "Content-Type: application/json" \
-d '{
  "utente": "Luigi"
}'
*/
