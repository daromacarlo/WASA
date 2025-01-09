package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) VediProfili(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Ottieni il nickname dall'URL
	chiamante := ps.ByName("utente")
	lista, err := rt.db.VediProfili(chiamante)
	if err != nil {
		http.Error(w, "Errore durante il recupero degli utenti: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(lista) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Nessun utente trovato"}`))
		return
	}

	// Converte la lista in JSON e restituisce la risposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(lista); err != nil {
		http.Error(w, "Errore durante la codifica della risposta in JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) ImpostaFotoProfilo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	UtenteChiamante := ps.ByName("utente")

	var input struct {
		PercorsoFoto string `json:"foto"`
	}

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

	// Imposta la foto profilo dell'utente
	err = rt.db.ImpostaFotoProfilo(UtenteChiamante, idFoto)
	if err != nil {
		http.Error(w, "Errore durante l'aggiornamento della foto profilo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Risposta di successo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Foto profilo aggiornata con successo",
	})
}

func (rt *_router) ImpostaNome(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	UtenteChiamante := ps.ByName("utente")

	var input struct {
		Nome string `json:"nome"` // Il campo deve essere maiuscolo per essere esportato
	}

	// Decodifica il corpo della richiesta
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	// Controlla se il nome è già in uso
	existe, err := rt.db.EsistenzaUtente(input.Nome)
	if err != nil {
		http.Error(w, "Errore durante il controllo dell'esistenza del nome utente: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if existe {
		http.Error(w, "Il nome utente è già in uso", http.StatusConflict)
		return
	}

	// Chiama la funzione per impostare il nome dell'utente
	err = rt.db.ImpostaNome(UtenteChiamante, input.Nome)
	if err != nil {
		http.Error(w, "Errore durante l'aggiornamento del nome utente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Risposta di successo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Nome utente aggiornato con successo",
	})
}
