package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createPrivateConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		User string `json:"user"`
	}

	CallingUser := ps.ByName("user")

	err := json.NewDecoder(r.Body).Decode(&input)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if len(input.User) == 0 {
		CreateJsonError(w, "The target username is required", http.StatusBadRequest)
		return
	}

	if len(CallingUser) == 0 {
		CreateJsonError(w, "Calling user not properly specified", http.StatusBadRequest)
		return
	}

	_, errorCode, err := rt.db.CreatePrivateChatDB(CallingUser, input.User)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while creating the private conversation: "+err.Error(), errorCode)
		return
	}

	CreateJsonResponse(w, "Private conversation successfully created", http.StatusOK)
}
