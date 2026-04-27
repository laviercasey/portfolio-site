package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/portfolio/backend/internal/database"
	"github.com/portfolio/backend/internal/model"
)

const MaxUploadSize = 50 << 20

var allowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"application/pdf": true,
	"video/mp4":       true,
	"video/webm":      true,
}

var mimeToExt = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/gif":       ".gif",
	"image/webp":      ".webp",
	"application/pdf": ".pdf",
	"video/mp4":       ".mp4",
	"video/webm":      ".webm",
}

type UploadService struct {
	db        *database.DB
	uploadDir string
}

func NewUploadService(db *database.DB, uploadDir string) *UploadService {
	return &UploadService{db: db, uploadDir: uploadDir}
}

func (s *UploadService) Upload(ctx context.Context, header *multipart.FileHeader) (*model.Media, error) {
	if header.Size > MaxUploadSize {
		return nil, fmt.Errorf("file too large: %d bytes exceeds limit of %d", header.Size, MaxUploadSize)
	}

	src, err := header.Open()
	if err != nil {
		return nil, fmt.Errorf("open uploaded file: %w", err)
	}
	defer src.Close()

	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("read file header: %w", err)
	}
	mimeType := http.DetectContentType(buf[:n])
	if !allowedMimeTypes[mimeType] {
		return nil, fmt.Errorf("unsupported file type: %s", mimeType)
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek file: %w", err)
	}

	if err := os.MkdirAll(s.uploadDir, 0o755); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}

	ext := mimeToExt[mimeType]
	fileID := uuid.New()
	filename := fileID.String() + ext

	destPath := filepath.Join(s.uploadDir, filepath.Base(filename))

	dst, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("create destination file: %w", err)
	}

	written, err := io.Copy(dst, src)
	if err != nil {
		_ = dst.Close()
		if rerr := os.Remove(destPath); rerr != nil && !errors.Is(rerr, os.ErrNotExist) {
			slog.Error("failed to remove orphaned upload", "path", destPath, "error", rerr)
		}
		return nil, fmt.Errorf("write file: %w", err)
	}
	if err := dst.Close(); err != nil {
		if rerr := os.Remove(destPath); rerr != nil && !errors.Is(rerr, os.ErrNotExist) {
			slog.Error("failed to remove orphaned upload", "path", destPath, "error", rerr)
		}
		return nil, fmt.Errorf("close destination file: %w", err)
	}

	urlPath := "/uploads/" + filename

	var m model.Media
	err = s.db.Pool.QueryRow(ctx,
		`INSERT INTO media (filename, original_name, mime_type, size, url)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, filename, original_name, mime_type, size, url, created_at`,
		filename, cleanOriginalName(header.Filename), mimeType, written, urlPath,
	).Scan(&m.ID, &m.Filename, &m.OriginalName, &m.MimeType, &m.Size, &m.URL, &m.CreatedAt)
	if err != nil {
		if rerr := os.Remove(destPath); rerr != nil && !errors.Is(rerr, os.ErrNotExist) {
			slog.Error("failed to remove orphaned upload", "path", destPath, "error", rerr)
		}
		return nil, fmt.Errorf("insert media record: %w", err)
	}

	return &m, nil
}

func (s *UploadService) List(ctx context.Context) ([]model.Media, error) {
	rows, err := s.db.Pool.Query(ctx,
		`SELECT id, filename, original_name, mime_type, size, url, created_at FROM media ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list media: %w", err)
	}
	defer rows.Close()

	var media []model.Media
	for rows.Next() {
		var m model.Media
		if err := rows.Scan(&m.ID, &m.Filename, &m.OriginalName, &m.MimeType, &m.Size, &m.URL, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan media: %w", err)
		}
		media = append(media, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate media: %w", err)
	}
	if media == nil {
		media = []model.Media{}
	}
	return media, nil
}

func (s *UploadService) Delete(ctx context.Context, id uuid.UUID) error {
	var filename string
	err := s.db.Pool.QueryRow(ctx,
		`SELECT filename FROM media WHERE id = $1`, id,
	).Scan(&filename)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("get media for delete: %w", err)
	}

	tag, err := s.db.Pool.Exec(ctx, `DELETE FROM media WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete media record: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	filePath := filepath.Join(s.uploadDir, filepath.Base(filename))
	if err := os.Remove(filePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Error("failed to remove orphaned upload", "path", filePath, "error", err)
	}

	return nil
}

func cleanOriginalName(name string) string {
	return filepath.Base(name)
}
