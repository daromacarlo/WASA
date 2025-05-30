package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	callingUser := ps.ByName("user")
	chatStr := ps.ByName("chat")

	chat, err := strconv.Atoi(chatStr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "The 'chat' parameter must be a valid integer", http.StatusBadRequest)
		return
	}

	list, errorCode, err := rt.db.GetChat(callingUser, chat)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while retrieving the conversation: "+err.Error(), errorCode)
		return
	}

	if len(list) == 0 {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		CreateJsonError(w, "Error while encoding the JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
