package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/portfolio/backend/internal/database"
)

const (
	githubAPIBaseURL    = "https://api.github.com"
	githubAPIVersion    = "2022-11-28"
	githubClientTimeout = 10 * time.Second
)

var githubURLRe = regexp.MustCompile(`(?i)github\.com/([^/?#\s]+)/([^/?#\s]+)`)

type starProject struct {
	ID        uuid.UUID
	Slug      string
	GithubURL string
}

type githubRepoResponse struct {
	StargazersCount int `json:"stargazers_count"`
}

type GithubStarsService struct {
	db       *database.DB
	token    string
	interval time.Duration
	logger   *slog.Logger
	client   *http.Client
	baseURL  string
}

func NewGithubStarsService(
	db *database.DB,
	token string,
	interval time.Duration,
	logger *slog.Logger,
) *GithubStarsService {
	if logger == nil {
		logger = slog.Default()
	}
	if interval <= 0 {
		interval = 6 * time.Hour
	}
	return &GithubStarsService{
		db:       db,
		token:    token,
		interval: interval,
		logger:   logger,
		client:   &http.Client{Timeout: githubClientTimeout},
		baseURL:  githubAPIBaseURL,
	}
}

func (s *GithubStarsService) Start(ctx context.Context) {
	go s.run(ctx)
}

func (s *GithubStarsService) run(ctx context.Context) {
	if err := s.SyncAll(ctx); err != nil {
		s.logger.Warn("github stars: initial sync failed", slog.String("error", err.Error()))
	}

	t := time.NewTicker(s.interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := s.SyncAll(ctx); err != nil {
				s.logger.Warn("github stars: periodic sync failed", slog.String("error", err.Error()))
			}
		}
	}
}

func (s *GithubStarsService) SyncAll(ctx context.Context) error {
	projects, err := s.listProjects(ctx)
	if err != nil {
		return fmt.Errorf("list projects: %w", err)
	}

	for _, p := range projects {
		if err := ctx.Err(); err != nil {
			return err
		}

		owner, repo, ok := parseGithubURL(p.GithubURL)
		if !ok {
			s.logger.Warn(
				"github stars: cannot parse url",
				slog.String("slug", p.Slug),
				slog.String("url", p.GithubURL),
			)
			continue
		}

		stars, err := s.fetchStars(ctx, owner, repo)
		if err != nil {
			s.logger.Warn(
				"github stars: fetch failed",
				slog.String("slug", p.Slug),
				slog.String("repo", owner+"/"+repo),
				slog.String("error", err.Error()),
			)
			continue
		}

		if err := s.updateStars(ctx, p.ID, stars); err != nil {
			s.logger.Warn(
				"github stars: update failed",
				slog.String("slug", p.Slug),
				slog.String("error", err.Error()),
			)
			continue
		}

		s.logger.Info(
			"github stars: synced",
			slog.String("slug", p.Slug),
			slog.String("repo", owner+"/"+repo),
			slog.Int("stars", stars),
		)
	}
	return nil
}

func (s *GithubStarsService) listProjects(ctx context.Context) ([]starProject, error) {
	rows, err := s.db.Pool.Query(
		ctx,
		`SELECT id, slug, github_url FROM projects WHERE github_url IS NOT NULL AND github_url <> ''`,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var out []starProject
	for rows.Next() {
		var p starProject
		if err := rows.Scan(&p.ID, &p.Slug, &p.GithubURL); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate: %w", err)
	}
	return out, nil
}

func (s *GithubStarsService) fetchStars(ctx context.Context, owner, repo string) (int, error) {
	apiURL := fmt.Sprintf("%s/repos/%s/%s", s.baseURL, owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return 0, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", githubAPIVersion)
	if s.token != "" {
		req.Header.Set("Authorization", "Bearer "+s.token)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("http: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	switch {
	case resp.StatusCode == http.StatusNotFound:
		return 0, nil
	case resp.StatusCode == http.StatusTooManyRequests:
		return 0, errors.New("rate limited (429)")
	case resp.StatusCode >= 500:
		return 0, fmt.Errorf("upstream %d", resp.StatusCode)
	case resp.StatusCode < 200 || resp.StatusCode >= 300:
		return 0, fmt.Errorf("status %d", resp.StatusCode)
	}

	const maxResponseBytes = 1 << 20
	var data githubRepoResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxResponseBytes)).Decode(&data); err != nil {
		return 0, fmt.Errorf("decode: %w", err)
	}
	if data.StargazersCount < 0 {
		data.StargazersCount = 0
	}
	return data.StargazersCount, nil
}

func (s *GithubStarsService) updateStars(ctx context.Context, id uuid.UUID, stars int) error {
	if _, err := s.db.Pool.Exec(
		ctx,
		`UPDATE projects SET stars = $1, updated_at = now() WHERE id = $2`,
		stars, id,
	); err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func parseGithubURL(raw string) (owner, repo string, ok bool) {
	if raw == "" {
		return "", "", false
	}
	m := githubURLRe.FindStringSubmatch(raw)
	if len(m) != 3 {
		return "", "", false
	}
	owner = strings.TrimSpace(m[1])
	repo = strings.TrimSuffix(strings.TrimSpace(m[2]), ".git")
	if owner == "" || repo == "" {
		return "", "", false
	}
	return owner, repo, true
}
