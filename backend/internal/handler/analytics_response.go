package handler

import "net/http"

type apiEnvelope struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func writeSuccess(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, apiEnvelope{Success: true, Data: data})
}

func writeAPIError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, apiEnvelope{Success: false, Error: msg})
}
