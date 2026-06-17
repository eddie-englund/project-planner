-- +goose Up
CREATE TABLE project_statuses (
  id         UUID         NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  project_id UUID         NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  name       VARCHAR(64)  NOT NULL,
  color      VARCHAR(32)  NOT NULL DEFAULT '#6B7C8E',
  position   INT          NOT NULL DEFAULT 0,
  is_terminal BOOLEAN     NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE topic_tickets (
  id         UUID         NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  topic_id   UUID         NOT NULL REFERENCES project_topics(id) ON DELETE CASCADE,
  status_id  UUID         REFERENCES project_statuses(id) ON DELETE SET NULL,
  title      VARCHAR(256) NOT NULL,
  body       TEXT         NOT NULL DEFAULT '',
  urls       TEXT[]       NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);
