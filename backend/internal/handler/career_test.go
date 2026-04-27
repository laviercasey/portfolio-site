package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/portfolio/backend/internal/service"
	"github.com/portfolio/backend/internal/testutil"
)

func TestCareer_GetAll_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)

	eduRows := pgxmock.NewRows([]string{
		"id", "institution", "degree", "field", "start_year", "end_year",
		"description", "logo_url", "related_project_slugs",
	})
	workRows := pgxmock.NewRows([]string{
		"id", "company", "position", "start_date", "end_date",
		"is_current", "description", "technologies", "logo_url",
		"achievements", "full_description",
	})
	certRows := pgxmock.NewRows([]string{
		"id", "title", "issuer", "issue_date", "credential_id", "url", "image_url",
	})
	pubRows := pgxmock.NewRows([]string{
		"id", "title", "journal", "year", "doi", "url", "abstract",
	})

	pool.ExpectQuery("FROM education").WillReturnRows(eduRows)
	pool.ExpectQuery("FROM work_history").WillReturnRows(workRows)
	pool.ExpectQuery("FROM certificates").WillReturnRows(certRows)
	pool.ExpectQuery("FROM publications").WillReturnRows(pubRows)

	h := NewCareerHandler(service.NewCareerService(db))

	req := httptest.NewRequest(http.MethodGet, "/api/career", nil)
	rec := httptest.NewRecorder()
	h.GetAll(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCareer_GetAll_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery("FROM education").WillReturnError(errors.New("db down"))

	h := NewCareerHandler(service.NewCareerService(db))

	req := httptest.NewRequest(http.MethodGet, "/api/career", nil)
	rec := httptest.NewRecorder()
	h.GetAll(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}

func TestCareer_Create_InvalidType(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewCareerHandler(service.NewCareerService(db))

	req := chiRequest(http.MethodPost, "/api/career/bogus",
		map[string]string{"type": "bogus"},
		`{}`)
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid career type") {
		t.Errorf("body: %s", rec.Body.String())
	}
}

func TestCareer_Create_InvalidJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewCareerHandler(service.NewCareerService(db))

	req := chiRequest(http.MethodPost, "/api/career/work",
		map[string]string{"type": "work"},
		`{not-json`)
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestCareer_Create_Work_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)

	id := uuid.New()
	row := pgxmock.NewRows([]string{
		"id", "company", "position", "start_date", "end_date",
		"is_current", "description", "technologies", "logo_url",
		"achievements", "full_description",
	}).AddRow(
		id, []byte(`{"en":"Acme"}`), []byte(`{"en":"Dev"}`),
		"2024-01-01", nil,
		false, []byte(`{"en":"d"}`), []string{}, nil,
		[]byte(`[]`), []byte(`{}`),
	)
	pool.ExpectQuery("INSERT INTO work_history").WithArgs(testutil.AnyArgs(11)...).WillReturnRows(row)

	h := NewCareerHandler(service.NewCareerService(db))

	body := `{"company":{"en":"Acme"},"position":{"en":"Dev"},"startDate":"2024-01-01","description":{"en":"d"}}`
	req := chiRequest(http.MethodPost, "/api/career/work",
		map[string]string{"type": "work"}, body)
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCareer_Create_Publication_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)

	id := uuid.New()
	row := pgxmock.NewRows([]string{
		"id", "title", "journal", "year", "doi", "url", "abstract",
	}).AddRow(id,
		[]byte(`{"en":"Paper"}`), []byte(`{"en":"Journal"}`),
		"2024", nil, nil, []byte(`{}`),
	)
	pool.ExpectQuery("INSERT INTO publications").WithArgs(testutil.AnyArgs(7)...).WillReturnRows(row)

	h := NewCareerHandler(service.NewCareerService(db))

	body := `{"title":{"en":"Paper"},"journal":{"en":"Journal"},"year":2024}`
	req := chiRequest(http.MethodPost, "/api/career/publication",
		map[string]string{"type": "publication"}, body)
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCareer_Create_Certificate_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)

	id := uuid.New()
	row := pgxmock.NewRows([]string{
		"id", "title", "issuer", "issue_date", "credential_id", "url", "image_url",
	}).AddRow(id,
		[]byte(`{"en":"Cert"}`), []byte(`{"en":"Org"}`),
		"2024-01-02", nil, nil, nil,
	)
	pool.ExpectQuery("INSERT INTO certificates").WithArgs(testutil.AnyArgs(7)...).WillReturnRows(row)

	h := NewCareerHandler(service.NewCareerService(db))

	body := `{"title":{"en":"Cert"},"issuer":{"en":"Org"},"date":"2024-01-02"}`
	req := chiRequest(http.MethodPost, "/api/career/certificate",
		map[string]string{"type": "certificate"}, body)
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCareer_Create_Education_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)

	id := uuid.New()
	row := pgxmock.NewRows([]string{
		"id", "institution", "degree", "field", "start_year", "end_year",
		"description", "logo_url", "related_project_slugs",
	}).AddRow(id,
		[]byte(`{"en":"MIT"}`), []byte(`{"en":"BS"}`), []byte(`{"en":"CS"}`),
		"2020", nil, []byte(`{}`), nil, []string{},
	)
	pool.ExpectQuery("INSERT INTO education").WithArgs(testutil.AnyArgs(9)...).WillReturnRows(row)

	h := NewCareerHandler(service.NewCareerService(db))

	body := `{"institution":{"en":"MIT"},"degree":{"en":"BS"},"field":{"en":"CS"},"startYear":2020}`
	req := chiRequest(http.MethodPost, "/api/career/education",
		map[string]string{"type": "education"}, body)
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCareer_Update_InvalidType(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewCareerHandler(service.NewCareerService(db))

	req := chiRequest(http.MethodPut, "/api/career/xyz/"+uuid.New().String(),
		map[string]string{"type": "xyz", "id": uuid.New().String()}, `{}`)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestCareer_Update_InvalidID(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewCareerHandler(service.NewCareerService(db))

	req := chiRequest(http.MethodPut, "/api/career/work/not-uuid",
		map[string]string{"type": "work", "id": "not-uuid"}, `{}`)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestCareer_Update_InvalidJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewCareerHandler(service.NewCareerService(db))

	id := uuid.New()
	req := chiRequest(http.MethodPut, "/api/career/work/"+id.String(),
		map[string]string{"type": "work", "id": id.String()}, `{not-json`)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestCareer_Update_Work_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)

	id := uuid.New()
	row := pgxmock.NewRows([]string{
		"id", "company", "position", "start_date", "end_date",
		"is_current", "description", "technologies", "logo_url",
		"achievements", "full_description",
	}).AddRow(
		id, []byte(`{"en":"Acme"}`), []byte(`{"en":"Dev"}`),
		"2024-01-01", nil,
		false, []byte(`{"en":"d"}`), []string{}, nil,
		[]byte(`[]`), []byte(`{}`),
	)
	pool.ExpectQuery("UPDATE work_history").WithArgs(testutil.AnyArgs(12)...).WillReturnRows(row)

	h := NewCareerHandler(service.NewCareerService(db))

	body := `{"company":{"en":"Acme"},"position":{"en":"Dev"},"startDate":"2024-01-01","description":{"en":"d"}}`
	req := chiRequest(http.MethodPut, "/api/career/work/"+id.String(),
		map[string]string{"type": "work", "id": id.String()}, body)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCareer_Update_NotFound(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery("UPDATE work_history").WithArgs(testutil.AnyArgs(12)...).WillReturnError(pgx.ErrNoRows)

	h := NewCareerHandler(service.NewCareerService(db))

	id := uuid.New()
	body := `{"company":{"en":"A"},"position":{"en":"B"},"startDate":"2024-01-01","description":{"en":"d"}}`
	req := chiRequest(http.MethodPut, "/api/career/work/"+id.String(),
		map[string]string{"type": "work", "id": id.String()}, body)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", rec.Code)
	}
}

func TestCareer_Update_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery("UPDATE publications").WithArgs(testutil.AnyArgs(8)...).WillReturnError(errors.New("boom"))

	h := NewCareerHandler(service.NewCareerService(db))

	id := uuid.New()
	body := `{"title":{"en":"P"},"journal":{"en":"J"},"year":2024}`
	req := chiRequest(http.MethodPut, "/api/career/publication/"+id.String(),
		map[string]string{"type": "publication", "id": id.String()}, body)
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}

func TestCareer_Delete_InvalidType(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewCareerHandler(service.NewCareerService(db))

	req := chiRequest(http.MethodDelete, "/api/career/xyz/"+uuid.New().String(),
		map[string]string{"type": "xyz", "id": uuid.New().String()}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestCareer_Delete_InvalidID(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	h := NewCareerHandler(service.NewCareerService(db))

	req := chiRequest(http.MethodDelete, "/api/career/work/not-uuid",
		map[string]string{"type": "work", "id": "not-uuid"}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestCareer_Delete_Work_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectExec("DELETE FROM work_history").
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	h := NewCareerHandler(service.NewCareerService(db))

	req := chiRequest(http.MethodDelete, "/api/career/work/"+id.String(),
		map[string]string{"type": "work", "id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: got %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCareer_Delete_Education_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectExec("DELETE FROM education").
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	h := NewCareerHandler(service.NewCareerService(db))

	req := chiRequest(http.MethodDelete, "/api/career/education/"+id.String(),
		map[string]string{"type": "education", "id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: got %d, want 204", rec.Code)
	}
}

func TestCareer_Delete_Certificate_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectExec("DELETE FROM certificates").
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	h := NewCareerHandler(service.NewCareerService(db))

	req := chiRequest(http.MethodDelete, "/api/career/certificate/"+id.String(),
		map[string]string{"type": "certificate", "id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: got %d, want 204", rec.Code)
	}
}

func TestCareer_Delete_Publication_OK(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectExec("DELETE FROM publications").
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	h := NewCareerHandler(service.NewCareerService(db))

	req := chiRequest(http.MethodDelete, "/api/career/publication/"+id.String(),
		map[string]string{"type": "publication", "id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: got %d, want 204", rec.Code)
	}
}

func TestCareer_Delete_NotFound(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectExec("DELETE FROM work_history").
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	h := NewCareerHandler(service.NewCareerService(db))

	req := chiRequest(http.MethodDelete, "/api/career/work/"+id.String(),
		map[string]string{"type": "work", "id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", rec.Code)
	}
}

func TestCareer_Delete_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()
	pool.ExpectExec("DELETE FROM work_history").
		WithArgs(id).
		WillReturnError(errors.New("boom"))

	h := NewCareerHandler(service.NewCareerService(db))

	req := chiRequest(http.MethodDelete, "/api/career/work/"+id.String(),
		map[string]string{"type": "work", "id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
}
