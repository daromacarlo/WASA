package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	nick := r.URL.Query().Get("nick")

	//formattazione dati in ingresso
	w.Header().Set("Content-Type", "application/json")

	// La variabile esistenza serve a controllare se esiste un utente con il nome 'nick' nel database
	esistenza, err := rt.db.SearchUser(nick)

	// Se c'è un errore nella ricerca dell'utente, restituisci un errore
	if err != nil {
		http.Error(w, fmt.Sprintf("Errore nella ricerca dell'utente: %s", err), http.StatusInternalServerError)
		return
	}

	// Se l'utente esiste, restituisci un messaggio di bentornato all'utente, al server chi ha effettuato il login
	if esistenza {
		fmt.Printf(`login effettuato da '%s' `, nick)
		fmt.Fprintf(w, `Bentornato %s!"`, nick)
		return
	}

	// Se l'utente non esiste, viene creato
	err = rt.db.CreateUser(nick)
	if err != nil {
		http.Error(w, fmt.Sprintf("Errore nella creazione dell'utente: %s", err), http.StatusInternalServerError)
		return
	}

	// Se l'utente è stato creato, restituisci un messaggio di benvenuto all'utente, al server viene scritto chi si è registrato
	fmt.Printf(`registrazione effettuata da '%s' `, nick)
	fmt.Fprintf(w, `Benvenuto, ciao %s!`, nick)
}
