package handler

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/portfolio/backend/internal/service"
	"github.com/portfolio/backend/internal/testutil"
)

var pngMagic = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
	0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
	0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
	0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
	0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
	0x42, 0x60, 0x82,
}

func buildMultipart(t *testing.T, fieldName, filename string, content []byte) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile(fieldName, filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write content: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	return &buf, w.FormDataContentType()
}

func TestUpload_Upload_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	dir := t.TempDir()

	pool.ExpectQuery("INSERT INTO media").
		WithArgs(testutil.AnyArgs(5)...).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "filename", "original_name", "mime_type", "size", "url", "created_at",
		}).AddRow(
			uuid.New(), "abc.png", "pic.png", "image/png", int64(68), "/uploads/abc.png", time.Now(),
		))

	h := NewUploadHandler(service.NewUploadService(db, dir))

	body, contentType := buildMultipart(t, "file", "pic.png", pngMagic)
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()
	h.Upload(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
}

func TestUpload_Upload_MissingFile(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	dir := t.TempDir()

	h := NewUploadHandler(service.NewUploadService(db, dir))

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if err := w.WriteField("notfile", "value"); err != nil {
		t.Fatalf("write field: %v", err)
	}
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	h.Upload(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestUpload_Upload_InvalidMultipart(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	dir := t.TempDir()

	h := NewUploadHandler(service.NewUploadService(db, dir))

	req := httptest.NewRequest(http.MethodPost, "/api/upload", bytes.NewReader([]byte("not multipart")))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=X")
	rec := httptest.NewRecorder()
	h.Upload(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestUpload_Upload_UnsupportedType(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	dir := t.TempDir()

	h := NewUploadHandler(service.NewUploadService(db, dir))

	body, contentType := buildMultipart(t, "file", "bad.exe", []byte("MZ\x00\x00 not an image"))
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()
	h.Upload(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}

func TestUpload_List_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	dir := t.TempDir()

	rows := pgxmock.NewRows([]string{
		"id", "filename", "original_name", "mime_type", "size", "url", "created_at",
	}).AddRow(
		uuid.New(), "a.png", "a.png", "image/png", int64(100), "/uploads/a.png", time.Now(),
	)
	pool.ExpectQuery("FROM media ORDER BY created_at DESC").WillReturnRows(rows)

	h := NewUploadHandler(service.NewUploadService(db, dir))

	req := httptest.NewRequest(http.MethodGet, "/api/media", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestUpload_List_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	dir := t.TempDir()

	pool.ExpectQuery("FROM media").WillReturnError(errors.New("boom"))

	h := NewUploadHandler(service.NewUploadService(db, dir))

	req := httptest.NewRequest(http.MethodGet, "/api/media", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}

func TestUpload_Delete_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	dir := t.TempDir()

	id := uuid.New()
	pool.ExpectQuery("SELECT filename FROM media").
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{"filename"}).AddRow("gone.png"))
	pool.ExpectExec("DELETE FROM media").
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	h := NewUploadHandler(service.NewUploadService(db, dir))

	req := chiRequest(http.MethodDelete, "/api/media/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: got %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
}

func TestUpload_Delete_InvalidID(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	dir := t.TempDir()

	h := NewUploadHandler(service.NewUploadService(db, dir))

	req := chiRequest(http.MethodDelete, "/api/media/bad",
		map[string]string{"id": "bad"}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestUpload_Delete_NotFound(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	dir := t.TempDir()

	id := uuid.New()
	pool.ExpectQuery("SELECT filename FROM media").
		WithArgs(id).
		WillReturnError(pgx.ErrNoRows)

	h := NewUploadHandler(service.NewUploadService(db, dir))

	req := chiRequest(http.MethodDelete, "/api/media/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", rec.Code)
	}
}

func TestUpload_Delete_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	dir := t.TempDir()

	id := uuid.New()
	pool.ExpectQuery("SELECT filename FROM media").
		WithArgs(id).
		WillReturnError(errors.New("boom"))

	h := NewUploadHandler(service.NewUploadService(db, dir))

	req := chiRequest(http.MethodDelete, "/api/media/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}
