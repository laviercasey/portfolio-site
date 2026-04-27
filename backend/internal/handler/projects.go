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

type ProjectHandler struct {
	svc      *service.ProjectService
	validate *validator.Validate
	rev      service.Revalidator
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		svc:      svc,
		validate: validator.New(),
	}
}

func (h *ProjectHandler) WithRevalidator(r service.Revalidator) *ProjectHandler {
	h.rev = r
	return h
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	projects, err := h.svc.List(r.Context(), category)
	if err != nil {
		slog.Error("failed to list projects", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list projects")
		return
	}

	writeJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	project, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "project not found")
			return
		}
		slog.Error("failed to get project", "slug", slug, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get project")
		return
	}

	writeJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input model.CreateProjectInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(input); err != nil {
		writeError(w, http.StatusBadRequest, formatValidationError(err))
		return
	}

	project, err := h.svc.Create(r.Context(), input)
	if err != nil {
		slog.Error("failed to create project", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create project")
		return
	}

	revalidateAsync(h.rev, []string{"/"})
	writeJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input model.UpdateProjectInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(input); err != nil {
		writeError(w, http.StatusBadRequest, formatValidationError(err))
		return
	}

	project, err := h.svc.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "project not found")
			return
		}
		slog.Error("failed to update project", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update project")
		return
	}

	revalidateAsync(h.rev, []string{"/"})
	writeJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "project not found")
			return
		}
		slog.Error("failed to delete project", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete project")
		return
	}

	revalidateAsync(h.rev, []string{"/"})
	w.WriteHeader(http.StatusNoContent)
}
