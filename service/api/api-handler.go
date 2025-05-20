package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/liveness", rt.liveness)
	rt.router.GET("/check/:chat", rt.isGroup)                                                                  // presente in API
	rt.router.POST("/wasachat/:utente/gruppi", rt.createGroup)                                                 // presente in API
	rt.router.PUT("/wasachat", rt.register)                                                                    // presente in API
	rt.router.POST("/wasachat", rt.doLogin)                                                                    // presente in API
	rt.router.POST("/wasachat/:utente/chats/:chat", rt.sendMessage)                                            // presente in API
	rt.router.PUT("/wasachat/:utente/chats/gruppi/:chat/aggiungi", rt.addToGroup)                              // presente in API
	rt.router.GET("/wasachat/:utente/chats/:chat", rt.getConversation)                                         // presente in API
	rt.router.DELETE("/wasachat/:utente/chats/:chat", rt.leaveGroup)                                           // ptesente in API
	rt.router.DELETE("/wasachat/:utente/chats/:chat/messaggi/:messaggio", rt.deleteMessage)                    // ptesente in API
	rt.router.PUT("/wasachat/:utente/gruppi/:chat/foto", rt.setGroupPhoto)                                     // ptesente in API
	rt.router.PUT("/wasachat/:utente/gruppi/:chat/nome", rt.setGroupName)                                      // presente in API
	rt.router.PUT("/wasachat/:utente/foto", rt.setMyPhoto)                                                     // presente in API
	rt.router.PUT("/wasachat/:utente/nome", rt.setMyUserName)                                                  // presente in API
	rt.router.GET("/wasachat/:utente/utenti/gruppi/:gruppo", rt.usersInGroup)                                  // presente in API
	rt.router.GET("/wasachat/:utente/chats", rt.getMyConversations)                                            // presente in API
	rt.router.DELETE("/wasachat/:utente/messaggi/:messaggio", rt.deleteComment)                                // presente in API
	rt.router.POST("/wasachat/:utente/messaggi/:messaggio", rt.commentMessage)                                 // presente in API
	rt.router.POST("/wasachat/:utente/risposta/chats/:chat/messaggi/:messaggio", rt.ansMessage)                // presente in API
	rt.router.POST("/wasachat/:utente/inoltro/:nuovachat/messaggi/:messaggio", rt.forwardMessage)              // presente in API
	rt.router.POST("/inoltro/:utente/a/:destinatario/inoltro/messaggi/:messaggio", rt.forwardMessageToNewChat) // presente in API
	rt.router.POST("/wasachat/:utente/conversazioniprivate", rt.createPrivateConversation)                     // presente in API
	return rt.router
}
