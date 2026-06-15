-- +goose Up
CREATE TABLE IF NOT EXISTS projects(
  id UUID NOT NULL,
  title VARCHAR(128) NOT NULL,
  created_by UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  color VARCHAR(32) NOT NULL,
  image_url text,
  PRIMARY KEY(id),
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS project_topics(
  id UUID NOT NULL,
  project_id UUID NOT NULL,
  index INT NOT NULL,
  title VARCHAR(128),
  created_at TIMESTAMPZ NOT NULL DEFAULT now(),
  image_url text,
  color VARCHAR(32) NOT NULL, 
  PRIMARY KEY(id),
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
