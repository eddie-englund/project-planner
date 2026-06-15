package api

import (
	"log/slog"
	"net/http"

	db "github.com/eddie-englund/project-planner/backend/internal/db/generated"
)

func NewRouter(logger *slog.Logger, queries *db.Queries) http.Handler {
	logger.Info("Registered router")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	ph := &ProjectHandler{db: queries, logger: logger}
	mux.HandleFunc("POST /projects", ph.Create)
	mux.HandleFunc("GET /projects", ph.List)
	mux.HandleFunc("GET /projects/{id}", ph.GetByID)
	mux.HandleFunc("PUT /projects/{id}", ph.Update)
	mux.HandleFunc("DELETE /projects/{id}", ph.Delete)

	th := &TopicHandler{db: queries, logger: logger}
	mux.HandleFunc("POST /projects/{projectId}/topics", th.Create)
	mux.HandleFunc("GET /projects/{projectId}/topics", th.List)
	mux.HandleFunc("GET /projects/{projectId}/topics/{id}", th.GetByID)
	mux.HandleFunc("PUT /projects/{projectId}/topics/{id}", th.Update)
	mux.HandleFunc("DELETE /projects/{projectId}/topics/{id}", th.Delete)

	return mux
}
