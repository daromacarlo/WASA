package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) RispondiAMessaggioPrivato(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	var input struct {
		Testo        string `json:"testo"`
		PercorsoFoto string `json:"foto"`
	}
	// Recupera i parametri utente e destinatario dalla URL
	UtenteChiamante := ps.ByName("utente")
	Destinatario := ps.ByName("destinatario")
	IdMessaggio, err := strconv.Atoi(ps.ByName("idMessaggio")) // Recupera l'ID del messaggio a cui si risponde
	if err != nil {
		http.Error(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}

	// Decodifica il corpo della richiesta JSON
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	// Verifica che il messaggio abbia almeno del testo o una foto
	if len(input.Testo) == 0 && len(input.PercorsoFoto) == 0 {
		http.Error(w, "Errore, il messaggio deve contenere almeno una foto o del testo", http.StatusBadRequest)
		return
	}

	// Se ci sono sia testo che foto, non è permesso
	if len(input.Testo) > 0 && len(input.PercorsoFoto) > 0 {
		http.Error(w, "Errore, il messaggio non può contenere sia foto che testo", http.StatusBadRequest)
		return
	}

	// Se il messaggio contiene solo testo
	if len(input.Testo) > 0 {
		// Rispondi con un messaggio testuale
		err = rt.db.RispondiMessaggioPrivatoTesto(UtenteChiamante, Destinatario, IdMessaggio, input.Testo)
		if err != nil {
			http.Error(w, "Errore durante la risposta al messaggio: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Risposta testuale inviata con successo"))
		return
	}

	// Se il messaggio contiene solo una foto
	if len(input.PercorsoFoto) > 0 {
		// Leggi il file della foto
		fileFoto, err := ReadImageFile(input.PercorsoFoto)
		if err != nil {
			http.Error(w, "Errore durante la lettura della foto: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Chiamata al database per creare la foto
		idFoto, err := rt.db.CreaFoto(input.PercorsoFoto, fileFoto)
		if err != nil {
			http.Error(w, "Errore durante l'inserimento della foto: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Crea il messaggio con la foto
		err = rt.db.RispondiMessaggioPrivatoFoto(UtenteChiamante, Destinatario, IdMessaggio, idFoto)
		if err != nil {
			http.Error(w, "Errore durante la creazione del messaggio con foto: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Risposta con foto inviata con successo"))
		return
	}
}

/*
curl -X POST http://localhost:3000/wasachat/utente123/risposta/chats/chatprivate/destinario456/789 \
-H "Content-Type: application/json" \
-d '{
  "testo": "Ciao! Come va?",
  "foto": ""
}'
*/

func (rt *_router) RispondiAMessaggioGruppo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	var input struct {
		Testo        string `json:"testo"`
		PercorsoFoto string `json:"foto"`
	}
	UtenteChiamante := ps.ByName("utente")

	IdMessaggio, err := strconv.Atoi(ps.ByName("idMessaggio")) // Recupera l'ID del messaggio a cui si risponde
	if err != nil {
		http.Error(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}
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
		// Chiamata al database per creare la foto profilo
		idFoto, err := rt.db.CreaFoto(input.PercorsoFoto, fileFoto)
		if err != nil {
			http.Error(w, "Errore durante l'inserimento della foto profilo: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = rt.db.RispondiMessaggioGruppoFoto(UtenteChiamante, Conversazione, IdMessaggio, idFoto)
		if err != nil {
			http.Error(w, "Errore durante la creazione del messaggio: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Risposta di successo
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Foto inviata con successo "))
		return
	}

	err = rt.db.RispondiMessaggioGruppoTesto(UtenteChiamante, Conversazione, IdMessaggio, input.Testo)
	if err != nil {
		http.Error(w, "Errore durante la creazione del messaggio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Risposta di successo
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Messaggio di testo inviato con successo "))
}

/*
curl -X POST http://localhost:3000/wasachat/utente123/risposta/chats/gruppi/453252/32143 \
-H "Content-Type: application/json" \
-d '{
  "testo": "Ciao! Come va?",
  "foto": ""
}'
*/
