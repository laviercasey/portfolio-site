package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/portfolio/backend/internal/model"
)

type workHistoryInput struct {
	Company         model.I18n   `json:"company"`
	Position        model.I18n   `json:"position"`
	StartDate       string       `json:"startDate"`
	EndDate         *string      `json:"endDate"`
	Current         bool         `json:"current"`
	Description     model.I18n   `json:"description"`
	Technologies    []string     `json:"technologies"`
	LogoURL         string       `json:"logoUrl"`
	SortOrder       int          `json:"sortOrder"`
	Achievements    []model.I18n `json:"achievements"`
	FullDescription model.I18n   `json:"fullDescription"`
}

func (s *CareerService) listWorkHistory(ctx context.Context) ([]model.WorkHistory, error) {
	rows, err := s.db.Pool.Query(ctx,
		`SELECT id, company, position, start_date, end_date, is_current, description, technologies, logo_url, achievements, full_description FROM work_history ORDER BY sort_order ASC, start_date DESC`)
	if err != nil {
		return nil, fmt.Errorf("list work history: %w", err)
	}
	defer rows.Close()

	var items []model.WorkHistory
	for rows.Next() {
		w, err := scanWorkHistory(rows)
		if err != nil {
			return nil, fmt.Errorf("scan work history: %w", err)
		}
		items = append(items, *w)
	}
	if items == nil {
		items = []model.WorkHistory{}
	}
	return items, nil
}

func (s *CareerService) createWorkHistory(ctx context.Context, body json.RawMessage) (*model.WorkHistory, error) {
	var in workHistoryInput
	if err := json.Unmarshal(body, &in); err != nil {
		return nil, fmt.Errorf("parse work history input: %w", err)
	}

	compJSON, err := marshalJSON(in.Company)
	if err != nil {
		return nil, err
	}
	posJSON, err := marshalJSON(in.Position)
	if err != nil {
		return nil, err
	}
	descJSON, err := marshalJSON(in.Description)
	if err != nil {
		return nil, err
	}

	if in.Technologies == nil {
		in.Technologies = []string{}
	}
	if in.Achievements == nil {
		in.Achievements = []model.I18n{}
	}
	achievementsJSON, err := marshalJSON(in.Achievements)
	if err != nil {
		return nil, err
	}
	fullDescJSON, err := marshalJSON(in.FullDescription)
	if err != nil {
		return nil, err
	}

	var w model.WorkHistory
	var compB, posB, descB, achB, fullDescB []byte
	var sd string
	var ed *string
	var logoURL *string

	err = s.db.Pool.QueryRow(ctx,
		`INSERT INTO work_history (company, position, start_date, end_date, is_current, description, technologies, logo_url, sort_order, achievements, full_description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, company, position, start_date, end_date, is_current, description, technologies, logo_url, achievements, full_description`,
		compJSON, posJSON, in.StartDate, in.EndDate, in.Current, descJSON, in.Technologies, nullableStr(in.LogoURL), in.SortOrder, achievementsJSON, fullDescJSON,
	).Scan(&w.ID, &compB, &posB, &sd, &ed, &w.Current, &descB, &w.Technologies, &logoURL, &achB, &fullDescB)
	if err != nil {
		return nil, fmt.Errorf("create work history: %w", err)
	}

	if err := parseWorkHistoryFields(&w, compB, posB, descB, sd, ed, logoURL, achB, fullDescB); err != nil {
		return nil, err
	}
	return &w, nil
}

func (s *CareerService) updateWorkHistory(ctx context.Context, id uuid.UUID, body json.RawMessage) (*model.WorkHistory, error) {
	var in workHistoryInput
	if err := json.Unmarshal(body, &in); err != nil {
		return nil, fmt.Errorf("parse work history input: %w", err)
	}

	compJSON, err := marshalJSON(in.Company)
	if err != nil {
		return nil, err
	}
	posJSON, err := marshalJSON(in.Position)
	if err != nil {
		return nil, err
	}
	descJSON, err := marshalJSON(in.Description)
	if err != nil {
		return nil, err
	}

	if in.Technologies == nil {
		in.Technologies = []string{}
	}
	if in.Achievements == nil {
		in.Achievements = []model.I18n{}
	}
	achievementsJSON, err := marshalJSON(in.Achievements)
	if err != nil {
		return nil, err
	}
	fullDescJSON, err := marshalJSON(in.FullDescription)
	if err != nil {
		return nil, err
	}

	var w model.WorkHistory
	var compB, posB, descB, achB, fullDescB []byte
	var sd string
	var ed *string
	var logoURL *string

	err = s.db.Pool.QueryRow(ctx,
		`UPDATE work_history SET company=$2, position=$3, start_date=$4, end_date=$5, is_current=$6, description=$7, technologies=$8, logo_url=$9, sort_order=$10, achievements=$11, full_description=$12
WHERE id = $1 RETURNING id, company, position, start_date, end_date, is_current, description, technologies, logo_url, achievements, full_description`,
		id, compJSON, posJSON, in.StartDate, in.EndDate, in.Current, descJSON, in.Technologies, nullableStr(in.LogoURL), in.SortOrder, achievementsJSON, fullDescJSON,
	).Scan(&w.ID, &compB, &posB, &sd, &ed, &w.Current, &descB, &w.Technologies, &logoURL, &achB, &fullDescB)
	if err != nil {
		return nil, fmt.Errorf("update work history: %w", err)
	}

	if err := parseWorkHistoryFields(&w, compB, posB, descB, sd, ed, logoURL, achB, fullDescB); err != nil {
		return nil, err
	}
	return &w, nil
}

func scanWorkHistory(rows pgx.Rows) (*model.WorkHistory, error) {
	var w model.WorkHistory
	var compB, posB, descB, achB, fullDescB []byte
	var sd string
	var ed *string
	var logoURL *string

	if err := rows.Scan(&w.ID, &compB, &posB, &sd, &ed, &w.Current, &descB, &w.Technologies, &logoURL, &achB, &fullDescB); err != nil {
		return nil, err
	}

	if err := parseWorkHistoryFields(&w, compB, posB, descB, sd, ed, logoURL, achB, fullDescB); err != nil {
		return nil, err
	}
	return &w, nil
}

func parseWorkHistoryFields(w *model.WorkHistory, compB, posB, descB []byte, sd string, ed *string, logoURL *string, achB, fullDescB []byte) error {
	if err := unmarshalJSON(compB, &w.Company); err != nil {
		return err
	}
	if err := unmarshalJSON(posB, &w.Position); err != nil {
		return err
	}
	if err := unmarshalJSON(descB, &w.Description); err != nil {
		return err
	}
	if err := unmarshalJSON(achB, &w.Achievements); err != nil {
		return err
	}
	if err := unmarshalJSON(fullDescB, &w.FullDescription); err != nil {
		return err
	}

	t, err := parseFlexibleDate(sd)
	if err != nil {
		return fmt.Errorf("work history start_date: %w", err)
	}
	w.StartDate = t

	if ed != nil {
		t, err := parseFlexibleDate(*ed)
		if err != nil {
			slog.Warn("work history end_date unparseable, leaving zero", "value", *ed, "err", err)
		} else {
			w.EndDate = &t
		}
	}

	if w.Technologies == nil {
		w.Technologies = []string{}
	}
	if w.Achievements == nil {
		w.Achievements = []model.I18n{}
	}
	if logoURL != nil {
		w.LogoURL = *logoURL
	}
	return nil
}
