-- +goose Up
ALTER TABLE projects ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE project_topics ALTER COLUMN created_at TYPE TIMESTAMPTZ
  USING created_at::TIMESTAMPTZ;
