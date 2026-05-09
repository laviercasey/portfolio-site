package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/portfolio/backend/internal/database"
	"github.com/portfolio/backend/internal/model"
)

const serviceColumns = `id, slug, num, icon_key, visual_key, accent, title, lead, bullets, stack, timeline, price_ru, price_en, case_projects, sort_order, created_at, updated_at`

const serviceFaqColumns = `id, question, answer, sort_order, created_at, updated_at`

const serviceProcessStepColumns = `id, num, title, description, sort_order, created_at, updated_at`

type ServicesService struct {
	db *database.DB
}

func NewServicesService(db *database.DB) *ServicesService {
	return &ServicesService{db: db}
}

func scanService(rs interface {
	Scan(...any) error
}) (*model.Service, error) {
	var s model.Service
	var titleRaw, leadRaw, bulletsRaw, timelineRaw, casesRaw []byte
	if err := rs.Scan(
		&s.ID, &s.Slug, &s.Num, &s.IconKey, &s.VisualKey, &s.Accent,
		&titleRaw, &leadRaw, &bulletsRaw, &s.Stack, &timelineRaw,
		&s.PriceRu, &s.PriceEn, &casesRaw, &s.Order,
		&s.CreatedAt, &s.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if err := unmarshalJSON(titleRaw, &s.Title); err != nil {
		return nil, fmt.Errorf("unmarshal title: %w", err)
	}
	if err := unmarshalJSON(leadRaw, &s.Lead); err != nil {
		return nil, fmt.Errorf("unmarshal lead: %w", err)
	}
	if err := unmarshalJSON(bulletsRaw, &s.Bullets); err != nil {
		return nil, fmt.Errorf("unmarshal bullets: %w", err)
	}
	if err := unmarshalJSON(timelineRaw, &s.Timeline); err != nil {
		return nil, fmt.Errorf("unmarshal timeline: %w", err)
	}
	if err := unmarshalJSON(casesRaw, &s.CaseProjects); err != nil {
		return nil, fmt.Errorf("unmarshal case_projects: %w", err)
	}
	if s.CaseProjects == nil {
		s.CaseProjects = []model.ServiceCaseProject{}
	}
	return &s, nil
}

func (s *ServicesService) ListServices(ctx context.Context) ([]model.Service, error) {
	query := `SELECT ` + serviceColumns + ` FROM services ORDER BY sort_order ASC, created_at ASC`
	rows, err := s.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}
	defer rows.Close()

	out := []model.Service{}
	for rows.Next() {
		svc, err := scanService(rows)
		if err != nil {
			return nil, fmt.Errorf("scan service: %w", err)
		}
		out = append(out, *svc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate services: %w", err)
	}
	return out, nil
}

func (s *ServicesService) CreateService(ctx context.Context, input model.CreateServiceInput) (*model.Service, error) {
	titleJSON, err := marshalJSON(input.Title)
	if err != nil {
		return nil, err
	}
	leadJSON, err := marshalJSON(input.Lead)
	if err != nil {
		return nil, err
	}
	bulletsJSON, err := marshalJSON(input.Bullets)
	if err != nil {
		return nil, err
	}
	timelineJSON, err := marshalJSON(input.Timeline)
	if err != nil {
		return nil, err
	}
	if input.CaseProjects == nil {
		input.CaseProjects = []model.ServiceCaseProject{}
	}
	casesJSON, err := marshalJSON(input.CaseProjects)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO services (slug, num, icon_key, visual_key, accent, title, lead, bullets, stack, timeline, price_ru, price_en, case_projects, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING ` + serviceColumns

	row := s.db.Pool.QueryRow(ctx, query,
		input.Slug, input.Num, input.IconKey, input.VisualKey, input.Accent,
		titleJSON, leadJSON, bulletsJSON, input.Stack, timelineJSON,
		input.PriceRu, input.PriceEn, casesJSON, input.Order,
	)
	svc, err := scanService(row)
	if err != nil {
		return nil, fmt.Errorf("create service: %w", err)
	}
	return svc, nil
}

func (s *ServicesService) UpdateService(ctx context.Context, id uuid.UUID, input model.UpdateServiceInput) (*model.Service, error) {
	getRow := s.db.Pool.QueryRow(ctx, `SELECT `+serviceColumns+` FROM services WHERE id = $1`, id)
	existing, err := scanService(getRow)
	if err != nil {
		return nil, fmt.Errorf("get service for update: %w", err)
	}

	if input.Slug != nil {
		existing.Slug = *input.Slug
	}
	if input.Num != nil {
		existing.Num = *input.Num
	}
	if input.IconKey != nil {
		existing.IconKey = *input.IconKey
	}
	if input.VisualKey != nil {
		existing.VisualKey = *input.VisualKey
	}
	if input.Accent != nil {
		existing.Accent = *input.Accent
	}
	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Lead != nil {
		existing.Lead = *input.Lead
	}
	if input.Bullets != nil {
		existing.Bullets = *input.Bullets
	}
	if input.Stack != nil {
		existing.Stack = *input.Stack
	}
	if input.Timeline != nil {
		existing.Timeline = *input.Timeline
	}
	if input.PriceRu != nil {
		existing.PriceRu = *input.PriceRu
	}
	if input.PriceEn != nil {
		existing.PriceEn = *input.PriceEn
	}
	if input.CaseProjects != nil {
		existing.CaseProjects = *input.CaseProjects
	}
	if input.Order != nil {
		existing.Order = *input.Order
	}

	titleJSON, err := marshalJSON(existing.Title)
	if err != nil {
		return nil, err
	}
	leadJSON, err := marshalJSON(existing.Lead)
	if err != nil {
		return nil, err
	}
	bulletsJSON, err := marshalJSON(existing.Bullets)
	if err != nil {
		return nil, err
	}
	timelineJSON, err := marshalJSON(existing.Timeline)
	if err != nil {
		return nil, err
	}
	casesJSON, err := marshalJSON(existing.CaseProjects)
	if err != nil {
		return nil, err
	}

	query := `UPDATE services SET
    slug = $2, num = $3, icon_key = $4, visual_key = $5, accent = $6,
    title = $7, lead = $8, bullets = $9, stack = $10, timeline = $11,
    price_ru = $12, price_en = $13, case_projects = $14, sort_order = $15
WHERE id = $1
RETURNING ` + serviceColumns

	row := s.db.Pool.QueryRow(ctx, query,
		id, existing.Slug, existing.Num, existing.IconKey, existing.VisualKey, existing.Accent,
		titleJSON, leadJSON, bulletsJSON, existing.Stack, timelineJSON,
		existing.PriceRu, existing.PriceEn, casesJSON, existing.Order,
	)
	svc, err := scanService(row)
	if err != nil {
		return nil, fmt.Errorf("update service: %w", err)
	}
	return svc, nil
}

func (s *ServicesService) DeleteService(ctx context.Context, id uuid.UUID) error {
	tag, err := s.db.Pool.Exec(ctx, `DELETE FROM services WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete service: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func scanFaq(rs interface {
	Scan(...any) error
}) (*model.ServiceFaq, error) {
	var f model.ServiceFaq
	var qRaw, aRaw []byte
	if err := rs.Scan(&f.ID, &qRaw, &aRaw, &f.Order, &f.CreatedAt, &f.UpdatedAt); err != nil {
		return nil, err
	}
	if err := unmarshalJSON(qRaw, &f.Question); err != nil {
		return nil, fmt.Errorf("unmarshal question: %w", err)
	}
	if err := unmarshalJSON(aRaw, &f.Answer); err != nil {
		return nil, fmt.Errorf("unmarshal answer: %w", err)
	}
	return &f, nil
}

func (s *ServicesService) ListFaqs(ctx context.Context) ([]model.ServiceFaq, error) {
	rows, err := s.db.Pool.Query(ctx, `SELECT `+serviceFaqColumns+` FROM service_faqs ORDER BY sort_order ASC, created_at ASC`)
	if err != nil {
		return nil, fmt.Errorf("list service faqs: %w", err)
	}
	defer rows.Close()

	out := []model.ServiceFaq{}
	for rows.Next() {
		f, err := scanFaq(rows)
		if err != nil {
			return nil, fmt.Errorf("scan service faq: %w", err)
		}
		out = append(out, *f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate service faqs: %w", err)
	}
	return out, nil
}

func (s *ServicesService) CreateFaq(ctx context.Context, input model.CreateServiceFaqInput) (*model.ServiceFaq, error) {
	qJSON, err := marshalJSON(input.Question)
	if err != nil {
		return nil, err
	}
	aJSON, err := marshalJSON(input.Answer)
	if err != nil {
		return nil, err
	}
	row := s.db.Pool.QueryRow(ctx,
		`INSERT INTO service_faqs (question, answer, sort_order) VALUES ($1, $2, $3) RETURNING `+serviceFaqColumns,
		qJSON, aJSON, input.Order,
	)
	f, err := scanFaq(row)
	if err != nil {
		return nil, fmt.Errorf("create service faq: %w", err)
	}
	return f, nil
}

func (s *ServicesService) UpdateFaq(ctx context.Context, id uuid.UUID, input model.UpdateServiceFaqInput) (*model.ServiceFaq, error) {
	getRow := s.db.Pool.QueryRow(ctx, `SELECT `+serviceFaqColumns+` FROM service_faqs WHERE id = $1`, id)
	existing, err := scanFaq(getRow)
	if err != nil {
		return nil, fmt.Errorf("get service faq for update: %w", err)
	}
	if input.Question != nil {
		existing.Question = *input.Question
	}
	if input.Answer != nil {
		existing.Answer = *input.Answer
	}
	if input.Order != nil {
		existing.Order = *input.Order
	}
	qJSON, err := marshalJSON(existing.Question)
	if err != nil {
		return nil, err
	}
	aJSON, err := marshalJSON(existing.Answer)
	if err != nil {
		return nil, err
	}
	row := s.db.Pool.QueryRow(ctx,
		`UPDATE service_faqs SET question = $2, answer = $3, sort_order = $4 WHERE id = $1 RETURNING `+serviceFaqColumns,
		id, qJSON, aJSON, existing.Order,
	)
	f, err := scanFaq(row)
	if err != nil {
		return nil, fmt.Errorf("update service faq: %w", err)
	}
	return f, nil
}

func (s *ServicesService) DeleteFaq(ctx context.Context, id uuid.UUID) error {
	tag, err := s.db.Pool.Exec(ctx, `DELETE FROM service_faqs WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete service faq: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func scanProcessStep(rs interface {
	Scan(...any) error
}) (*model.ServiceProcessStep, error) {
	var p model.ServiceProcessStep
	var titleRaw, descRaw []byte
	if err := rs.Scan(&p.ID, &p.Num, &titleRaw, &descRaw, &p.Order, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	if err := unmarshalJSON(titleRaw, &p.Title); err != nil {
		return nil, fmt.Errorf("unmarshal step title: %w", err)
	}
	if err := unmarshalJSON(descRaw, &p.Description); err != nil {
		return nil, fmt.Errorf("unmarshal step description: %w", err)
	}
	return &p, nil
}

func (s *ServicesService) ListProcessSteps(ctx context.Context) ([]model.ServiceProcessStep, error) {
	rows, err := s.db.Pool.Query(ctx, `SELECT `+serviceProcessStepColumns+` FROM service_process_steps ORDER BY sort_order ASC, created_at ASC`)
	if err != nil {
		return nil, fmt.Errorf("list process steps: %w", err)
	}
	defer rows.Close()

	out := []model.ServiceProcessStep{}
	for rows.Next() {
		p, err := scanProcessStep(rows)
		if err != nil {
			return nil, fmt.Errorf("scan process step: %w", err)
		}
		out = append(out, *p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate process steps: %w", err)
	}
	return out, nil
}

func (s *ServicesService) CreateProcessStep(ctx context.Context, input model.CreateServiceProcessStepInput) (*model.ServiceProcessStep, error) {
	titleJSON, err := marshalJSON(input.Title)
	if err != nil {
		return nil, err
	}
	descJSON, err := marshalJSON(input.Description)
	if err != nil {
		return nil, err
	}
	row := s.db.Pool.QueryRow(ctx,
		`INSERT INTO service_process_steps (num, title, description, sort_order) VALUES ($1, $2, $3, $4) RETURNING `+serviceProcessStepColumns,
		input.Num, titleJSON, descJSON, input.Order,
	)
	p, err := scanProcessStep(row)
	if err != nil {
		return nil, fmt.Errorf("create process step: %w", err)
	}
	return p, nil
}

func (s *ServicesService) UpdateProcessStep(ctx context.Context, id uuid.UUID, input model.UpdateServiceProcessStepInput) (*model.ServiceProcessStep, error) {
	getRow := s.db.Pool.QueryRow(ctx, `SELECT `+serviceProcessStepColumns+` FROM service_process_steps WHERE id = $1`, id)
	existing, err := scanProcessStep(getRow)
	if err != nil {
		return nil, fmt.Errorf("get process step for update: %w", err)
	}
	if input.Num != nil {
		existing.Num = *input.Num
	}
	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.Order != nil {
		existing.Order = *input.Order
	}
	titleJSON, err := marshalJSON(existing.Title)
	if err != nil {
		return nil, err
	}
	descJSON, err := marshalJSON(existing.Description)
	if err != nil {
		return nil, err
	}
	row := s.db.Pool.QueryRow(ctx,
		`UPDATE service_process_steps SET num = $2, title = $3, description = $4, sort_order = $5 WHERE id = $1 RETURNING `+serviceProcessStepColumns,
		id, existing.Num, titleJSON, descJSON, existing.Order,
	)
	p, err := scanProcessStep(row)
	if err != nil {
		return nil, fmt.Errorf("update process step: %w", err)
	}
	return p, nil
}

func (s *ServicesService) DeleteProcessStep(ctx context.Context, id uuid.UUID) error {
	tag, err := s.db.Pool.Exec(ctx, `DELETE FROM service_process_steps WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete process step: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (s *ServicesService) GetPageData(ctx context.Context) (*model.ServicesPageData, error) {
	services, err := s.ListServices(ctx)
	if err != nil {
		return nil, err
	}
	faqs, err := s.ListFaqs(ctx)
	if err != nil {
		return nil, err
	}
	steps, err := s.ListProcessSteps(ctx)
	if err != nil {
		return nil, err
	}
	return &model.ServicesPageData{
		Services:     services,
		Faqs:         faqs,
		ProcessSteps: steps,
	}, nil
}
