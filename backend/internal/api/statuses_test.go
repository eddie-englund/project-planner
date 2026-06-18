package api

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	db "github.com/eddie-englund/project-planner/backend/internal/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStatusStore struct {
	statuses    []db.ProjectStatus
	status      db.ProjectStatus
	err         error
	seedErr     error
	// seededStatuses are returned after SeedDefaultStatuses is called
	seededStatuses []db.ProjectStatus
	seedCalled     bool
}

func (m *mockStatusStore) ListProjectStatuses(_ context.Context, _ pgtype.UUID) ([]db.ProjectStatus, error) {
	if m.seedCalled && m.seededStatuses != nil {
		return m.seededStatuses, m.err
	}
	return m.statuses, m.err
}

func (m *mockStatusStore) SeedDefaultStatuses(_ context.Context, _ pgtype.UUID) error {
	m.seedCalled = true
	return m.seedErr
}

func (m *mockStatusStore) CreateStatus(_ context.Context, _ db.CreateStatusParams) (db.ProjectStatus, error) {
	return m.status, m.err
}

func newTestStatusHandler(store statusStore) *StatusHandler {
	return &StatusHandler{db: store, logger: slog.Default()}
}

func TestStatusHandler_List_ReturnsStatuses(t *testing.T) {
	mock := &mockStatusStore{statuses: []db.ProjectStatus{
		{Name: "Open", Color: "#6B8E7C", Position: 0, IsTerminal: false},
		{Name: "Closed", Color: "#8E7C6B", Position: 1, IsTerminal: true},
	}}
	h := newTestStatusHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/projects/00000000-0000-0000-0000-000000000001/statuses", nil)
	req.SetPathValue("projectId", "00000000-0000-0000-0000-000000000001")
	w := httptest.NewRecorder()
	h.List(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Open")
	assert.Contains(t, w.Body.String(), "Closed")
}

func TestStatusHandler_List_EmptyAutoSeeds(t *testing.T) {
	mock := &mockStatusStore{
		statuses: []db.ProjectStatus{},
		seededStatuses: []db.ProjectStatus{
			{Name: "Open", Color: "#6B8E7C", Position: 0, IsTerminal: false},
			{Name: "Closed", Color: "#8E7C6B", Position: 1, IsTerminal: true},
		},
	}
	h := newTestStatusHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/projects/00000000-0000-0000-0000-000000000001/statuses", nil)
	req.SetPathValue("projectId", "00000000-0000-0000-0000-000000000001")
	w := httptest.NewRecorder()
	h.List(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.True(t, mock.seedCalled)
	assert.Contains(t, w.Body.String(), "Open")
}

func TestStatusHandler_List_InvalidProjectID_Returns400(t *testing.T) {
	mock := &mockStatusStore{}
	h := newTestStatusHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/projects/bad/statuses", nil)
	req.SetPathValue("projectId", "bad")
	w := httptest.NewRecorder()
	h.List(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStatusHandler_Create_ValidBody_Returns201(t *testing.T) {
	mock := &mockStatusStore{status: db.ProjectStatus{Name: "In Progress", Color: "#888888"}}
	h := newTestStatusHandler(mock)

	body := `{"name":"In Progress","color":"#888888","position":2}`
	req := httptest.NewRequest(http.MethodPost, "/projects/00000000-0000-0000-0000-000000000001/statuses", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("projectId", "00000000-0000-0000-0000-000000000001")
	w := httptest.NewRecorder()
	h.Create(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "In Progress")
}

func TestStatusHandler_Create_MissingName_Returns400(t *testing.T) {
	mock := &mockStatusStore{}
	h := newTestStatusHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/projects/00000000-0000-0000-0000-000000000001/statuses", strings.NewReader(`{"color":"#888888"}`))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("projectId", "00000000-0000-0000-0000-000000000001")
	w := httptest.NewRecorder()
	h.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
