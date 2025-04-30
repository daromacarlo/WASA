package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Funzione per inoltrare messaggi
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	NuovaConversazioneStr := ps.ByName("nuovachat")
	IdMessaggioStr := ps.ByName("messaggio")

	IdMessaggio, err := strconv.Atoi(IdMessaggioStr)
	if err != nil {
		CreaErroreJson(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}

	NuovaConversazione, err := strconv.Atoi(NuovaConversazioneStr)
	if err != nil {
		CreaErroreJson(w, "ID nuova conversazione non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.InoltraMessaggio(UtenteChiamante, NuovaConversazione, IdMessaggio)
	if err != nil {
		CreaErroreJson(w, "Errore durante l'inoltro del messaggio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	CreaRispostaJson(w, "Messaggio inoltrato con successo", http.StatusOK)
}
