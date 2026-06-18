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

type mockTopicStore struct {
	topic  db.ProjectTopic
	topics []db.ProjectTopic
	err    error
}

func (m *mockTopicStore) CreateProjectTopic(_ context.Context, _ db.CreateProjectTopicParams) (db.ProjectTopic, error) {
	return m.topic, m.err
}
func (m *mockTopicStore) ListProjectTopics(_ context.Context, _ pgtype.UUID) ([]db.ProjectTopic, error) {
	return m.topics, m.err
}
func (m *mockTopicStore) GetProjectTopicByID(_ context.Context, _ db.GetProjectTopicByIDParams) (db.ProjectTopic, error) {
	return m.topic, m.err
}
func (m *mockTopicStore) UpdateProjectTopic(_ context.Context, _ db.UpdateProjectTopicParams) (db.ProjectTopic, error) {
	return m.topic, m.err
}
func (m *mockTopicStore) DeleteProjectTopic(_ context.Context, _ pgtype.UUID) error {
	return m.err
}

func newTestTopicHandler(store topicStore) *TopicHandler {
	return &TopicHandler{db: store, logger: slog.Default()}
}

const validProjectID = "00000000-0000-0000-0000-000000000001"
const validTopicID = "00000000-0000-0000-0000-000000000002"

func TestTopicHandler_List_ReturnsTopics(t *testing.T) {
	mock := &mockTopicStore{topics: []db.ProjectTopic{
		{Title: "Frontend", Color: "#aabbcc"},
		{Title: "Backend", Color: "#ccbbaa"},
	}}
	h := newTestTopicHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/projects/"+validProjectID+"/topics", nil)
	req.SetPathValue("projectId", validProjectID)
	w := httptest.NewRecorder()
	h.List(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Frontend")
	assert.Contains(t, w.Body.String(), "Backend")
}

func TestTopicHandler_List_DBError_Returns500(t *testing.T) {
	mock := &mockTopicStore{err: errors.New("db error")}
	h := newTestTopicHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/projects/"+validProjectID+"/topics", nil)
	req.SetPathValue("projectId", validProjectID)
	w := httptest.NewRecorder()
	h.List(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTopicHandler_Create_ValidBody_Returns201(t *testing.T) {
	mock := &mockTopicStore{topic: db.ProjectTopic{Title: "New Topic", Color: "#ff0000"}}
	h := newTestTopicHandler(mock)

	body := `{"title":"New Topic","color":"#ff0000","index":0}`
	req := httptest.NewRequest(http.MethodPost, "/projects/"+validProjectID+"/topics", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("projectId", validProjectID)
	w := httptest.NewRecorder()
	h.Create(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "New Topic")
}

func TestTopicHandler_Create_MissingTitle_Returns400(t *testing.T) {
	mock := &mockTopicStore{}
	h := newTestTopicHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/projects/"+validProjectID+"/topics", strings.NewReader(`{"color":"#ff0000"}`))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("projectId", validProjectID)
	w := httptest.NewRecorder()
	h.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTopicHandler_GetByID_Found_Returns200(t *testing.T) {
	mock := &mockTopicStore{topic: db.ProjectTopic{Title: "Found Topic", Color: "#123456"}}
	h := newTestTopicHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/projects/"+validProjectID+"/topics/"+validTopicID, nil)
	req.SetPathValue("projectId", validProjectID)
	req.SetPathValue("id", validTopicID)
	w := httptest.NewRecorder()
	h.GetByID(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Found Topic")
}

func TestTopicHandler_GetByID_NotFound_Returns404(t *testing.T) {
	mock := &mockTopicStore{err: pgx.ErrNoRows}
	h := newTestTopicHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/projects/"+validProjectID+"/topics/"+validTopicID, nil)
	req.SetPathValue("projectId", validProjectID)
	req.SetPathValue("id", validTopicID)
	w := httptest.NewRecorder()
	h.GetByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTopicHandler_Update_Valid_Returns200(t *testing.T) {
	mock := &mockTopicStore{topic: db.ProjectTopic{Title: "Updated Topic", Color: "#abcdef"}}
	h := newTestTopicHandler(mock)

	body := `{"title":"Updated Topic","color":"#abcdef","index":1}`
	req := httptest.NewRequest(http.MethodPut, "/projects/"+validProjectID+"/topics/"+validTopicID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", validTopicID)
	w := httptest.NewRecorder()
	h.Update(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated Topic")
}

func TestTopicHandler_Delete_Success_Returns204(t *testing.T) {
	mock := &mockTopicStore{}
	h := newTestTopicHandler(mock)

	req := httptest.NewRequest(http.MethodDelete, "/projects/"+validProjectID+"/topics/"+validTopicID, nil)
	req.SetPathValue("id", validTopicID)
	w := httptest.NewRecorder()
	h.Delete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
