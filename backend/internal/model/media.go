package model

import (
	"time"

	"github.com/google/uuid"
)

type I18n map[string]string

type Media struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Filename     string    `json:"filename" db:"filename"`
	OriginalName string    `json:"originalName" db:"original_name"`
	MimeType     string    `json:"mimeType" db:"mime_type"`
	Size         int64     `json:"size" db:"size"`
	URL          string    `json:"url" db:"url"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}
