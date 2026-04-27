package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/portfolio/backend/internal/model"
	"github.com/portfolio/backend/internal/service"
)

type InquiryHandler struct {
	svc      *service.InquiryService
	validate *validator.Validate
}

func NewInquiryHandler(svc *service.InquiryService) *InquiryHandler {
	return &InquiryHandler{
		svc:      svc,
		validate: validator.New(),
	}
}

func (h *InquiryHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	inquiries, err := h.svc.List(r.Context(), status)
	if err != nil {
		slog.Error("failed to list inquiries", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list inquiries")
		return
	}

	writeJSON(w, http.StatusOK, inquiries)
}

func (h *InquiryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid inquiry id")
		return
	}

	inquiry, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "inquiry not found")
			return
		}
		slog.Error("failed to get inquiry", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get inquiry")
		return
	}

	writeJSON(w, http.StatusOK, inquiry)
}

func (h *InquiryHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input model.CreateInquiryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(input); err != nil {
		writeError(w, http.StatusBadRequest, formatValidationError(err))
		return
	}

	inquiry, err := h.svc.Create(r.Context(), input)
	if err != nil {
		slog.Error("failed to create inquiry", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create inquiry")
		return
	}

	writeJSON(w, http.StatusCreated, inquiry)
}

func (h *InquiryHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid inquiry id")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input model.UpdateInquiryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(input); err != nil {
		writeError(w, http.StatusBadRequest, formatValidationError(err))
		return
	}

	inquiry, err := h.svc.UpdateStatus(r.Context(), id, input.Status, input.AdminNotes)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "inquiry not found")
			return
		}
		slog.Error("failed to update inquiry status", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update inquiry status")
		return
	}

	writeJSON(w, http.StatusOK, inquiry)
}
