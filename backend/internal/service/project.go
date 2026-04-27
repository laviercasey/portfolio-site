package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/portfolio/backend/internal/database"
	"github.com/portfolio/backend/internal/model"
)

const projectColumns = `id, slug, title, short_desc, description, category, status, tags, tech_stack, goal_desc, github_url, demo_url, site_url, video_url, thumbnail, images, stars, featured, sort_order, created_at, updated_at, problem, approach, outcome, tech_choices, highlights, timeline_started, timeline_shipped, demo_credentials`

type ProjectService struct {
	db *database.DB
}

func NewProjectService(db *database.DB) *ProjectService {
	return &ProjectService{db: db}
}

func (s *ProjectService) List(ctx context.Context, category string) ([]model.Project, error) {
	var query string
	var args []any

	if category != "" {
		query = `SELECT ` + projectColumns + ` FROM projects WHERE category = $1 ORDER BY sort_order ASC, created_at DESC`
		args = []any{category}
	} else {
		query = `SELECT ` + projectColumns + ` FROM projects ORDER BY sort_order ASC, created_at DESC`
	}

	rows, err := s.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, *p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate projects: %w", err)
	}

	if projects == nil {
		projects = []model.Project{}
	}
	return projects, nil
}

func (s *ProjectService) GetBySlug(ctx context.Context, slug string) (*model.Project, error) {
	query := `SELECT ` + projectColumns + ` FROM projects WHERE slug = $1`

	row := s.db.Pool.QueryRow(ctx, query, slug)
	p, err := scanProjectRow(row)
	if err != nil {
		return nil, fmt.Errorf("get project by slug: %w", err)
	}
	return p, nil
}

func (s *ProjectService) Create(ctx context.Context, input model.CreateProjectInput) (*model.Project, error) {
	titleJSON, err := marshalJSON(input.Title)
	if err != nil {
		return nil, err
	}
	shortDescJSON, err := marshalJSON(input.ShortDescription)
	if err != nil {
		return nil, err
	}
	descJSON, err := marshalJSON(input.Desc)
	if err != nil {
		return nil, err
	}
	goalDescJSON, err := marshalJSON(input.GoalDescription)
	if err != nil {
		return nil, err
	}

	if input.Tags == nil {
		input.Tags = []string{}
	}
	if input.TechStack == nil {
		input.TechStack = []string{}
	}
	if input.Images == nil {
		input.Images = []string{}
	}
	if input.TechChoices == nil {
		input.TechChoices = []model.TechChoice{}
	}
	if input.Highlights == nil {
		input.Highlights = []model.Highlight{}
	}
	if input.DemoCredentials == nil {
		input.DemoCredentials = []model.DemoCredential{}
	}

	problemJSON, err := marshalJSON(input.Problem)
	if err != nil {
		return nil, err
	}
	approachJSON, err := marshalJSON(input.Approach)
	if err != nil {
		return nil, err
	}
	outcomeJSON, err := marshalJSON(input.Outcome)
	if err != nil {
		return nil, err
	}
	techChoicesJSON, err := marshalJSON(input.TechChoices)
	if err != nil {
		return nil, err
	}
	highlightsJSON, err := marshalJSON(input.Highlights)
	if err != nil {
		return nil, err
	}
	credsJSON, err := marshalJSON(input.DemoCredentials)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO projects (slug, title, short_desc, description, category, status, tags, tech_stack, goal_desc, github_url, demo_url, site_url, video_url, thumbnail, images, featured, sort_order, problem, approach, outcome, tech_choices, highlights, timeline_started, timeline_shipped, demo_credentials)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
RETURNING ` + projectColumns

	row := s.db.Pool.QueryRow(ctx, query,
		input.Slug, titleJSON, shortDescJSON, descJSON,
		input.Category, input.Status, input.Tags, input.TechStack,
		goalDescJSON, input.GithubURL, input.DemoURL,
		input.SiteURL, input.VideoURL, input.ThumbnailURL,
		input.Images, input.Featured, input.Order,
		problemJSON, approachJSON, outcomeJSON,
		techChoicesJSON, highlightsJSON,
		input.TimelineStarted, input.TimelineShipped, credsJSON,
	)

	p, err := scanProjectRow(row)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return p, nil
}

func (s *ProjectService) Update(ctx context.Context, id uuid.UUID, input model.UpdateProjectInput) (*model.Project, error) {
	getQuery := `SELECT ` + projectColumns + ` FROM projects WHERE id = $1`
	row := s.db.Pool.QueryRow(ctx, getQuery, id)
	existing, err := scanProjectRow(row)
	if err != nil {
		return nil, fmt.Errorf("get project for update: %w", err)
	}

	if input.Slug != nil {
		existing.Slug = *input.Slug
	}
	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.ShortDescription != nil {
		existing.ShortDescription = *input.ShortDescription
	}
	if input.Desc != nil {
		existing.Description = *input.Desc
	}
	if input.Category != nil {
		existing.Category = *input.Category
	}
	if input.Status != nil {
		existing.Status = *input.Status
	}
	if input.Tags != nil {
		existing.Tags = *input.Tags
	}
	if input.TechStack != nil {
		existing.TechStack = *input.TechStack
	}
	if input.GoalDescription != nil {
		existing.GoalDescription = *input.GoalDescription
	}
	if input.GithubURL != nil {
		existing.GithubURL = *input.GithubURL
	}
	if input.DemoURL != nil {
		existing.DemoURL = *input.DemoURL
	}
	if input.SiteURL != nil {
		existing.SiteURL = *input.SiteURL
	}
	if input.VideoURL != nil {
		existing.VideoURL = *input.VideoURL
	}
	if input.ThumbnailURL != nil {
		existing.ThumbnailURL = *input.ThumbnailURL
	}
	if input.Images != nil {
		existing.Images = *input.Images
	}
	if input.Featured != nil {
		existing.Featured = *input.Featured
	}
	if input.Order != nil {
		existing.Order = *input.Order
	}
	if input.Problem != nil {
		existing.Problem = *input.Problem
	}
	if input.Approach != nil {
		existing.Approach = *input.Approach
	}
	if input.Outcome != nil {
		existing.Outcome = *input.Outcome
	}
	if input.TechChoices != nil {
		existing.TechChoices = *input.TechChoices
	}
	if input.Highlights != nil {
		existing.Highlights = *input.Highlights
	}
	if input.TimelineStarted != nil {
		existing.TimelineStarted = input.TimelineStarted
	}
	if input.TimelineShipped != nil {
		existing.TimelineShipped = input.TimelineShipped
	}
	if input.DemoCredentials != nil {
		existing.DemoCredentials = *input.DemoCredentials
	}

	titleJSON, err := marshalJSON(existing.Title)
	if err != nil {
		return nil, err
	}
	shortDescJSON, err := marshalJSON(existing.ShortDescription)
	if err != nil {
		return nil, err
	}
	descJSON, err := marshalJSON(existing.Description)
	if err != nil {
		return nil, err
	}
	goalDescJSON, err := marshalJSON(existing.GoalDescription)
	if err != nil {
		return nil, err
	}
	problemJSON, err := marshalJSON(existing.Problem)
	if err != nil {
		return nil, err
	}
	approachJSON, err := marshalJSON(existing.Approach)
	if err != nil {
		return nil, err
	}
	outcomeJSON, err := marshalJSON(existing.Outcome)
	if err != nil {
		return nil, err
	}
	techChoicesJSON, err := marshalJSON(existing.TechChoices)
	if err != nil {
		return nil, err
	}
	highlightsJSON, err := marshalJSON(existing.Highlights)
	if err != nil {
		return nil, err
	}
	credsJSON, err := marshalJSON(existing.DemoCredentials)
	if err != nil {
		return nil, err
	}

	updateQuery := `UPDATE projects SET
    slug = $2, title = $3, short_desc = $4, description = $5,
    category = $6, status = $7, tags = $8, tech_stack = $9,
    goal_desc = $10, github_url = $11, demo_url = $12,
    site_url = $13, video_url = $14, thumbnail = $15,
    images = $16, featured = $17, sort_order = $18,
    problem = $19, approach = $20, outcome = $21,
    tech_choices = $22, highlights = $23,
    timeline_started = $24, timeline_shipped = $25,
    demo_credentials = $26
WHERE id = $1
RETURNING ` + projectColumns

	updRow := s.db.Pool.QueryRow(ctx, updateQuery,
		id, existing.Slug, titleJSON, shortDescJSON, descJSON,
		existing.Category, existing.Status, existing.Tags, existing.TechStack,
		goalDescJSON, existing.GithubURL, existing.DemoURL,
		existing.SiteURL, existing.VideoURL, existing.ThumbnailURL,
		existing.Images, existing.Featured, existing.Order,
		problemJSON, approachJSON, outcomeJSON,
		techChoicesJSON, highlightsJSON,
		existing.TimelineStarted, existing.TimelineShipped, credsJSON,
	)

	p, err := scanProjectRow(updRow)
	if err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}
	return p, nil
}

func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := s.db.Pool.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
