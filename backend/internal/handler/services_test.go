package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/portfolio/backend/internal/service"
	"github.com/portfolio/backend/internal/testutil"
)

var pgxNoRows = pgx.ErrNoRows

func strReader(s string) io.Reader { return strings.NewReader(s) }

func newSvcHandler(t *testing.T) (*ServicesHandler, pgxmock.PgxPoolIface) {
	t.Helper()
	db, pool := testutil.NewMockDB(t)
	return NewServicesHandler(service.NewServicesService(db)), pool
}

func svcCols() []string {
	return []string{
		"id", "slug", "num", "icon_key", "visual_key", "accent",
		"title", "lead", "bullets", "stack", "timeline",
		"price_ru", "price_en", "case_projects", "sort_order",
		"created_at", "updated_at",
	}
}

func svcRow(id uuid.UUID, slug string) []any {
	now := time.Now().UTC()
	return []any{
		id, slug, "01", "bot", "terminal", "#5eb3ff",
		[]byte(`{"ru":"T","en":"T"}`),
		[]byte(`{"ru":"L","en":"L"}`),
		[]byte(`{"ru":["b"],"en":["b"]}`),
		"stack", []byte(`{"ru":"1w","en":"1w"}`),
		"100", "$100", []byte(`[]`), 10, now, now,
	}
}

func faqColsH() []string {
	return []string{"id", "question", "answer", "sort_order", "created_at", "updated_at"}
}
func faqRowH(id uuid.UUID) []any {
	now := time.Now().UTC()
	return []any{id, []byte(`{"ru":"q","en":"q"}`), []byte(`{"ru":"a","en":"a"}`), 10, now, now}
}

func stepColsH() []string {
	return []string{"id", "num", "title", "description", "sort_order", "created_at", "updated_at"}
}
func stepRowH(id uuid.UUID) []any {
	now := time.Now().UTC()
	return []any{id, "01", []byte(`{"ru":"t","en":"t"}`), []byte(`{"ru":"d","en":"d"}`), 10, now, now}
}

func TestServicesHandler_NewWithRevalidator(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	h2 := h.WithRevalidator(nil)
	if h2 == nil || h2 != h {
		t.Fatal("WithRevalidator should return same handler")
	}
}

func TestServicesHandler_GetPageData_OK(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)

	pool.ExpectQuery(`FROM services`).
		WillReturnRows(pgxmock.NewRows(svcCols()).AddRow(svcRow(uuid.New(), "x")...))
	pool.ExpectQuery(`FROM service_faqs`).
		WillReturnRows(pgxmock.NewRows(faqColsH()).AddRow(faqRowH(uuid.New())...))
	pool.ExpectQuery(`FROM service_process_steps`).
		WillReturnRows(pgxmock.NewRows(stepColsH()).AddRow(stepRowH(uuid.New())...))

	req := httptest.NewRequest(http.MethodGet, "/api/services", nil)
	rec := httptest.NewRecorder()
	h.GetPageData(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestServicesHandler_GetPageData_DBError(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	pool.ExpectQuery(`FROM services`).WillReturnError(errors.New("db"))

	req := httptest.NewRequest(http.MethodGet, "/api/services", nil)
	rec := httptest.NewRecorder()
	h.GetPageData(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_ListServices_OK(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	pool.ExpectQuery(`FROM services ORDER BY sort_order`).
		WillReturnRows(pgxmock.NewRows(svcCols()).AddRow(svcRow(uuid.New(), "x")...))

	req := httptest.NewRequest(http.MethodGet, "/api/services/list", nil)
	rec := httptest.NewRecorder()
	h.ListServices(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: %d", rec.Code)
	}
	var out []map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	if len(out) != 1 {
		t.Errorf("len: %d", len(out))
	}
}

func TestServicesHandler_ListServices_DBError(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	pool.ExpectQuery(`FROM services`).WillReturnError(errors.New("db"))

	req := httptest.NewRequest(http.MethodGet, "/api/services/list", nil)
	rec := httptest.NewRecorder()
	h.ListServices(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_CreateService_OK(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	pool.ExpectQuery(`INSERT INTO services`).
		WithArgs(testutil.AnyArgs(14)...).
		WillReturnRows(pgxmock.NewRows(svcCols()).AddRow(svcRow(uuid.New(), "telegram")...))

	body := `{"slug":"telegram","num":"01","iconKey":"bot","visualKey":"terminal","accent":"#000","title":{"ru":"t","en":"t"},"lead":{"ru":"l","en":"l"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/services", strReader(body))
	rec := httptest.NewRecorder()
	h.CreateService(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status: %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestServicesHandler_CreateService_BadJSON(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/services", strReader(`{not json`))
	rec := httptest.NewRecorder()
	h.CreateService(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_CreateService_ValidationError(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/services", strReader(`{"slug":"a"}`))
	rec := httptest.NewRecorder()
	h.CreateService(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_CreateService_DBError(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	pool.ExpectQuery(`INSERT INTO services`).
		WithArgs(testutil.AnyArgs(14)...).
		WillReturnError(errors.New("db"))

	body := `{"slug":"abc","num":"01","iconKey":"bot","visualKey":"terminal","accent":"#000","title":{"ru":"t","en":"t"},"lead":{"ru":"l","en":"l"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/services", strReader(body))
	rec := httptest.NewRecorder()
	h.CreateService(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_UpdateService_OK(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM services WHERE id`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows(svcCols()).AddRow(svcRow(id, "old")...))
	pool.ExpectQuery(`UPDATE services SET`).
		WithArgs(testutil.AnyArgs(15)...).
		WillReturnRows(pgxmock.NewRows(svcCols()).AddRow(svcRow(id, "new")...))

	req := chiRequest(http.MethodPut, "/api/services/"+id.String(),
		map[string]string{"id": id.String()},
		`{"slug":"new"}`)
	rec := httptest.NewRecorder()
	h.UpdateService(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestServicesHandler_UpdateService_BadID(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	req := chiRequest(http.MethodPut, "/api/services/not-uuid",
		map[string]string{"id": "not-uuid"}, `{}`)
	rec := httptest.NewRecorder()
	h.UpdateService(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_UpdateService_BadJSON(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	id := uuid.New()
	req := chiRequest(http.MethodPut, "/api/services/"+id.String(),
		map[string]string{"id": id.String()}, `{not json`)
	rec := httptest.NewRecorder()
	h.UpdateService(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_UpdateService_NotFound(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM services WHERE id`).
		WithArgs(id).
		WillReturnError(pgxNoRows)

	req := chiRequest(http.MethodPut, "/api/services/"+id.String(),
		map[string]string{"id": id.String()}, `{"slug":"new"}`)
	rec := httptest.NewRecorder()
	h.UpdateService(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_DeleteService_OK(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectExec(`DELETE FROM services`).
		WithArgs(id).
		WillReturnResult(pgconn.NewCommandTag("DELETE 1"))

	req := chiRequest(http.MethodDelete, "/api/services/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.DeleteService(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_DeleteService_BadID(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	req := chiRequest(http.MethodDelete, "/api/services/x",
		map[string]string{"id": "x"}, "")
	rec := httptest.NewRecorder()
	h.DeleteService(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_DeleteService_NotFound(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectExec(`DELETE FROM services`).
		WithArgs(id).
		WillReturnResult(pgconn.NewCommandTag("DELETE 0"))

	req := chiRequest(http.MethodDelete, "/api/services/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.DeleteService(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_DeleteService_DBError(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectExec(`DELETE FROM services`).WillReturnError(errors.New("db"))

	req := chiRequest(http.MethodDelete, "/api/services/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.DeleteService(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_ListFaqs_OK(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	pool.ExpectQuery(`FROM service_faqs ORDER BY sort_order`).
		WillReturnRows(pgxmock.NewRows(faqColsH()).AddRow(faqRowH(uuid.New())...))

	req := httptest.NewRequest(http.MethodGet, "/api/services/faqs", nil)
	rec := httptest.NewRecorder()
	h.ListFaqs(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_ListFaqs_DBError(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	pool.ExpectQuery(`FROM service_faqs`).WillReturnError(errors.New("db"))

	req := httptest.NewRequest(http.MethodGet, "/api/services/faqs", nil)
	rec := httptest.NewRecorder()
	h.ListFaqs(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_CreateFaq_OK(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	pool.ExpectQuery(`INSERT INTO service_faqs`).
		WithArgs(testutil.AnyArgs(3)...).
		WillReturnRows(pgxmock.NewRows(faqColsH()).AddRow(faqRowH(uuid.New())...))

	body := `{"question":{"ru":"q","en":"q"},"answer":{"ru":"a","en":"a"},"order":10}`
	req := httptest.NewRequest(http.MethodPost, "/api/services/faqs", strReader(body))
	rec := httptest.NewRecorder()
	h.CreateFaq(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status: %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestServicesHandler_CreateFaq_BadJSON(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/services/faqs", strReader(`{not`))
	rec := httptest.NewRecorder()
	h.CreateFaq(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_CreateFaq_Validation(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/services/faqs", strReader(`{}`))
	rec := httptest.NewRecorder()
	h.CreateFaq(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_CreateFaq_DBError(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	pool.ExpectQuery(`INSERT INTO service_faqs`).WillReturnError(errors.New("db"))

	body := `{"question":{"ru":"q"},"answer":{"ru":"a"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/services/faqs", strReader(body))
	rec := httptest.NewRecorder()
	h.CreateFaq(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_UpdateFaq_OK(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM service_faqs WHERE id`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows(faqColsH()).AddRow(faqRowH(id)...))
	pool.ExpectQuery(`UPDATE service_faqs`).
		WithArgs(testutil.AnyArgs(4)...).
		WillReturnRows(pgxmock.NewRows(faqColsH()).AddRow(faqRowH(id)...))

	req := chiRequest(http.MethodPut, "/api/services/faqs/"+id.String(),
		map[string]string{"id": id.String()}, `{"order":99}`)
	rec := httptest.NewRecorder()
	h.UpdateFaq(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_UpdateFaq_BadID(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	req := chiRequest(http.MethodPut, "/api/services/faqs/x",
		map[string]string{"id": "x"}, `{}`)
	rec := httptest.NewRecorder()
	h.UpdateFaq(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_UpdateFaq_BadJSON(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	id := uuid.New()
	req := chiRequest(http.MethodPut, "/api/services/faqs/"+id.String(),
		map[string]string{"id": id.String()}, `{not`)
	rec := httptest.NewRecorder()
	h.UpdateFaq(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_UpdateFaq_NotFound(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM service_faqs`).
		WithArgs(id).
		WillReturnError(pgxNoRows)

	req := chiRequest(http.MethodPut, "/api/services/faqs/"+id.String(),
		map[string]string{"id": id.String()}, `{"order":1}`)
	rec := httptest.NewRecorder()
	h.UpdateFaq(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestServicesHandler_DeleteFaq_OK(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectExec(`DELETE FROM service_faqs`).WithArgs(id).
		WillReturnResult(pgconn.NewCommandTag("DELETE 1"))

	req := chiRequest(http.MethodDelete, "/api/services/faqs/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.DeleteFaq(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_DeleteFaq_BadID(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	req := chiRequest(http.MethodDelete, "/api/services/faqs/x",
		map[string]string{"id": "x"}, "")
	rec := httptest.NewRecorder()
	h.DeleteFaq(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_DeleteFaq_NotFound(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectExec(`DELETE FROM service_faqs`).WithArgs(id).
		WillReturnResult(pgconn.NewCommandTag("DELETE 0"))

	req := chiRequest(http.MethodDelete, "/api/services/faqs/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.DeleteFaq(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_DeleteFaq_DBError(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectExec(`DELETE FROM service_faqs`).WillReturnError(errors.New("db"))

	req := chiRequest(http.MethodDelete, "/api/services/faqs/"+id.String(),
		map[string]string{"id": id.String()}, "")
	rec := httptest.NewRecorder()
	h.DeleteFaq(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_ListProcessSteps_OK(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	pool.ExpectQuery(`FROM service_process_steps ORDER BY sort_order`).
		WillReturnRows(pgxmock.NewRows(stepColsH()).AddRow(stepRowH(uuid.New())...))
	rec := httptest.NewRecorder()
	h.ListProcessSteps(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_ListProcessSteps_DBError(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	pool.ExpectQuery(`FROM service_process_steps`).WillReturnError(errors.New("db"))
	rec := httptest.NewRecorder()
	h.ListProcessSteps(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_CreateProcessStep_OK(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	pool.ExpectQuery(`INSERT INTO service_process_steps`).
		WithArgs(testutil.AnyArgs(4)...).
		WillReturnRows(pgxmock.NewRows(stepColsH()).AddRow(stepRowH(uuid.New())...))
	body := `{"num":"01","title":{"ru":"t","en":"t"},"description":{"ru":"d","en":"d"},"order":10}`
	rec := httptest.NewRecorder()
	h.CreateProcessStep(rec, httptest.NewRequest(http.MethodPost, "/", strReader(body)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("status: %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestServicesHandler_CreateProcessStep_BadJSON(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	rec := httptest.NewRecorder()
	h.CreateProcessStep(rec, httptest.NewRequest(http.MethodPost, "/", strReader(`{not`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_CreateProcessStep_Validation(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	rec := httptest.NewRecorder()
	h.CreateProcessStep(rec, httptest.NewRequest(http.MethodPost, "/", strReader(`{}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_CreateProcessStep_DBError(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	pool.ExpectQuery(`INSERT INTO service_process_steps`).
		WithArgs(testutil.AnyArgs(4)...).
		WillReturnError(errors.New("db"))
	body := `{"num":"01","title":{"ru":"t"},"description":{"ru":"d"}}`
	rec := httptest.NewRecorder()
	h.CreateProcessStep(rec, httptest.NewRequest(http.MethodPost, "/", strReader(body)))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_UpdateProcessStep_OK(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM service_process_steps WHERE id`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows(stepColsH()).AddRow(stepRowH(id)...))
	pool.ExpectQuery(`UPDATE service_process_steps`).
		WithArgs(testutil.AnyArgs(5)...).
		WillReturnRows(pgxmock.NewRows(stepColsH()).AddRow(stepRowH(id)...))
	rec := httptest.NewRecorder()
	h.UpdateProcessStep(rec, chiRequest(http.MethodPut, "/", map[string]string{"id": id.String()}, `{"order":2}`))
	if rec.Code != http.StatusOK {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_UpdateProcessStep_BadID(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	rec := httptest.NewRecorder()
	h.UpdateProcessStep(rec, chiRequest(http.MethodPut, "/", map[string]string{"id": "x"}, `{}`))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_UpdateProcessStep_BadJSON(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	id := uuid.New()
	rec := httptest.NewRecorder()
	h.UpdateProcessStep(rec, chiRequest(http.MethodPut, "/", map[string]string{"id": id.String()}, `{not`))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_UpdateProcessStep_NotFound(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM service_process_steps`).
		WithArgs(id).
		WillReturnError(pgxNoRows)
	rec := httptest.NewRecorder()
	h.UpdateProcessStep(rec, chiRequest(http.MethodPut, "/", map[string]string{"id": id.String()}, `{"order":1}`))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_DeleteProcessStep_OK(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectExec(`DELETE FROM service_process_steps`).
		WithArgs(id).
		WillReturnResult(pgconn.NewCommandTag("DELETE 1"))
	rec := httptest.NewRecorder()
	h.DeleteProcessStep(rec, chiRequest(http.MethodDelete, "/", map[string]string{"id": id.String()}, ""))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_DeleteProcessStep_BadID(t *testing.T) {
	t.Parallel()
	h, _ := newSvcHandler(t)
	rec := httptest.NewRecorder()
	h.DeleteProcessStep(rec, chiRequest(http.MethodDelete, "/", map[string]string{"id": "x"}, ""))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_DeleteProcessStep_NotFound(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectExec(`DELETE FROM service_process_steps`).
		WithArgs(id).
		WillReturnResult(pgconn.NewCommandTag("DELETE 0"))
	rec := httptest.NewRecorder()
	h.DeleteProcessStep(rec, chiRequest(http.MethodDelete, "/", map[string]string{"id": id.String()}, ""))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestServicesHandler_DeleteProcessStep_DBError(t *testing.T) {
	t.Parallel()
	h, pool := newSvcHandler(t)
	id := uuid.New()
	pool.ExpectExec(`DELETE FROM service_process_steps`).WillReturnError(errors.New("db"))
	rec := httptest.NewRecorder()
	h.DeleteProcessStep(rec, chiRequest(http.MethodDelete, "/", map[string]string{"id": id.String()}, ""))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", rec.Code)
	}
}
