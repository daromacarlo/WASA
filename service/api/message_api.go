package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var input struct {
		Text  string `json:"text"`
		Photo string `json:"photo"`
	}

	CallingUser := ps.ByName("user")
	ConvStr := ps.ByName("chat")

	Conv, err := strconv.Atoi(ConvStr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if len(input.Text) == 0 && len(input.Photo) == 0 {
		CreateJsonError(w, "Error: the message must contain at least a photo or some text", http.StatusBadRequest)
		return
	}

	if len(input.Text) != 0 && len(input.Photo) != 0 {
		idPhoto, err := rt.db.CreaFoto(input.Photo)
		if !errors.Is(err, nil) {
			CreateJsonError(w, "Error while saving the photo in the message: "+err.Error(), http.StatusInternalServerError)
			return
		}
		errorCode, _, err := rt.db.CreatePhotoTextMessageDB(CallingUser, Conv, idPhoto, input.Text)
		if !errors.Is(err, nil) {
			CreateJsonError(w, "Error while saving the message with photo and text: "+err.Error(), errorCode)
			return
		}
		CreateJsonResponse(w, "Photo and text message sent successfully", http.StatusCreated)
		return
	}

	if len(input.Text) == 0 {
		idPhoto, err := rt.db.CreaFoto(input.Photo)
		if !errors.Is(err, nil) {
			CreateJsonError(w, "Error while saving the photo in the message: "+err.Error(), http.StatusInternalServerError)
			return
		}
		errorCode, _, err := rt.db.CreatePhotoMessageDB(CallingUser, Conv, idPhoto)
		if !errors.Is(err, nil) {
			CreateJsonError(w, "Error while saving the message with photo: "+err.Error(), errorCode)
			return
		}
		CreateJsonResponse(w, "Photo message sent successfully", http.StatusCreated)
		return
	}

	errorCode, _, err := rt.db.CreateTextMessageDB(CallingUser, Conv, input.Text)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while saving the text message: "+err.Error(), errorCode)
		return
	}

	CreateJsonResponse(w, "Text message sent successfully", http.StatusCreated)
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	CallingUser := ps.ByName("user")
	IDMessagestr := ps.ByName("message")
	IDChatstr := ps.ByName("chat")

	IDMessaggio, err := strconv.Atoi(IDMessagestr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	IDChat, err := strconv.Atoi(IDChatstr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteMessage(CallingUser, IDMessaggio, IDChat)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while deleting the message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	CreateJsonResponse(w, "Message successfully deleted", http.StatusOK)
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var input struct {
		Reaction string `json:"reaction"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while parsing the request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	CallingUser := ps.ByName("user")
	IDMessagestr := ps.ByName("message")

	IDMessaggio, err := strconv.Atoi(IDMessagestr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	errorCode, err := rt.db.AddComment(CallingUser, IDMessaggio, input.Reaction)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while adding the comment: "+err.Error(), errorCode)
		return
	}

	CreateJsonResponse(w, "Comment successfully added", http.StatusOK)
}

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	CallingUser := ps.ByName("user")
	IDMessagestr := ps.ByName("message")

	IDMessage, err := strconv.Atoi(IDMessagestr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteComment(CallingUser, IDMessage)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while deleting the comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	CreateJsonResponse(w, "Comment successfully deleted", http.StatusOK)
}
