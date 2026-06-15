package api

import (
	"errors"
	"log/slog"
	"net/http"

	db "github.com/eddie-englund/project-planner/backend/internal/db/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProjectHandler struct {
	db     *db.Queries
	logger *slog.Logger
}

type createProjectRequest struct {
	Title    string  `json:"title"     validate:"required"`
	Color    string  `json:"color"     validate:"required"`
	ImageURL *string `json:"image_url" validate:"omitempty"`
}

type updateProjectRequest struct {
	Title    string  `json:"title"     validate:"required"`
	Color    string  `json:"color"     validate:"required"`
	ImageURL *string `json:"image_url" validate:"omitempty"`
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := decodeAndValidate[createProjectRequest](r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	params := db.CreateProjectParams{
		Title:    req.Title,
		Color:    req.Color,
		ImageUrl: optText(req.ImageURL),
	}

	project, err := h.db.CreateProject(r.Context(), params)
	if err != nil {
		h.logger.Error("create project", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	projects, err := h.db.ListProjects(r.Context())
	if err != nil {
		h.logger.Error("list projects", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if projects == nil {
		projects = []db.Project{}
	}
	writeJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	project, err := h.db.GetProjectByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.logger.Error("get project", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	req, err := decodeAndValidate[updateProjectRequest](r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	params := db.UpdateProjectParams{
		ID:       id,
		Title:    req.Title,
		Color:    req.Color,
		ImageUrl: optText(req.ImageURL),
	}

	project, err := h.db.UpdateProject(r.Context(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.logger.Error("update project", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteProject(r.Context(), id); err != nil {
		h.logger.Error("delete project", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
