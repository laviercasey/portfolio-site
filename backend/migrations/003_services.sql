-- +goose Up
-- +goose StatementBegin
CREATE TABLE services (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug          TEXT UNIQUE NOT NULL,
    num           TEXT NOT NULL,
    icon_key      TEXT NOT NULL,
    visual_key    TEXT NOT NULL,
    accent        TEXT NOT NULL,
    title         JSONB NOT NULL,
    lead          JSONB NOT NULL,
    bullets       JSONB NOT NULL DEFAULT '{"ru":[],"en":[]}'::jsonb,
    stack         TEXT NOT NULL DEFAULT '',
    timeline      JSONB NOT NULL DEFAULT '{}'::jsonb,
    price_ru      TEXT NOT NULL DEFAULT '',
    price_en      TEXT NOT NULL DEFAULT '',
    case_projects JSONB NOT NULL DEFAULT '[]'::jsonb,
    sort_order    INTEGER NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_services_sort_order ON services (sort_order);

CREATE TRIGGER trg_services_updated_at
    BEFORE UPDATE ON services
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE service_faqs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question    JSONB NOT NULL,
    answer      JSONB NOT NULL,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_service_faqs_sort_order ON service_faqs (sort_order);

CREATE TRIGGER trg_service_faqs_updated_at
    BEFORE UPDATE ON service_faqs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE service_process_steps (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    num          TEXT NOT NULL,
    title        JSONB NOT NULL,
    description  JSONB NOT NULL,
    sort_order   INTEGER NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_service_process_steps_sort_order ON service_process_steps (sort_order);

CREATE TRIGGER trg_service_process_steps_updated_at
    BEFORE UPDATE ON service_process_steps
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS service_process_steps;
DROP TABLE IF EXISTS service_faqs;
DROP TABLE IF EXISTS services;
