package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) usersInGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := ps.ByName("user")
	groupstr := ps.ByName("group")
	group, err := strconv.Atoi(groupstr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error during the name conversion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	list, errorCode, err := rt.db.UsersInGroup(user, group)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error during retriving users: "+err.Error(), errorCode)
		return
	}
	if len(list) == 0 {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(list); err != nil {
		CreateJsonError(w, "Error during the creation of JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := ps.ByName("user")
	var input struct {
		Photo string `json:"photo"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error: not valid request", http.StatusBadRequest)
		return
	}

	if len(input.Photo) == 0 {
		CreateJsonError(w, "photo is requested", http.StatusBadRequest)
		return
	}

	idPhoto, err := rt.db.CreaFoto(input.Photo)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error during uploading your photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = rt.db.SetPhotoDB(user, idPhoto)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error during updating your photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	CreateJsonResponse(w, "Photo updated ", http.StatusOK)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("user")
	var input struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error: not valid request", http.StatusBadRequest)
		return
	}

	errorCode, err := rt.db.SetNameDB(UtenteChiamante, input.Name)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error during updating your name: "+err.Error(), errorCode)
		return
	}

	CreateJsonResponse(w, "Name updated ", http.StatusOK)
}

func (rt *_router) idFromName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	nickname := ps.ByName("user")

	id, errorCode, err := rt.db.IDFromNICK(nickname)
	if err != nil {
		CreateJsonError(w, "Error during retrieving the id: "+err.Error(), errorCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		ID int `json:"id"`
	}{ID: id})
	if err != nil {
		http.Error(w, "Failed to encode JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
