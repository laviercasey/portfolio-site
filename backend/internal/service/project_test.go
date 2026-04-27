package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/portfolio/backend/internal/model"
	"github.com/portfolio/backend/internal/testutil"
)

func projectColumnsSlice() []string {
	return []string{
		"id", "slug", "title", "short_desc", "description",
		"category", "status", "tags", "tech_stack", "goal_desc",
		"github_url", "demo_url", "site_url", "video_url",
		"thumbnail", "images", "stars", "featured", "sort_order",
		"created_at", "updated_at",
		"problem", "approach", "outcome",
		"tech_choices", "highlights",
		"timeline_started", "timeline_shipped", "demo_credentials",
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func fullProjectRow(t *testing.T, id uuid.UUID, slug string) []any {
	t.Helper()
	title := mustJSON(t, model.I18n{"en": "T"})
	shortD := mustJSON(t, model.I18n{"en": "S"})
	desc := mustJSON(t, model.I18n{"en": "D"})
	goal := mustJSON(t, model.I18n{"en": "G"})
	problem := mustJSON(t, model.I18n{"en": "P"})
	approach := mustJSON(t, model.I18n{"en": "A"})
	outcome := mustJSON(t, model.I18n{"en": "O"})
	techChoices := mustJSON(t, []model.TechChoice{{Tech: "Go", Reason: model.I18n{"en": "fast"}}})
	highlights := mustJSON(t, []model.Highlight{{Label: model.I18n{"en": "users"}, Value: "1000"}})
	creds := mustJSON(t, []model.DemoCredential{{Role: "admin", Login: "u", Password: "p"}})
	now := time.Now().UTC()

	return []any{
		id, slug, title, shortD, desc,
		"web", "completed", []string{"go", "web"}, []string{"Go", "Postgres"}, goal,
		(*string)(nil), (*string)(nil), (*string)(nil), (*string)(nil),
		(*string)(nil), []string{"img1.png"}, int(5), true, 1,
		now, now,
		problem, approach, outcome,
		techChoices, highlights,
		(*time.Time)(nil), (*time.Time)(nil), creds,
	}
}

func TestProjectService_List_AllCategories(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	id := uuid.New()
	row := fullProjectRow(t, id, "slug-1")

	pool.ExpectQuery(`SELECT .+ FROM projects ORDER BY sort_order`).
		WillReturnRows(pgxmock.NewRows(projectColumnsSlice()).AddRow(row...))

	got, err := svc.List(context.Background(), "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len: got %d", len(got))
	}
	if got[0].Slug != "slug-1" {
		t.Errorf("slug: %q", got[0].Slug)
	}
	if len(got[0].TechChoices) != 1 {
		t.Errorf("techChoices: %v", got[0].TechChoices)
	}
}

func TestProjectService_List_FilteredByCategory(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	pool.ExpectQuery(`SELECT .+ FROM projects WHERE category = \$1`).
		WithArgs("mobile").
		WillReturnRows(pgxmock.NewRows(projectColumnsSlice()))

	got, err := svc.List(context.Background(), "mobile")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil")
	}
	if len(got) != 0 {
		t.Errorf("len: %d", len(got))
	}
}

func TestProjectService_List_QueryError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	pool.ExpectQuery(`SELECT .+ FROM projects`).WillReturnError(errors.New("db"))

	if _, err := svc.List(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	}
}

func TestProjectService_List_ScanError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	id := uuid.New()
	row := fullProjectRow(t, id, "s")
	rows := pgxmock.NewRows(projectColumnsSlice()).AddRow(row...)
	rows.RowError(0, errors.New("scan boom"))
	pool.ExpectQuery(`SELECT .+ FROM projects`).WillReturnRows(rows)

	if _, err := svc.List(context.Background(), ""); err == nil {
		t.Fatal("expected scan error")
	}
}

func TestProjectService_List_RowsErr(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	id := uuid.New()
	row := fullProjectRow(t, id, "s")
	rows := pgxmock.NewRows(projectColumnsSlice()).
		AddRow(row...).
		RowError(0, errors.New("iter fail"))
	pool.ExpectQuery(`SELECT .+ FROM projects`).WillReturnRows(rows)

	if _, err := svc.List(context.Background(), ""); err == nil {
		t.Fatal("expected iterate error")
	}
}

func TestProjectService_GetBySlug_Success(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	id := uuid.New()
	row := fullProjectRow(t, id, "my-slug")
	pool.ExpectQuery(`SELECT .+ FROM projects WHERE slug = \$1`).
		WithArgs("my-slug").
		WillReturnRows(pgxmock.NewRows(projectColumnsSlice()).AddRow(row...))

	got, err := svc.GetBySlug(context.Background(), "my-slug")
	if err != nil {
		t.Fatalf("GetBySlug: %v", err)
	}
	if got.Slug != "my-slug" {
		t.Errorf("slug: %q", got.Slug)
	}
}

func TestProjectService_GetBySlug_NotFound(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	pool.ExpectQuery(`SELECT .+ FROM projects WHERE slug = \$1`).
		WithArgs("missing").
		WillReturnError(pgx.ErrNoRows)

	_, err := svc.GetBySlug(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestProjectService_Delete_Success(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	id := uuid.New()
	pool.ExpectExec(`DELETE FROM projects WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	if err := svc.Delete(context.Background(), id); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestProjectService_Delete_NoRows(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	id := uuid.New()
	pool.ExpectExec(`DELETE FROM projects WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err := svc.Delete(context.Background(), id)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("got %v, want pgx.ErrNoRows", err)
	}
}

func TestProjectService_Delete_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	id := uuid.New()
	pool.ExpectExec(`DELETE FROM projects WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(errors.New("db"))

	if err := svc.Delete(context.Background(), id); err == nil {
		t.Fatal("expected error")
	}
}

func TestProjectService_Create_HappyPath(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	input := model.CreateProjectInput{
		Slug:             "new-slug",
		Title:            model.I18n{"en": "New"},
		ShortDescription: model.I18n{"en": "short"},
		Desc:             model.I18n{"en": "full"},
		Category:         "web",
		Status:           "completed",
		Tags:             []string{"a"},
		TechStack:        []string{"Go"},
	}

	id := uuid.New()
	row := fullProjectRow(t, id, "new-slug")
	pool.ExpectQuery(`INSERT INTO projects`).
		WithArgs(testutil.AnyArgs(25)...).
		WillReturnRows(pgxmock.NewRows(projectColumnsSlice()).AddRow(row...))

	got, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != id {
		t.Errorf("ID mismatch")
	}
}

func TestProjectService_Create_NilSlicesNormalized(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	input := model.CreateProjectInput{
		Slug:     "s",
		Title:    model.I18n{"en": "T"},
		Category: "web",
		Status:   "completed",
	}

	id := uuid.New()
	row := fullProjectRow(t, id, "s")
	pool.ExpectQuery(`INSERT INTO projects`).
		WithArgs(testutil.AnyArgs(25)...).
		WillReturnRows(pgxmock.NewRows(projectColumnsSlice()).AddRow(row...))

	got, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.Tags == nil {
		t.Error("Tags should be non-nil after populate")
	}
	if got.TechStack == nil {
		t.Error("TechStack should be non-nil")
	}
}

func TestProjectService_Create_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	pool.ExpectQuery(`INSERT INTO projects`).WillReturnError(errors.New("constraint"))

	input := model.CreateProjectInput{Slug: "s", Title: model.I18n{"en": "T"}, Category: "web", Status: "completed"}
	if _, err := svc.Create(context.Background(), input); err == nil {
		t.Fatal("expected error")
	}
}

func TestProjectService_Update_HappyPath(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	id := uuid.New()
	row := fullProjectRow(t, id, "old-slug")

	pool.ExpectQuery(`SELECT .+ FROM projects WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows(projectColumnsSlice()).AddRow(row...))

	newSlug := "new-slug"
	newTitle := model.I18n{"en": "NewTitle"}
	featured := false
	input := model.UpdateProjectInput{
		Slug:     &newSlug,
		Title:    &newTitle,
		Featured: &featured,
	}

	updRow := fullProjectRow(t, id, "new-slug")
	pool.ExpectQuery(`UPDATE projects SET`).
		WithArgs(testutil.AnyArgs(26)...).
		WillReturnRows(pgxmock.NewRows(projectColumnsSlice()).AddRow(updRow...))

	got, err := svc.Update(context.Background(), id, input)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got.Slug != "new-slug" {
		t.Errorf("slug: %q", got.Slug)
	}
}

func TestProjectService_Update_SelectError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	id := uuid.New()
	pool.ExpectQuery(`SELECT .+ FROM projects WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(pgx.ErrNoRows)

	if _, err := svc.Update(context.Background(), id, model.UpdateProjectInput{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestProjectService_Update_UpdateError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	id := uuid.New()
	row := fullProjectRow(t, id, "s")

	pool.ExpectQuery(`SELECT .+ FROM projects WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows(projectColumnsSlice()).AddRow(row...))
	pool.ExpectQuery(`UPDATE projects SET`).
		WithArgs(testutil.AnyArgs(26)...).
		WillReturnError(errors.New("update fail"))

	if _, err := svc.Update(context.Background(), id, model.UpdateProjectInput{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestProjectService_Update_AllFieldsApplied(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	id := uuid.New()
	row := fullProjectRow(t, id, "old")

	pool.ExpectQuery(`SELECT .+ FROM projects WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows(projectColumnsSlice()).AddRow(row...))

	updRow := fullProjectRow(t, id, "new")
	pool.ExpectQuery(`UPDATE projects SET`).
		WithArgs(testutil.AnyArgs(26)...).
		WillReturnRows(pgxmock.NewRows(projectColumnsSlice()).AddRow(updRow...))

	slug := "new"
	title := model.I18n{"en": "T2"}
	short := model.I18n{"en": "S2"}
	desc := model.I18n{"en": "D2"}
	cat := "mobile"
	status := "in_development"
	tags := []string{"x"}
	stack := []string{"Rust"}
	goal := model.I18n{"en": "G2"}
	gh := "https://gh"
	demo := "https://demo"
	site := "https://site"
	video := "https://video"
	thumb := "https://thumb"
	images := []string{"a.png"}
	featured := true
	order := 99
	problem := model.I18n{"en": "P2"}
	approach := model.I18n{"en": "A2"}
	outcome := model.I18n{"en": "O2"}
	techChoices := []model.TechChoice{{Tech: "Rust"}}
	highlights := []model.Highlight{{Value: "x"}}
	started := time.Now()
	shipped := time.Now()
	creds := []model.DemoCredential{{Role: "admin"}}

	input := model.UpdateProjectInput{
		Slug: &slug, Title: &title, ShortDescription: &short, Desc: &desc,
		Category: &cat, Status: &status, Tags: &tags, TechStack: &stack,
		GoalDescription: &goal, GithubURL: &gh, DemoURL: &demo, SiteURL: &site,
		VideoURL: &video, ThumbnailURL: &thumb, Images: &images,
		Featured: &featured, Order: &order,
		Problem: &problem, Approach: &approach, Outcome: &outcome,
		TechChoices: &techChoices, Highlights: &highlights,
		TimelineStarted: &started, TimelineShipped: &shipped,
		DemoCredentials: &creds,
	}

	if _, err := svc.Update(context.Background(), id, input); err != nil {
		t.Fatalf("Update: %v", err)
	}
}

func TestProjectService_PopulateNullables(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewProjectService(db)

	id := uuid.New()
	title := mustJSON(t, model.I18n{"en": "T"})
	gh := "https://gh.com"
	demo := "https://demo.com"
	site := "https://site.com"
	video := "https://video.com"
	thumb := "https://thumb.png"
	started := time.Now().UTC()
	shipped := started.Add(24 * time.Hour)

	row := []any{
		id, "slug", title, title, title,
		"web", "completed", []string{}, []string{}, title,
		&gh, &demo, &site, &video,
		&thumb, []string{}, int(0), false, 0,
		time.Now(), time.Now(),
		title, title, title,
		[]byte(`[]`), []byte(`[]`),
		&started, &shipped, []byte(`[]`),
	}

	pool.ExpectQuery(`SELECT .+ FROM projects WHERE slug = \$1`).
		WithArgs("slug").
		WillReturnRows(pgxmock.NewRows(projectColumnsSlice()).AddRow(row...))

	got, err := svc.GetBySlug(context.Background(), "slug")
	if err != nil {
		t.Fatalf("GetBySlug: %v", err)
	}
	if got.GithubURL != gh {
		t.Errorf("gh: %q", got.GithubURL)
	}
	if got.DemoURL != demo {
		t.Errorf("demo: %q", got.DemoURL)
	}
	if got.SiteURL != site {
		t.Errorf("site: %q", got.SiteURL)
	}
	if got.VideoURL != video {
		t.Errorf("video: %q", got.VideoURL)
	}
	if got.ThumbnailURL != thumb {
		t.Errorf("thumb: %q", got.ThumbnailURL)
	}
	if got.TimelineStarted == nil {
		t.Error("timelineStarted nil")
	}
	if got.TimelineShipped == nil {
		t.Error("timelineShipped nil")
	}
}
