package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/portfolio/backend/internal/testutil"
)

func TestParseGithubURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantOwner string
		wantRepo  string
		wantOK    bool
	}{
		{"plain", "https://github.com/owner/repo", "owner", "repo", true},
		{"with .git", "https://github.com/owner/repo.git", "owner", "repo", true},
		{"with query", "https://github.com/owner/repo?tab=readme", "owner", "repo", true},
		{"with fragment", "https://github.com/owner/repo#section", "owner", "repo", true},
		{"trailing slash", "https://github.com/owner/repo/", "owner", "repo", true},
		{"http scheme", "http://github.com/owner/repo", "owner", "repo", true},
		{"no scheme", "github.com/owner/repo", "owner", "repo", true},
		{"www prefix", "https://www.github.com/owner/repo", "owner", "repo", true},
		{"uppercase host", "https://GitHub.com/Foo/Bar", "Foo", "Bar", true},
		{"subpath", "https://github.com/owner/repo/issues/1", "owner", "repo", true},
		{"empty", "", "", "", false},
		{"non-github", "https://gitlab.com/owner/repo", "", "", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			owner, repo, ok := parseGithubURL(tc.input)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v (input=%q)", ok, tc.wantOK, tc.input)
			}
			if !ok {
				return
			}
			if owner != tc.wantOwner {
				t.Errorf("owner = %q, want %q", owner, tc.wantOwner)
			}
			if repo != tc.wantRepo {
				t.Errorf("repo = %q, want %q", repo, tc.wantRepo)
			}
		})
	}
}

func TestGithubStarsService_FetchStars_HappyPath(t *testing.T) {
	t.Parallel()

	var gotAuth, gotAccept, gotAPIVer, gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotAccept = r.Header.Get("Accept")
		gotAPIVer = r.Header.Get("X-GitHub-Api-Version")
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"stargazers_count": 42}`)
	}))
	defer ts.Close()

	db, _ := testutil.NewMockDB(t)
	svc := &GithubStarsService{
		db:      db,
		token:   "tok-123",
		logger:  testutil.SilentLogger(),
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: ts.URL,
	}

	stars, err := svc.fetchStars(context.Background(), "laviercasey", "portfolio-site")
	if err != nil {
		t.Fatalf("fetchStars: %v", err)
	}
	if stars != 42 {
		t.Errorf("stars = %d, want 42", stars)
	}
	if gotPath != "/repos/laviercasey/portfolio-site" {
		t.Errorf("path = %q", gotPath)
	}
	if gotAuth != "Bearer tok-123" {
		t.Errorf("auth = %q, want Bearer tok-123", gotAuth)
	}
	if gotAccept != "application/vnd.github+json" {
		t.Errorf("accept = %q", gotAccept)
	}
	if gotAPIVer != "2022-11-28" {
		t.Errorf("api version = %q", gotAPIVer)
	}
}

func TestGithubStarsService_FetchStars_NoTokenOmitsHeader(t *testing.T) {
	t.Parallel()

	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		_, _ = fmt.Fprint(w, `{"stargazers_count": 5}`)
	}))
	defer ts.Close()

	db, _ := testutil.NewMockDB(t)
	svc := &GithubStarsService{
		db:      db,
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: ts.URL,
		logger:  testutil.SilentLogger(),
	}

	if _, err := svc.fetchStars(context.Background(), "a", "b"); err != nil {
		t.Fatalf("fetchStars: %v", err)
	}
	if gotAuth != "" {
		t.Errorf("auth header should be empty when no token, got %q", gotAuth)
	}
}

func TestGithubStarsService_FetchStars_StatusCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		status     int
		wantStars  int
		wantErrSub string
	}{
		{"404 returns 0 no error", http.StatusNotFound, 0, ""},
		{"429 errors", http.StatusTooManyRequests, 0, "rate limited"},
		{"500 errors", http.StatusInternalServerError, 0, "upstream 500"},
		{"503 errors", http.StatusServiceUnavailable, 0, "upstream 503"},
		{"403 errors", http.StatusForbidden, 0, "status 403"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
			}))
			defer ts.Close()

			db, _ := testutil.NewMockDB(t)
			svc := &GithubStarsService{
				db:      db,
				client:  &http.Client{Timeout: 2 * time.Second},
				baseURL: ts.URL,
				logger:  testutil.SilentLogger(),
			}

			stars, err := svc.fetchStars(context.Background(), "o", "r")
			if tc.wantErrSub == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if stars != tc.wantStars {
					t.Errorf("stars = %d, want %d", stars, tc.wantStars)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErrSub) {
				t.Fatalf("error = %v, want contains %q", err, tc.wantErrSub)
			}
		})
	}
}

func TestGithubStarsService_FetchStars_NetworkError(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := &GithubStarsService{
		db:      db,
		client:  &http.Client{Timeout: 100 * time.Millisecond},
		baseURL: "http://127.0.0.1:1",
		logger:  testutil.SilentLogger(),
	}
	if _, err := svc.fetchStars(context.Background(), "o", "r"); err == nil {
		t.Fatal("expected network error")
	}
}

func TestGithubStarsService_FetchStars_DecodeError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `not json`)
	}))
	defer ts.Close()

	db, _ := testutil.NewMockDB(t)
	svc := &GithubStarsService{
		db:      db,
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: ts.URL,
		logger:  testutil.SilentLogger(),
	}
	if _, err := svc.fetchStars(context.Background(), "o", "r"); err == nil {
		t.Fatal("expected decode error")
	}
}

func TestGithubStarsService_SyncAll_HappyPath(t *testing.T) {
	t.Parallel()

	hits := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		switch r.URL.Path {
		case "/repos/laviercasey/portfolio-site":
			_, _ = fmt.Fprint(w, `{"stargazers_count": 42}`)
		case "/repos/laviercasey/another":
			_, _ = fmt.Fprint(w, `{"stargazers_count": 7}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	db, pool := testutil.NewMockDB(t)
	id1, id2 := uuid.New(), uuid.New()

	pool.ExpectQuery(`SELECT id, slug, github_url FROM projects`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "slug", "github_url"}).
			AddRow(id1, "portfolio", "https://github.com/laviercasey/portfolio-site").
			AddRow(id2, "another", "https://github.com/laviercasey/another.git"))
	pool.ExpectExec(`UPDATE projects SET stars = \$1, updated_at = now\(\) WHERE id = \$2`).
		WithArgs(42, id1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	pool.ExpectExec(`UPDATE projects SET stars = \$1, updated_at = now\(\) WHERE id = \$2`).
		WithArgs(7, id2).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	svc := &GithubStarsService{
		db:      db,
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: ts.URL,
		logger:  testutil.SilentLogger(),
	}

	if err := svc.SyncAll(context.Background()); err != nil {
		t.Fatalf("SyncAll: %v", err)
	}
	if hits != 2 {
		t.Errorf("hits = %d, want 2", hits)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}

func TestGithubStarsService_SyncAll_404SetsZero(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer ts.Close()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()

	pool.ExpectQuery(`SELECT id, slug, github_url FROM projects`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "slug", "github_url"}).
			AddRow(id, "deleted", "https://github.com/x/y"))
	pool.ExpectExec(`UPDATE projects SET stars = \$1, updated_at = now\(\) WHERE id = \$2`).
		WithArgs(0, id).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	svc := &GithubStarsService{
		db:      db,
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: ts.URL,
		logger:  testutil.SilentLogger(),
	}

	if err := svc.SyncAll(context.Background()); err != nil {
		t.Fatalf("SyncAll: %v", err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}

func TestGithubStarsService_SyncAll_5xxSkipsUpdate(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer ts.Close()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()

	pool.ExpectQuery(`SELECT id, slug, github_url FROM projects`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "slug", "github_url"}).
			AddRow(id, "p", "https://github.com/x/y"))

	svc := &GithubStarsService{
		db:      db,
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: ts.URL,
		logger:  testutil.SilentLogger(),
	}

	if err := svc.SyncAll(context.Background()); err != nil {
		t.Fatalf("SyncAll: %v", err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}

func TestGithubStarsService_SyncAll_429SkipsUpdate(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer ts.Close()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()

	pool.ExpectQuery(`SELECT id, slug, github_url FROM projects`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "slug", "github_url"}).
			AddRow(id, "p", "https://github.com/x/y"))

	svc := &GithubStarsService{
		db:      db,
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: ts.URL,
		logger:  testutil.SilentLogger(),
	}

	if err := svc.SyncAll(context.Background()); err != nil {
		t.Fatalf("SyncAll: %v", err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}

func TestGithubStarsService_SyncAll_BadURLSkipped(t *testing.T) {
	t.Parallel()

	hit := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hit = true
		_, _ = fmt.Fprint(w, `{"stargazers_count": 0}`)
	}))
	defer ts.Close()

	db, pool := testutil.NewMockDB(t)
	id := uuid.New()

	pool.ExpectQuery(`SELECT id, slug, github_url FROM projects`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "slug", "github_url"}).
			AddRow(id, "p", "not-a-github-url"))

	svc := &GithubStarsService{
		db:      db,
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: ts.URL,
		logger:  testutil.SilentLogger(),
	}

	if err := svc.SyncAll(context.Background()); err != nil {
		t.Fatalf("SyncAll: %v", err)
	}
	if hit {
		t.Error("API should not be hit for unparseable URLs")
	}
}

func TestGithubStarsService_SyncAll_QueryError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery(`SELECT id, slug, github_url FROM projects`).
		WillReturnError(errors.New("db down"))

	svc := &GithubStarsService{
		db:      db,
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: "http://unused",
		logger:  testutil.SilentLogger(),
	}

	if err := svc.SyncAll(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}

func TestGithubStarsService_SyncAll_UpdateErrorContinues(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/x/y":
			_, _ = fmt.Fprint(w, `{"stargazers_count": 1}`)
		case "/repos/a/b":
			_, _ = fmt.Fprint(w, `{"stargazers_count": 2}`)
		}
	}))
	defer ts.Close()

	db, pool := testutil.NewMockDB(t)
	id1, id2 := uuid.New(), uuid.New()

	pool.ExpectQuery(`SELECT id, slug, github_url FROM projects`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "slug", "github_url"}).
			AddRow(id1, "p1", "https://github.com/x/y").
			AddRow(id2, "p2", "https://github.com/a/b"))
	pool.ExpectExec(`UPDATE projects SET stars`).
		WithArgs(1, id1).
		WillReturnError(errors.New("update failed"))
	pool.ExpectExec(`UPDATE projects SET stars`).
		WithArgs(2, id2).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	svc := &GithubStarsService{
		db:      db,
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: ts.URL,
		logger:  testutil.SilentLogger(),
	}

	if err := svc.SyncAll(context.Background()); err != nil {
		t.Fatalf("SyncAll: %v", err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}

func TestGithubStarsService_SyncAll_ContextCanceled(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	id1, id2 := uuid.New(), uuid.New()
	pool.ExpectQuery(`SELECT id, slug, github_url FROM projects`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "slug", "github_url"}).
			AddRow(id1, "p1", "https://github.com/x/y").
			AddRow(id2, "p2", "https://github.com/a/b"))

	svc := &GithubStarsService{
		db:      db,
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: "http://unused",
		logger:  testutil.SilentLogger(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := svc.SyncAll(ctx); err == nil {
		t.Fatal("expected ctx canceled error")
	}
}

func TestNewGithubStarsService_DefaultsApplied(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewGithubStarsService(db, "", 0, nil)
	if svc.interval != 6*time.Hour {
		t.Errorf("interval = %v, want 6h", svc.interval)
	}
	if svc.logger == nil {
		t.Error("logger should default to slog.Default()")
	}
	if svc.baseURL != githubAPIBaseURL {
		t.Errorf("baseURL = %q, want %q", svc.baseURL, githubAPIBaseURL)
	}
	if svc.client == nil {
		t.Error("client should be initialized")
	}
}

func TestNewGithubStarsService_CustomInterval(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewGithubStarsService(db, "tok", 30*time.Minute, testutil.SilentLogger())
	if svc.interval != 30*time.Minute {
		t.Errorf("interval = %v, want 30m", svc.interval)
	}
	if svc.token != "tok" {
		t.Errorf("token = %q", svc.token)
	}
}

func TestGithubStarsService_Start_StopsOnContextCancel(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	pool.ExpectQuery(`SELECT id, slug, github_url FROM projects`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "slug", "github_url"}))

	svc := &GithubStarsService{
		db:       db,
		interval: 24 * time.Hour,
		logger:   testutil.SilentLogger(),
		client:   &http.Client{Timeout: time.Second},
		baseURL:  "http://unused",
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		svc.run(ctx)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("run() did not return after context cancel")
	}
}
