package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) SendPrivateMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	nickname := ps.ByName("nickname")

	// Struttura per ricevere i dati dal body
	var input struct {
		Destinatario string `json:"destinatario"`
		Messaggio    string `json:"messaggio"`
	}

	// Decodifica il corpo della richiesta
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || len(input.Destinatario) == 0 || len(input.Messaggio) == 0 {
		http.Error(w, "Parametro destinatario o messaggio mancante o non valido", http.StatusBadRequest)
		return
	}

	// Chiamata al database per creare il messaggio privato
	err = rt.db.CreateMessaggioPrivato(nickname, input.Destinatario, input.Messaggio)
	if err != nil {
		http.Error(w, "Errore durante la creazione del messaggio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Risposta di successo
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
