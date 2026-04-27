-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION update_updated_at() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

CREATE TABLE projects (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug              TEXT UNIQUE NOT NULL,
    title             JSONB NOT NULL,
    short_desc        JSONB NOT NULL DEFAULT '{}'::jsonb,
    description       JSONB NOT NULL DEFAULT '{}'::jsonb,
    category          TEXT NOT NULL CHECK (category IN ('web','mobile','data','research','other')),
    status            TEXT NOT NULL CHECK (status IN ('completed','in_development','open_source')),
    tags              TEXT[] NOT NULL DEFAULT '{}',
    tech_stack        TEXT[] NOT NULL DEFAULT '{}',
    goal_desc         JSONB NOT NULL DEFAULT '{}'::jsonb,
    github_url        TEXT NOT NULL DEFAULT '',
    demo_url          TEXT NOT NULL DEFAULT '',
    site_url          TEXT NOT NULL DEFAULT '',
    video_url         TEXT NOT NULL DEFAULT '',
    thumbnail         TEXT NOT NULL DEFAULT '',
    images            TEXT[] NOT NULL DEFAULT '{}',
    stars             INTEGER NOT NULL DEFAULT 0,
    featured          BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order        INTEGER NOT NULL DEFAULT 0,
    problem           JSONB NOT NULL DEFAULT '{}'::jsonb,
    approach          JSONB NOT NULL DEFAULT '{}'::jsonb,
    outcome           JSONB NOT NULL DEFAULT '{}'::jsonb,
    tech_choices      JSONB NOT NULL DEFAULT '[]'::jsonb,
    highlights        JSONB NOT NULL DEFAULT '[]'::jsonb,
    timeline_started  DATE,
    timeline_shipped  DATE,
    demo_credentials  JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_projects_category    ON projects (category);
CREATE INDEX idx_projects_featured    ON projects (featured);
CREATE INDEX idx_projects_sort_order  ON projects (sort_order);
CREATE INDEX idx_projects_title_gin   ON projects USING GIN (title);

CREATE TRIGGER trg_projects_updated_at
    BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE education (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    institution            JSONB NOT NULL,
    degree                 JSONB NOT NULL,
    field                  JSONB NOT NULL,
    start_year             TEXT NOT NULL,
    end_year               TEXT,
    description            JSONB NOT NULL DEFAULT '{}'::jsonb,
    logo_url               TEXT NOT NULL DEFAULT '',
    related_project_slugs  TEXT[] NOT NULL DEFAULT '{}',
    sort_order             INTEGER NOT NULL DEFAULT 0,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_education_sort_order ON education (sort_order);

CREATE TRIGGER trg_education_updated_at
    BEFORE UPDATE ON education
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE work_history (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company           JSONB NOT NULL,
    position          JSONB NOT NULL,
    start_date        TEXT NOT NULL,
    end_date          TEXT,
    is_current        BOOLEAN NOT NULL DEFAULT FALSE,
    description       JSONB NOT NULL DEFAULT '{}'::jsonb,
    full_description  JSONB NOT NULL DEFAULT '{}'::jsonb,
    technologies      TEXT[] NOT NULL DEFAULT '{}',
    achievements      JSONB NOT NULL DEFAULT '[]'::jsonb,
    logo_url          TEXT NOT NULL DEFAULT '',
    sort_order        INTEGER NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_work_history_is_current ON work_history (is_current);
CREATE INDEX idx_work_history_sort_order ON work_history (sort_order);

CREATE TRIGGER trg_work_history_updated_at
    BEFORE UPDATE ON work_history
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE certificates (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title          JSONB NOT NULL,
    issuer         JSONB NOT NULL,
    issue_date     TEXT NOT NULL,
    credential_id  TEXT,
    url            TEXT,
    image_url      TEXT NOT NULL DEFAULT '',
    sort_order     INTEGER NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_certificates_sort_order ON certificates (sort_order);

CREATE TRIGGER trg_certificates_updated_at
    BEFORE UPDATE ON certificates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE publications (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       JSONB NOT NULL,
    journal     JSONB NOT NULL,
    year        TEXT NOT NULL,
    doi         TEXT,
    url         TEXT NOT NULL DEFAULT '',
    abstract    JSONB NOT NULL DEFAULT '{}'::jsonb,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_publications_sort_order ON publications (sort_order);
CREATE INDEX idx_publications_year       ON publications (year);

CREATE TRIGGER trg_publications_updated_at
    BEFORE UPDATE ON publications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE content (
    section     TEXT PRIMARY KEY,
    data        JSONB NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_content_updated_at
    BEFORE UPDATE ON content
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE inquiries (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT NOT NULL,
    email         TEXT NOT NULL,
    telegram      TEXT,
    company       TEXT,
    inquiry_type  TEXT NOT NULL CHECK (inquiry_type IN ('freelance','fulltime','collaboration','other')),
    budget        TEXT,
    message       TEXT NOT NULL,
    status        TEXT NOT NULL DEFAULT 'new' CHECK (status IN ('new','read','replied','archived')),
    admin_notes   TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_inquiries_created_at ON inquiries (created_at DESC);
CREATE INDEX idx_inquiries_email      ON inquiries (email);
CREATE INDEX idx_inquiries_status     ON inquiries (status);

CREATE TRIGGER trg_inquiries_updated_at
    BEFORE UPDATE ON inquiries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE media (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    filename       TEXT NOT NULL,
    original_name  TEXT NOT NULL,
    mime_type      TEXT NOT NULL,
    size           BIGINT NOT NULL,
    url            TEXT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS media;
DROP TABLE IF EXISTS inquiries;
DROP TABLE IF EXISTS content;
DROP TABLE IF EXISTS publications;
DROP TABLE IF EXISTS certificates;
DROP TABLE IF EXISTS work_history;
DROP TABLE IF EXISTS education;
DROP TABLE IF EXISTS projects;
DROP FUNCTION IF EXISTS update_updated_at();
