package model

import (
	"time"

	"github.com/google/uuid"
)

type Inquiry struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Email      string    `json:"email" db:"email"`
	Company    string    `json:"company,omitempty" db:"company"`
	Telegram   string    `json:"telegram,omitempty" db:"telegram"`
	Type       string    `json:"type" db:"inquiry_type"`
	Budget     string    `json:"budget,omitempty" db:"budget"`
	Message    string    `json:"message" db:"message"`
	Status     string    `json:"status" db:"status"`
	AdminNotes string    `json:"adminNotes,omitempty" db:"admin_notes"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateInquiryInput struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Company  string `json:"company" validate:"max=100"`
	Telegram string `json:"telegram" validate:"max=100"`
	Type     string `json:"type" validate:"required,oneof=freelance fulltime collaboration other"`
	Budget   string `json:"budget" validate:"max=50"`
	Message  string `json:"message" validate:"required,min=10,max=5000"`
}

type UpdateInquiryInput struct {
	Status     string `json:"status" validate:"required,oneof=new read replied archived"`
	AdminNotes string `json:"adminNotes" validate:"max=5000"`
}
