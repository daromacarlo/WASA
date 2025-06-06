package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {

	// GET
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/liveness", rt.liveness)
	rt.router.GET("/check/:chat", rt.isGroup)
	rt.router.GET("/wasachat/:user/chats/:chat", rt.getConversation)
	rt.router.GET("/wasachat/:user/usercheck/groups/:group", rt.usersInGroup)
	rt.router.GET("/wasachat/:user/chats", rt.getMyConversations)
	rt.router.GET("/wasachat/:user", rt.idFromName)

	// POST
	rt.router.POST("/wasachat/:user/groups", rt.createGroup)
	rt.router.POST("/wasachat", rt.doLogin)
	rt.router.POST("/wasachat/:user/chats/:chat", rt.sendMessage)
	rt.router.POST("/wasachat/:user/chats/:chat/messages/:message", rt.ansMessage)
	rt.router.POST("/wasachat/:user/privateconversation", rt.createPrivateConversation)
	rt.router.POST("/wasachat/:user/forw/:target/messages/:message", rt.forwardMessage)
	rt.router.POST("/wasachat/:user/forwnew/:target/messages/:message", rt.forwardMessageToNewChat)
	rt.router.POST("/wasachat/:user/messages/:message", rt.commentMessage)

	// PUT
	rt.router.PUT("/wasachat/:user/groups/:chat/add", rt.addToGroup)
	rt.router.PUT("/wasachat/:user/groups/:chat/photo", rt.setGroupPhoto)
	rt.router.PUT("/wasachat/:user/groups/:chat/name", rt.setGroupName)
	rt.router.PUT("/wasachat/:user/usersettings/photo", rt.setMyPhoto)
	rt.router.PUT("/wasachat/:user/usersettings/name", rt.setMyUserName)

	// DELETE
	rt.router.DELETE("/wasachat/:user/chats/:chat", rt.leaveGroup)
	rt.router.DELETE("/wasachat/:user/messages/:message", rt.uncommentMessage)
	rt.router.DELETE("/wasachat/:user/chats/:chat/messages/:message", rt.deleteMessage)

	return rt.router
}
