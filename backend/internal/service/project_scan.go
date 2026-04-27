package service

import (
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/portfolio/backend/internal/model"
)

type projectJSONFields struct {
	titleJSON, shortDescJSON, descJSON, goalDescJSON []byte
	problemJSON, approachJSON, outcomeJSON           []byte
	techChoicesJSON, highlightsJSON, credsJSON       []byte
}

type projectNullableFields struct {
	githubURL, demoURL, siteURL, videoURL, thumbnail *string
	timelineStarted, timelineShipped                 *time.Time
}

func scanProject(rows pgx.Rows) (*model.Project, error) {
	var p model.Project
	var j projectJSONFields
	var n projectNullableFields

	err := rows.Scan(
		&p.ID, &p.Slug, &j.titleJSON, &j.shortDescJSON, &j.descJSON,
		&p.Category, &p.Status, &p.Tags, &p.TechStack, &j.goalDescJSON,
		&n.githubURL, &n.demoURL, &n.siteURL, &n.videoURL,
		&n.thumbnail, &p.Images, &p.Stars,
		&p.Featured, &p.Order, &p.CreatedAt, &p.UpdatedAt,
		&j.problemJSON, &j.approachJSON, &j.outcomeJSON,
		&j.techChoicesJSON, &j.highlightsJSON,
		&n.timelineStarted, &n.timelineShipped, &j.credsJSON,
	)
	if err != nil {
		return nil, err
	}

	if err := populateProjectFields(&p, j, n); err != nil {
		return nil, err
	}
	return &p, nil
}

func scanProjectRow(row pgx.Row) (*model.Project, error) {
	var p model.Project
	var j projectJSONFields
	var n projectNullableFields

	err := row.Scan(
		&p.ID, &p.Slug, &j.titleJSON, &j.shortDescJSON, &j.descJSON,
		&p.Category, &p.Status, &p.Tags, &p.TechStack, &j.goalDescJSON,
		&n.githubURL, &n.demoURL, &n.siteURL, &n.videoURL,
		&n.thumbnail, &p.Images, &p.Stars,
		&p.Featured, &p.Order, &p.CreatedAt, &p.UpdatedAt,
		&j.problemJSON, &j.approachJSON, &j.outcomeJSON,
		&j.techChoicesJSON, &j.highlightsJSON,
		&n.timelineStarted, &n.timelineShipped, &j.credsJSON,
	)
	if err != nil {
		return nil, err
	}

	if err := populateProjectFields(&p, j, n); err != nil {
		return nil, err
	}
	return &p, nil
}

func populateProjectFields(p *model.Project, j projectJSONFields, n projectNullableFields) error {
	if err := unmarshalJSON(j.titleJSON, &p.Title); err != nil {
		return err
	}
	if err := unmarshalJSON(j.shortDescJSON, &p.ShortDescription); err != nil {
		return err
	}
	if err := unmarshalJSON(j.descJSON, &p.Description); err != nil {
		return err
	}
	if err := unmarshalJSON(j.goalDescJSON, &p.GoalDescription); err != nil {
		return err
	}
	if err := unmarshalJSON(j.problemJSON, &p.Problem); err != nil {
		return err
	}
	if err := unmarshalJSON(j.approachJSON, &p.Approach); err != nil {
		return err
	}
	if err := unmarshalJSON(j.outcomeJSON, &p.Outcome); err != nil {
		return err
	}
	if err := unmarshalJSON(j.techChoicesJSON, &p.TechChoices); err != nil {
		return err
	}
	if err := unmarshalJSON(j.highlightsJSON, &p.Highlights); err != nil {
		return err
	}
	if err := unmarshalJSON(j.credsJSON, &p.DemoCredentials); err != nil {
		return err
	}

	if n.githubURL != nil {
		p.GithubURL = *n.githubURL
	}
	if n.demoURL != nil {
		p.DemoURL = *n.demoURL
	}
	if n.siteURL != nil {
		p.SiteURL = *n.siteURL
	}
	if n.videoURL != nil {
		p.VideoURL = *n.videoURL
	}
	if n.thumbnail != nil {
		p.ThumbnailURL = *n.thumbnail
	}
	p.TimelineStarted = n.timelineStarted
	p.TimelineShipped = n.timelineShipped

	if p.Tags == nil {
		p.Tags = []string{}
	}
	if p.TechStack == nil {
		p.TechStack = []string{}
	}
	if p.Images == nil {
		p.Images = []string{}
	}
	if p.TechChoices == nil {
		p.TechChoices = []model.TechChoice{}
	}
	if p.Highlights == nil {
		p.Highlights = []model.Highlight{}
	}
	if p.DemoCredentials == nil {
		p.DemoCredentials = []model.DemoCredential{}
	}
	return nil
}
