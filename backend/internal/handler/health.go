package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/portfolio/backend/internal/database"
)

type HealthHandler struct {
	db *database.DB
}

func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	dbStatus := "ok"
	code := http.StatusOK

	if err := h.db.Health(r.Context()); err != nil {
		dbStatus = "unhealthy"
		status = "degraded"
		code = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status":   status,
		"database": dbStatus,
	}); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}
