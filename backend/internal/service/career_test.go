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

func educationCols() []string {
	return []string{
		"id", "institution", "degree", "field",
		"start_year", "end_year", "description",
		"logo_url", "related_project_slugs",
	}
}

func workCols() []string {
	return []string{
		"id", "company", "position", "start_date", "end_date",
		"is_current", "description", "technologies", "logo_url",
		"achievements", "full_description",
	}
}

func certCols() []string {
	return []string{
		"id", "title", "issuer", "issue_date",
		"credential_id", "url", "image_url",
	}
}

func pubCols() []string {
	return []string{
		"id", "title", "journal", "year", "doi", "url", "abstract",
	}
}

func i18nJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal i18n: %v", err)
	}
	return b
}

func TestCareerService_Delete_InvalidType(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	err := svc.Delete(context.Background(), "bogus", uuid.New())
	if !errors.Is(err, ErrInvalidCareerType) {
		t.Errorf("got %v, want ErrInvalidCareerType", err)
	}
}

func TestCareerService_Delete_Success(t *testing.T) {
	t.Parallel()

	cases := []struct {
		careerType string
		tableRe    string
	}{
		{"education", `DELETE FROM education WHERE id = \$1`},
		{"work", `DELETE FROM work_history WHERE id = \$1`},
		{"certificate", `DELETE FROM certificates WHERE id = \$1`},
		{"publication", `DELETE FROM publications WHERE id = \$1`},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.careerType, func(t *testing.T) {
			t.Parallel()
			db, pool := testutil.NewMockDB(t)
			svc := NewCareerService(db)
			id := uuid.New()

			pool.ExpectExec(tc.tableRe).
				WithArgs(id).
				WillReturnResult(pgxmock.NewResult("DELETE", 1))

			if err := svc.Delete(context.Background(), tc.careerType, id); err != nil {
				t.Fatalf("Delete: %v", err)
			}
		})
	}
}

func TestCareerService_Delete_NoRows(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)
	id := uuid.New()

	pool.ExpectExec(`DELETE FROM education WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err := svc.Delete(context.Background(), "education", id)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("got %v, want pgx.ErrNoRows", err)
	}
}

func TestCareerService_Delete_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)
	id := uuid.New()

	pool.ExpectExec(`DELETE FROM work_history WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(errors.New("db fail"))

	if err := svc.Delete(context.Background(), "work", id); err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_Create_InvalidType(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	_, err := svc.Create(context.Background(), "bogus", json.RawMessage(`{}`))
	if !errors.Is(err, ErrInvalidCareerType) {
		t.Errorf("got %v, want ErrInvalidCareerType", err)
	}
}

func TestCareerService_Update_InvalidType(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	_, err := svc.Update(context.Background(), "bogus", uuid.New(), json.RawMessage(`{}`))
	if !errors.Is(err, ErrInvalidCareerType) {
		t.Errorf("got %v, want ErrInvalidCareerType", err)
	}
}

func TestCareerService_GetAll_Success(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	edID := uuid.New()
	wID := uuid.New()
	cID := uuid.New()
	pID := uuid.New()
	inst := i18nJSON(t, model.I18n{"en": "MIT"})
	deg := i18nJSON(t, model.I18n{"en": "BSc"})
	field := i18nJSON(t, model.I18n{"en": "CS"})
	desc := i18nJSON(t, model.I18n{"en": "Studied"})

	pool.ExpectQuery(`SELECT .+ FROM education ORDER BY`).
		WillReturnRows(pgxmock.NewRows(educationCols()).
			AddRow(edID, inst, deg, field, "2018", (*string)(nil), desc, (*string)(nil), []string{}))

	comp := i18nJSON(t, model.I18n{"en": "Acme"})
	pos := i18nJSON(t, model.I18n{"en": "Eng"})
	wDesc := i18nJSON(t, model.I18n{"en": "work desc"})
	ach := i18nJSON(t, []model.I18n{{"en": "shipped"}})
	fullDesc := i18nJSON(t, model.I18n{"en": "long"})

	pool.ExpectQuery(`SELECT .+ FROM work_history ORDER BY`).
		WillReturnRows(pgxmock.NewRows(workCols()).
			AddRow(wID, comp, pos, "2020-01-01", (*string)(nil), true, wDesc, []string{"Go"}, (*string)(nil), ach, fullDesc))

	title := i18nJSON(t, model.I18n{"en": "Cert"})
	issuer := i18nJSON(t, model.I18n{"en": "Org"})
	pool.ExpectQuery(`SELECT .+ FROM certificates ORDER BY`).
		WillReturnRows(pgxmock.NewRows(certCols()).
			AddRow(cID, title, issuer, "2022-05-01", (*string)(nil), (*string)(nil), (*string)(nil)))

	pTitle := i18nJSON(t, model.I18n{"en": "Paper"})
	pJournal := i18nJSON(t, model.I18n{"en": "Journal"})
	pAbs := i18nJSON(t, model.I18n{"en": "abstract"})
	pool.ExpectQuery(`SELECT .+ FROM publications ORDER BY`).
		WillReturnRows(pgxmock.NewRows(pubCols()).
			AddRow(pID, pTitle, pJournal, "2024", (*string)(nil), (*string)(nil), pAbs))

	got, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if len(got.Education) != 1 {
		t.Errorf("education len: got %d", len(got.Education))
	}
	if got.Education[0].StartYear != 2018 {
		t.Errorf("start year: got %d", got.Education[0].StartYear)
	}
	if len(got.WorkHistory) != 1 {
		t.Errorf("work len: got %d", len(got.WorkHistory))
	}
	if !got.WorkHistory[0].Current {
		t.Errorf("current flag")
	}
	if len(got.Certificates) != 1 {
		t.Errorf("cert len: got %d", len(got.Certificates))
	}
	if len(got.Publications) != 1 {
		t.Errorf("pub len: got %d", len(got.Publications))
	}
	if got.Publications[0].Year != 2024 {
		t.Errorf("year: got %d", got.Publications[0].Year)
	}
}

func TestCareerService_GetAll_EducationError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`SELECT .+ FROM education`).WillReturnError(errors.New("db"))

	if _, err := svc.GetAll(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_GetAll_WorkError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`SELECT .+ FROM education`).
		WillReturnRows(pgxmock.NewRows(educationCols()))
	pool.ExpectQuery(`SELECT .+ FROM work_history`).WillReturnError(errors.New("db"))

	if _, err := svc.GetAll(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_GetAll_CertError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`SELECT .+ FROM education`).
		WillReturnRows(pgxmock.NewRows(educationCols()))
	pool.ExpectQuery(`SELECT .+ FROM work_history`).
		WillReturnRows(pgxmock.NewRows(workCols()))
	pool.ExpectQuery(`SELECT .+ FROM certificates`).WillReturnError(errors.New("db"))

	if _, err := svc.GetAll(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_GetAll_PubError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`SELECT .+ FROM education`).
		WillReturnRows(pgxmock.NewRows(educationCols()))
	pool.ExpectQuery(`SELECT .+ FROM work_history`).
		WillReturnRows(pgxmock.NewRows(workCols()))
	pool.ExpectQuery(`SELECT .+ FROM certificates`).
		WillReturnRows(pgxmock.NewRows(certCols()))
	pool.ExpectQuery(`SELECT .+ FROM publications`).WillReturnError(errors.New("db"))

	if _, err := svc.GetAll(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_CreateEducation_HappyPath(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	body := json.RawMessage(`{
      "institution": {"en": "MIT"},
      "degree": {"en": "BSc"},
      "field": {"en": "CS"},
      "startYear": 2015,
      "endYear": 2019,
      "description": {"en": "Studied hard"},
      "logoUrl": "https://x/logo.png",
      "sortOrder": 1,
      "relatedProjectSlugs": ["p1", "p2"]
    }`)

	id := uuid.New()
	inst := i18nJSON(t, model.I18n{"en": "MIT"})
	deg := i18nJSON(t, model.I18n{"en": "BSc"})
	field := i18nJSON(t, model.I18n{"en": "CS"})
	desc := i18nJSON(t, model.I18n{"en": "Studied hard"})
	endYear := "2019"
	logo := "https://x/logo.png"

	pool.ExpectQuery(`INSERT INTO education`).
		WithArgs(testutil.AnyArgs(9)...).
		WillReturnRows(pgxmock.NewRows(educationCols()).
			AddRow(id, inst, deg, field, "2015", &endYear, desc, &logo, []string{"p1", "p2"}))

	got, err := svc.Create(context.Background(), "education", body)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	e, ok := got.(*model.Education)
	if !ok {
		t.Fatalf("want *model.Education, got %T", got)
	}
	if e.StartYear != 2015 {
		t.Errorf("start: %d", e.StartYear)
	}
	if e.EndYear == nil || *e.EndYear != 2019 {
		t.Errorf("end year")
	}
	if e.LogoURL != logo {
		t.Errorf("logo: %q", e.LogoURL)
	}
	if len(e.RelatedProjectSlugs) != 2 {
		t.Errorf("slugs: %v", e.RelatedProjectSlugs)
	}
}

func TestCareerService_CreateEducation_NilEndYearAndSlugs(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	body := json.RawMessage(`{
      "institution": {"en": "MIT"},
      "degree": {"en": "BSc"},
      "field": {"en": "CS"},
      "startYear": 2015,
      "description": {"en": "d"},
      "sortOrder": 0
    }`)

	id := uuid.New()
	inst := i18nJSON(t, model.I18n{"en": "MIT"})
	deg := i18nJSON(t, model.I18n{"en": "BSc"})
	field := i18nJSON(t, model.I18n{"en": "CS"})
	desc := i18nJSON(t, model.I18n{"en": "d"})

	pool.ExpectQuery(`INSERT INTO education`).
		WithArgs(testutil.AnyArgs(9)...).
		WillReturnRows(pgxmock.NewRows(educationCols()).
			AddRow(id, inst, deg, field, "2015", (*string)(nil), desc, (*string)(nil), []string(nil)))

	got, err := svc.Create(context.Background(), "education", body)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	e := got.(*model.Education)
	if e.EndYear != nil {
		t.Errorf("expected nil EndYear, got %d", *e.EndYear)
	}
	if e.LogoURL != "" {
		t.Errorf("expected empty logo, got %q", e.LogoURL)
	}
	if e.RelatedProjectSlugs == nil {
		t.Error("expected non-nil slice")
	}
}

func TestCareerService_CreateEducation_BadJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	_, err := svc.Create(context.Background(), "education", json.RawMessage(`{bad}`))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_CreateEducation_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`INSERT INTO education`).
		WithArgs(testutil.AnyArgs(9)...).
		WillReturnError(errors.New("db"))

	body := json.RawMessage(`{"institution":{"en":"M"},"degree":{"en":"B"},"field":{"en":"C"},"startYear":2020,"description":{}}`)
	if _, err := svc.Create(context.Background(), "education", body); err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_UpdateEducation_HappyPath(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	id := uuid.New()
	body := json.RawMessage(`{
      "institution": {"en": "MIT"},
      "degree": {"en": "BSc"},
      "field": {"en": "CS"},
      "startYear": 2015,
      "description": {"en": "d"},
      "sortOrder": 0
    }`)

	inst := i18nJSON(t, model.I18n{"en": "MIT"})
	deg := i18nJSON(t, model.I18n{"en": "BSc"})
	field := i18nJSON(t, model.I18n{"en": "CS"})
	desc := i18nJSON(t, model.I18n{"en": "d"})

	pool.ExpectQuery(`UPDATE education SET`).
		WithArgs(testutil.AnyArgs(10)...).
		WillReturnRows(pgxmock.NewRows(educationCols()).
			AddRow(id, inst, deg, field, "2015", (*string)(nil), desc, (*string)(nil), []string{}))

	got, err := svc.Update(context.Background(), "education", id, body)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if e := got.(*model.Education); e.ID != id {
		t.Errorf("ID mismatch")
	}
}

func TestCareerService_UpdateEducation_BadJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	_, err := svc.Update(context.Background(), "education", uuid.New(), json.RawMessage(`{bad}`))
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestCareerService_UpdateEducation_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`UPDATE education SET`).
		WithArgs(testutil.AnyArgs(10)...).
		WillReturnError(errors.New("db"))

	body := json.RawMessage(`{"institution":{"en":"M"},"degree":{"en":"B"},"field":{"en":"C"},"startYear":2020,"description":{}}`)
	if _, err := svc.Update(context.Background(), "education", uuid.New(), body); err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_ListEducation_ScanError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`SELECT .+ FROM education`).
		WillReturnRows(pgxmock.NewRows(educationCols()).
			AddRow(nil, nil, nil, nil, nil, nil, nil, nil, nil))

	_, err := svc.GetAll(context.Background())
	if err == nil {
		t.Fatal("expected scan error")
	}
}

func TestCareerService_CreateCertificate_HappyPath(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	body := json.RawMessage(`{
      "title": {"en": "AWS"},
      "issuer": {"en": "Amazon"},
      "date": "2022-05-01",
      "credentialId": "CRED-123",
      "url": "https://aws.com/cred",
      "imageUrl": "https://x/img.png",
      "sortOrder": 1
    }`)

	id := uuid.New()
	title := i18nJSON(t, model.I18n{"en": "AWS"})
	issuer := i18nJSON(t, model.I18n{"en": "Amazon"})
	cred := "CRED-123"
	url := "https://aws.com/cred"
	img := "https://x/img.png"

	pool.ExpectQuery(`INSERT INTO certificates`).
		WithArgs(testutil.AnyArgs(7)...).
		WillReturnRows(pgxmock.NewRows(certCols()).
			AddRow(id, title, issuer, "2022-05-01", &cred, &url, &img))

	got, err := svc.Create(context.Background(), "certificate", body)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	c := got.(*model.Certificate)
	if c.CredentialID != cred {
		t.Errorf("cred: %q", c.CredentialID)
	}
	if c.Date.Year() != 2022 || c.Date.Month() != time.May {
		t.Errorf("date: %v", c.Date)
	}
}

func TestCareerService_CreateCertificate_BadJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	_, err := svc.Create(context.Background(), "certificate", json.RawMessage(`{bad}`))
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestCareerService_CreateCertificate_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`INSERT INTO certificates`).WithArgs(testutil.AnyArgs(7)...).WillReturnError(errors.New("db"))

	body := json.RawMessage(`{"title":{},"issuer":{},"date":"2022-05-01"}`)
	if _, err := svc.Create(context.Background(), "certificate", body); err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_UpdateCertificate_HappyPath(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	id := uuid.New()
	body := json.RawMessage(`{"title":{"en":"T"},"issuer":{"en":"I"},"date":"2023-01-01"}`)
	title := i18nJSON(t, model.I18n{"en": "T"})
	issuer := i18nJSON(t, model.I18n{"en": "I"})

	pool.ExpectQuery(`UPDATE certificates SET`).
		WithArgs(testutil.AnyArgs(8)...).
		WillReturnRows(pgxmock.NewRows(certCols()).
			AddRow(id, title, issuer, "2023-01-01", (*string)(nil), (*string)(nil), (*string)(nil)))

	if _, err := svc.Update(context.Background(), "certificate", id, body); err != nil {
		t.Fatalf("Update: %v", err)
	}
}

func TestCareerService_UpdateCertificate_BadJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	_, err := svc.Update(context.Background(), "certificate", uuid.New(), json.RawMessage(`{bad}`))
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestCareerService_UpdateCertificate_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`UPDATE certificates SET`).WithArgs(testutil.AnyArgs(8)...).WillReturnError(errors.New("db"))
	body := json.RawMessage(`{"title":{},"issuer":{},"date":"2023-01-01"}`)
	if _, err := svc.Update(context.Background(), "certificate", uuid.New(), body); err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_CertBadDate(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	id := uuid.New()
	title := i18nJSON(t, model.I18n{"en": "T"})
	issuer := i18nJSON(t, model.I18n{"en": "I"})

	pool.ExpectQuery(`INSERT INTO certificates`).
		WithArgs(testutil.AnyArgs(7)...).
		WillReturnRows(pgxmock.NewRows(certCols()).
			AddRow(id, title, issuer, "not-a-date", (*string)(nil), (*string)(nil), (*string)(nil)))

	body := json.RawMessage(`{"title":{"en":"T"},"issuer":{"en":"I"},"date":"not-a-date"}`)
	if _, err := svc.Create(context.Background(), "certificate", body); err == nil {
		t.Fatal("expected date parse error")
	}
}

func TestCareerService_CreatePublication_HappyPath(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	body := json.RawMessage(`{
      "title": {"en": "Paper"},
      "journal": {"en": "Nature"},
      "year": 2024,
      "doi": "10.1/x",
      "url": "https://doi.org/10.1/x",
      "abstract": {"en": "abs"},
      "sortOrder": 0
    }`)

	id := uuid.New()
	title := i18nJSON(t, model.I18n{"en": "Paper"})
	journal := i18nJSON(t, model.I18n{"en": "Nature"})
	abs := i18nJSON(t, model.I18n{"en": "abs"})
	doi := "10.1/x"
	url := "https://doi.org/10.1/x"

	pool.ExpectQuery(`INSERT INTO publications`).
		WithArgs(testutil.AnyArgs(7)...).
		WillReturnRows(pgxmock.NewRows(pubCols()).
			AddRow(id, title, journal, "2024", &doi, &url, abs))

	got, err := svc.Create(context.Background(), "publication", body)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	p := got.(*model.Publication)
	if p.Year != 2024 {
		t.Errorf("year: %d", p.Year)
	}
	if p.DOI != doi {
		t.Errorf("doi: %q", p.DOI)
	}
}

func TestCareerService_CreatePublication_BadJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	_, err := svc.Create(context.Background(), "publication", json.RawMessage(`{bad}`))
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestCareerService_CreatePublication_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`INSERT INTO publications`).WithArgs(testutil.AnyArgs(7)...).WillReturnError(errors.New("db"))
	body := json.RawMessage(`{"title":{},"journal":{},"year":2024,"abstract":{}}`)
	if _, err := svc.Create(context.Background(), "publication", body); err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_UpdatePublication_HappyPath(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	id := uuid.New()
	title := i18nJSON(t, model.I18n{"en": "T"})
	journal := i18nJSON(t, model.I18n{"en": "J"})
	abs := i18nJSON(t, model.I18n{"en": "A"})

	pool.ExpectQuery(`UPDATE publications SET`).
		WithArgs(testutil.AnyArgs(8)...).
		WillReturnRows(pgxmock.NewRows(pubCols()).
			AddRow(id, title, journal, "2024", (*string)(nil), (*string)(nil), abs))

	body := json.RawMessage(`{"title":{"en":"T"},"journal":{"en":"J"},"year":2024,"abstract":{"en":"A"}}`)
	if _, err := svc.Update(context.Background(), "publication", id, body); err != nil {
		t.Fatalf("Update: %v", err)
	}
}

func TestCareerService_UpdatePublication_BadJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	_, err := svc.Update(context.Background(), "publication", uuid.New(), json.RawMessage(`{bad}`))
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestCareerService_UpdatePublication_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`UPDATE publications SET`).WithArgs(testutil.AnyArgs(8)...).WillReturnError(errors.New("db"))
	body := json.RawMessage(`{"title":{},"journal":{},"year":2024,"abstract":{}}`)
	if _, err := svc.Update(context.Background(), "publication", uuid.New(), body); err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_CreateWorkHistory_HappyPath(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	body := json.RawMessage(`{
      "company": {"en": "Acme"},
      "position": {"en": "Eng"},
      "startDate": "2020-01-01",
      "endDate": "2023-06-01",
      "current": false,
      "description": {"en": "desc"},
      "technologies": ["Go", "React"],
      "logoUrl": "https://x/logo.png",
      "sortOrder": 1,
      "achievements": [{"en": "shipped v1"}],
      "fullDescription": {"en": "full"}
    }`)

	id := uuid.New()
	comp := i18nJSON(t, model.I18n{"en": "Acme"})
	pos := i18nJSON(t, model.I18n{"en": "Eng"})
	desc := i18nJSON(t, model.I18n{"en": "desc"})
	ach := i18nJSON(t, []model.I18n{{"en": "shipped v1"}})
	full := i18nJSON(t, model.I18n{"en": "full"})
	endDate := "2023-06-01"
	logo := "https://x/logo.png"

	pool.ExpectQuery(`INSERT INTO work_history`).
		WithArgs(testutil.AnyArgs(11)...).
		WillReturnRows(pgxmock.NewRows(workCols()).
			AddRow(id, comp, pos, "2020-01-01", &endDate, false, desc, []string{"Go", "React"}, &logo, ach, full))

	got, err := svc.Create(context.Background(), "work", body)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	w := got.(*model.WorkHistory)
	if w.StartDate.Year() != 2020 {
		t.Errorf("start year: %d", w.StartDate.Year())
	}
	if w.EndDate == nil || w.EndDate.Year() != 2023 {
		t.Errorf("end date: %v", w.EndDate)
	}
	if len(w.Technologies) != 2 {
		t.Errorf("tech: %v", w.Technologies)
	}
}

func TestCareerService_CreateWorkHistory_NullOptionals(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	body := json.RawMessage(`{
      "company": {"en": "A"},
      "position": {"en": "B"},
      "startDate": "2021-01-01",
      "current": true,
      "description": {"en": "d"},
      "sortOrder": 0
    }`)

	id := uuid.New()
	comp := i18nJSON(t, model.I18n{"en": "A"})
	pos := i18nJSON(t, model.I18n{"en": "B"})
	desc := i18nJSON(t, model.I18n{"en": "d"})

	pool.ExpectQuery(`INSERT INTO work_history`).
		WithArgs(testutil.AnyArgs(11)...).
		WillReturnRows(pgxmock.NewRows(workCols()).
			AddRow(id, comp, pos, "2021-01-01", (*string)(nil), true, desc, []string(nil), (*string)(nil), []byte(`[]`), []byte(`{}`)))

	got, err := svc.Create(context.Background(), "work", body)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	w := got.(*model.WorkHistory)
	if w.EndDate != nil {
		t.Errorf("end date should be nil, got %v", w.EndDate)
	}
	if w.Technologies == nil {
		t.Errorf("technologies should be non-nil")
	}
	if w.Achievements == nil {
		t.Errorf("achievements should be non-nil")
	}
}

func TestCareerService_CreateWorkHistory_BadEndDateIsLogged(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	body := json.RawMessage(`{"company":{"en":"A"},"position":{"en":"B"},"startDate":"2020-01-01","current":false,"description":{"en":"d"},"sortOrder":0}`)

	id := uuid.New()
	comp := i18nJSON(t, model.I18n{"en": "A"})
	pos := i18nJSON(t, model.I18n{"en": "B"})
	desc := i18nJSON(t, model.I18n{"en": "d"})
	badEnd := "garbage-date"

	pool.ExpectQuery(`INSERT INTO work_history`).
		WithArgs(testutil.AnyArgs(11)...).
		WillReturnRows(pgxmock.NewRows(workCols()).
			AddRow(id, comp, pos, "2020-01-01", &badEnd, false, desc, []string{}, (*string)(nil), []byte(`[]`), []byte(`{}`)))

	got, err := svc.Create(context.Background(), "work", body)
	if err != nil {
		t.Fatalf("Create should tolerate bad end date, got err: %v", err)
	}
	w := got.(*model.WorkHistory)
	if w.EndDate != nil {
		t.Errorf("EndDate should stay nil on unparseable, got %v", w.EndDate)
	}
}

func TestCareerService_CreateWorkHistory_BadStartDate(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	body := json.RawMessage(`{"company":{"en":"A"},"position":{"en":"B"},"startDate":"bad","current":false,"description":{"en":"d"},"sortOrder":0}`)

	id := uuid.New()
	comp := i18nJSON(t, model.I18n{"en": "A"})
	pos := i18nJSON(t, model.I18n{"en": "B"})
	desc := i18nJSON(t, model.I18n{"en": "d"})

	pool.ExpectQuery(`INSERT INTO work_history`).
		WithArgs(testutil.AnyArgs(11)...).
		WillReturnRows(pgxmock.NewRows(workCols()).
			AddRow(id, comp, pos, "bad", (*string)(nil), false, desc, []string{}, (*string)(nil), []byte(`[]`), []byte(`{}`)))

	if _, err := svc.Create(context.Background(), "work", body); err == nil {
		t.Fatal("expected parse error on start date")
	}
}

func TestCareerService_CreateWorkHistory_BadJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	_, err := svc.Create(context.Background(), "work", json.RawMessage(`{bad}`))
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestCareerService_CreateWorkHistory_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`INSERT INTO work_history`).WithArgs(testutil.AnyArgs(11)...).WillReturnError(errors.New("db"))
	body := json.RawMessage(`{"company":{},"position":{},"startDate":"2020-01-01","description":{}}`)
	if _, err := svc.Create(context.Background(), "work", body); err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_UpdateWorkHistory_HappyPath(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	id := uuid.New()
	comp := i18nJSON(t, model.I18n{"en": "A"})
	pos := i18nJSON(t, model.I18n{"en": "B"})
	desc := i18nJSON(t, model.I18n{"en": "d"})

	pool.ExpectQuery(`UPDATE work_history SET`).
		WithArgs(testutil.AnyArgs(12)...).
		WillReturnRows(pgxmock.NewRows(workCols()).
			AddRow(id, comp, pos, "2020-01-01", (*string)(nil), true, desc, []string{}, (*string)(nil), []byte(`[]`), []byte(`{}`)))

	body := json.RawMessage(`{"company":{"en":"A"},"position":{"en":"B"},"startDate":"2020-01-01","current":true,"description":{"en":"d"},"sortOrder":0}`)
	if _, err := svc.Update(context.Background(), "work", id, body); err != nil {
		t.Fatalf("Update: %v", err)
	}
}

func TestCareerService_UpdateWorkHistory_BadJSON(t *testing.T) {
	t.Parallel()

	db, _ := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	_, err := svc.Update(context.Background(), "work", uuid.New(), json.RawMessage(`{bad}`))
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestCareerService_UpdateWorkHistory_DBError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`UPDATE work_history SET`).WithArgs(testutil.AnyArgs(12)...).WillReturnError(errors.New("db"))
	body := json.RawMessage(`{"company":{},"position":{},"startDate":"2020-01-01","description":{}}`)
	if _, err := svc.Update(context.Background(), "work", uuid.New(), body); err == nil {
		t.Fatal("expected error")
	}
}

func TestCareerService_ListWorkHistory_ScanError(t *testing.T) {
	t.Parallel()

	db, pool := testutil.NewMockDB(t)
	svc := NewCareerService(db)

	pool.ExpectQuery(`SELECT .+ FROM education`).
		WillReturnRows(pgxmock.NewRows(educationCols()))
	pool.ExpectQuery(`SELECT .+ FROM work_history`).
		WillReturnRows(pgxmock.NewRows(workCols()).
			AddRow(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil))

	_, err := svc.GetAll(context.Background())
	if err == nil {
		t.Fatal("expected scan error")
	}
}
