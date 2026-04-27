package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		slog.Error("json marshal failed", "error", err)
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(b); err != nil {
		slog.Error("write response failed", "error", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func formatValidationError(err error) string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return "validation failed"
	}

	var msgs []string
	for _, fe := range ve {
		msgs = append(msgs, fmt.Sprintf("field '%s' failed on '%s' validation", fe.Field(), fe.Tag()))
	}
	return strings.Join(msgs, "; ")
}
