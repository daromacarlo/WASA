package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Funzione che serve a creare un gruppo dato un nome e una foto
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Nome string `json:"nome"`
		Foto string `json:"foto"`
	}
	UtenteChiamante := ps.ByName("utente")

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		CreaErroreJson(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	if len(input.Nome) == 0 {
		CreaErroreJson(w, "Il nome è obbligatorio", http.StatusBadRequest)
		return
	}
	if len(input.Foto) == 0 {
		CreaErroreJson(w, "La foto è obbligatoria", http.StatusBadRequest)
		return
	}

	idFoto, err := rt.db.CreaFoto(input.Foto)
	if err != nil {
		CreaErroreJson(w, "Errore durante l'inserimento della foto del gruppo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	codiceErrore, err := rt.db.CreaGruppoDB(UtenteChiamante, input.Nome, idFoto)
	if err != nil {
		CreaErroreJson(w, "Errore durante la creazione del gruppo: "+err.Error(), codiceErrore)
		return
	}
}

// test
/*
curl -X POST http://localhost:3000/wasachat/:utente/gruppi \
-H "Content-Type: application/json" \
-d '{
  "nome": "Gruppo1",
  "foto": "/home/carlo/Scrivania/wasachat/immagini/prova.png"
}'
*/

// funzione che serve ad aggiungere un utente ad un gruppo, se l'utente è già presente nel gruppo non verrà aggiunto
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Utente string `json:"utente_da_aggiungere"`
	}
	UtenteChiamante := ps.ByName("utente")
	idConversazioneStr := ps.ByName("chat")

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		CreaErroreJson(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	if len(input.Utente) == 0 {
		CreaErroreJson(w, "Il nome è obbligatorio", http.StatusBadRequest)
		return
	}

	idConversazione, err := strconv.Atoi(idConversazioneStr)
	if err != nil {
		CreaErroreJson(w, "ID della conversazione non valido", http.StatusBadRequest)
		return
	}

	codiceErrore, err := rt.db.AggiungiAGruppoDB(idConversazione, UtenteChiamante, input.Utente)
	if err != nil {
		CreaErroreJson(w, "Errore durante l'aggiunta dell'utente: "+err.Error(), codiceErrore)
		return
	}
}

// test
/*
curl -X PUT http://localhost:3000/wasachat/:utente/chats/gruppi/:chat/aggiungi \
-H "Content-Type: application/json" \
-d '{
  "utente": "Luigi"
}'
*/

// Funzione che serve a lasciare un gruppo, l'utente deve essere presente nel gruppo
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	idConversazioneStr := ps.ByName("chat")
	idConversazione, err := strconv.Atoi(idConversazioneStr)
	if err != nil {
		CreaErroreJson(w, "ID della conversazione non valido", http.StatusBadRequest)
		return
	}
	codiceErrore, err := rt.db.LasciaGruppo(idConversazione, UtenteChiamante)
	if err != nil {
		CreaErroreJson(w, "Errore durante la rimozione dell'utente: "+err.Error(), codiceErrore)
		return
	}
}

// Funzione per impostare una nuova foto al gruppo
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idConversazioneStr := ps.ByName("chat")
	var input struct {
		Foto string `json:"foto"`
	}
	UtenteChiamante := ps.ByName("utente")
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		CreaErroreJson(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}
	if len(input.Foto) == 0 {
		CreaErroreJson(w, "La foto è obbligatoria", http.StatusBadRequest)
		return
	}

	idFoto, err := rt.db.CreaFoto(input.Foto)
	if err != nil {
		CreaErroreJson(w, "Errore durante l'inserimento della foto del gruppo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	idConversazione, err := strconv.Atoi(idConversazioneStr)
	if err != nil {
		CreaErroreJson(w, "ID della conversazione non valido", http.StatusBadRequest)
		return
	}

	codiceErrore, err := rt.db.ImpostaFotoGruppo(UtenteChiamante, idFoto, idConversazione)
	if err != nil {
		CreaErroreJson(w, "Errore durante l'aggiornamento della foto del gruppo: "+err.Error(), codiceErrore)
		return
	}
}

// Funzione per impostare un nome ad un gruppo
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	idConversazioneStr := ps.ByName("chat")

	var input struct {
		Nome string `json:"nome"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		CreaErroreJson(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	idConversazione, err := strconv.Atoi(idConversazioneStr)
	if err != nil {
		CreaErroreJson(w, "ID della conversazione non valido", http.StatusBadRequest)
		return
	}

	codiceErrore, err := rt.db.ImpostaNomeGruppo(UtenteChiamante, input.Nome, idConversazione)
	if err != nil {
		CreaErroreJson(w, "Errore durante l'aggiornamento del nome del gruppo: "+err.Error(), codiceErrore)
		return
	}
}

// IsGroup verifica se la chat è un gruppo
func (rt *_router) isGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	chatParam := ps.ByName("chat")

	chatID, err := strconv.Atoi(chatParam)
	if err != nil {
		CreaErroreJson(w, "ID chat non valido", http.StatusBadRequest)
		return
	}

	idGruppo, codiceErrore, err := rt.db.CercaConversazioneGruppo(chatID)
	if err != nil {
		CreaErroreJson(w, "Errore interno del server: "+err.Error(), codiceErrore)
		return
	}
	// Se è un gruppo invio l'id gruppo
	risposta := struct {
		IsGroup bool `json:"is_group"`
		GroupID int  `json:"group_id,omitempty"`
	}{
		IsGroup: idGruppo > 0,
	}

	if idGruppo > 0 {
		risposta.GroupID = idGruppo
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(risposta); err != nil {
		CreaErroreJson(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}
}
