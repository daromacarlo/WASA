package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
// interfaccia utente ad alto livello
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/liveness", rt.liveness)
	rt.router.POST("/wasachat/:utente/gruppi", rt.CreaGruppo)
	rt.router.POST("/wasachat", rt.Registrazione)
	rt.router.POST("/wasachat/:utente/chats/:chat", rt.InviaMessaggio)
	rt.router.PUT("/wasachat/:utente/chats/gruppi/:chat/aggiungi", rt.AggiungiAGruppo)
	rt.router.GET("/wasachat/:utente/chats/:chat", rt.GetConversazione)
	rt.router.DELETE("/wasachat/:utente/chats/:chat", rt.LasciaGruppo)
	rt.router.DELETE("/wasachat/:utente/chats/:chat/messaggi/:messaggio", rt.EliminaMessaggio)
	rt.router.PUT("/wasachat/:utente/gruppi/:chat/foto", rt.ImpostaFotoGruppo)
	rt.router.PUT("/wasachat/:utente/gruppi/:chat/nome", rt.ImpostaNomeGruppo)
	rt.router.PUT("/wasachat/:utente/foto", rt.ImpostaFotoProfilo)
	rt.router.PUT("/wasachat/:utente/nome", rt.ImpostaNome)
	rt.router.GET("/wasachat/:utente/utenti", rt.VediProfili) //funzione di test
	rt.router.GET("/wasachat/:utente/chats", rt.GetConversazioni)
	rt.router.DELETE("/wasachat/:utente/messaggi/:commento", rt.EliminaCommento)
	rt.router.POST("/wasachat/:utente/messaggi/:messaggio", rt.AggiungiCommento)
	rt.router.POST("/wasachat/:utente/risposta/chats/:chat/:idMessaggio", rt.RispondiAMessaggio)
	rt.router.POST("/wasachat/:utente/inoltra/:nuovachat/:messaggio", rt.InoltraMessaggio)
	rt.router.POST("/wasachat/:utente/conversazioniprivate", rt.CreaConversazionePrivata)
	return rt.router
}
