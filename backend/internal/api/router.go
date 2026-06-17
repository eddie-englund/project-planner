package api

import (
	"log/slog"
	"net/http"
	"os"

	db "github.com/eddie-englund/project-planner/backend/internal/db/generated"
)

func corsMiddleware(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

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

	sh := &StatusHandler{db: queries, logger: logger}
	mux.HandleFunc("GET /projects/{projectId}/statuses", sh.List)
	mux.HandleFunc("POST /projects/{projectId}/statuses", sh.Create)

	tkh := &TicketHandler{db: queries, logger: logger}
	mux.HandleFunc("POST /projects/{projectId}/topics/{topicId}/tickets", tkh.Create)
	mux.HandleFunc("GET /projects/{projectId}/topics/{topicId}/tickets", tkh.List)
	mux.HandleFunc("GET /projects/{projectId}/topics/{topicId}/tickets/{id}", tkh.GetByID)
	mux.HandleFunc("PUT /projects/{projectId}/topics/{topicId}/tickets/{id}", tkh.Update)
	mux.HandleFunc("DELETE /projects/{projectId}/topics/{topicId}/tickets/{id}", tkh.Delete)

	origin := os.Getenv("CORS_ORIGIN")
	if origin == "" {
		origin = "http://localhost:5173"
	}
	return corsMiddleware(origin, mux)
}
