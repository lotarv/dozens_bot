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
}

func (t *DocumentsTransport) getRules(w http.ResponseWriter, r *http.Request) {
	rules, err := t.service.GetRules()
	if err != nil {
		http.Error(w, "failed to fetch rules", http.StatusInternalServerError)
		slog.Error("failed to get rules", "error", err)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(rules)
	if err != nil {
		http.Error(w, "failed to encode rules", http.StatusInternalServerError)
		slog.Error("failed to encode rules", "error", err)
		return
	}
}
