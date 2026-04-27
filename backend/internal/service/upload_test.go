package service

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/portfolio/backend/internal/testutil"
)

func pngHeader() []byte {
	return []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
}

func jpegHeader() []byte {
	return []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 'J', 'F', 'I', 'F'}
}

func makeMultipartHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
	part, err := mw.CreatePart(h)
	if err != nil {
		t.Fatalf("create part: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write part: %v", err)
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close mw: %v", err)
	}

	mr := multipart.NewReader(&buf, mw.Boundary())
	form, err := mr.ReadForm(10 << 20)
	if err != nil {
		t.Fatalf("read form: %v", err)
	}
	t.Cleanup(func() { _ = form.RemoveAll() })
	files := form.File["file"]
	if len(files) == 0 {
		t.Fatal("no files parsed")
	}
	return files[0]
}

func TestUploadService_Upload_HappyPathPNG(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	db, pool := testutil.NewMockDB(t)
	svc := NewUploadService(db, dir)

	pngBytes := append(pngHeader(), make([]byte, 200)...)
	header := makeMultipartHeader(t, "my pic.png", pngBytes)

	id := uuid.New()
	now := time.Now().UTC()
	pool.ExpectQuery(`INSERT INTO media`).
		WithArgs(testutil.AnyArgs(5)...).
		WillReturnRows(pgxmock.NewRows([]string{"id", "filename", "original_name", "mime_type", "size", "url", "created_at"}).
			AddRow(id, "stored.png", "my pic.png", "image/png", int64(len(pngBytes)), "/uploads/stored.png", now))

	got, err := svc.Upload(context.Background(), header)
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if got.MimeType != "image/png" {
		t.Errorf("mime: got %q", got.MimeType)
	}
	if got.OriginalName != "my pic.png" {
		t.Errorf("original: got %q", got.OriginalName)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 file on disk, got %d", len(entries))
	}
	if filepath.Ext(entries[0].Name()) != ".png" {
		t.Errorf("ext: got %q, want .png", filepath.Ext(entries[0].Name()))
	}
}

func TestUploadService_Upload_TooLarge(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	db, _ := testutil.NewMockDB(t)
	svc := NewUploadService(db, dir)

	header := makeMultipartHeader(t, "big.png", pngHeader())
	header.Size = MaxUploadSize + 1

	_, err := svc.Upload(context.Background(), header)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUploadService_Upload_UnsupportedMIME(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	db, _ := testutil.NewMockDB(t)
	svc := NewUploadService(db, dir)

	header := makeMultipartHeader(t, "script.sh", []byte("#!/bin/sh\necho hi\n"))

	_, err := svc.Upload(context.Background(), header)
	if err == nil {
		t.Fatal("expected error")
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("expected no files on disk, got %d", len(entries))
	}
}

func TestUploadService_Upload_DBInsertErrorCleansUpFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	db, pool := testutil.NewMockDB(t)
	svc := NewUploadService(db, dir)

	header := makeMultipartHeader(t, "x.jpg", append(jpegHeader(), make([]byte, 100)...))

	pool.ExpectQuery(`INSERT INTO media`).
		WithArgs(testutil.AnyArgs(5)...).
		WillReturnError(errors.New("db write fail"))

	_, err := svc.Upload(context.Background(), header)
	if err == nil {
		t.Fatal("expected error")
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("expected cleanup on db error; found %d files", len(entries))
	}
}

func TestUploadService_List_Empty(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewUploadService(db, t.TempDir())

	pool.ExpectQuery(`SELECT id, filename, original_name, mime_type, size, url, created_at FROM media`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "filename", "original_name", "mime_type", "size", "url", "created_at"}))

	got, err := svc.List(context.Background())
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

func TestUploadService_List_Rows(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewUploadService(db, t.TempDir())

	now := time.Now().UTC()
	pool.ExpectQuery(`SELECT id, filename, original_name, mime_type, size, url, created_at FROM media`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "filename", "original_name", "mime_type", "size", "url", "created_at"}).
			AddRow(uuid.New(), "a.png", "a.png", "image/png", int64(1), "/uploads/a.png", now).
			AddRow(uuid.New(), "b.jpg", "b.jpg", "image/jpeg", int64(2), "/uploads/b.jpg", now))

	got, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len: got %d, want 2", len(got))
	}
}

func TestUploadService_List_QueryError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewUploadService(db, t.TempDir())

	pool.ExpectQuery(`SELECT id, filename, original_name, mime_type, size, url, created_at FROM media`).
		WillReturnError(errors.New("db err"))

	_, err := svc.List(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUploadService_List_ScanError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewUploadService(db, t.TempDir())

	rows := pgxmock.NewRows([]string{"id", "filename", "original_name", "mime_type", "size", "url", "created_at"}).
		AddRow(uuid.New(), "a.png", "a.png", "image/png", int64(1), "/uploads/a.png", time.Now().UTC())
	rows.RowError(0, errors.New("scan boom"))
	pool.ExpectQuery(`SELECT id, filename, original_name, mime_type, size, url, created_at FROM media`).
		WillReturnRows(rows)

	_, err := svc.List(context.Background())
	if err == nil {
		t.Fatal("expected scan error")
	}
}

func TestUploadService_List_RowsErr(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewUploadService(db, t.TempDir())

	now := time.Now().UTC()
	rows := pgxmock.NewRows([]string{"id", "filename", "original_name", "mime_type", "size", "url", "created_at"}).
		AddRow(uuid.New(), "a", "a", "image/png", int64(1), "/a", now).
		RowError(0, errors.New("iter fail"))

	pool.ExpectQuery(`SELECT id, filename, original_name, mime_type, size, url, created_at FROM media`).
		WillReturnRows(rows)

	_, err := svc.List(context.Background())
	if err == nil {
		t.Fatal("expected iterate error")
	}
}

func TestUploadService_Delete_Success(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	db, pool := testutil.NewMockDB(t)
	svc := NewUploadService(db, dir)

	filename := "abc.png"
	filePath := filepath.Join(dir, filename)
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	id := uuid.New()
	pool.ExpectQuery(`SELECT filename FROM media WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{"filename"}).AddRow(filename))
	pool.ExpectExec(`DELETE FROM media WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	if err := svc.Delete(context.Background(), id); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Errorf("file should be removed, stat err=%v", err)
	}
}

func TestUploadService_Delete_FileAlreadyMissing(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	db, pool := testutil.NewMockDB(t)
	svc := NewUploadService(db, dir)

	id := uuid.New()
	pool.ExpectQuery(`SELECT filename FROM media WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{"filename"}).AddRow("ghost.png"))
	pool.ExpectExec(`DELETE FROM media WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	if err := svc.Delete(context.Background(), id); err != nil {
		t.Errorf("Delete when file missing should still succeed, got %v", err)
	}
}

func TestUploadService_Delete_NotFoundOnSelect(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewUploadService(db, t.TempDir())

	id := uuid.New()
	pool.ExpectQuery(`SELECT filename FROM media WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(pgx.ErrNoRows)

	err := svc.Delete(context.Background(), id)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("got %v, want pgx.ErrNoRows", err)
	}
}

func TestUploadService_Delete_SelectOtherError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewUploadService(db, t.TempDir())

	id := uuid.New()
	pool.ExpectQuery(`SELECT filename FROM media WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(errors.New("db issue"))

	err := svc.Delete(context.Background(), id)
	if err == nil {
		t.Fatal("expected error")
	}
	if errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("should not be pgx.ErrNoRows")
	}
}

func TestUploadService_Delete_ExecError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewUploadService(db, t.TempDir())

	id := uuid.New()
	pool.ExpectQuery(`SELECT filename FROM media WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{"filename"}).AddRow("x.png"))
	pool.ExpectExec(`DELETE FROM media WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(errors.New("delete fail"))

	if err := svc.Delete(context.Background(), id); err == nil {
		t.Fatal("expected error")
	}
}

func TestUploadService_Delete_NoRowsAffected(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewUploadService(db, t.TempDir())

	id := uuid.New()
	pool.ExpectQuery(`SELECT filename FROM media WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{"filename"}).AddRow("x.png"))
	pool.ExpectExec(`DELETE FROM media WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err := svc.Delete(context.Background(), id)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("got %v, want pgx.ErrNoRows", err)
	}
}

func TestCleanOriginalName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no path", "file.png", "file.png"},
		{"unix path", "/tmp/foo/bar.png", "bar.png"},
		{"traversal", "../../etc/passwd", "passwd"},
		{"trailing slash", "dir/", "dir"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := cleanOriginalName(tc.in)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
