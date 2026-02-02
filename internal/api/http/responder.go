package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

func SendError(w http.ResponseWriter, statusCode int, err error) {

	log.Printf("Error:%v", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
