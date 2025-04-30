package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Funzione per inviare un messaggio
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Testo string `json:"testo"`
		Foto  string `json:"foto"`
	}
	UtenteChiamante := ps.ByName("utente")
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
			CreaErroreJson(w, "Errore durante l'inserimento della foto del messaggio: "+err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = rt.db.CreaMessaggioFotoDB(UtenteChiamante, Conversazione, idFoto)
		if err != nil {
			CreaErroreJson(w, "Errore durante l'inserimento del messaggio con foto: "+err.Error(), http.StatusInternalServerError)
			return
		}
		CreaRispostaJson(w, "Messaggio con foto inviato con successo", http.StatusCreated)
		return
	}

	_, err = rt.db.CreaMessaggioTestualeDB(UtenteChiamante, Conversazione, input.Testo)
	if err != nil {
		CreaErroreJson(w, "Errore durante l'inserimento del messaggio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	CreaRispostaJson(w, "Messaggio inviato con successo", http.StatusCreated)
}

// Funzione per eliminare un messaggio
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	IDMessaggiostr := ps.ByName("messaggio")
	IDChatstr := ps.ByName("chat")

	IDMessaggio, err := strconv.Atoi(IDMessaggiostr)
	if err != nil {
		CreaErroreJson(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}
	IDChat, err := strconv.Atoi(IDChatstr)
	if err != nil {
		CreaErroreJson(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.EliminaMessaggio(UtenteChiamante, IDMessaggio, IDChat)
	if err != nil {
		CreaErroreJson(w, "Errore durante l'eliminazione del messaggio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	CreaRispostaJson(w, "Messaggio eliminato con successo", http.StatusOK)
}

// Funzione per commentare un messaggio
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Reazione string `json:"reazione"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		CreaErroreJson(w, "Errore nel parsing del corpo della richiesta: "+err.Error(), http.StatusBadRequest)
		return
	}

	UtenteChiamante := ps.ByName("utente")
	IDMessaggiostr := ps.ByName("messaggio")

	IDMessaggio, err := strconv.Atoi(IDMessaggiostr)
	if err != nil {
		CreaErroreJson(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.AggiungiCommento(UtenteChiamante, IDMessaggio, input.Reazione)
	if err != nil {
		CreaErroreJson(w, "Errore durante la scrittura del commento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	CreaRispostaJson(w, "Commento aggiunto con successo", http.StatusOK)
}

// Funzione per eliminare un commento
func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	IDMessaggiostr := ps.ByName("messaggio")

	IDMessaggio, err := strconv.Atoi(IDMessaggiostr)
	if err != nil {
		CreaErroreJson(w, "ID commento non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.EliminaCommento(UtenteChiamante, IDMessaggio)
	if err != nil {
		CreaErroreJson(w, "Errore durante l'eliminazione del commento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	CreaRispostaJson(w, "Commento eliminato con successo", http.StatusOK)
}
