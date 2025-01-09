package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	//route speciali
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/liveness", rt.liveness)

	//funzione di login e registrazione
	rt.router.POST("/wasachat", rt.Login)

	//fuzione per cambiare nickname (vedi esempio 1 di example-list per esempio di utilizzo)
	rt.router.PUT("/wasachat/:nickname/impostazioni/impostazioninickname", rt.SetMyUsername)

	//funzione per cambiare o aggiungere la propria foto profilo
	rt.router.PUT("/wasachat/:nickname/impostazioni/impostazioniimmagineprofilo", rt.SetMyPhoto)

	//funzione per inviare messaggi in una chat privata
	rt.router.POST("/wasachat/:nickname", rt.SendPrivateMessage)

	rt.router.GET("/wasachat/:nickname/chats", rt.GetMyConversation)

	//funzioni di test
	rt.router.GET("/wasachat/:nickname/test/user", rt.GetAllUsers)
	rt.router.GET("/wasachat/:nickname/test/messaggi", rt.GetAllMessaggi)
	return rt.router
}
