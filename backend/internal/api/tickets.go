package api

import (
	"errors"
	"log/slog"
	"net/http"

	db "github.com/eddie-englund/project-planner/backend/internal/db/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type TicketHandler struct {
	db     ticketStore
	logger *slog.Logger
}

type createTicketRequest struct {
	Title    string   `json:"title"    validate:"required"`
	Body     string   `json:"body"`
	URLs     []string `json:"urls"`
	StatusID *string  `json:"statusId"`
}

type updateTicketRequest struct {
	Title    string   `json:"title"    validate:"required"`
	Body     string   `json:"body"`
	URLs     []string `json:"urls"`
	StatusID *string  `json:"statusId"`
}

func parseTopicID(r *http.Request) (pgtype.UUID, error) {
	var id pgtype.UUID
	return id, id.Scan(r.PathValue("topicId"))
}

func optUUID(s *string) pgtype.UUID {
	if s == nil {
		return pgtype.UUID{}
	}
	var u pgtype.UUID
	_ = u.Scan(*s)
	return u
}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	topicID, err := parseTopicID(r)
	if err != nil {
		http.Error(w, "invalid topic id", http.StatusBadRequest)
		return
	}

	req, err := decodeAndValidate[createTicketRequest](r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	urls := req.URLs
	if urls == nil {
		urls = []string{}
	}

	ticket, err := h.db.CreateTicket(r.Context(), db.CreateTicketParams{
		TopicID:  topicID,
		StatusID: optUUID(req.StatusID),
		Title:    req.Title,
		Body:     req.Body,
		Urls:     urls,
	})
	if err != nil {
		h.logger.Error("create ticket", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, ticketToResponse(ticket))
}

func (h *TicketHandler) List(w http.ResponseWriter, r *http.Request) {
	topicID, err := parseTopicID(r)
	if err != nil {
		http.Error(w, "invalid topic id", http.StatusBadRequest)
		return
	}

	tickets, err := h.db.ListTicketsByTopic(r.Context(), topicID)
	if err != nil {
		h.logger.Error("list tickets", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := make([]TicketResponse, len(tickets))
	for i, t := range tickets {
		resp[i] = ticketToResponse(t)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *TicketHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	topicID, err := parseTopicID(r)
	if err != nil {
		http.Error(w, "invalid topic id", http.StatusBadRequest)
		return
	}

	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ticket, err := h.db.GetTicketByID(r.Context(), db.GetTicketByIDParams{ID: id, TopicID: topicID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.logger.Error("get ticket", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, ticketToResponse(ticket))
}

func (h *TicketHandler) Update(w http.ResponseWriter, r *http.Request) {
	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	req, err := decodeAndValidate[updateTicketRequest](r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	urls := req.URLs
	if urls == nil {
		urls = []string{}
	}

	ticket, err := h.db.UpdateTicket(r.Context(), db.UpdateTicketParams{
		ID:       id,
		Title:    req.Title,
		Body:     req.Body,
		Urls:     urls,
		StatusID: optUUID(req.StatusID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.logger.Error("update ticket", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, ticketToResponse(ticket))
}

func (h *TicketHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteTicket(r.Context(), id); err != nil {
		h.logger.Error("delete ticket", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
