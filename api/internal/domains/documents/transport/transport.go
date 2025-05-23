package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/documents/types"
)

type service interface {
	GetRules() (types.Document, error)
	GetReports(username string) (*types.ReportsResponse, error)
}

type DocumentsTransport struct {
	router  *chi.Mux
	service service
}

func New(router *chi.Mux, service service) *DocumentsTransport {
	return &DocumentsTransport{
		router:  router,
		service: service,
	}
}

func (t *DocumentsTransport) RegisterRoutes() {
	t.router.Get("/api/rules", t.getRules)
	t.router.Get("/api/reports/{username}", t.getReports)
}

func (t *DocumentsTransport) getRules(w http.ResponseWriter, r *http.Request) {
	rules, err := t.service.GetRules()
	if err != nil {
		http.Error(w, "failed to fetch rules", http.StatusInternalServerError)
		slog.Error("failed to get rules", "error", err)
		return
	}

	// w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(rules)
	if err != nil {
		http.Error(w, "failed to encode rules", http.StatusInternalServerError)
		slog.Error("failed to encode rules", "error", err)
		return
	}
}

func (t *DocumentsTransport) getReports(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	reports, err := t.service.GetReports(username)
	if err != nil {
		http.Error(w, "failed to fetch reports", http.StatusInternalServerError)
		slog.Error("failed to fetch reports", "error", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reports)
}
