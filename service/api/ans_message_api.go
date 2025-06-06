package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) ansMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	user := ps.ByName("user")
	IdMessageStr := ps.ByName("message")

	IdMessage, err := strconv.Atoi(IdMessageStr)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

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
		CreateJsonError(w, "Error: The message must contain at least a photo or some text", http.StatusBadRequest)
		return
	}

	if len(input.Text) > 0 && len(input.Photo) > 0 {
		idPhoto, err := rt.db.CreaFoto(input.Photo)
		if !errors.Is(err, nil) {
			CreateJsonError(w, "Error while uploading the photo for the reply: "+err.Error(), http.StatusInternalServerError)
			return
		}
		errorCode, err := rt.db.AnsPhotoTextMessage(user, Conv, IdMessage, idPhoto, input.Text)
		if !errors.Is(err, nil) {
			CreateJsonError(w, "Error while replying with photo and text: "+err.Error(), errorCode)
			return
		}
		CreateJsonResponse(w, "Reply with photo and text sent successfully", http.StatusOK)
		return
	}

	if len(input.Text) == 0 {
		idPhoto, err := rt.db.CreaFoto(input.Photo)
		if !errors.Is(err, nil) {
			CreateJsonError(w, "Error while uploading the photo for the reply: "+err.Error(), http.StatusInternalServerError)
			return
		}
		errorCode, err := rt.db.AnsPhotoMessage(user, Conv, IdMessage, idPhoto)
		if !errors.Is(err, nil) {
			CreateJsonError(w, "Error while replying with photo: "+err.Error(), errorCode)
			return
		}
		CreateJsonResponse(w, "Reply with photo sent successfully", http.StatusOK)
		return
	}

	errorCode, err := rt.db.AnsTextMessage(user, Conv, IdMessage, input.Text)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while sending the reply message: "+err.Error(), errorCode)
		return
	}
	CreateJsonResponse(w, "Reply sent successfully", http.StatusOK)
}
