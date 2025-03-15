package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	var input struct {
		Testo        string `json:"testo"`
		PercorsoFoto string `json:"foto"`
	}
	UtenteChiamante := ps.ByName("utente")
	ConversazioneStr := ps.ByName("chat")

	// Converti Conversazione in intero
	Conversazione, err := strconv.Atoi(ConversazioneStr)
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	// Decodifica il corpo della richiesta
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}
	// Controlliamo che la richiesta sia valida
	if len(input.Testo) == 0 && len(input.PercorsoFoto) == 0 {
		http.Error(w, "Errore, il messaggio deve contenere almeno una foto o del testo", http.StatusBadRequest)
		return
	}

	if len(input.Testo) > 0 && len(input.PercorsoFoto) > 0 {
		http.Error(w, "Errore, il messaggio non può contenere sia foto che testo", http.StatusBadRequest)
		return
	}

	if len(input.Testo) == 0 {
		var fileFoto []byte

		if len(input.PercorsoFoto) > 0 {
			fileFoto, err = ReadImageFile(input.PercorsoFoto)
			if err != nil {
				http.Error(w, "Errore durante la lettura della foto: "+err.Error(), http.StatusBadRequest)
				return
			}
		}
		// Chiamata al database per creare la foto
		idFoto, err := rt.db.CreaFoto(input.PercorsoFoto, fileFoto)
		if err != nil {
			http.Error(w, "Errore durante l'inserimento della foto profilo: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = rt.db.CreaMessaggioFotoDB(UtenteChiamante, Conversazione, idFoto)
		if err != nil {
			http.Error(w, "Errore durante la creazione del messaggio: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Risposta di successo
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Foto inviata con successo "))
		return
	}

	err = rt.db.CreaMessaggioTestualeDB(UtenteChiamante, Conversazione, input.Testo)
	if err != nil {
		http.Error(w, "Errore durante la creazione del messaggio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Risposta di successo
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Messaggio di testo inviato con successo "))
}

// test
/*
curl -X POST http://localhost:3000/wasachat/:utente/chats/:chat \
-H "Content-Type: application/json" \
-d '{
  "testo": "ciao come va?",
  "foto": ""
}'

curl -X POST http://localhost:3000/wasachat/:utente/chats/:chat \
-H "Content-Type: application/json" \
-d '{
  "testo": "",
  "foto": "/home/carlo/Scrivania/wasachat/immagini/prova.png"
}'
*/

func (rt *_router) EliminaMessaggio(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	IDMessaggiostr := ps.ByName("messaggio")
	IDChatstr := ps.ByName("chat")

	// Converti Conversazione in intero
	IDMessaggio, err := strconv.Atoi(IDMessaggiostr)
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	// Converti Conversazione in intero
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

	w.Write([]byte("Messaggio eliminato con successo "))
}

func (rt *_router) commentMessagge(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Reazione string `json:"reazione"`
	}

	//decodifica corpo del json
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Errore nel parsing del corpo della richiesta: "+err.Error(), http.StatusBadRequest)
		return
	}

	UtenteChiamante := ps.ByName("utente")
	IDMessaggiostr := ps.ByName("messaggio")

	// Converti IDMessaggio in intero
	IDMessaggio, err := strconv.Atoi(IDMessaggiostr)
	if err != nil {
		http.Error(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}

	// Aggiungi il commento al database
	err = rt.db.AggiungiCommento(UtenteChiamante, IDMessaggio, input.Reazione)
	if err != nil {
		http.Error(w, "Errore durante la scrittura del commento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Risposta positiva
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Commento aggiunto con successo"))
}

func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	IDCommentostr := ps.ByName("commento")

	// Converti Conversazione in intero
	IDCommento, err := strconv.Atoi(IDCommentostr)
	if err != nil {
		http.Error(w, "ID commento non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.EliminaCommento(UtenteChiamante, IDCommento)
	if err != nil {
		http.Error(w, "Errore durante l'eliminazione del commento "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Commento eliminato con successo "))
}
