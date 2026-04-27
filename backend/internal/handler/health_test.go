package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/portfolio/backend/internal/database"
)

func newPingablePool(t *testing.T) (*database.DB, pgxmock.PgxPoolIface) {
	t.Helper()
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool: %v", err)
	}
	t.Cleanup(pool.Close)
	return database.NewWithPool(pool), pool
}

func TestHealth_OK(t *testing.T) {
	t.Parallel()

	db, pool := newPingablePool(t)
	pool.ExpectPing()

	h := NewHealthHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.Check(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var out map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out["status"] != "ok" {
		t.Errorf("status: got %q, want ok", out["status"])
	}
	if out["database"] != "ok" {
		t.Errorf("database: got %q, want ok", out["database"])
	}
}

func TestHealth_Degraded(t *testing.T) {
	t.Parallel()

	db, pool := newPingablePool(t)
	pool.ExpectPing().WillReturnError(errors.New("db down"))

	h := NewHealthHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.Check(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status: got %d, want 503", rec.Code)
	}
	var out map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out["status"] != "degraded" {
		t.Errorf("status: got %q, want degraded", out["status"])
	}
	if out["database"] != "unhealthy" {
		t.Errorf("database: got %q, want unhealthy", out["database"])
	}
}
