package api

import (
	"os"
)

// (questa funzione non è chiamabile direttamente dall'utente)
// Funzione di supporto per leggere un file immagine
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

	return fileData, nil
}
