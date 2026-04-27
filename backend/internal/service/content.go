package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/portfolio/backend/internal/database"
	"github.com/portfolio/backend/internal/model"
)

type ContentService struct {
	db *database.DB
}

func NewContentService(db *database.DB) *ContentService {
	return &ContentService{db: db}
}

func (s *ContentService) GetAll(ctx context.Context) ([]model.Content, error) {
	rows, err := s.db.Pool.Query(ctx, `SELECT section, data, updated_at FROM content ORDER BY section`)
	if err != nil {
		return nil, fmt.Errorf("list content: %w", err)
	}
	defer rows.Close()

	var contents []model.Content
	for rows.Next() {
		var c model.Content
		if err := rows.Scan(&c.Section, &c.Data, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan content: %w", err)
		}
		contents = append(contents, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate content: %w", err)
	}
	if contents == nil {
		contents = []model.Content{}
	}
	return contents, nil
}

func (s *ContentService) Upsert(ctx context.Context, section string, data json.RawMessage) (*model.Content, error) {
	var c model.Content
	err := s.db.Pool.QueryRow(ctx,
		`INSERT INTO content (section, data)
VALUES ($1, $2)
ON CONFLICT (section) DO UPDATE SET data = $2, updated_at = now()
RETURNING section, data, updated_at`,
		section, data,
	).Scan(&c.Section, &c.Data, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("upsert content: %w", err)
	}
	return &c, nil
}
