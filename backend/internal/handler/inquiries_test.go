package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/portfolio/backend/internal/service"
	"github.com/portfolio/backend/internal/testutil"
)

func inquiryRow(id uuid.UUID) *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id", "name", "email", "company", "telegram",
		"inquiry_type", "budget", "message",
		"status", "admin_notes", "created_at", "updated_at",
	}).AddRow(
		id, "Jane", "jane@example.com", nil, nil,
		"freelance", nil, "Hello there, this is a test message.",
		"new", nil, time.Now(), time.Now(),
	)
}

func TestInquiries_List_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectQuery("FROM inquiries ORDER BY created_at DESC").
		WillReturnRows(inquiryRow(id))

	h := NewInquiryHandler(service.NewInquiryService(db))

	req := httptest.NewRequest(http.MethodGet, "/api/inquiries", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestInquiries_List_WithStatus(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectQuery("FROM inquiries WHERE status").
		WithArgs("new").
		WillReturnRows(inquiryRow(id))

	h := NewInquiryHandler(service.NewInquiryService(db))

	req := httptest.NewRequest(http.MethodGet, "/api/inquiries?status=new", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
}

func TestInquiries_List_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery("FROM inquiries").WillReturnError(errors.New("db down"))

	h := NewInquiryHandler(service.NewInquiryService(db))

	req := httptest.NewRequest(http.MethodGet, "/api/inquiries", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}

func TestInquiries_GetByID_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectQuery("FROM inquiries WHERE id").
		WithArgs(id).
		WillReturnRows(inquiryRow(id))

	h := NewInquiryHandler(service.NewInquiryService(db))

	req := chiRequest(http.MethodGet, "/api/inquiries/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.GetByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestInquiries_GetByID_InvalidUUID(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewInquiryHandler(service.NewInquiryService(db))

	req := chiRequest(http.MethodGet, "/api/inquiries/not-uuid",
		map[string]string{"id": "not-uuid"}, "")
	rec := httptest.NewRecorder()
	h.GetByID(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestInquiries_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectQuery("FROM inquiries WHERE id").
		WithArgs(id).
		WillReturnError(pgx.ErrNoRows)

	h := NewInquiryHandler(service.NewInquiryService(db))

	req := chiRequest(http.MethodGet, "/api/inquiries/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.GetByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", rec.Code)
	}
}

func TestInquiries_GetByID_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectQuery("FROM inquiries WHERE id").
		WithArgs(id).
		WillReturnError(errors.New("boom"))

	h := NewInquiryHandler(service.NewInquiryService(db))

	req := chiRequest(http.MethodGet, "/api/inquiries/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.GetByID(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}

func TestInquiries_Create_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectQuery("INSERT INTO inquiries").
		WithArgs(testutil.AnyArgs(7)...).
		WillReturnRows(inquiryRow(id))

	h := NewInquiryHandler(service.NewInquiryService(db))

	body := `{"name":"Jane","email":"jane@example.com","type":"freelance","message":"Hello there, this is a test message."}`
	req := httptest.NewRequest(http.MethodPost, "/api/inquiries", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out["email"] != "jane@example.com" {
		t.Errorf("email: got %v", out["email"])
	}
}

func TestInquiries_Create_InvalidJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewInquiryHandler(service.NewInquiryService(db))

	req := httptest.NewRequest(http.MethodPost, "/api/inquiries", strings.NewReader(`{not-json`))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestInquiries_Create_ValidationError(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewInquiryHandler(service.NewInquiryService(db))

	body := `{"name":"x","email":"not-email","type":"weird","message":"short"}`
	req := httptest.NewRequest(http.MethodPost, "/api/inquiries", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "validation") && !strings.Contains(rec.Body.String(), "failed on") {
		t.Errorf("expected validation error, got %s", rec.Body.String())
	}
}

func TestInquiries_Create_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery("INSERT INTO inquiries").
		WithArgs(testutil.AnyArgs(7)...).
		WillReturnError(errors.New("boom"))

	h := NewInquiryHandler(service.NewInquiryService(db))

	body := `{"name":"Jane","email":"jane@example.com","type":"freelance","message":"Hello there, this is a test message."}`
	req := httptest.NewRequest(http.MethodPost, "/api/inquiries", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}

func TestInquiries_UpdateStatus_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectQuery("UPDATE inquiries SET status").
		WithArgs(testutil.AnyArgs(3)...).
		WillReturnRows(inquiryRow(id))

	h := NewInquiryHandler(service.NewInquiryService(db))

	body := `{"status":"read","adminNotes":"seen it"}`
	req := chiRequest(http.MethodPatch, "/api/inquiries/"+id.String(),
		map[string]string{"id": id.String()}, body)
	rec := httptest.NewRecorder()
	h.UpdateStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestInquiries_UpdateStatus_InvalidUUID(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewInquiryHandler(service.NewInquiryService(db))

	req := chiRequest(http.MethodPatch, "/api/inquiries/bad",
		map[string]string{"id": "bad"}, `{"status":"read"}`)
	rec := httptest.NewRecorder()
	h.UpdateStatus(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestInquiries_UpdateStatus_InvalidJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewInquiryHandler(service.NewInquiryService(db))

	id := uuid.New()
	req := chiRequest(http.MethodPatch, "/api/inquiries/"+id.String(),
		map[string]string{"id": id.String()}, `{not-json`)
	rec := httptest.NewRecorder()
	h.UpdateStatus(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestInquiries_UpdateStatus_ValidationError(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewInquiryHandler(service.NewInquiryService(db))

	id := uuid.New()
	body := `{"status":"not-valid"}`
	req := chiRequest(http.MethodPatch, "/api/inquiries/"+id.String(),
		map[string]string{"id": id.String()}, body)
	rec := httptest.NewRecorder()
	h.UpdateStatus(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestInquiries_UpdateStatus_NotFound(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectQuery("UPDATE inquiries SET status").
		WithArgs(testutil.AnyArgs(3)...).
		WillReturnError(pgx.ErrNoRows)

	h := NewInquiryHandler(service.NewInquiryService(db))

	body := `{"status":"read"}`
	req := chiRequest(http.MethodPatch, "/api/inquiries/"+id.String(),
		map[string]string{"id": id.String()}, body)
	rec := httptest.NewRecorder()
	h.UpdateStatus(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", rec.Code)
	}
}
