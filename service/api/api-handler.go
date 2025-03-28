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
	rt.router.GET("/check/:chat", rt.IsGroup)
	rt.router.POST("/wasachat/:utente/gruppi", rt.CreaGruppo)
	rt.router.PUT("/wasachat", rt.registrare)
	rt.router.POST("/wasachat", rt.doLogin)
	rt.router.POST("/wasachat/:utente/chats/:chat", rt.sendMessage)
	rt.router.PUT("/wasachat/:utente/chats/gruppi/:chat/aggiungi", rt.addToGroup)
	rt.router.GET("/wasachat/:utente/chats/:chat", rt.getConversation)
	rt.router.DELETE("/wasachat/:utente/chats/:chat", rt.leaveGroup)
	rt.router.DELETE("/wasachat/:utente/chats/:chat/messaggi/:messaggio", rt.EliminaMessaggio)
	rt.router.PUT("/wasachat/:utente/gruppi/:chat/foto", rt.setGroupPhoto)
	rt.router.PUT("/wasachat/:utente/gruppi/:chat/nome", rt.setGroupName)
	rt.router.PUT("/wasachat/:utente/foto", rt.setMyPhoto)
	rt.router.PUT("/wasachat/:utente/nome", rt.setMyUserName)
	rt.router.GET("/wasachat/:utente/utenti", rt.VediProfili) //funzione di test
	rt.router.GET("/wasachat/:utente/chats", rt.getMyConversation)
	rt.router.DELETE("/wasachat/:utente/messaggi/:commento", rt.deleteComment)
	rt.router.POST("/wasachat/:utente/messaggi/:messaggio", rt.commentMessagge)
	rt.router.POST("/wasachat/:utente/risposta/chats/:chat/:idMessaggio", rt.RispondiAMessaggio)
	rt.router.POST("/wasachat/:utente/inoltra/:nuovachat/:messaggio", rt.forwardMessagge)
	rt.router.POST("/wasachat/:utente/conversazioniprivate", rt.CreaConversazionePrivata)
	return rt.router
}
