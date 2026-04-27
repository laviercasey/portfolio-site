package model

import (
	"time"

	"github.com/google/uuid"
)

type Education struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	Institution         I18n      `json:"institution" db:"institution"`
	Degree              I18n      `json:"degree" db:"degree"`
	Field               I18n      `json:"field" db:"field"`
	StartYear           int       `json:"startYear" db:"start_year"`
	EndYear             *int      `json:"endYear,omitempty" db:"end_year"`
	Description         I18n      `json:"description,omitempty" db:"description"`
	LogoURL             string    `json:"logoUrl,omitempty" db:"logo_url"`
	RelatedProjectSlugs []string  `json:"relatedProjectSlugs" db:"related_project_slugs"`
}

type WorkHistory struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Company         I18n       `json:"company" db:"company"`
	Position        I18n       `json:"position" db:"position"`
	StartDate       time.Time  `json:"startDate" db:"start_date"`
	EndDate         *time.Time `json:"endDate,omitempty" db:"end_date"`
	Current         bool       `json:"current" db:"current"`
	Description     I18n       `json:"description,omitempty" db:"description"`
	Technologies    []string   `json:"technologies" db:"technologies"`
	LogoURL         string     `json:"logoUrl,omitempty" db:"logo_url"`
	Achievements    []I18n     `json:"achievements" db:"achievements"`
	FullDescription I18n       `json:"fullDescription,omitempty" db:"full_description"`
}

type Certificate struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Title        I18n      `json:"title" db:"title"`
	Issuer       I18n      `json:"issuer" db:"issuer"`
	Date         time.Time `json:"date" db:"date"`
	CredentialID string    `json:"credentialId,omitempty" db:"credential_id"`
	URL          string    `json:"url,omitempty" db:"url"`
	ImageURL     string    `json:"imageUrl,omitempty" db:"image_url"`
}

type Publication struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Title    I18n      `json:"title" db:"title"`
	Journal  I18n      `json:"journal" db:"journal"`
	Year     int       `json:"year" db:"year"`
	DOI      string    `json:"doi,omitempty" db:"doi"`
	URL      string    `json:"url,omitempty" db:"url"`
	Abstract I18n      `json:"abstract,omitempty" db:"abstract"`
}

type CareerData struct {
	Education    []Education   `json:"education"`
	WorkHistory  []WorkHistory `json:"workHistory"`
	Certificates []Certificate `json:"certificates"`
	Publications []Publication `json:"publications"`
}
