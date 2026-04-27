package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/portfolio/backend/internal/database"
	"github.com/portfolio/backend/internal/model"
)

var ErrInvalidCareerType = errors.New("invalid career type: must be education, work, certificate, or publication")

type CareerService struct {
	db *database.DB
}

func NewCareerService(db *database.DB) *CareerService {
	return &CareerService{db: db}
}

func (s *CareerService) GetAll(ctx context.Context) (*model.CareerData, error) {
	data := &model.CareerData{}

	edu, err := s.listEducation(ctx)
	if err != nil {
		return nil, err
	}
	data.Education = edu

	work, err := s.listWorkHistory(ctx)
	if err != nil {
		return nil, err
	}
	data.WorkHistory = work

	certs, err := s.listCertificates(ctx)
	if err != nil {
		return nil, err
	}
	data.Certificates = certs

	pubs, err := s.listPublications(ctx)
	if err != nil {
		return nil, err
	}
	data.Publications = pubs

	return data, nil
}

func (s *CareerService) Create(ctx context.Context, careerType string, body json.RawMessage) (any, error) {
	switch careerType {
	case "education":
		return s.createEducation(ctx, body)
	case "work":
		return s.createWorkHistory(ctx, body)
	case "certificate":
		return s.createCertificate(ctx, body)
	case "publication":
		return s.createPublication(ctx, body)
	default:
		return nil, ErrInvalidCareerType
	}
}

func (s *CareerService) Update(ctx context.Context, careerType string, id uuid.UUID, body json.RawMessage) (any, error) {
	switch careerType {
	case "education":
		return s.updateEducation(ctx, id, body)
	case "work":
		return s.updateWorkHistory(ctx, id, body)
	case "certificate":
		return s.updateCertificate(ctx, id, body)
	case "publication":
		return s.updatePublication(ctx, id, body)
	default:
		return nil, ErrInvalidCareerType
	}
}

func (s *CareerService) Delete(ctx context.Context, careerType string, id uuid.UUID) error {
	var query string
	switch careerType {
	case "education":
		query = `DELETE FROM education WHERE id = $1`
	case "work":
		query = `DELETE FROM work_history WHERE id = $1`
	case "certificate":
		query = `DELETE FROM certificates WHERE id = $1`
	case "publication":
		query = `DELETE FROM publications WHERE id = $1`
	default:
		return ErrInvalidCareerType
	}

	tag, err := s.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete %s: %w", careerType, err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (s *CareerService) listEducation(ctx context.Context) ([]model.Education, error) {
	rows, err := s.db.Pool.Query(ctx,
		`SELECT id, institution, degree, field, start_year, end_year, description, logo_url, related_project_slugs FROM education ORDER BY sort_order ASC, start_year DESC`)
	if err != nil {
		return nil, fmt.Errorf("list education: %w", err)
	}
	defer rows.Close()

	var items []model.Education
	for rows.Next() {
		e, err := scanEducation(rows)
		if err != nil {
			return nil, fmt.Errorf("scan education: %w", err)
		}
		items = append(items, *e)
	}
	if items == nil {
		items = []model.Education{}
	}
	return items, nil
}

type educationInput struct {
	Institution         model.I18n `json:"institution"`
	Degree              model.I18n `json:"degree"`
	Field               model.I18n `json:"field"`
	StartYear           int        `json:"startYear"`
	EndYear             *int       `json:"endYear"`
	Description         model.I18n `json:"description"`
	LogoURL             string     `json:"logoUrl"`
	SortOrder           int        `json:"sortOrder"`
	RelatedProjectSlugs []string   `json:"relatedProjectSlugs"`
}

func (s *CareerService) createEducation(ctx context.Context, body json.RawMessage) (*model.Education, error) {
	var in educationInput
	if err := json.Unmarshal(body, &in); err != nil {
		return nil, fmt.Errorf("parse education input: %w", err)
	}

	instJSON, err := marshalJSON(in.Institution)
	if err != nil {
		return nil, err
	}
	degJSON, err := marshalJSON(in.Degree)
	if err != nil {
		return nil, err
	}
	fieldJSON, err := marshalJSON(in.Field)
	if err != nil {
		return nil, err
	}
	descJSON, err := marshalJSON(in.Description)
	if err != nil {
		return nil, err
	}

	startYear := fmt.Sprintf("%d", in.StartYear)
	var endYear *string
	if in.EndYear != nil {
		s := fmt.Sprintf("%d", *in.EndYear)
		endYear = &s
	}

	if in.RelatedProjectSlugs == nil {
		in.RelatedProjectSlugs = []string{}
	}

	var e model.Education
	var instB, degB, fieldB, descB []byte
	var sy string
	var ey *string
	var logoURL *string
	var rps []string

	err = s.db.Pool.QueryRow(ctx,
		`INSERT INTO education (institution, degree, field, start_year, end_year, description, logo_url, sort_order, related_project_slugs)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, institution, degree, field, start_year, end_year, description, logo_url, related_project_slugs`,
		instJSON, degJSON, fieldJSON, startYear, endYear, descJSON, nullableStr(in.LogoURL), in.SortOrder, in.RelatedProjectSlugs,
	).Scan(&e.ID, &instB, &degB, &fieldB, &sy, &ey, &descB, &logoURL, &rps)
	if err != nil {
		return nil, fmt.Errorf("create education: %w", err)
	}

	if err := parseEducationFields(&e, instB, degB, fieldB, descB, sy, ey, logoURL, rps); err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *CareerService) updateEducation(ctx context.Context, id uuid.UUID, body json.RawMessage) (*model.Education, error) {
	var in educationInput
	if err := json.Unmarshal(body, &in); err != nil {
		return nil, fmt.Errorf("parse education input: %w", err)
	}

	instJSON, err := marshalJSON(in.Institution)
	if err != nil {
		return nil, err
	}
	degJSON, err := marshalJSON(in.Degree)
	if err != nil {
		return nil, err
	}
	fieldJSON, err := marshalJSON(in.Field)
	if err != nil {
		return nil, err
	}
	descJSON, err := marshalJSON(in.Description)
	if err != nil {
		return nil, err
	}

	startYear := fmt.Sprintf("%d", in.StartYear)
	var endYear *string
	if in.EndYear != nil {
		s := fmt.Sprintf("%d", *in.EndYear)
		endYear = &s
	}

	if in.RelatedProjectSlugs == nil {
		in.RelatedProjectSlugs = []string{}
	}

	var e model.Education
	var instB, degB, fieldB, descB []byte
	var sy string
	var ey *string
	var logoURL *string
	var rps []string

	err = s.db.Pool.QueryRow(ctx,
		`UPDATE education SET institution=$2, degree=$3, field=$4, start_year=$5, end_year=$6, description=$7, logo_url=$8, sort_order=$9, related_project_slugs=$10
WHERE id = $1 RETURNING id, institution, degree, field, start_year, end_year, description, logo_url, related_project_slugs`,
		id, instJSON, degJSON, fieldJSON, startYear, endYear, descJSON, nullableStr(in.LogoURL), in.SortOrder, in.RelatedProjectSlugs,
	).Scan(&e.ID, &instB, &degB, &fieldB, &sy, &ey, &descB, &logoURL, &rps)
	if err != nil {
		return nil, fmt.Errorf("update education: %w", err)
	}

	if err := parseEducationFields(&e, instB, degB, fieldB, descB, sy, ey, logoURL, rps); err != nil {
		return nil, err
	}
	return &e, nil
}

func scanEducation(rows pgx.Rows) (*model.Education, error) {
	var e model.Education
	var instB, degB, fieldB, descB []byte
	var sy string
	var ey *string
	var logoURL *string
	var rps []string

	if err := rows.Scan(&e.ID, &instB, &degB, &fieldB, &sy, &ey, &descB, &logoURL, &rps); err != nil {
		return nil, err
	}

	if err := parseEducationFields(&e, instB, degB, fieldB, descB, sy, ey, logoURL, rps); err != nil {
		return nil, err
	}
	return &e, nil
}

func parseEducationFields(e *model.Education, instB, degB, fieldB, descB []byte, sy string, ey *string, logoURL *string, rps []string) error {
	if err := unmarshalJSON(instB, &e.Institution); err != nil {
		return err
	}
	if err := unmarshalJSON(degB, &e.Degree); err != nil {
		return err
	}
	if err := unmarshalJSON(fieldB, &e.Field); err != nil {
		return err
	}
	if err := unmarshalJSON(descB, &e.Description); err != nil {
		return err
	}

	fmt.Sscanf(sy, "%d", &e.StartYear)
	if ey != nil {
		var v int
		fmt.Sscanf(*ey, "%d", &v)
		e.EndYear = &v
	}
	if logoURL != nil {
		e.LogoURL = *logoURL
	}
	if rps == nil {
		rps = []string{}
	}
	e.RelatedProjectSlugs = rps
	return nil
}
