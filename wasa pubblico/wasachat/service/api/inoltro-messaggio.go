package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// funzione per inoltrare messaggi ad un gruppo
func (rt *_router) InoltraMessaggioGruppo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	err = rt.db.InoltraMessaggioGruppo(UtenteChiamante, NuovaConversazione, IdMessaggio)
	if err != nil {
		http.Error(w, "Errore durante la creazione del messaggio: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Messaggio inoltrato con successo "))
	return
}

// funzione per inoltrare messaggi ad una chat privata
func (rt *_router) InoltraMessaggioPrivato(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	destinatario := ps.ByName("destinatario")
	IdMessaggioStr := ps.ByName("messaggio")
	IdMessaggio, err := strconv.Atoi(IdMessaggioStr)
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}
	err = rt.db.InoltraMessaggioPrivato(UtenteChiamante, destinatario, IdMessaggio)
	if err != nil {
		http.Error(w, "Errore durante la creazione del messaggio: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Messaggio inoltrato con successo "))
	return
}
