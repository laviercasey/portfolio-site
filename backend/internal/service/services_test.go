package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/portfolio/backend/internal/model"
	"github.com/portfolio/backend/internal/testutil"
)

func newServiceRow(id uuid.UUID, slug string, sortOrder int) []any {
	now := time.Now().UTC()
	return []any{
		id, slug, "01", "bot", "terminal", "#5eb3ff",
		[]byte(`{"ru":"T","en":"T"}`),
		[]byte(`{"ru":"L","en":"L"}`),
		[]byte(`{"ru":["b1"],"en":["b1"]}`),
		"stack", []byte(`{"ru":"1w","en":"1w"}`),
		"100", "$100",
		[]byte(`[{"slug":"x","name":"X"}]`),
		sortOrder, now, now,
	}
}

func serviceCols() []string {
	return []string{
		"id", "slug", "num", "icon_key", "visual_key", "accent",
		"title", "lead", "bullets", "stack", "timeline",
		"price_ru", "price_en", "case_projects", "sort_order",
		"created_at", "updated_at",
	}
}

func TestServicesService_ListServices_Empty(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	pool.ExpectQuery(`SELECT .+ FROM services ORDER BY sort_order`).
		WillReturnRows(pgxmock.NewRows(serviceCols()))

	got, err := svc.ListServices(context.Background())
	if err != nil {
		t.Fatalf("ListServices: %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestServicesService_ListServices_Rows(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	id1, id2 := uuid.New(), uuid.New()
	rows := pgxmock.NewRows(serviceCols()).
		AddRow(newServiceRow(id1, "a", 10)...).
		AddRow(newServiceRow(id2, "b", 20)...)
	pool.ExpectQuery(`SELECT .+ FROM services ORDER BY sort_order`).WillReturnRows(rows)

	got, err := svc.ListServices(context.Background())
	if err != nil {
		t.Fatalf("ListServices: %v", err)
	}
	if len(got) != 2 || got[0].Slug != "a" || got[1].Slug != "b" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestServicesService_ListServices_QueryError(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	pool.ExpectQuery(`SELECT .+ FROM services`).WillReturnError(errors.New("boom"))
	_, err := svc.ListServices(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestServicesService_ListServices_RowsErr(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	id := uuid.New()
	rows := pgxmock.NewRows(serviceCols()).
		AddRow(newServiceRow(id, "a", 10)...).
		RowError(0, errors.New("iter fail"))
	pool.ExpectQuery(`SELECT .+ FROM services`).WillReturnRows(rows)

	_, err := svc.ListServices(context.Background())
	if err == nil {
		t.Fatal("expected iterate error")
	}
}

func TestServicesService_CreateService_Success(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	id := uuid.New()
	pool.ExpectQuery(`INSERT INTO services`).
		WithArgs(testutil.AnyArgs(14)...).
		WillReturnRows(pgxmock.NewRows(serviceCols()).
			AddRow(newServiceRow(id, "telegram-bots", 10)...))

	got, err := svc.CreateService(context.Background(), model.CreateServiceInput{
		Slug: "telegram-bots", Num: "01",
		IconKey: "bot", VisualKey: "terminal", Accent: "#5eb3ff",
		Title:   model.I18n{"ru": "T", "en": "T"},
		Lead:    model.I18n{"ru": "L", "en": "L"},
		Bullets: model.ServiceBullets{Ru: []string{"b1"}, En: []string{"b1"}},
		Stack:   "stack", Timeline: model.I18n{"ru": "1w", "en": "1w"},
		PriceRu: "100", PriceEn: "$100",
		CaseProjects: []model.ServiceCaseProject{{Slug: "x", Name: "X"}},
		Order:        10,
	})
	if err != nil {
		t.Fatalf("CreateService: %v", err)
	}
	if got.Slug != "telegram-bots" {
		t.Errorf("slug: %q", got.Slug)
	}
}

func TestServicesService_CreateService_NilCaseProjects(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	id := uuid.New()
	pool.ExpectQuery(`INSERT INTO services`).
		WithArgs(testutil.AnyArgs(14)...).
		WillReturnRows(pgxmock.NewRows(serviceCols()).
			AddRow(newServiceRow(id, "x", 0)...))

	_, err := svc.CreateService(context.Background(), model.CreateServiceInput{
		Slug: "x", Num: "01", IconKey: "bot", VisualKey: "terminal", Accent: "#000",
		Title: model.I18n{"ru": "x", "en": "x"},
		Lead:  model.I18n{"ru": "l", "en": "l"},
	})
	if err != nil {
		t.Fatalf("CreateService: %v", err)
	}
}

func TestServicesService_UpdateService_Success(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM services WHERE id`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows(serviceCols()).AddRow(newServiceRow(id, "old", 5)...))
	pool.ExpectQuery(`UPDATE services SET`).
		WithArgs(testutil.AnyArgs(15)...).
		WillReturnRows(pgxmock.NewRows(serviceCols()).AddRow(newServiceRow(id, "new", 9)...))

	newSlug := "new"
	newOrder := 9
	got, err := svc.UpdateService(context.Background(), id, model.UpdateServiceInput{
		Slug:  &newSlug,
		Order: &newOrder,
	})
	if err != nil {
		t.Fatalf("UpdateService: %v", err)
	}
	if got.Slug != "new" {
		t.Errorf("slug: %q", got.Slug)
	}
}

func TestServicesService_UpdateService_AllFields(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM services WHERE id`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows(serviceCols()).AddRow(newServiceRow(id, "old", 5)...))
	pool.ExpectQuery(`UPDATE services SET`).
		WithArgs(testutil.AnyArgs(15)...).
		WillReturnRows(pgxmock.NewRows(serviceCols()).AddRow(newServiceRow(id, "new", 9)...))

	slug, num, ik, vk, accent, stack, priceRu, priceEn := "new", "02", "code", "browser", "#fff", "x", "1₽", "$1"
	order := 9
	title := model.I18n{"ru": "x"}
	lead := model.I18n{"ru": "y"}
	timeline := model.I18n{"ru": "1w"}
	bullets := model.ServiceBullets{Ru: []string{"a"}, En: []string{"a"}}
	cases := []model.ServiceCaseProject{{Slug: "p", Name: "P"}}

	_, err := svc.UpdateService(context.Background(), id, model.UpdateServiceInput{
		Slug: &slug, Num: &num, IconKey: &ik, VisualKey: &vk, Accent: &accent,
		Title: &title, Lead: &lead, Bullets: &bullets, Stack: &stack,
		Timeline: &timeline, PriceRu: &priceRu, PriceEn: &priceEn,
		CaseProjects: &cases, Order: &order,
	})
	if err != nil {
		t.Fatalf("UpdateService: %v", err)
	}
}

func TestServicesService_UpdateService_NotFound(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM services WHERE id`).WillReturnError(errors.New("no row"))

	_, err := svc.UpdateService(context.Background(), id, model.UpdateServiceInput{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestServicesService_DeleteService_Success(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	id := uuid.New()
	pool.ExpectExec(`DELETE FROM services WHERE id`).WithArgs(id).
		WillReturnResult(pgconn.NewCommandTag("DELETE 1"))

	if err := svc.DeleteService(context.Background(), id); err != nil {
		t.Fatalf("DeleteService: %v", err)
	}
}

func TestServicesService_DeleteService_NotFound(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	id := uuid.New()
	pool.ExpectExec(`DELETE FROM services WHERE id`).WithArgs(id).
		WillReturnResult(pgconn.NewCommandTag("DELETE 0"))

	if err := svc.DeleteService(context.Background(), id); err == nil {
		t.Fatal("expected ErrNoRows")
	}
}

func TestServicesService_DeleteService_DBError(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	pool.ExpectExec(`DELETE FROM services`).WillReturnError(errors.New("db down"))
	if err := svc.DeleteService(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}

func faqRow(id uuid.UUID, sortOrder int) []any {
	now := time.Now().UTC()
	return []any{id, []byte(`{"ru":"q","en":"q"}`), []byte(`{"ru":"a","en":"a"}`), sortOrder, now, now}
}
func faqCols() []string {
	return []string{"id", "question", "answer", "sort_order", "created_at", "updated_at"}
}

func TestServicesService_ListFaqs(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM service_faqs ORDER BY sort_order`).
		WillReturnRows(pgxmock.NewRows(faqCols()).AddRow(faqRow(id, 10)...))

	got, err := svc.ListFaqs(context.Background())
	if err != nil {
		t.Fatalf("ListFaqs: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len: %d", len(got))
	}
}

func TestServicesService_ListFaqs_QueryError(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	pool.ExpectQuery(`FROM service_faqs`).WillReturnError(errors.New("err"))
	_, err := svc.ListFaqs(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestServicesService_CreateFaq_Success(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	id := uuid.New()

	pool.ExpectQuery(`INSERT INTO service_faqs`).
		WithArgs(testutil.AnyArgs(3)...).
		WillReturnRows(pgxmock.NewRows(faqCols()).AddRow(faqRow(id, 10)...))

	_, err := svc.CreateFaq(context.Background(), model.CreateServiceFaqInput{
		Question: model.I18n{"ru": "q"}, Answer: model.I18n{"ru": "a"}, Order: 10,
	})
	if err != nil {
		t.Fatalf("CreateFaq: %v", err)
	}
}

func TestServicesService_UpdateFaq_Success(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	id := uuid.New()

	pool.ExpectQuery(`SELECT .+ FROM service_faqs WHERE id`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows(faqCols()).AddRow(faqRow(id, 5)...))
	pool.ExpectQuery(`UPDATE service_faqs`).
		WithArgs(testutil.AnyArgs(4)...).
		WillReturnRows(pgxmock.NewRows(faqCols()).AddRow(faqRow(id, 9)...))

	q := model.I18n{"ru": "q2"}
	a := model.I18n{"ru": "a2"}
	order := 9
	_, err := svc.UpdateFaq(context.Background(), id, model.UpdateServiceFaqInput{
		Question: &q, Answer: &a, Order: &order,
	})
	if err != nil {
		t.Fatalf("UpdateFaq: %v", err)
	}
}

func TestServicesService_UpdateFaq_NotFound(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	pool.ExpectQuery(`SELECT .+ FROM service_faqs`).WillReturnError(errors.New("nope"))
	_, err := svc.UpdateFaq(context.Background(), uuid.New(), model.UpdateServiceFaqInput{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestServicesService_DeleteFaq_Success(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	id := uuid.New()
	pool.ExpectExec(`DELETE FROM service_faqs`).WithArgs(id).
		WillReturnResult(pgconn.NewCommandTag("DELETE 1"))
	if err := svc.DeleteFaq(context.Background(), id); err != nil {
		t.Fatalf("DeleteFaq: %v", err)
	}
}

func TestServicesService_DeleteFaq_NotFound(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	pool.ExpectExec(`DELETE FROM service_faqs`).
		WillReturnResult(pgconn.NewCommandTag("DELETE 0"))
	if err := svc.DeleteFaq(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected ErrNoRows")
	}
}

func TestServicesService_DeleteFaq_DBError(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	pool.ExpectExec(`DELETE FROM service_faqs`).WillReturnError(errors.New("db"))
	if err := svc.DeleteFaq(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}

func stepRow(id uuid.UUID, sortOrder int) []any {
	now := time.Now().UTC()
	return []any{id, "01", []byte(`{"ru":"t","en":"t"}`), []byte(`{"ru":"d","en":"d"}`), sortOrder, now, now}
}
func stepCols() []string {
	return []string{"id", "num", "title", "description", "sort_order", "created_at", "updated_at"}
}

func TestServicesService_ListProcessSteps(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM service_process_steps ORDER BY sort_order`).
		WillReturnRows(pgxmock.NewRows(stepCols()).AddRow(stepRow(id, 10)...))
	got, err := svc.ListProcessSteps(context.Background())
	if err != nil {
		t.Fatalf("ListProcessSteps: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len: %d", len(got))
	}
}

func TestServicesService_ListProcessSteps_QueryError(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	pool.ExpectQuery(`FROM service_process_steps`).WillReturnError(errors.New("err"))
	_, err := svc.ListProcessSteps(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestServicesService_CreateProcessStep(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	id := uuid.New()
	pool.ExpectQuery(`INSERT INTO service_process_steps`).
		WithArgs(testutil.AnyArgs(4)...).
		WillReturnRows(pgxmock.NewRows(stepCols()).AddRow(stepRow(id, 10)...))

	_, err := svc.CreateProcessStep(context.Background(), model.CreateServiceProcessStepInput{
		Num: "01", Title: model.I18n{"ru": "t"}, Description: model.I18n{"ru": "d"}, Order: 10,
	})
	if err != nil {
		t.Fatalf("CreateProcessStep: %v", err)
	}
}

func TestServicesService_UpdateProcessStep(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM service_process_steps WHERE id`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows(stepCols()).AddRow(stepRow(id, 5)...))
	pool.ExpectQuery(`UPDATE service_process_steps`).
		WithArgs(testutil.AnyArgs(5)...).
		WillReturnRows(pgxmock.NewRows(stepCols()).AddRow(stepRow(id, 9)...))

	num, order := "02", 9
	title := model.I18n{"ru": "x"}
	desc := model.I18n{"ru": "y"}
	_, err := svc.UpdateProcessStep(context.Background(), id, model.UpdateServiceProcessStepInput{
		Num: &num, Title: &title, Description: &desc, Order: &order,
	})
	if err != nil {
		t.Fatalf("UpdateProcessStep: %v", err)
	}
}

func TestServicesService_UpdateProcessStep_NotFound(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	pool.ExpectQuery(`SELECT .+ FROM service_process_steps`).WillReturnError(errors.New("nope"))
	_, err := svc.UpdateProcessStep(context.Background(), uuid.New(), model.UpdateServiceProcessStepInput{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestServicesService_DeleteProcessStep_Success(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	id := uuid.New()
	pool.ExpectExec(`DELETE FROM service_process_steps`).WithArgs(id).
		WillReturnResult(pgconn.NewCommandTag("DELETE 1"))
	if err := svc.DeleteProcessStep(context.Background(), id); err != nil {
		t.Fatalf("DeleteProcessStep: %v", err)
	}
}

func TestServicesService_DeleteProcessStep_NotFound(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	pool.ExpectExec(`DELETE FROM service_process_steps`).
		WillReturnResult(pgconn.NewCommandTag("DELETE 0"))
	if err := svc.DeleteProcessStep(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected ErrNoRows")
	}
}

func TestServicesService_DeleteProcessStep_DBError(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)
	pool.ExpectExec(`DELETE FROM service_process_steps`).WillReturnError(errors.New("db"))
	if err := svc.DeleteProcessStep(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}

func TestServicesService_GetPageData(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	pool.ExpectQuery(`FROM services ORDER BY sort_order`).
		WillReturnRows(pgxmock.NewRows(serviceCols()).AddRow(newServiceRow(uuid.New(), "x", 1)...))
	pool.ExpectQuery(`FROM service_faqs ORDER BY sort_order`).
		WillReturnRows(pgxmock.NewRows(faqCols()).AddRow(faqRow(uuid.New(), 1)...))
	pool.ExpectQuery(`FROM service_process_steps ORDER BY sort_order`).
		WillReturnRows(pgxmock.NewRows(stepCols()).AddRow(stepRow(uuid.New(), 1)...))

	got, err := svc.GetPageData(context.Background())
	if err != nil {
		t.Fatalf("GetPageData: %v", err)
	}
	if len(got.Services) != 1 || len(got.Faqs) != 1 || len(got.ProcessSteps) != 1 {
		t.Errorf("counts: %d/%d/%d", len(got.Services), len(got.Faqs), len(got.ProcessSteps))
	}
}

func TestServicesService_GetPageData_ServicesError(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	pool.ExpectQuery(`FROM services`).WillReturnError(errors.New("db"))
	_, err := svc.GetPageData(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestServicesService_GetPageData_FaqsError(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	pool.ExpectQuery(`FROM services`).
		WillReturnRows(pgxmock.NewRows(serviceCols()))
	pool.ExpectQuery(`FROM service_faqs`).WillReturnError(errors.New("db"))

	_, err := svc.GetPageData(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestServicesService_GetPageData_StepsError(t *testing.T) {
	t.Parallel()
	db, pool := testutil.NewMockDB(t)
	svc := NewServicesService(db)

	pool.ExpectQuery(`FROM services`).WillReturnRows(pgxmock.NewRows(serviceCols()))
	pool.ExpectQuery(`FROM service_faqs`).WillReturnRows(pgxmock.NewRows(faqCols()))
	pool.ExpectQuery(`FROM service_process_steps`).WillReturnError(errors.New("db"))

	_, err := svc.GetPageData(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
