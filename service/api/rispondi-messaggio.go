package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) RispondiAMessaggio(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Testo string `json:"testo"`
		Foto  string `json:"foto"`
	}
	UtenteChiamante := ps.ByName("utente")

	IdMessaggio, err := strconv.Atoi(ps.ByName("idMessaggio"))
	if err != nil {
		http.Error(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}
	ConversazioneStr := ps.ByName("chat")
	Conversazione, err := strconv.Atoi(ConversazioneStr)
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}
	if len(input.Testo) == 0 && len(input.Foto) == 0 {
		http.Error(w, "Errore, il messaggio deve contenere almeno una foto o del testo", http.StatusBadRequest)
		return
	}

	if len(input.Testo) > 0 && len(input.Foto) > 0 {
		http.Error(w, "Errore, il messaggio non può contenere sia foto che testo", http.StatusBadRequest)
		return
	}

	if len(input.Testo) == 0 {
		idFoto, err := rt.db.CreaFoto(input.Foto)
		if err != nil {
			http.Error(w, "Errore durante l'inserimento della foto profilo: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = rt.db.RispondiMessaggioFoto(UtenteChiamante, Conversazione, IdMessaggio, idFoto)
		if err != nil {
			http.Error(w, "Errore durante la creazione del messaggio: "+err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	err = rt.db.RispondiMessaggioTesto(UtenteChiamante, Conversazione, IdMessaggio, input.Testo)
	if err != nil {
		http.Error(w, "Errore durante la creazione del messaggio: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
