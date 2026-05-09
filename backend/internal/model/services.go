package model

import (
	"time"

	"github.com/google/uuid"
)

type ServiceBullets struct {
	Ru []string `json:"ru"`
	En []string `json:"en"`
}

type ServiceCaseProject struct {
	Slug string `json:"slug" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type Service struct {
	ID           uuid.UUID            `json:"id" db:"id"`
	Slug         string               `json:"slug" db:"slug"`
	Num          string               `json:"num" db:"num"`
	IconKey      string               `json:"iconKey" db:"icon_key"`
	VisualKey    string               `json:"visualKey" db:"visual_key"`
	Accent       string               `json:"accent" db:"accent"`
	Title        I18n                 `json:"title" db:"title"`
	Lead         I18n                 `json:"lead" db:"lead"`
	Bullets      ServiceBullets       `json:"bullets" db:"bullets"`
	Stack        string               `json:"stack" db:"stack"`
	Timeline     I18n                 `json:"timeline" db:"timeline"`
	PriceRu      string               `json:"priceRu" db:"price_ru"`
	PriceEn      string               `json:"priceEn" db:"price_en"`
	CaseProjects []ServiceCaseProject `json:"caseProjects" db:"case_projects"`
	Order        int                  `json:"order" db:"sort_order"`
	CreatedAt    time.Time            `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time            `json:"updatedAt" db:"updated_at"`
}

type CreateServiceInput struct {
	Slug         string               `json:"slug" validate:"required,min=2,max=100"`
	Num          string               `json:"num" validate:"required,max=10"`
	IconKey      string               `json:"iconKey" validate:"required,oneof=bot code layers database workflow"`
	VisualKey    string               `json:"visualKey" validate:"required,oneof=terminal browser editor dataframe pipeline"`
	Accent       string               `json:"accent" validate:"required"`
	Title        I18n                 `json:"title" validate:"required"`
	Lead         I18n                 `json:"lead" validate:"required"`
	Bullets      ServiceBullets       `json:"bullets"`
	Stack        string               `json:"stack"`
	Timeline     I18n                 `json:"timeline"`
	PriceRu      string               `json:"priceRu"`
	PriceEn      string               `json:"priceEn"`
	CaseProjects []ServiceCaseProject `json:"caseProjects" validate:"dive"`
	Order        int                  `json:"order"`
}

type UpdateServiceInput struct {
	Slug         *string               `json:"slug" validate:"omitempty,min=2,max=100"`
	Num          *string               `json:"num" validate:"omitempty,max=10"`
	IconKey      *string               `json:"iconKey" validate:"omitempty,oneof=bot code layers database workflow"`
	VisualKey    *string               `json:"visualKey" validate:"omitempty,oneof=terminal browser editor dataframe pipeline"`
	Accent       *string               `json:"accent"`
	Title        *I18n                 `json:"title"`
	Lead         *I18n                 `json:"lead"`
	Bullets      *ServiceBullets       `json:"bullets"`
	Stack        *string               `json:"stack"`
	Timeline     *I18n                 `json:"timeline"`
	PriceRu      *string               `json:"priceRu"`
	PriceEn      *string               `json:"priceEn"`
	CaseProjects *[]ServiceCaseProject `json:"caseProjects" validate:"omitempty,dive"`
	Order        *int                  `json:"order"`
}

type ServiceFaq struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Question  I18n      `json:"question" db:"question"`
	Answer    I18n      `json:"answer" db:"answer"`
	Order     int       `json:"order" db:"sort_order"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateServiceFaqInput struct {
	Question I18n `json:"question" validate:"required"`
	Answer   I18n `json:"answer" validate:"required"`
	Order    int  `json:"order"`
}

type UpdateServiceFaqInput struct {
	Question *I18n `json:"question"`
	Answer   *I18n `json:"answer"`
	Order    *int  `json:"order"`
}

type ServiceProcessStep struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Num         string    `json:"num" db:"num"`
	Title       I18n      `json:"title" db:"title"`
	Description I18n      `json:"description" db:"description"`
	Order       int       `json:"order" db:"sort_order"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateServiceProcessStepInput struct {
	Num         string `json:"num" validate:"required,max=10"`
	Title       I18n   `json:"title" validate:"required"`
	Description I18n   `json:"description" validate:"required"`
	Order       int    `json:"order"`
}

type UpdateServiceProcessStepInput struct {
	Num         *string `json:"num" validate:"omitempty,max=10"`
	Title       *I18n   `json:"title"`
	Description *I18n   `json:"description"`
	Order       *int    `json:"order"`
}

type ServicesPageData struct {
	Services     []Service            `json:"services"`
	Faqs         []ServiceFaq         `json:"faqs"`
	ProcessSteps []ServiceProcessStep `json:"processSteps"`
}
