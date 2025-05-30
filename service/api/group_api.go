package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Name  string `json:"name"`
		Photo string `json:"photo"`
	}
	CallingUser := ps.ByName("user")

	err := json.NewDecoder(r.Body).Decode(&input)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if len(input.Name) == 0 {
		CreateJsonError(w, "Group name is required", http.StatusBadRequest)
		return
	}
	if len(input.Photo) == 0 {
		CreateJsonError(w, "Group photo is required", http.StatusBadRequest)
		return
	}

	idPhoto, err := rt.db.CreaFoto(input.Photo)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while uploading group photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	errorCode, err := rt.db.CreateGroupDB(CallingUser, input.Name, idPhoto)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while creating the group: "+err.Error(), errorCode)
		return
	}
	CreateJsonResponse(w, "Group created", http.StatusOK)
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		User string `json:"user_to_add"`
	}
	CallingUser := ps.ByName("user")
	chatIdStr := ps.ByName("chat")

	err := json.NewDecoder(r.Body).Decode(&input)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if len(input.User) == 0 {
		CreateJsonError(w, "Username to add is required", http.StatusBadRequest)
		return
	}

	chatId, err := strconv.Atoi(chatIdStr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	errorCode, err := rt.db.AddToGroupDB(chatId, CallingUser, input.User)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while adding the user: "+err.Error(), errorCode)
		return
	}
	CreateJsonResponse(w, "User successfully added to the group", http.StatusOK)
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	CallingUser := ps.ByName("user")
	chatIdStr := ps.ByName("chat")
	chatId, err := strconv.Atoi(chatIdStr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}
	errorCode, err := rt.db.LeaveGroup(chatId, CallingUser)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while leaving the group: "+err.Error(), errorCode)
		return
	}
	CreateJsonResponse(w, "Successfully left the group", http.StatusOK)
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	chatIdStr := ps.ByName("chat")
	var input struct {
		Photo string `json:"photo"`
	}
	CallingUser := ps.ByName("user")
	err := json.NewDecoder(r.Body).Decode(&input)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	if len(input.Photo) == 0 {
		CreateJsonError(w, "Group photo is required", http.StatusBadRequest)
		return
	}

	idPhoto, err := rt.db.CreaFoto(input.Photo)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while uploading group photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	idConv, err := strconv.Atoi(chatIdStr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	errorCode, err := rt.db.SetGroupPhotoDB(CallingUser, idPhoto, idConv)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while updating the group photo: "+err.Error(), errorCode)
		return
	}
	CreateJsonResponse(w, "Group photo successfully updated", http.StatusOK)
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	CallingUser := ps.ByName("user")
	idConvStr := ps.ByName("chat")

	var input struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	idConv, err := strconv.Atoi(idConvStr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	errorCode, err := rt.db.SetGroupNameDB(CallingUser, input.Name, idConv)
	if err != nil {
		CreateJsonError(w, "Error while updating the group name: "+err.Error(), errorCode)
		return
	}
	CreateJsonResponse(w, "Group name successfully updated", http.StatusOK)
}

func (rt *_router) isGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	chatParam := ps.ByName("chat")

	chatID, err := strconv.Atoi(chatParam)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	idGroup, errorCode, err := rt.db.SearchGroup(chatID)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Internal server error: "+err.Error(), errorCode)
		return
	}

	response := struct {
		IsGroup bool `json:"is_group"`
		GroupID int  `json:"group_id,omitempty"`
	}{
		IsGroup: idGroup > 0,
	}

	if idGroup > 0 {
		response.GroupID = idGroup
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		CreateJsonError(w, "Error encoding the response", http.StatusInternalServerError)
		return
	}
}
