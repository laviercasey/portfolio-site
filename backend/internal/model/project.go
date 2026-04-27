package model

import (
	"time"

	"github.com/google/uuid"
)

type TechChoice struct {
	Tech   string `json:"tech"`
	Reason I18n   `json:"reason"`
}

type Highlight struct {
	Label I18n   `json:"label"`
	Value string `json:"value"`
}

type DemoCredential struct {
	Role     string `json:"role"`
	Label    I18n   `json:"label"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Note     I18n   `json:"note"`
}

type Project struct {
	ID               uuid.UUID        `json:"id" db:"id"`
	Slug             string           `json:"slug" db:"slug"`
	Title            I18n             `json:"title" db:"title"`
	ShortDescription I18n             `json:"shortDescription" db:"short_desc"`
	Description      I18n             `json:"description" db:"description"`
	Category         string           `json:"category" db:"category"`
	Status           string           `json:"status" db:"status"`
	Tags             []string         `json:"tags" db:"tags"`
	TechStack        []string         `json:"techStack" db:"tech_stack"`
	GoalDescription  I18n             `json:"goalDescription" db:"goal_desc"`
	GithubURL        string           `json:"githubUrl,omitempty" db:"github_url"`
	DemoURL          string           `json:"demoUrl,omitempty" db:"demo_url"`
	SiteURL          string           `json:"siteUrl,omitempty" db:"site_url"`
	VideoURL         string           `json:"videoUrl,omitempty" db:"video_url"`
	ThumbnailURL     string           `json:"thumbnailUrl,omitempty" db:"thumbnail"`
	Images           []string         `json:"images,omitempty" db:"images"`
	Stars            int              `json:"stars,omitempty" db:"stars"`
	Featured         bool             `json:"featured" db:"featured"`
	Order            int              `json:"order" db:"sort_order"`
	Problem          I18n             `json:"problem,omitempty" db:"problem"`
	Approach         I18n             `json:"approach,omitempty" db:"approach"`
	Outcome          I18n             `json:"outcome,omitempty" db:"outcome"`
	TechChoices      []TechChoice     `json:"techChoices,omitempty" db:"tech_choices"`
	Highlights       []Highlight      `json:"highlights,omitempty" db:"highlights"`
	TimelineStarted  *time.Time       `json:"timelineStarted,omitempty" db:"timeline_started"`
	TimelineShipped  *time.Time       `json:"timelineShipped,omitempty" db:"timeline_shipped"`
	DemoCredentials  []DemoCredential `json:"demoCredentials,omitempty" db:"demo_credentials"`
	CreatedAt        time.Time        `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time        `json:"updatedAt" db:"updated_at"`
}

type CreateProjectInput struct {
	Slug             string           `json:"slug" validate:"required,min=2,max=100"`
	Title            I18n             `json:"title" validate:"required"`
	ShortDescription I18n             `json:"shortDescription"`
	Desc             I18n             `json:"description"`
	Category         string           `json:"category" validate:"required"`
	Status           string           `json:"status" validate:"required,oneof=completed in_development open_source"`
	Tags             []string         `json:"tags"`
	TechStack        []string         `json:"techStack"`
	GoalDescription  I18n             `json:"goalDescription"`
	GithubURL        string           `json:"githubUrl" validate:"omitempty,url"`
	DemoURL          string           `json:"demoUrl" validate:"omitempty,url"`
	SiteURL          string           `json:"siteUrl" validate:"omitempty,url"`
	VideoURL         string           `json:"videoUrl"`
	ThumbnailURL     string           `json:"thumbnailUrl"`
	Images           []string         `json:"images"`
	Featured         bool             `json:"featured"`
	Order            int              `json:"order"`
	Problem          I18n             `json:"problem"`
	Approach         I18n             `json:"approach"`
	Outcome          I18n             `json:"outcome"`
	TechChoices      []TechChoice     `json:"techChoices"`
	Highlights       []Highlight      `json:"highlights"`
	TimelineStarted  *time.Time       `json:"timelineStarted"`
	TimelineShipped  *time.Time       `json:"timelineShipped"`
	DemoCredentials  []DemoCredential `json:"demoCredentials" validate:"dive"`
}

type UpdateProjectInput struct {
	Slug             *string           `json:"slug" validate:"omitempty,min=2,max=100"`
	Title            *I18n             `json:"title"`
	ShortDescription *I18n             `json:"shortDescription"`
	Desc             *I18n             `json:"description"`
	Category         *string           `json:"category"`
	Status           *string           `json:"status" validate:"omitempty,oneof=completed in_development open_source"`
	Tags             *[]string         `json:"tags"`
	TechStack        *[]string         `json:"techStack"`
	GoalDescription  *I18n             `json:"goalDescription"`
	GithubURL        *string           `json:"githubUrl" validate:"omitempty,url"`
	DemoURL          *string           `json:"demoUrl" validate:"omitempty,url"`
	SiteURL          *string           `json:"siteUrl" validate:"omitempty,url"`
	VideoURL         *string           `json:"videoUrl"`
	ThumbnailURL     *string           `json:"thumbnailUrl"`
	Images           *[]string         `json:"images"`
	Featured         *bool             `json:"featured"`
	Order            *int              `json:"order"`
	Problem          *I18n             `json:"problem"`
	Approach         *I18n             `json:"approach"`
	Outcome          *I18n             `json:"outcome"`
	TechChoices      *[]TechChoice     `json:"techChoices"`
	Highlights       *[]Highlight      `json:"highlights"`
	TimelineStarted  *time.Time        `json:"timelineStarted"`
	TimelineShipped  *time.Time        `json:"timelineShipped"`
	DemoCredentials  *[]DemoCredential `json:"demoCredentials"`
}
