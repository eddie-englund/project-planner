package api

import (
	"context"

	db "github.com/eddie-englund/project-planner/backend/internal/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type projectStore interface {
	CreateProject(ctx context.Context, arg db.CreateProjectParams) (db.Project, error)
	ListProjects(ctx context.Context) ([]db.Project, error)
	GetProjectByID(ctx context.Context, id pgtype.UUID) (db.Project, error)
	UpdateProject(ctx context.Context, arg db.UpdateProjectParams) (db.Project, error)
	DeleteProject(ctx context.Context, id pgtype.UUID) error
	SeedDefaultStatuses(ctx context.Context, projectID pgtype.UUID) error
}

type topicStore interface {
	CreateProjectTopic(ctx context.Context, arg db.CreateProjectTopicParams) (db.ProjectTopic, error)
	ListProjectTopics(ctx context.Context, projectID pgtype.UUID) ([]db.ProjectTopic, error)
	GetProjectTopicByID(ctx context.Context, arg db.GetProjectTopicByIDParams) (db.ProjectTopic, error)
	UpdateProjectTopic(ctx context.Context, arg db.UpdateProjectTopicParams) (db.ProjectTopic, error)
	DeleteProjectTopic(ctx context.Context, id pgtype.UUID) error
}

type statusStore interface {
	ListProjectStatuses(ctx context.Context, projectID pgtype.UUID) ([]db.ProjectStatus, error)
	SeedDefaultStatuses(ctx context.Context, projectID pgtype.UUID) error
	CreateStatus(ctx context.Context, arg db.CreateStatusParams) (db.ProjectStatus, error)
}

type ticketStore interface {
	CreateTicket(ctx context.Context, arg db.CreateTicketParams) (db.TopicTicket, error)
	ListTicketsByTopic(ctx context.Context, topicID pgtype.UUID) ([]db.TopicTicket, error)
	GetTicketByID(ctx context.Context, arg db.GetTicketByIDParams) (db.TopicTicket, error)
	UpdateTicket(ctx context.Context, arg db.UpdateTicketParams) (db.TopicTicket, error)
	DeleteTicket(ctx context.Context, id pgtype.UUID) error
}
