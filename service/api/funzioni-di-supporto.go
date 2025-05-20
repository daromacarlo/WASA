package api

import (
	"encoding/json"
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

func CreaErroreJson(w http.ResponseWriter, messaggio string, codice int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codice)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"errore":       messaggio,
		"codiceErrore": strconv.Itoa(codice),
	})
}

func CreaRispostaJson(w http.ResponseWriter, messaggio string, codice int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codice)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"risposta": messaggio,
		"codice":   strconv.Itoa(codice),
	})
}
