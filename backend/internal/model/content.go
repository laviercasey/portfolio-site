package model

import (
	"encoding/json"
	"time"
)

type Content struct {
	Section   string          `json:"section" db:"section"`
	Data      json.RawMessage `json:"data" db:"data"`
	UpdatedAt time.Time       `json:"updatedAt" db:"updated_at"`
}

type UpdateContentInput struct {
	Data json.RawMessage `json:"data" validate:"required"`
}
