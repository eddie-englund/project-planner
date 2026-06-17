package api

import (
	"log/slog"
	"net/http"

	db "github.com/eddie-englund/project-planner/backend/internal/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type StatusHandler struct {
	db     *db.Queries
	logger *slog.Logger
}

type createStatusRequest struct {
	Name       string `json:"name"       validate:"required"`
	Color      string `json:"color"      validate:"required"`
	Position   int32  `json:"position"`
	IsTerminal bool   `json:"isTerminal"`
}

func (h *StatusHandler) List(w http.ResponseWriter, r *http.Request) {
	var projectID pgtype.UUID
	if err := projectID.Scan(r.PathValue("projectId")); err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	statuses, err := h.db.ListProjectStatuses(r.Context(), projectID)
	if err != nil {
		h.logger.Error("list statuses", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := make([]StatusResponse, len(statuses))
	for i, s := range statuses {
		resp[i] = statusToResponse(s)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *StatusHandler) Create(w http.ResponseWriter, r *http.Request) {
	var projectID pgtype.UUID
	if err := projectID.Scan(r.PathValue("projectId")); err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	req, err := decodeAndValidate[createStatusRequest](r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	status, err := h.db.CreateStatus(r.Context(), db.CreateStatusParams{
		ProjectID:  projectID,
		Name:       req.Name,
		Color:      req.Color,
		Position:   req.Position,
		IsTerminal: req.IsTerminal,
	})
	if err != nil {
		h.logger.Error("create status", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, statusToResponse(status))
}
