package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	db "github.com/eddie-englund/project-planner/backend/internal/db/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTicketStore struct {
	ticket         db.TopicTicket
	tickets        []db.TopicTicket
	projectTickets []db.ListTicketsByProjectRow
	err            error
}

func (m *mockTicketStore) CreateTicket(_ context.Context, _ db.CreateTicketParams) (db.TopicTicket, error) {
	return m.ticket, m.err
}
func (m *mockTicketStore) ListTicketsByTopic(_ context.Context, _ pgtype.UUID) ([]db.TopicTicket, error) {
	return m.tickets, m.err
}
func (m *mockTicketStore) ListTicketsByProject(_ context.Context, _ pgtype.UUID) ([]db.ListTicketsByProjectRow, error) {
	return m.projectTickets, m.err
}
func (m *mockTicketStore) GetTicketByID(_ context.Context, _ db.GetTicketByIDParams) (db.TopicTicket, error) {
	return m.ticket, m.err
}
func (m *mockTicketStore) UpdateTicket(_ context.Context, _ db.UpdateTicketParams) (db.TopicTicket, error) {
	return m.ticket, m.err
}
func (m *mockTicketStore) DeleteTicket(_ context.Context, _ pgtype.UUID) error {
	return m.err
}

func newTestTicketHandler(store ticketStore) *TicketHandler {
	return &TicketHandler{db: store, logger: slog.Default()}
}

const validTicketID = "00000000-0000-0000-0000-000000000003"

func ticketReq(method, url, body string, pathValues map[string]string) *http.Request {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, url, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	for k, v := range pathValues {
		req.SetPathValue(k, v)
	}
	return req
}

func TestTicketHandler_List_ReturnsTickets(t *testing.T) {
	mock := &mockTicketStore{tickets: []db.TopicTicket{
		{Title: "Fix bug", Urls: []string{}},
		{Title: "Add feature", Urls: []string{}},
	}}
	h := newTestTicketHandler(mock)

	req := ticketReq(http.MethodGet, "/", "", map[string]string{"topicId": validTopicID})
	w := httptest.NewRecorder()
	h.List(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Fix bug")
	assert.Contains(t, w.Body.String(), "Add feature")
}

func TestTicketHandler_List_DBError_Returns500(t *testing.T) {
	mock := &mockTicketStore{err: errors.New("db error")}
	h := newTestTicketHandler(mock)

	req := ticketReq(http.MethodGet, "/", "", map[string]string{"topicId": validTopicID})
	w := httptest.NewRecorder()
	h.List(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTicketHandler_Create_ValidBody_Returns201(t *testing.T) {
	mock := &mockTicketStore{ticket: db.TopicTicket{Title: "New ticket", Urls: []string{}}}
	h := newTestTicketHandler(mock)

	req := ticketReq(http.MethodPost, "/", `{"title":"New ticket"}`, map[string]string{"topicId": validTopicID})
	w := httptest.NewRecorder()
	h.Create(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "New ticket")
}

func TestTicketHandler_Create_MissingTitle_Returns400(t *testing.T) {
	mock := &mockTicketStore{}
	h := newTestTicketHandler(mock)

	req := ticketReq(http.MethodPost, "/", `{"body":"no title"}`, map[string]string{"topicId": validTopicID})
	w := httptest.NewRecorder()
	h.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTicketHandler_GetByID_Found_Returns200(t *testing.T) {
	mock := &mockTicketStore{ticket: db.TopicTicket{Title: "Found ticket", Urls: []string{}}}
	h := newTestTicketHandler(mock)

	req := ticketReq(http.MethodGet, "/", "", map[string]string{
		"topicId": validTopicID,
		"id":      validTicketID,
	})
	w := httptest.NewRecorder()
	h.GetByID(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Found ticket")
}

func TestTicketHandler_GetByID_NotFound_Returns404(t *testing.T) {
	mock := &mockTicketStore{err: pgx.ErrNoRows}
	h := newTestTicketHandler(mock)

	req := ticketReq(http.MethodGet, "/", "", map[string]string{
		"topicId": validTopicID,
		"id":      validTicketID,
	})
	w := httptest.NewRecorder()
	h.GetByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTicketHandler_Update_Valid_Returns200(t *testing.T) {
	mock := &mockTicketStore{ticket: db.TopicTicket{Title: "Updated", Urls: []string{}}}
	h := newTestTicketHandler(mock)

	req := ticketReq(http.MethodPut, "/", `{"title":"Updated"}`, map[string]string{"id": validTicketID})
	w := httptest.NewRecorder()
	h.Update(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated")
}

func TestTicketHandler_Update_NotFound_Returns404(t *testing.T) {
	mock := &mockTicketStore{err: pgx.ErrNoRows}
	h := newTestTicketHandler(mock)

	req := ticketReq(http.MethodPut, "/", `{"title":"X"}`, map[string]string{"id": validTicketID})
	w := httptest.NewRecorder()
	h.Update(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTicketHandler_Delete_Success_Returns204(t *testing.T) {
	mock := &mockTicketStore{}
	h := newTestTicketHandler(mock)

	req := ticketReq(http.MethodDelete, "/", "", map[string]string{"id": validTicketID})
	w := httptest.NewRecorder()
	h.Delete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestTicketHandler_ListByProject_ReturnsTickets(t *testing.T) {
	mock := &mockTicketStore{projectTickets: []db.ListTicketsByProjectRow{
		{Title: "Fix bug", TopicColor: "#10b981", TopicTitle: "Backend", Urls: []string{}},
		{Title: "Design UI", TopicColor: "#3b82f6", TopicTitle: "Frontend", Urls: []string{}},
	}}
	h := newTestTicketHandler(mock)

	req := ticketReq(http.MethodGet, "/", "", map[string]string{"projectId": validProjectID})
	w := httptest.NewRecorder()
	h.ListByProject(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Fix bug")
	assert.Contains(t, w.Body.String(), "#10b981")
	assert.Contains(t, w.Body.String(), "Backend")
}

func TestTicketHandler_ListByProject_DBError_Returns500(t *testing.T) {
	mock := &mockTicketStore{err: errors.New("db error")}
	h := newTestTicketHandler(mock)

	req := ticketReq(http.MethodGet, "/", "", map[string]string{"projectId": validProjectID})
	w := httptest.NewRecorder()
	h.ListByProject(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTicketHandler_ListByProject_InvalidProjectID_Returns400(t *testing.T) {
	mock := &mockTicketStore{}
	h := newTestTicketHandler(mock)

	req := ticketReq(http.MethodGet, "/", "", map[string]string{"projectId": "not-a-uuid"})
	w := httptest.NewRecorder()
	h.ListByProject(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
