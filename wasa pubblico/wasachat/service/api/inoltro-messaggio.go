package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// funzione per inoltrare messaggi ad un gruppo
func (rt *_router) InoltraMessaggio(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	NuovaConversazioneStr := ps.ByName("nuovachat")
	IdMessaggioStr := ps.ByName("messaggio")
	IdMessaggio, err := strconv.Atoi(IdMessaggioStr)
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}
	NuovaConversazione, err := strconv.Atoi(NuovaConversazioneStr)
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}
	err = rt.db.InoltraMessaggio(UtenteChiamante, NuovaConversazione, IdMessaggio)
	if err != nil {
		http.Error(w, "Errore durante la creazione del messaggio: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Messaggio inoltrato con successo "))
	return
}
