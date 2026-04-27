package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/portfolio/backend/internal/database"
	"github.com/portfolio/backend/internal/model"
)

const inquiryColumns = `id, name, email, company, telegram, inquiry_type, budget, message, status, admin_notes, created_at, updated_at`

type InquiryService struct {
	db *database.DB
}

func NewInquiryService(db *database.DB) *InquiryService {
	return &InquiryService{db: db}
}

func (s *InquiryService) List(ctx context.Context, status string) ([]model.Inquiry, error) {
	var query string
	var args []any

	if status != "" {
		query = `SELECT ` + inquiryColumns + ` FROM inquiries WHERE status = $1 ORDER BY created_at DESC`
		args = []any{status}
	} else {
		query = `SELECT ` + inquiryColumns + ` FROM inquiries ORDER BY created_at DESC`
	}

	rows, err := s.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list inquiries: %w", err)
	}
	defer rows.Close()

	var inquiries []model.Inquiry
	for rows.Next() {
		inq, err := scanInquiry(rows)
		if err != nil {
			return nil, fmt.Errorf("scan inquiry: %w", err)
		}
		inquiries = append(inquiries, *inq)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate inquiries: %w", err)
	}
	if inquiries == nil {
		inquiries = []model.Inquiry{}
	}
	return inquiries, nil
}

func (s *InquiryService) GetByID(ctx context.Context, id uuid.UUID) (*model.Inquiry, error) {
	row := s.db.Pool.QueryRow(ctx,
		`SELECT `+inquiryColumns+` FROM inquiries WHERE id = $1`, id)

	inq, err := scanInquiryRow(row)
	if err != nil {
		return nil, fmt.Errorf("get inquiry: %w", err)
	}
	return inq, nil
}

func (s *InquiryService) Create(ctx context.Context, input model.CreateInquiryInput) (*model.Inquiry, error) {
	row := s.db.Pool.QueryRow(ctx,
		`INSERT INTO inquiries (name, email, company, telegram, inquiry_type, budget, message)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING `+inquiryColumns,
		input.Name, input.Email, nullableStr(input.Company), nullableStr(input.Telegram),
		input.Type, nullableStr(input.Budget), input.Message,
	)

	inq, err := scanInquiryRow(row)
	if err != nil {
		return nil, fmt.Errorf("create inquiry: %w", err)
	}
	return inq, nil
}

func (s *InquiryService) UpdateStatus(ctx context.Context, id uuid.UUID, status string, adminNotes string) (*model.Inquiry, error) {
	switch status {
	case "new", "read", "replied", "archived":
	default:
		return nil, fmt.Errorf("invalid status %q", status)
	}

	row := s.db.Pool.QueryRow(ctx,
		`UPDATE inquiries SET status = $2, admin_notes = $3, updated_at = NOW() WHERE id = $1
RETURNING `+inquiryColumns,
		id, status, adminNotes,
	)

	inq, err := scanInquiryRow(row)
	if err != nil {
		return nil, fmt.Errorf("update inquiry status: %w", err)
	}
	return inq, nil
}

func scanInquiry(rows pgx.Rows) (*model.Inquiry, error) {
	var inq model.Inquiry
	var company, telegram, budget, adminNotes *string

	err := rows.Scan(
		&inq.ID, &inq.Name, &inq.Email, &company, &telegram,
		&inq.Type, &budget, &inq.Message,
		&inq.Status, &adminNotes, &inq.CreatedAt, &inq.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if company != nil {
		inq.Company = *company
	}
	if telegram != nil {
		inq.Telegram = *telegram
	}
	if budget != nil {
		inq.Budget = *budget
	}
	if adminNotes != nil {
		inq.AdminNotes = *adminNotes
	}
	return &inq, nil
}

func scanInquiryRow(row pgx.Row) (*model.Inquiry, error) {
	var inq model.Inquiry
	var company, telegram, budget, adminNotes *string

	err := row.Scan(
		&inq.ID, &inq.Name, &inq.Email, &company, &telegram,
		&inq.Type, &budget, &inq.Message,
		&inq.Status, &adminNotes, &inq.CreatedAt, &inq.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if company != nil {
		inq.Company = *company
	}
	if telegram != nil {
		inq.Telegram = *telegram
	}
	if budget != nil {
		inq.Budget = *budget
	}
	if adminNotes != nil {
		inq.AdminNotes = *adminNotes
	}
	return &inq, nil
}
