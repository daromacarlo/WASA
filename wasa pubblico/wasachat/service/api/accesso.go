package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// funzione che serve a registrare un nuovo utente
func (rt *_router) registrare(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	var input struct {
		Nickname string `json:"nickname"`
		Foto     string `json:"foto"`
	}

	// Decodifica il corpo della richiesta
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	// Validazione: il nickname è obbligatorio
	if len(input.Nickname) == 0 {
		http.Error(w, "Il nickname è obbligatorio", http.StatusBadRequest)
		return
	}

	// Validazione: la foto è obbligatoria
	if len(input.Foto) == 0 {
		http.Error(w, "La foto è obbligatoria", http.StatusBadRequest)
		return
	}

	// Chiamata al database per creare la foto profilo
	idFoto, err := rt.db.CreaFoto(input.Foto)
	if err != nil {
		http.Error(w, "Errore durante l'inserimento della foto profilo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Chiamata al database per creare l'utente, passo l'id dell'ultima foto creata
	err = rt.db.CreaUtente(input.Nickname, idFoto)
	if err != nil {
		http.Error(w, "Errore durante la creazione dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Utente creato con successo "))
}

// test
/*
curl -X PUT http://localhost:3000/wasachat \
-H "Content-Type: application/json" \
-d '{
  "nickname": "Mario",
  "foto": "/home/carlo/Scrivania/wasachat/immagini/prova.png"
}'
*/
// funzione che serve a fare il login di un utente
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	var input struct {
		Nickname string `json:"nickname"`
	}

	// Decodifica il corpo della richiesta
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	// Validazione: il nickname è obbligatorio
	if len(input.Nickname) == 0 {
		http.Error(w, "Il nickname è obbligatorio", http.StatusBadRequest)
		return
	}

	// Chiamata al database per verificare se l'utente esiste
	utente, err := rt.db.Login(input.Nickname)
	if err != nil {
		http.Error(w, "Errore durante la verifica dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if utente == "" {
		// Utente non trovato
		http.Error(w, "Credenziali non valide", http.StatusUnauthorized)
		return
	}

	// Se l'utente esiste, rispondi con un messaggio di successo
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login effettuato con successo"))
}
