package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/portfolio/backend/internal/model"
	"github.com/portfolio/backend/internal/testutil"
)

func inquiryColumnsSlice() []string {
	return []string{
		"id", "name", "email", "company", "telegram",
		"inquiry_type", "budget", "message", "status",
		"admin_notes", "created_at", "updated_at",
	}
}

func TestInquiryService_List_AllStatuses(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewInquiryService(db)

	id := uuid.New()
	now := time.Now().UTC()
	company := "Acme"
	telegram := "@user"
	budget := "5k"
	notes := "note"

	pool.ExpectQuery(`SELECT .+ FROM inquiries ORDER BY created_at DESC`).
		WillReturnRows(pgxmock.NewRows(inquiryColumnsSlice()).
			AddRow(id, "John", "j@x.com", &company, &telegram, "freelance", &budget, "Hi there!", "new", &notes, now, now).
			AddRow(uuid.New(), "Jane", "j2@x.com", nil, nil, "fulltime", nil, "Msg", "read", nil, now, now))

	got, err := svc.List(context.Background(), "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len: got %d, want 2", len(got))
	}
	if got[0].Name != "John" {
		t.Errorf("[0].Name: got %q", got[0].Name)
	}
	if got[0].Company != "Acme" {
		t.Errorf("[0].Company: got %q, want Acme", got[0].Company)
	}
	if got[0].Telegram != "@user" {
		t.Errorf("[0].Telegram: got %q", got[0].Telegram)
	}
	if got[0].Budget != "5k" {
		t.Errorf("[0].Budget: got %q", got[0].Budget)
	}
	if got[0].AdminNotes != "note" {
		t.Errorf("[0].AdminNotes: got %q", got[0].AdminNotes)
	}
	if got[1].Company != "" {
		t.Errorf("[1].Company nullable: got %q", got[1].Company)
	}
	if got[1].Telegram != "" {
		t.Errorf("[1].Telegram nullable: got %q", got[1].Telegram)
	}
	if got[1].Budget != "" {
		t.Errorf("[1].Budget nullable: got %q", got[1].Budget)
	}
	if got[1].AdminNotes != "" {
		t.Errorf("[1].AdminNotes nullable: got %q", got[1].AdminNotes)
	}
}

func TestInquiryService_List_FilteredByStatus(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewInquiryService(db)

	pool.ExpectQuery(`SELECT .+ FROM inquiries WHERE status = \$1 ORDER BY created_at DESC`).
		WithArgs("new").
		WillReturnRows(pgxmock.NewRows(inquiryColumnsSlice()))

	got, err := svc.List(context.Background(), "new")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil empty slice")
	}
	if len(got) != 0 {
		t.Errorf("len: got %d", len(got))
	}
}

func TestInquiryService_List_QueryError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewInquiryService(db)

	pool.ExpectQuery(`SELECT .+ FROM inquiries`).
		WillReturnError(errors.New("db down"))

	_, err := svc.List(context.Background(), "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInquiryService_List_ScanError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewInquiryService(db)

	rows := pgxmock.NewRows(inquiryColumnsSlice()).
		AddRow(uuid.New(), "n", "e@x", (*string)(nil), (*string)(nil), "t", (*string)(nil), "m", "new", (*string)(nil), time.Now().UTC(), time.Now().UTC())
	rows.RowError(0, errors.New("scan boom"))
	pool.ExpectQuery(`SELECT .+ FROM inquiries`).WillReturnRows(rows)

	_, err := svc.List(context.Background(), "")
	if err == nil {
		t.Fatal("expected scan error")
	}
}

func TestInquiryService_List_RowsErr(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewInquiryService(db)

	id := uuid.New()
	now := time.Now().UTC()
	rows := pgxmock.NewRows(inquiryColumnsSlice()).
		AddRow(id, "n", "e@x.com", nil, nil, "freelance", nil, "m", "new", nil, now, now).
		RowError(0, errors.New("iter fail"))

	pool.ExpectQuery(`SELECT .+ FROM inquiries`).WillReturnRows(rows)

	_, err := svc.List(context.Background(), "")
	if err == nil {
		t.Fatal("expected iterate error")
	}
}

func TestInquiryService_GetByID_Success(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewInquiryService(db)

	id := uuid.New()
	now := time.Now().UTC()
	pool.ExpectQuery(`SELECT .+ FROM inquiries WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows(inquiryColumnsSlice()).
			AddRow(id, "X", "x@x.com", nil, nil, "other", nil, "hey", "new", nil, now, now))

	got, err := svc.GetByID(context.Background(), id)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.ID != id {
		t.Errorf("ID mismatch")
	}
}

func TestInquiryService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewInquiryService(db)

	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM inquiries WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(errors.New("no rows"))

	_, err := svc.GetByID(context.Background(), id)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInquiryService_Create_Success(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewInquiryService(db)

	input := model.CreateInquiryInput{
		Name:     "Alice",
		Email:    "a@x.com",
		Company:  "Acme",
		Telegram: "@a",
		Type:     "freelance",
		Budget:   "10k",
		Message:  "Hello, I'd like to work together.",
	}
	id := uuid.New()
	now := time.Now().UTC()
	company := "Acme"
	telegram := "@a"
	budget := "10k"

	pool.ExpectQuery(`INSERT INTO inquiries`).
		WithArgs("Alice", "a@x.com", &company, &telegram, "freelance", &budget, "Hello, I'd like to work together.").
		WillReturnRows(pgxmock.NewRows(inquiryColumnsSlice()).
			AddRow(id, "Alice", "a@x.com", &company, &telegram, "freelance", &budget, "Hello, I'd like to work together.", "new", nil, now, now))

	got, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != id {
		t.Errorf("ID")
	}
	if got.Status != "new" {
		t.Errorf("status: got %q", got.Status)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestInquiryService_Create_EmptyOptionalsAsNull(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewInquiryService(db)

	input := model.CreateInquiryInput{
		Name:    "Bob",
		Email:   "b@x.com",
		Type:    "other",
		Message: "hi there longer message",
	}
	id := uuid.New()
	now := time.Now().UTC()

	pool.ExpectQuery(`INSERT INTO inquiries`).
		WithArgs("Bob", "b@x.com", (*string)(nil), (*string)(nil), "other", (*string)(nil), "hi there longer message").
		WillReturnRows(pgxmock.NewRows(inquiryColumnsSlice()).
			AddRow(id, "Bob", "b@x.com", nil, nil, "other", nil, "hi there longer message", "new", nil, now, now))

	if _, err := svc.Create(context.Background(), input); err != nil {
		t.Fatalf("Create: %v", err)
	}
}

func TestInquiryService_Create_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewInquiryService(db)

	pool.ExpectQuery(`INSERT INTO inquiries`).
		WillReturnError(errors.New("constraint violation"))

	_, err := svc.Create(context.Background(), model.CreateInquiryInput{
		Name: "x", Email: "x@x.com", Type: "other", Message: "hi",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInquiryService_UpdateStatus_ValidStatuses(t *testing.T) {
	t.Parallel()

	statuses := []string{"new", "read", "replied", "archived"}
	for _, s := range statuses {
		s := s
		t.Run(s, func(t *testing.T) {
			t.Parallel()
			db, pool := testutil.NewMockDB(t)
			svc := NewInquiryService(db)

			id := uuid.New()
			now := time.Now().UTC()
			notes := "admin note"

			pool.ExpectQuery(`UPDATE inquiries SET status = \$2`).
				WithArgs(id, s, "admin note").
				WillReturnRows(pgxmock.NewRows(inquiryColumnsSlice()).
					AddRow(id, "n", "e@x.com", nil, nil, "other", nil, "m", s, &notes, now, now))

			got, err := svc.UpdateStatus(context.Background(), id, s, "admin note")
			if err != nil {
				t.Fatalf("UpdateStatus: %v", err)
			}
			if got.Status != s {
				t.Errorf("status: got %q, want %q", got.Status, s)
			}
			if got.AdminNotes != "admin note" {
				t.Errorf("notes: got %q", got.AdminNotes)
			}
		})
	}
}

func TestInquiryService_UpdateStatus_InvalidStatus(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewInquiryService(db)

	bads := []string{"", "pending", "unknown", "NEW", "deleted"}
	for _, s := range bads {
		s := s
		t.Run(s, func(t *testing.T) {
			t.Parallel()
			_, err := svc.UpdateStatus(context.Background(), uuid.New(), s, "")
			if err == nil {
				t.Fatalf("expected error for status %q", s)
			}
		})
	}
}

func TestInquiryService_UpdateStatus_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewInquiryService(db)

	id := uuid.New()
	pool.ExpectQuery(`UPDATE inquiries SET status`).
		WithArgs(id, "new", "").
		WillReturnError(errors.New("no rows"))

	_, err := svc.UpdateStatus(context.Background(), id, "new", "")
	if err == nil {
		t.Fatal("expected error")
	}
}
