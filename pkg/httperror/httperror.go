package httperror

import (
	"encoding/json"
	"net/http"
)

type HttpError struct {
	Code         int    `json:"-"`
	ErrorMessage string `json:"error"`
}

func SendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(&HttpError{ErrorMessage: message})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
