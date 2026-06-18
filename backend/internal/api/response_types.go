package api

import (
	"fmt"
	"time"

	db "github.com/eddie-englund/project-planner/backend/internal/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProjectResponse struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Color     string  `json:"color"`
	ImageUrl  *string `json:"imageUrl"`
	CreatedBy string  `json:"createdBy"`
	CreatedAt string  `json:"createdAt"`
}

func uuidToString(u pgtype.UUID) string {
	b := u.Bytes
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

type TopicResponse struct {
	ID        string  `json:"id"`
	ProjectID string  `json:"projectId"`
	Index     int32   `json:"index"`
	Title     string  `json:"title"`
	Color     string  `json:"color"`
	ImageUrl  *string `json:"imageUrl"`
	CreatedAt string  `json:"createdAt"`
}

func topicToResponse(t db.ProjectTopic) TopicResponse {
	var imageUrl *string
	if t.ImageUrl.Valid {
		imageUrl = &t.ImageUrl.String
	}
	return TopicResponse{
		ID:        uuidToString(t.ID),
		ProjectID: uuidToString(t.ProjectID),
		Index:     t.Index,
		Title:     t.Title,
		Color:     t.Color,
		ImageUrl:  imageUrl,
		CreatedAt: t.CreatedAt.Time.Format(time.RFC3339),
	}
}

type StatusResponse struct {
	ID         string `json:"id"`
	ProjectID  string `json:"projectId"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	Position   int32  `json:"position"`
	IsTerminal bool   `json:"isTerminal"`
}

func statusToResponse(s db.ProjectStatus) StatusResponse {
	return StatusResponse{
		ID:         uuidToString(s.ID),
		ProjectID:  uuidToString(s.ProjectID),
		Name:       s.Name,
		Color:      s.Color,
		Position:   s.Position,
		IsTerminal: s.IsTerminal,
	}
}

type TicketResponse struct {
	ID        string   `json:"id"`
	TopicID   string   `json:"topicId"`
	StatusID  *string  `json:"statusId"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	URLs      []string `json:"urls"`
	CreatedAt string   `json:"createdAt"`
}

func ticketToResponse(t db.TopicTicket) TicketResponse {
	var statusID *string
	if t.StatusID.Valid {
		s := uuidToString(t.StatusID)
		statusID = &s
	}
	urls := t.Urls
	if urls == nil {
		urls = []string{}
	}
	return TicketResponse{
		ID:        uuidToString(t.ID),
		TopicID:   uuidToString(t.TopicID),
		StatusID:  statusID,
		Title:     t.Title,
		Body:      t.Body,
		URLs:      urls,
		CreatedAt: t.CreatedAt.Time.Format(time.RFC3339),
	}
}

type TicketWithTopicResponse struct {
	ID         string   `json:"id"`
	TopicID    string   `json:"topicId"`
	TopicColor string   `json:"topicColor"`
	TopicTitle string   `json:"topicTitle"`
	StatusID   *string  `json:"statusId"`
	Title      string   `json:"title"`
	Body       string   `json:"body"`
	URLs       []string `json:"urls"`
	CreatedAt  string   `json:"createdAt"`
}

func ticketWithTopicToResponse(t db.ListTicketsByProjectRow) TicketWithTopicResponse {
	var statusID *string
	if t.StatusID.Valid {
		s := uuidToString(t.StatusID)
		statusID = &s
	}
	urls := t.Urls
	if urls == nil {
		urls = []string{}
	}
	return TicketWithTopicResponse{
		ID:         uuidToString(t.ID),
		TopicID:    uuidToString(t.TopicID),
		TopicColor: t.TopicColor,
		TopicTitle: t.TopicTitle,
		StatusID:   statusID,
		Title:      t.Title,
		Body:       t.Body,
		URLs:       urls,
		CreatedAt:  t.CreatedAt.Time.Format(time.RFC3339),
	}
}

func projectToResponse(p db.Project) ProjectResponse {
	var imageUrl *string
	if p.ImageUrl.Valid {
		imageUrl = &p.ImageUrl.String
	}
	return ProjectResponse{
		ID:        uuidToString(p.ID),
		Title:     p.Title,
		Color:     p.Color,
		ImageUrl:  imageUrl,
		CreatedBy: uuidToString(p.CreatedBy),
		CreatedAt: p.CreatedAt.Time.Format(time.RFC3339),
	}
}
