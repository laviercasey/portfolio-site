package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/portfolio/backend/internal/service"
)

type CareerHandler struct {
	svc *service.CareerService
	rev service.Revalidator
}

func NewCareerHandler(svc *service.CareerService) *CareerHandler {
	return &CareerHandler{svc: svc}
}

func (h *CareerHandler) WithRevalidator(r service.Revalidator) *CareerHandler {
	h.rev = r
	return h
}

func (h *CareerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetAll(r.Context())
	if err != nil {
		slog.Error("failed to get career data", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get career data")
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (h *CareerHandler) Create(w http.ResponseWriter, r *http.Request) {
	careerType := chi.URLParam(r, "type")
	if !isValidCareerType(careerType) {
		writeError(w, http.StatusBadRequest, "invalid career type: must be education, work, certificate, or publication")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	if !json.Valid(body) {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	result, err := h.svc.Create(r.Context(), careerType, json.RawMessage(body))
	if err != nil {
		if errors.Is(err, service.ErrInvalidCareerType) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		slog.Error("failed to create career entry", "type", careerType, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create career entry")
		return
	}

	revalidateAsync(h.rev, []string{"/"})
	writeJSON(w, http.StatusCreated, result)
}

func (h *CareerHandler) Update(w http.ResponseWriter, r *http.Request) {
	careerType := chi.URLParam(r, "type")
	if !isValidCareerType(careerType) {
		writeError(w, http.StatusBadRequest, "invalid career type: must be education, work, certificate, or publication")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	if !json.Valid(body) {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	result, err := h.svc.Update(r.Context(), careerType, id, json.RawMessage(body))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "career entry not found")
			return
		}
		if errors.Is(err, service.ErrInvalidCareerType) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		slog.Error("failed to update career entry", "type", careerType, "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update career entry")
		return
	}

	revalidateAsync(h.rev, []string{"/"})
	writeJSON(w, http.StatusOK, result)
}

func (h *CareerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	careerType := chi.URLParam(r, "type")
	if !isValidCareerType(careerType) {
		writeError(w, http.StatusBadRequest, "invalid career type: must be education, work, certificate, or publication")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(r.Context(), careerType, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "career entry not found")
			return
		}
		if errors.Is(err, service.ErrInvalidCareerType) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		slog.Error("failed to delete career entry", "type", careerType, "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete career entry")
		return
	}

	revalidateAsync(h.rev, []string{"/"})
	w.WriteHeader(http.StatusNoContent)
}

func isValidCareerType(t string) bool {
	switch t {
	case "education", "work", "certificate", "publication":
		return true
	}
	return false
}
