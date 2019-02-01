package goauthserver

import (
	"encoding/json"
	"net/http"
)

// JSON will parse payload and send the response back to client
func JSON(w http.ResponseWriter, code int, payload interface{}) {
	// Parse the payload into byte, error has been resolved at the upper level, ignore here
	response, _ := json.Marshal(payload)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Error will be used to out put error
func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, code, map[string]string{"error": message})
}
