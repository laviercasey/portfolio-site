package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/portfolio/backend/internal/testutil"
)

func TestContentService_GetAll_Empty(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewContentService(db)

	pool.ExpectQuery(`SELECT section, data, updated_at FROM content`).
		WillReturnRows(pgxmock.NewRows([]string{"section", "data", "updated_at"}))

	got, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil empty slice")
	}
	if len(got) != 0 {
		t.Errorf("len: got %d, want 0", len(got))
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestContentService_GetAll_ManyRows(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewContentService(db)
	now := time.Now().UTC()

	pool.ExpectQuery(`SELECT section, data, updated_at FROM content`).
		WillReturnRows(pgxmock.NewRows([]string{"section", "data", "updated_at"}).
			AddRow("about", json.RawMessage(`{"title":"A"}`), now).
			AddRow("hero", json.RawMessage(`{"title":"H"}`), now))

	got, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len: got %d, want 2", len(got))
	}
	if got[0].Section != "about" || got[1].Section != "hero" {
		t.Errorf("sections: got %q and %q", got[0].Section, got[1].Section)
	}
}

func TestContentService_GetAll_QueryError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewContentService(db)

	pool.ExpectQuery(`SELECT section, data, updated_at FROM content`).
		WillReturnError(errors.New("connection refused"))

	_, err := svc.GetAll(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() == "" {
		t.Error("empty error message")
	}
}

func TestContentService_GetAll_ScanError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewContentService(db)

	rows := pgxmock.NewRows([]string{"section", "data", "updated_at"}).
		AddRow("homepage", []byte(`{}`), time.Now().UTC())
	rows.RowError(0, errors.New("scan boom"))
	pool.ExpectQuery(`SELECT section, data, updated_at FROM content`).WillReturnRows(rows)

	_, err := svc.GetAll(context.Background())
	if err == nil {
		t.Fatal("expected scan error")
	}
}

func TestContentService_GetAll_RowsErr(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewContentService(db)

	rows := pgxmock.NewRows([]string{"section", "data", "updated_at"}).
		AddRow("hero", json.RawMessage(`{}`), time.Now()).
		RowError(0, errors.New("iter fail"))

	pool.ExpectQuery(`SELECT section, data, updated_at FROM content`).WillReturnRows(rows)

	_, err := svc.GetAll(context.Background())
	if err == nil {
		t.Fatal("expected iterate error")
	}
}

func TestContentService_Upsert_Success(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewContentService(db)
	now := time.Now().UTC()
	data := json.RawMessage(`{"title":"Hello"}`)

	pool.ExpectQuery(`INSERT INTO content`).
		WithArgs("hero", data).
		WillReturnRows(pgxmock.NewRows([]string{"section", "data", "updated_at"}).
			AddRow("hero", data, now))

	got, err := svc.Upsert(context.Background(), "hero", data)
	if err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	if got.Section != "hero" {
		t.Errorf("section: got %q, want hero", got.Section)
	}
	if string(got.Data) != string(data) {
		t.Errorf("data mismatch: got %s, want %s", string(got.Data), string(data))
	}
	if !got.UpdatedAt.Equal(now) {
		t.Errorf("updatedAt: got %v, want %v", got.UpdatedAt, now)
	}
}

func TestContentService_Upsert_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewContentService(db)

	pool.ExpectQuery(`INSERT INTO content`).
		WithArgs("hero", json.RawMessage(`{}`)).
		WillReturnError(errors.New("constraint fail"))

	_, err := svc.Upsert(context.Background(), "hero", json.RawMessage(`{}`))
	if err == nil {
		t.Fatal("expected error")
	}
}
