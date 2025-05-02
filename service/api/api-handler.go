package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/liveness", rt.liveness)
	rt.router.GET("/check/:chat", rt.IsGroup)
	rt.router.POST("/wasachat/:utente/gruppi", rt.createGroup)                                    // presente in API
	rt.router.PUT("/wasachat", rt.register)                                                       // presente
	rt.router.POST("/wasachat", rt.doLogin)                                                       // presente
	rt.router.POST("/wasachat/:utente/chats/:chat", rt.sendMessage)                               // presente
	rt.router.PUT("/wasachat/:utente/chats/gruppi/:chat/aggiungi", rt.addToGroup)                 // presente
	rt.router.GET("/wasachat/:utente/chats/:chat", rt.getConversation)                            // presente
	rt.router.DELETE("/wasachat/:utente/chats/:chat", rt.leaveGroup)                              // ptesente
	rt.router.DELETE("/wasachat/:utente/chats/:chat/messaggi/:messaggio", rt.deleteMessage)       // ptesente
	rt.router.PUT("/wasachat/:utente/gruppi/:chat/foto", rt.setGroupPhoto)                        // ptesente
	rt.router.PUT("/wasachat/:utente/gruppi/:chat/nome", rt.setGroupName)                         // presente
	rt.router.PUT("/wasachat/:utente/foto", rt.setMyPhoto)                                        // presente
	rt.router.PUT("/wasachat/:utente/nome", rt.setMyUserName)                                     // presente
	rt.router.GET("/wasachat/:utente/utenti/gruppi/:gruppo", rt.usersInGroup)                     // presente
	rt.router.GET("/wasachat/:utente/chats", rt.getMyConversations)                               // presente
	rt.router.DELETE("/wasachat/:utente/messaggi/:messaggio", rt.deleteComment)                   // presente
	rt.router.POST("/wasachat/:utente/messaggi/:messaggio", rt.commentMessage)                    // presente
	rt.router.POST("/wasachat/:utente/risposta/chats/:chat/messaggi/:messaggio", rt.ansMessage)   // presente
	rt.router.POST("/wasachat/:utente/inoltro/:nuovachat/messaggi/:messaggio", rt.forwardMessage) // presente
	rt.router.POST("/wasachat/:utente/conversazioniprivate", rt.createPrivateConversation)        // presente
	return rt.router
}
