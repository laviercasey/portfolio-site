package handler

import (
	"errors"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/portfolio/backend/internal/service"
)

type UploadHandler struct {
	svc *service.UploadService
}

func NewUploadHandler(svc *service.UploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, service.MaxUploadSize)

	_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || params["boundary"] == "" {
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	mr := multipart.NewReader(r.Body, params["boundary"])
	form, err := mr.ReadForm(service.MaxUploadSize)
	if err != nil {
		writeError(w, http.StatusBadRequest, "file too large or invalid multipart form")
		return
	}
	defer func() {
		_ = form.RemoveAll()
	}()

	files := form.File["file"]
	if len(files) == 0 {
		writeError(w, http.StatusBadRequest, "missing 'file' field in form")
		return
	}

	media, err := h.svc.Upload(r.Context(), files[0])
	if err != nil {
		slog.Error("failed to upload file", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to upload file")
		return
	}

	writeJSON(w, http.StatusCreated, media)
}

func (h *UploadHandler) List(w http.ResponseWriter, r *http.Request) {
	media, err := h.svc.List(r.Context())
	if err != nil {
		slog.Error("failed to list media", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list media")
		return
	}
	writeJSON(w, http.StatusOK, media)
}

func (h *UploadHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "media not found")
			return
		}
		slog.Error("failed to delete media", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete media")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
