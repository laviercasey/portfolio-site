package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/portfolio/backend/internal/service"
)

type ContentHandler struct {
	svc *service.ContentService
	rev service.Revalidator
}

func NewContentHandler(svc *service.ContentService) *ContentHandler {
	return &ContentHandler{svc: svc}
}

func (h *ContentHandler) WithRevalidator(r service.Revalidator) *ContentHandler {
	h.rev = r
	return h
}

func isValidContentSection(s string) bool {
	switch s {
	case "homepage", "contact", "career":
		return true
	}
	return false
}

func (h *ContentHandler) List(w http.ResponseWriter, r *http.Request) {
	contents, err := h.svc.GetAll(r.Context())
	if err != nil {
		slog.Error("failed to list content", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list content")
		return
	}

	writeJSON(w, http.StatusOK, contents)
}

func (h *ContentHandler) Update(w http.ResponseWriter, r *http.Request) {
	section := chi.URLParam(r, "section")
	if !isValidContentSection(section) {
		writeError(w, http.StatusBadRequest, "invalid section: must be homepage, contact, or career")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var body json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(body) == 0 {
		writeError(w, http.StatusBadRequest, "data is required")
		return
	}

	content, err := h.svc.Upsert(r.Context(), section, body)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "content section not found")
			return
		}
		slog.Error("failed to update content", "section", section, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update content")
		return
	}

	revalidateAsync(h.rev, []string{"/"})
	writeJSON(w, http.StatusOK, content)
}
