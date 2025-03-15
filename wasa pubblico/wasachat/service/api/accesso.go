package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// funzione che serve a registrare un nuovo utente
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Struttura per ricevere i dati dal body
	var input struct {
		Nickname     string `json:"nickname"`
		PercorsoFoto string `json:"foto"`
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
curl -X POST http://localhost:3000/wasachat \
-H "Content-Type: application/json" \
-d '{
  "nickname": "Mario",
  "foto": "/home/carlo/Scrivania/wasachat/immagini/prova.png"
}'
*/
