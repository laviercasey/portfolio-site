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

func projectColumnNames() []string {
	return []string{
		"id", "slug", "title", "short_desc", "description",
		"category", "status", "tags", "tech_stack", "goal_desc",
		"github_url", "demo_url", "site_url", "video_url",
		"thumbnail", "images", "stars",
		"featured", "sort_order", "created_at", "updated_at",
		"problem", "approach", "outcome",
		"tech_choices", "highlights",
		"timeline_started", "timeline_shipped", "demo_credentials",
	}
}

func addProjectRow(rows *pgxmock.Rows, id uuid.UUID, slug string) *pgxmock.Rows {
	now := time.Now()
	return rows.AddRow(
		id, slug,
		[]byte(`{"en":"T"}`), []byte(`{"en":"SD"}`), []byte(`{"en":"D"}`),
		"web", "completed", []string{}, []string{}, []byte(`{}`),
		nil, nil, nil, nil,
		nil, []string{}, 0,
		false, 0, now, now,
		[]byte(`{}`), []byte(`{}`), []byte(`{}`),
		[]byte(`[]`), []byte(`[]`),
		nil, nil, []byte(`[]`),
	)
}

func TestProjects_List_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	rows := pgxmock.NewRows(projectColumnNames())
	addProjectRow(rows, id, "foo")

	pool.ExpectQuery("FROM projects ORDER BY").WillReturnRows(rows)

	h := NewProjectHandler(service.NewProjectService(db))

	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var out []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("len: %d", len(out))
	}
}

func TestProjects_List_WithCategory(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	rows := pgxmock.NewRows(projectColumnNames())
	addProjectRow(rows, id, "foo")

	pool.ExpectQuery("FROM projects WHERE category").
		WithArgs("web").
		WillReturnRows(rows)

	h := NewProjectHandler(service.NewProjectService(db))

	req := httptest.NewRequest(http.MethodGet, "/api/projects?category=web", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
}

func TestProjects_List_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery("FROM projects").WillReturnError(errors.New("boom"))

	h := NewProjectHandler(service.NewProjectService(db))

	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}

func TestProjects_GetBySlug_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	rows := pgxmock.NewRows(projectColumnNames())
	addProjectRow(rows, id, "foo")

	pool.ExpectQuery("FROM projects WHERE slug").
		WithArgs("foo").
		WillReturnRows(rows)

	h := NewProjectHandler(service.NewProjectService(db))

	req := chiRequest(http.MethodGet, "/api/projects/foo",
		map[string]string{"slug": "foo"}, "")
	rec := httptest.NewRecorder()
	h.GetBySlug(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestProjects_GetBySlug_Empty(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewProjectHandler(service.NewProjectService(db))

	req := chiRequest(http.MethodGet, "/api/projects/",
		map[string]string{"slug": ""}, "")
	rec := httptest.NewRecorder()
	h.GetBySlug(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestProjects_GetBySlug_NotFound(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery("FROM projects WHERE slug").
		WithArgs("missing").
		WillReturnError(pgx.ErrNoRows)

	h := NewProjectHandler(service.NewProjectService(db))

	req := chiRequest(http.MethodGet, "/api/projects/missing",
		map[string]string{"slug": "missing"}, "")
	rec := httptest.NewRecorder()
	h.GetBySlug(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", rec.Code)
	}
}

func TestProjects_GetBySlug_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery("FROM projects WHERE slug").
		WithArgs("foo").
		WillReturnError(errors.New("boom"))

	h := NewProjectHandler(service.NewProjectService(db))

	req := chiRequest(http.MethodGet, "/api/projects/foo",
		map[string]string{"slug": "foo"}, "")
	rec := httptest.NewRecorder()
	h.GetBySlug(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}

func TestProjects_Create_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	rows := pgxmock.NewRows(projectColumnNames())
	addProjectRow(rows, id, "new-slug")

	pool.ExpectQuery("INSERT INTO projects").WithArgs(testutil.AnyArgs(25)...).WillReturnRows(rows)

	h := NewProjectHandler(service.NewProjectService(db))

	body := `{"slug":"new-slug","title":{"en":"T"},"category":"web","status":"completed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/projects", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
}

func TestProjects_Create_InvalidJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewProjectHandler(service.NewProjectService(db))

	req := httptest.NewRequest(http.MethodPost, "/api/projects", strings.NewReader(`{not-json`))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestProjects_Create_ValidationError(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewProjectHandler(service.NewProjectService(db))

	body := `{"slug":"x"}`
	req := httptest.NewRequest(http.MethodPost, "/api/projects", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestProjects_Create_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery("INSERT INTO projects").WithArgs(testutil.AnyArgs(25)...).WillReturnError(errors.New("boom"))

	h := NewProjectHandler(service.NewProjectService(db))

	body := `{"slug":"new-slug","title":{"en":"T"},"category":"web","status":"completed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/projects", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}

func TestProjects_Update_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()

	selRows := pgxmock.NewRows(projectColumnNames())
	addProjectRow(selRows, id, "foo")
	pool.ExpectQuery("FROM projects WHERE id").
		WithArgs(id).
		WillReturnRows(selRows)

	updRows := pgxmock.NewRows(projectColumnNames())
	addProjectRow(updRows, id, "foo")
	pool.ExpectQuery("UPDATE projects SET").WithArgs(testutil.AnyArgs(26)...).WillReturnRows(updRows)

	h := NewProjectHandler(service.NewProjectService(db))

	body := `{"featured":true}`
	req := chiRequest(http.MethodPut, "/api/projects/"+id.String(),
		map[string]string{"id": id.String()}, body)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestProjects_Update_InvalidID(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewProjectHandler(service.NewProjectService(db))

	req := chiRequest(http.MethodPut, "/api/projects/bad",
		map[string]string{"id": "bad"}, `{}`)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestProjects_Update_InvalidJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewProjectHandler(service.NewProjectService(db))

	id := uuid.New()
	req := chiRequest(http.MethodPut, "/api/projects/"+id.String(),
		map[string]string{"id": id.String()}, `{not-json`)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestProjects_Update_ValidationError(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewProjectHandler(service.NewProjectService(db))

	id := uuid.New()
	body := `{"githubUrl":"not-a-url"}`
	req := chiRequest(http.MethodPut, "/api/projects/"+id.String(),
		map[string]string{"id": id.String()}, body)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestProjects_Update_NotFound(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectQuery("FROM projects WHERE id").
		WithArgs(id).
		WillReturnError(pgx.ErrNoRows)

	h := NewProjectHandler(service.NewProjectService(db))

	body := `{}`
	req := chiRequest(http.MethodPut, "/api/projects/"+id.String(),
		map[string]string{"id": id.String()}, body)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", rec.Code)
	}
}

func TestProjects_Delete_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectExec("DELETE FROM projects").
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	h := NewProjectHandler(service.NewProjectService(db))

	req := chiRequest(http.MethodDelete, "/api/projects/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: got %d, want 204", rec.Code)
	}
}

func TestProjects_Delete_InvalidID(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewProjectHandler(service.NewProjectService(db))

	req := chiRequest(http.MethodDelete, "/api/projects/bad",
		map[string]string{"id": "bad"}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestProjects_Delete_NotFound(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectExec("DELETE FROM projects").
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	h := NewProjectHandler(service.NewProjectService(db))

	req := chiRequest(http.MethodDelete, "/api/projects/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", rec.Code)
	}
}

func TestProjects_Delete_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectExec("DELETE FROM projects").
		WithArgs(id).
		WillReturnError(errors.New("boom"))

	h := NewProjectHandler(service.NewProjectService(db))

	req := chiRequest(http.MethodDelete, "/api/projects/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}
