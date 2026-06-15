package api

import (
	"errors"
	"log/slog"
	"net/http"

	db "github.com/eddie-englund/project-planner/backend/internal/db/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type TopicHandler struct {
	db     *db.Queries
	logger *slog.Logger
}

type createTopicRequest struct {
	Index    int32   `json:"index"`
	Title    *string `json:"title"     validate:"omitempty"`
	Color    string  `json:"color"     validate:"required"`
	ImageURL *string `json:"image_url" validate:"omitempty"`
}

type updateTopicRequest struct {
	Index    int32   `json:"index"`
	Title    *string `json:"title"     validate:"omitempty"`
	Color    string  `json:"color"     validate:"required"`
	ImageURL *string `json:"image_url" validate:"omitempty"`
}

func parseProjectID(r *http.Request) (pgtype.UUID, error) {
	var id pgtype.UUID
	return id, id.Scan(r.PathValue("projectId"))
}

func (h *TopicHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseProjectID(r)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	req, err := decodeAndValidate[createTopicRequest](r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	params := db.CreateProjectTopicParams{
		ProjectID: projectID,
		Index:     req.Index,
		Color:     req.Color,
		Title:     optText(req.Title),
		ImageUrl:  optText(req.ImageURL),
	}

	topic, err := h.db.CreateProjectTopic(r.Context(), params)
	if err != nil {
		h.logger.Error("create topic", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, topic)
}

func (h *TopicHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseProjectID(r)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	topics, err := h.db.ListProjectTopics(r.Context(), projectID)
	if err != nil {
		h.logger.Error("list topics", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if topics == nil {
		topics = []db.ProjectTopic{}
	}
	writeJSON(w, http.StatusOK, topics)
}

func (h *TopicHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseProjectID(r)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	topic, err := h.db.GetProjectTopicByID(r.Context(), db.GetProjectTopicByIDParams{
		ID:        id,
		ProjectID: projectID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.logger.Error("get topic", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, topic)
}

func (h *TopicHandler) Update(w http.ResponseWriter, r *http.Request) {
	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	req, err := decodeAndValidate[updateTopicRequest](r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	params := db.UpdateProjectTopicParams{
		ID:       id,
		Color:    req.Color,
		Index:    req.Index,
		Title:    optText(req.Title),
		ImageUrl: optText(req.ImageURL),
	}

	topic, err := h.db.UpdateProjectTopic(r.Context(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.logger.Error("update topic", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, topic)
}

func (h *TopicHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteProjectTopic(r.Context(), id); err != nil {
		h.logger.Error("delete topic", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
