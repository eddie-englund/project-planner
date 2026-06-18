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

type mockProjectStore struct {
	project  db.Project
	projects []db.Project
	err      error
}

func (m *mockProjectStore) CreateProject(_ context.Context, _ db.CreateProjectParams) (db.Project, error) {
	return m.project, m.err
}
func (m *mockProjectStore) ListProjects(_ context.Context) ([]db.Project, error) {
	return m.projects, m.err
}
func (m *mockProjectStore) GetProjectByID(_ context.Context, _ pgtype.UUID) (db.Project, error) {
	return m.project, m.err
}
func (m *mockProjectStore) UpdateProject(_ context.Context, _ db.UpdateProjectParams) (db.Project, error) {
	return m.project, m.err
}
func (m *mockProjectStore) DeleteProject(_ context.Context, _ pgtype.UUID) error {
	return m.err
}
func (m *mockProjectStore) SeedDefaultStatuses(_ context.Context, _ pgtype.UUID) error {
	return m.err
}

func newTestProjectHandler(store projectStore) *ProjectHandler {
	return &ProjectHandler{db: store, logger: slog.Default()}
}

func TestProjectHandler_List_ReturnsProjects(t *testing.T) {
	mock := &mockProjectStore{projects: []db.Project{
		{Title: "Alpha", Color: "#aaaaaa"},
		{Title: "Beta", Color: "#bbbbbb"},
	}}
	h := newTestProjectHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Alpha")
	assert.Contains(t, w.Body.String(), "Beta")
}

func TestProjectHandler_List_DBError_Returns500(t *testing.T) {
	mock := &mockProjectStore{err: errors.New("db down")}
	h := newTestProjectHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProjectHandler_Create_ValidBody_Returns201(t *testing.T) {
	mock := &mockProjectStore{project: db.Project{Title: "New", Color: "#ff0000"}}
	h := newTestProjectHandler(mock)

	body := `{"title":"New","color":"#ff0000"}`
	req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "New")
}

func TestProjectHandler_Create_MissingTitle_Returns400(t *testing.T) {
	mock := &mockProjectStore{}
	h := newTestProjectHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(`{"color":"#ff0000"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProjectHandler_GetByID_Found_Returns200(t *testing.T) {
	mock := &mockProjectStore{project: db.Project{Title: "Found", Color: "#cccccc"}}
	h := newTestProjectHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/projects/00000000-0000-0000-0000-000000000001", nil)
	req.SetPathValue("id", "00000000-0000-0000-0000-000000000001")
	w := httptest.NewRecorder()
	h.GetByID(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Found")
}

func TestProjectHandler_GetByID_NotFound_Returns404(t *testing.T) {
	mock := &mockProjectStore{err: pgx.ErrNoRows}
	h := newTestProjectHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/projects/00000000-0000-0000-0000-000000000001", nil)
	req.SetPathValue("id", "00000000-0000-0000-0000-000000000001")
	w := httptest.NewRecorder()
	h.GetByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProjectHandler_GetByID_InvalidID_Returns400(t *testing.T) {
	mock := &mockProjectStore{}
	h := newTestProjectHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/projects/not-a-uuid", nil)
	req.SetPathValue("id", "not-a-uuid")
	w := httptest.NewRecorder()
	h.GetByID(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProjectHandler_Update_Valid_Returns200(t *testing.T) {
	mock := &mockProjectStore{project: db.Project{Title: "Updated", Color: "#111111"}}
	h := newTestProjectHandler(mock)

	body := `{"title":"Updated","color":"#111111"}`
	req := httptest.NewRequest(http.MethodPut, "/projects/00000000-0000-0000-0000-000000000001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "00000000-0000-0000-0000-000000000001")
	w := httptest.NewRecorder()
	h.Update(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated")
}

func TestProjectHandler_Update_NotFound_Returns404(t *testing.T) {
	mock := &mockProjectStore{err: pgx.ErrNoRows}
	h := newTestProjectHandler(mock)

	body := `{"title":"X","color":"#000000"}`
	req := httptest.NewRequest(http.MethodPut, "/projects/00000000-0000-0000-0000-000000000001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "00000000-0000-0000-0000-000000000001")
	w := httptest.NewRecorder()
	h.Update(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProjectHandler_Delete_Success_Returns204(t *testing.T) {
	mock := &mockProjectStore{}
	h := newTestProjectHandler(mock)

	req := httptest.NewRequest(http.MethodDelete, "/projects/00000000-0000-0000-0000-000000000001", nil)
	req.SetPathValue("id", "00000000-0000-0000-0000-000000000001")
	w := httptest.NewRecorder()
	h.Delete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestProjectHandler_Delete_InvalidID_Returns400(t *testing.T) {
	mock := &mockProjectStore{}
	h := newTestProjectHandler(mock)

	req := httptest.NewRequest(http.MethodDelete, "/projects/bad", nil)
	req.SetPathValue("id", "bad")
	w := httptest.NewRecorder()
	h.Delete(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
