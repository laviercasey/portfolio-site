package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/portfolio/backend/internal/model"
)

type publicationInput struct {
	Title     model.I18n `json:"title"`
	Journal   model.I18n `json:"journal"`
	Year      int        `json:"year"`
	DOI       string     `json:"doi"`
	URL       string     `json:"url"`
	Abstract  model.I18n `json:"abstract"`
	SortOrder int        `json:"sortOrder"`
}

func (s *CareerService) listPublications(ctx context.Context) ([]model.Publication, error) {
	rows, err := s.db.Pool.Query(ctx,
		`SELECT id, title, journal, year, doi, url, abstract FROM publications ORDER BY sort_order ASC, year DESC`)
	if err != nil {
		return nil, fmt.Errorf("list publications: %w", err)
	}
	defer rows.Close()

	var items []model.Publication
	for rows.Next() {
		p, err := scanPublication(rows)
		if err != nil {
			return nil, fmt.Errorf("scan publication: %w", err)
		}
		items = append(items, *p)
	}
	if items == nil {
		items = []model.Publication{}
	}
	return items, nil
}

func (s *CareerService) createPublication(ctx context.Context, body json.RawMessage) (*model.Publication, error) {
	var in publicationInput
	if err := json.Unmarshal(body, &in); err != nil {
		return nil, fmt.Errorf("parse publication input: %w", err)
	}

	titleJSON, err := marshalJSON(in.Title)
	if err != nil {
		return nil, err
	}
	journalJSON, err := marshalJSON(in.Journal)
	if err != nil {
		return nil, err
	}
	abstractJSON, err := marshalJSON(in.Abstract)
	if err != nil {
		return nil, err
	}

	yearStr := fmt.Sprintf("%d", in.Year)

	var p model.Publication
	var titleB, journalB, abstractB []byte
	var ys string
	var doi, url *string

	err = s.db.Pool.QueryRow(ctx,
		`INSERT INTO publications (title, journal, year, doi, url, abstract, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, title, journal, year, doi, url, abstract`,
		titleJSON, journalJSON, yearStr, nullableStr(in.DOI), nullableStr(in.URL), abstractJSON, in.SortOrder,
	).Scan(&p.ID, &titleB, &journalB, &ys, &doi, &url, &abstractB)
	if err != nil {
		return nil, fmt.Errorf("create publication: %w", err)
	}

	if err := parsePublicationFields(&p, titleB, journalB, abstractB, ys, doi, url); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *CareerService) updatePublication(ctx context.Context, id uuid.UUID, body json.RawMessage) (*model.Publication, error) {
	var in publicationInput
	if err := json.Unmarshal(body, &in); err != nil {
		return nil, fmt.Errorf("parse publication input: %w", err)
	}

	titleJSON, err := marshalJSON(in.Title)
	if err != nil {
		return nil, err
	}
	journalJSON, err := marshalJSON(in.Journal)
	if err != nil {
		return nil, err
	}
	abstractJSON, err := marshalJSON(in.Abstract)
	if err != nil {
		return nil, err
	}

	yearStr := fmt.Sprintf("%d", in.Year)

	var p model.Publication
	var titleB, journalB, abstractB []byte
	var ys string
	var doi, url *string

	err = s.db.Pool.QueryRow(ctx,
		`UPDATE publications SET title=$2, journal=$3, year=$4, doi=$5, url=$6, abstract=$7, sort_order=$8
WHERE id = $1 RETURNING id, title, journal, year, doi, url, abstract`,
		id, titleJSON, journalJSON, yearStr, nullableStr(in.DOI), nullableStr(in.URL), abstractJSON, in.SortOrder,
	).Scan(&p.ID, &titleB, &journalB, &ys, &doi, &url, &abstractB)
	if err != nil {
		return nil, fmt.Errorf("update publication: %w", err)
	}

	if err := parsePublicationFields(&p, titleB, journalB, abstractB, ys, doi, url); err != nil {
		return nil, err
	}
	return &p, nil
}

func scanPublication(rows pgx.Rows) (*model.Publication, error) {
	var p model.Publication
	var titleB, journalB, abstractB []byte
	var ys string
	var doi, url *string

	if err := rows.Scan(&p.ID, &titleB, &journalB, &ys, &doi, &url, &abstractB); err != nil {
		return nil, err
	}

	if err := parsePublicationFields(&p, titleB, journalB, abstractB, ys, doi, url); err != nil {
		return nil, err
	}
	return &p, nil
}

func parsePublicationFields(p *model.Publication, titleB, journalB, abstractB []byte, ys string, doi, url *string) error {
	if err := unmarshalJSON(titleB, &p.Title); err != nil {
		return err
	}
	if err := unmarshalJSON(journalB, &p.Journal); err != nil {
		return err
	}
	if err := unmarshalJSON(abstractB, &p.Abstract); err != nil {
		return err
	}

	fmt.Sscanf(ys, "%d", &p.Year)
	if doi != nil {
		p.DOI = *doi
	}
	if url != nil {
		p.URL = *url
	}
	return nil
}
