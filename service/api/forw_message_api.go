package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token := r.Header.Get("Authorization")
	exist, err := rt.VerifyToken(token)
	if exist == false {
		CreateJsonError(w, "Error: Token not valid", 401)
		return
	}
	if err != nil {
		CreateJsonError(w, "Error checking token", http.StatusInternalServerError)
		return
	}

	User := ps.ByName("user")
	NewConvStr := ps.ByName("target")
	IdMessageStr := ps.ByName("message")

	IdMessage, err := strconv.Atoi(IdMessageStr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	NewConv, err := strconv.Atoi(NewConvStr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid new conversation ID", http.StatusBadRequest)
		return
	}

	errorCode, err := rt.db.ForwardMessage(User, NewConv, IdMessage)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while forwarding the message: "+err.Error(), errorCode)
		return
	}

	CreateJsonResponse(w, "Message successfully forwarded", http.StatusOK)
	return
}

func (rt *_router) forwardMessageToNewChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token := r.Header.Get("Authorization")
	exist, err := rt.VerifyToken(token)
	if exist == false {
		CreateJsonError(w, "Error: Token not valid", 401)
		return
	}
	if err != nil {
		CreateJsonError(w, "Error checking token", http.StatusInternalServerError)
		return
	}

	User := ps.ByName("user")
	Target := ps.ByName("target")
	IdMessageStr := ps.ByName("message")

	messageId, err := strconv.Atoi(IdMessageStr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	errorCode, err := rt.db.ForwardMessageToNewChat(User, Target, messageId)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while forwarding the message: "+err.Error(), errorCode)
		return
	}

	CreateJsonResponse(w, "Message successfully forwarded", http.StatusOK)
	return
}
