package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) CreaGruppo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	var input struct {
		Nome         string `json:"nome"`
		PercorsoFoto string `json:"foto"`
	}
	UtenteChiamante := ps.ByName("utente")

	// Decodifica il corpo della richiesta
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	if len(input.Nome) == 0 {
		http.Error(w, "Il nome è obbligatorio", http.StatusBadRequest)
		return
	}

	// Validazione: la foto è obbligatoria
	if len(input.PercorsoFoto) == 0 {
		http.Error(w, "La foto è obbligatoria", http.StatusBadRequest)
		return
	}

	var fileFoto []byte

	if len(input.PercorsoFoto) > 0 {
		fileFoto, err = ReadImageFile(input.PercorsoFoto)
		if err != nil {
			http.Error(w, "Errore durante la lettura della foto: "+err.Error(), http.StatusBadRequest)
			return
		}
	}
	// Chiamata al database per creare la foto profilo
	idFoto, err := rt.db.CreaFoto(input.PercorsoFoto, fileFoto)
	if err != nil {
		http.Error(w, "Errore durante l'inserimento della foto profilo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = rt.db.CreaGruppoDB(UtenteChiamante, input.Nome, idFoto)
	if err != nil {
		http.Error(w, "Errore durante la creazione dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Risposta di successo
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Gruppo creato con successo "))
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

func (rt *_router) AggiungiAGruppo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	var input struct {
		Utente string `json:"utente"`
	}
	UtenteChiamante := ps.ByName("utente")
	idConversazioneStr := ps.ByName("chat")

	// Decodifica il corpo della richiesta
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	if len(input.Utente) == 0 {
		http.Error(w, "Il nome è obbligatorio", http.StatusBadRequest)
		return
	}

	// Converti idConversazione in un intero
	idConversazione, err := strconv.Atoi(idConversazioneStr)
	if err != nil {
		http.Error(w, "ID della conversazione non valido", http.StatusBadRequest)
		return
	}

	// Chiamata alla funzione AggiungiAGruppoDB
	err = rt.db.AggiungiAGruppoDB(idConversazione, UtenteChiamante, input.Utente)
	if err != nil {
		http.Error(w, "Errore durante l'aggiunta dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Risposta di successo
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Utente aggiunto con successo "))
}

// test
/*
curl -X PUT http://localhost:3000/wasachat/:utente/chats/gruppi/:chat/aggiungi \
-H "Content-Type: application/json" \
-d '{
  "utente": "Luigi"
}'
*/

func (rt *_router) LasciaGruppo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	idConversazioneStr := ps.ByName("chat")

	// Converti idConversazione in un intero
	idConversazione, err := strconv.Atoi(idConversazioneStr)
	if err != nil {
		http.Error(w, "ID della conversazione non valido", http.StatusBadRequest)
		return
	}

	// Chiamata alla funzione AggiungiAGruppoDB
	err = rt.db.LasciaGruppo(idConversazione, UtenteChiamante)
	if err != nil {
		http.Error(w, "Errore durante la rimozione dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Utente rimosso con successo "))
}

func (rt *_router) ImpostaFotoGruppo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	idConversazioneStr := ps.ByName("chat")
	var input struct {
		PercorsoFoto string `json:"foto"`
	}
	UtenteChiamante := ps.ByName("utente")

	// Decodifica il corpo della richiesta
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	// Validazione: la foto è obbligatoria
	if len(input.PercorsoFoto) == 0 {
		http.Error(w, "La foto è obbligatoria", http.StatusBadRequest)
		return
	}

	var fileFoto []byte

	if len(input.PercorsoFoto) > 0 {
		fileFoto, err = ReadImageFile(input.PercorsoFoto)
		if err != nil {
			http.Error(w, "Errore durante la lettura della foto: "+err.Error(), http.StatusBadRequest)
			return
		}
	}
	// Chiamata al database per creare la foto profilo
	idFoto, err := rt.db.CreaFoto(input.PercorsoFoto, fileFoto)
	if err != nil {
		http.Error(w, "Errore durante l'inserimento della foto profilo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Converti idConversazione in un intero
	idConversazione, err := strconv.Atoi(idConversazioneStr)
	if err != nil {
		http.Error(w, "ID della conversazione non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.ImpostaFotoGruppo(UtenteChiamante, idFoto, idConversazione)
	if err != nil {
		http.Error(w, "Errore durante la creazione dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Gruppo modificato con successo "))
}
func (rt *_router) ImpostaNomeGruppo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	UtenteChiamante := ps.ByName("utente")
	idConversazioneStr := ps.ByName("chat")

	var input struct {
		Nome string `json:"nome"` // Il campo deve essere maiuscolo per essere esportato
	}

	// Decodifica il corpo della richiesta
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	// Converti idConversazione in un intero
	idConversazione, err := strconv.Atoi(idConversazioneStr)
	if err != nil {
		http.Error(w, "ID della conversazione non valido", http.StatusBadRequest)
		return
	}

	// Chiama la funzione per impostare il nome del gruppo
	err = rt.db.ImpostaNomeGruppo(UtenteChiamante, input.Nome, idConversazione)
	if err != nil {
		http.Error(w, "Errore durante l'aggiornamento del nome del gruppo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Risposta di successo
	w.Write([]byte("Gruppo modificato con successo "))
}
