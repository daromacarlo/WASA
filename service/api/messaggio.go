package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Testo string `json:"testo"`
		Foto  string `json:"foto"`
	}
	UtenteChiamante := ps.ByName("utente")
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
			http.Error(w, "Errore durante l'inserimento della foto del messaggio: "+err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = rt.db.CreaMessaggioFotoDB(UtenteChiamante, Conversazione, idFoto)
		if err != nil {
			http.Error(w, " durante l'inserimento del messaggio: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}

	_, err = rt.db.CreaMessaggioTestualeDB(UtenteChiamante, Conversazione, input.Testo)
	if err != nil {
		http.Error(w, "Errore durante l'inserimento del messaggio nella tabella messaggi: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	IDMessaggiostr := ps.ByName("messaggio")
	IDChatstr := ps.ByName("chat")

	IDMessaggio, err := strconv.Atoi(IDMessaggiostr)
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}
	IDChat, err := strconv.Atoi(IDChatstr)
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.EliminaMessaggio(UtenteChiamante, IDMessaggio, IDChat)
	if err != nil {
		http.Error(w, "Errore durante l'eliminazione del messaggio: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Reazione string `json:"reazione"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Errore nel parsing del corpo della richiesta: "+err.Error(), http.StatusBadRequest)
		return
	}

	UtenteChiamante := ps.ByName("utente")
	IDMessaggiostr := ps.ByName("messaggio")

	IDMessaggio, err := strconv.Atoi(IDMessaggiostr)
	if err != nil {
		http.Error(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.AggiungiCommento(UtenteChiamante, IDMessaggio, input.Reazione)
	if err != nil {
		http.Error(w, "Errore durante la scrittura del commento: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	IDMessaggiostr := ps.ByName("messaggio")

	IDMessaggio, err := strconv.Atoi(IDMessaggiostr)
	if err != nil {
		http.Error(w, "ID commento non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.EliminaCommento(UtenteChiamante, IDMessaggio)
	if err != nil {
		http.Error(w, "Errore durante l'eliminazione del commento "+err.Error(), http.StatusInternalServerError)
		return
	}
}
