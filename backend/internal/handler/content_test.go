package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/portfolio/backend/internal/service"
	"github.com/portfolio/backend/internal/testutil"
)

func chiRequest(method, target string, params map[string]string, body string) *http.Request {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, target, strings.NewReader(body))
	} else {
		req = httptest.NewRequest(method, target, nil)
	}
	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestContent_List_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	rows := pgxmock.NewRows([]string{"section", "data", "updated_at"}).
		AddRow("homepage", []byte(`{"hero":"hi"}`), time.Now()).
		AddRow("contact", []byte(`{"email":"a@b"}`), time.Now())
	pool.ExpectQuery("SELECT section, data, updated_at FROM content").WillReturnRows(rows)

	h := NewContentHandler(service.NewContentService(db))

	req := httptest.NewRequest(http.MethodGet, "/api/content", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var out []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("len: got %d, want 2", len(out))
	}
}

func TestContent_List_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery("SELECT section").WillReturnError(errors.New("db down"))

	h := NewContentHandler(service.NewContentService(db))

	req := httptest.NewRequest(http.MethodGet, "/api/content", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}

func TestContent_Update_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	row := pgxmock.NewRows([]string{"section", "data", "updated_at"}).
		AddRow("homepage", []byte(`{"hero":"new"}`), time.Now())
	pool.ExpectQuery("INSERT INTO content").
		WithArgs("homepage", pgxmock.AnyArg()).
		WillReturnRows(row)

	h := NewContentHandler(service.NewContentService(db))

	req := chiRequest(http.MethodPut, "/api/content/homepage",
		map[string]string{"section": "homepage"},
		`{"hero":"new"}`)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestContent_Update_InvalidSection(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewContentHandler(service.NewContentService(db))

	req := chiRequest(http.MethodPut, "/api/content/bogus",
		map[string]string{"section": "bogus"},
		`{"a":1}`)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid section") {
		t.Errorf("body: %s", rec.Body.String())
	}
}

func TestContent_Update_InvalidJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewContentHandler(service.NewContentService(db))

	req := chiRequest(http.MethodPut, "/api/content/homepage",
		map[string]string{"section": "homepage"},
		`{not-json`)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestContent_Update_NotFound(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery("INSERT INTO content").
		WithArgs("homepage", pgxmock.AnyArg()).
		WillReturnError(pgx.ErrNoRows)

	h := NewContentHandler(service.NewContentService(db))

	req := chiRequest(http.MethodPut, "/api/content/homepage",
		map[string]string{"section": "homepage"},
		`{"hero":"x"}`)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", rec.Code)
	}
}

func TestContent_Update_ServerError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery("INSERT INTO content").
		WithArgs("career", pgxmock.AnyArg()).
		WillReturnError(errors.New("boom"))

	h := NewContentHandler(service.NewContentService(db))

	req := chiRequest(http.MethodPut, "/api/content/career",
		map[string]string{"section": "career"},
		`{"x":1}`)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}
