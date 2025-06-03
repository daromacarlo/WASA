package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

/*
func ReadImageFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileData := make([]byte, fileInfo.Size())
	_, err = file.Read(fileData)
	if err != nil {
		return nil, err
	}

	return
}
*/

func CreateJsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":     message,
		"errorCode": strconv.Itoa(code),
	})
}

func CreateJsonResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"response": message,
		"code":     strconv.Itoa(code),
	})
}

func CreateJsonAccessResponse(w http.ResponseWriter, message string, userId int, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"response": message,
		"userid":   strconv.Itoa(userId),
		"code":     strconv.Itoa(code),
	})
}
func (rt *_router) VerifyToken(token string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("Error: token not valid - empty token")
	}
	token_i, err := strconv.Atoi(token)
	if err != nil {
		return false, fmt.Errorf("Error: invalid token format - %v is not a valid integer", token)
	}

	exists, err := rt.db.UserExistenceId(token_i)
	if err != nil {
		return false, fmt.Errorf("Error: failed to verify user existence - DB error: %v", err)
	}

	if !exists {
		return false, fmt.Errorf("Error: user with ID %d does not exist", token_i)
	}
	return true, nil
}
