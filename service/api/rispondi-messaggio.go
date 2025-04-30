package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// funzione per rispondere ad un messaggio
func (rt *_router) ansMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Testo string `json:"testo"`
		Foto  string `json:"foto"`
	}

	UtenteChiamante := ps.ByName("utente")
	IdMessaggioStr := ps.ByName("messaggio")

	IdMessaggio, err := strconv.Atoi(IdMessaggioStr)
	if err != nil {
		CreaErroreJson(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}

	ConversazioneStr := ps.ByName("chat")
	Conversazione, err := strconv.Atoi(ConversazioneStr)
	if err != nil {
		CreaErroreJson(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		CreaErroreJson(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	if len(input.Testo) == 0 && len(input.Foto) == 0 {
		CreaErroreJson(w, "Errore, il messaggio deve contenere almeno una foto o del testo", http.StatusBadRequest)
		return
	}

	if len(input.Testo) > 0 && len(input.Foto) > 0 {
		CreaErroreJson(w, "Errore, il messaggio non può contenere sia foto che testo", http.StatusBadRequest)
		return
	}

	if len(input.Testo) == 0 {
		idFoto, err := rt.db.CreaFoto(input.Foto)
		if err != nil {
			CreaErroreJson(w, "Errore durante l'inserimento della foto per la risposta: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = rt.db.RispondiMessaggioFoto(UtenteChiamante, Conversazione, IdMessaggio, idFoto)
		if err != nil {
			CreaErroreJson(w, "Errore durante la risposta con foto: "+err.Error(), http.StatusInternalServerError)
			return
		}
		CreaRispostaJson(w, "Risposta con foto inviata con successo", http.StatusOK)
		return
	}

	err = rt.db.RispondiMessaggioTesto(UtenteChiamante, Conversazione, IdMessaggio, input.Testo)
	if err != nil {
		CreaErroreJson(w, "Errore durante l'inserimento del messaggio di risposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
	CreaRispostaJson(w, "Risposta inviata con successo", http.StatusOK)
}
