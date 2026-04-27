package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/portfolio/backend/internal/model"
)

type certificateInput struct {
	Title        model.I18n `json:"title"`
	Issuer       model.I18n `json:"issuer"`
	Date         string     `json:"date"`
	CredentialID string     `json:"credentialId"`
	URL          string     `json:"url"`
	ImageURL     string     `json:"imageUrl"`
	SortOrder    int        `json:"sortOrder"`
}

func (s *CareerService) listCertificates(ctx context.Context) ([]model.Certificate, error) {
	rows, err := s.db.Pool.Query(ctx,
		`SELECT id, title, issuer, issue_date, credential_id, url, image_url FROM certificates ORDER BY sort_order ASC, issue_date DESC`)
	if err != nil {
		return nil, fmt.Errorf("list certificates: %w", err)
	}
	defer rows.Close()

	var items []model.Certificate
	for rows.Next() {
		c, err := scanCertificate(rows)
		if err != nil {
			return nil, fmt.Errorf("scan certificate: %w", err)
		}
		items = append(items, *c)
	}
	if items == nil {
		items = []model.Certificate{}
	}
	return items, nil
}

func (s *CareerService) createCertificate(ctx context.Context, body json.RawMessage) (*model.Certificate, error) {
	var in certificateInput
	if err := json.Unmarshal(body, &in); err != nil {
		return nil, fmt.Errorf("parse certificate input: %w", err)
	}

	titleJSON, err := marshalJSON(in.Title)
	if err != nil {
		return nil, err
	}
	issuerJSON, err := marshalJSON(in.Issuer)
	if err != nil {
		return nil, err
	}

	var c model.Certificate
	var titleB, issuerB []byte
	var dateStr string
	var credID, url, imageURL *string

	err = s.db.Pool.QueryRow(ctx,
		`INSERT INTO certificates (title, issuer, issue_date, credential_id, url, image_url, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, title, issuer, issue_date, credential_id, url, image_url`,
		titleJSON, issuerJSON, in.Date, nullableStr(in.CredentialID), nullableStr(in.URL), nullableStr(in.ImageURL), in.SortOrder,
	).Scan(&c.ID, &titleB, &issuerB, &dateStr, &credID, &url, &imageURL)
	if err != nil {
		return nil, fmt.Errorf("create certificate: %w", err)
	}

	if err := parseCertificateFields(&c, titleB, issuerB, dateStr, credID, url, imageURL); err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *CareerService) updateCertificate(ctx context.Context, id uuid.UUID, body json.RawMessage) (*model.Certificate, error) {
	var in certificateInput
	if err := json.Unmarshal(body, &in); err != nil {
		return nil, fmt.Errorf("parse certificate input: %w", err)
	}

	titleJSON, err := marshalJSON(in.Title)
	if err != nil {
		return nil, err
	}
	issuerJSON, err := marshalJSON(in.Issuer)
	if err != nil {
		return nil, err
	}

	var c model.Certificate
	var titleB, issuerB []byte
	var dateStr string
	var credID, urlP, imageURL *string

	err = s.db.Pool.QueryRow(ctx,
		`UPDATE certificates SET title=$2, issuer=$3, issue_date=$4, credential_id=$5, url=$6, image_url=$7, sort_order=$8
WHERE id = $1 RETURNING id, title, issuer, issue_date, credential_id, url, image_url`,
		id, titleJSON, issuerJSON, in.Date, nullableStr(in.CredentialID), nullableStr(in.URL), nullableStr(in.ImageURL), in.SortOrder,
	).Scan(&c.ID, &titleB, &issuerB, &dateStr, &credID, &urlP, &imageURL)
	if err != nil {
		return nil, fmt.Errorf("update certificate: %w", err)
	}

	if err := parseCertificateFields(&c, titleB, issuerB, dateStr, credID, urlP, imageURL); err != nil {
		return nil, err
	}
	return &c, nil
}

func scanCertificate(rows pgx.Rows) (*model.Certificate, error) {
	var c model.Certificate
	var titleB, issuerB []byte
	var dateStr string
	var credID, url, imageURL *string

	if err := rows.Scan(&c.ID, &titleB, &issuerB, &dateStr, &credID, &url, &imageURL); err != nil {
		return nil, err
	}

	if err := parseCertificateFields(&c, titleB, issuerB, dateStr, credID, url, imageURL); err != nil {
		return nil, err
	}
	return &c, nil
}

func parseCertificateFields(c *model.Certificate, titleB, issuerB []byte, dateStr string, credID, url, imageURL *string) error {
	if err := unmarshalJSON(titleB, &c.Title); err != nil {
		return err
	}
	if err := unmarshalJSON(issuerB, &c.Issuer); err != nil {
		return err
	}

	t, err := parseFlexibleDate(dateStr)
	if err != nil {
		return fmt.Errorf("certificate issue_date: %w", err)
	}
	c.Date = t

	if credID != nil {
		c.CredentialID = *credID
	}
	if url != nil {
		c.URL = *url
	}
	if imageURL != nil {
		c.ImageURL = *imageURL
	}
	return nil
}
