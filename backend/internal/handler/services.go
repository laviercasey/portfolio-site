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

type ServicesHandler struct {
	svc      *service.ServicesService
	validate *validator.Validate
	rev      service.Revalidator
}

func NewServicesHandler(svc *service.ServicesService) *ServicesHandler {
	return &ServicesHandler{
		svc:      svc,
		validate: validator.New(),
	}
}

func (h *ServicesHandler) WithRevalidator(r service.Revalidator) *ServicesHandler {
	h.rev = r
	return h
}

func (h *ServicesHandler) revalidateServicesPage() {
	revalidateAsync(h.rev, []string{"/ru/services", "/en/services"})
}

func (h *ServicesHandler) GetPageData(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetPageData(r.Context())
	if err != nil {
		slog.Error("failed to get services page data", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get services page data")
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (h *ServicesHandler) ListServices(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.ListServices(r.Context())
	if err != nil {
		slog.Error("failed to list services", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list services")
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *ServicesHandler) CreateService(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input model.CreateServiceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(input); err != nil {
		writeError(w, http.StatusBadRequest, formatValidationError(err))
		return
	}
	svc, err := h.svc.CreateService(r.Context(), input)
	if err != nil {
		slog.Error("failed to create service", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create service")
		return
	}
	h.revalidateServicesPage()
	writeJSON(w, http.StatusCreated, svc)
}

func (h *ServicesHandler) UpdateService(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid service id")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input model.UpdateServiceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(input); err != nil {
		writeError(w, http.StatusBadRequest, formatValidationError(err))
		return
	}
	svc, err := h.svc.UpdateService(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "service not found")
			return
		}
		slog.Error("failed to update service", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update service")
		return
	}
	h.revalidateServicesPage()
	writeJSON(w, http.StatusOK, svc)
}

func (h *ServicesHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid service id")
		return
	}
	if err := h.svc.DeleteService(r.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "service not found")
			return
		}
		slog.Error("failed to delete service", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete service")
		return
	}
	h.revalidateServicesPage()
	w.WriteHeader(http.StatusNoContent)
}

func (h *ServicesHandler) ListFaqs(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.ListFaqs(r.Context())
	if err != nil {
		slog.Error("failed to list faqs", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list faqs")
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *ServicesHandler) CreateFaq(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input model.CreateServiceFaqInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(input); err != nil {
		writeError(w, http.StatusBadRequest, formatValidationError(err))
		return
	}
	f, err := h.svc.CreateFaq(r.Context(), input)
	if err != nil {
		slog.Error("failed to create faq", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create faq")
		return
	}
	h.revalidateServicesPage()
	writeJSON(w, http.StatusCreated, f)
}

func (h *ServicesHandler) UpdateFaq(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid faq id")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input model.UpdateServiceFaqInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(input); err != nil {
		writeError(w, http.StatusBadRequest, formatValidationError(err))
		return
	}
	f, err := h.svc.UpdateFaq(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "faq not found")
			return
		}
		slog.Error("failed to update faq", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update faq")
		return
	}
	h.revalidateServicesPage()
	writeJSON(w, http.StatusOK, f)
}

func (h *ServicesHandler) DeleteFaq(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid faq id")
		return
	}
	if err := h.svc.DeleteFaq(r.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "faq not found")
			return
		}
		slog.Error("failed to delete faq", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete faq")
		return
	}
	h.revalidateServicesPage()
	w.WriteHeader(http.StatusNoContent)
}

func (h *ServicesHandler) ListProcessSteps(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.ListProcessSteps(r.Context())
	if err != nil {
		slog.Error("failed to list process steps", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list process steps")
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *ServicesHandler) CreateProcessStep(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input model.CreateServiceProcessStepInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(input); err != nil {
		writeError(w, http.StatusBadRequest, formatValidationError(err))
		return
	}
	p, err := h.svc.CreateProcessStep(r.Context(), input)
	if err != nil {
		slog.Error("failed to create process step", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create process step")
		return
	}
	h.revalidateServicesPage()
	writeJSON(w, http.StatusCreated, p)
}

func (h *ServicesHandler) UpdateProcessStep(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid process step id")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input model.UpdateServiceProcessStepInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(input); err != nil {
		writeError(w, http.StatusBadRequest, formatValidationError(err))
		return
	}
	p, err := h.svc.UpdateProcessStep(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "process step not found")
			return
		}
		slog.Error("failed to update process step", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update process step")
		return
	}
	h.revalidateServicesPage()
	writeJSON(w, http.StatusOK, p)
}

func (h *ServicesHandler) DeleteProcessStep(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid process step id")
		return
	}
	if err := h.svc.DeleteProcessStep(r.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "process step not found")
			return
		}
		slog.Error("failed to delete process step", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete process step")
		return
	}
	h.revalidateServicesPage()
	w.WriteHeader(http.StatusNoContent)
}
