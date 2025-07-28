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
	GetDeclarations(username string) ([]types.DeclarationDocument, error)
	GetDeclarationByID(id string) (*types.DeclarationDocument, error)
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
	t.router.Get("/api/declarations/{username}", t.getDeclarations)
	t.router.Get("/api/declaration/{id}", t.getDeclaration)
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
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reports)
}

func (t *DocumentsTransport) getDeclarations(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	declarations, err := t.service.GetDeclarations(username)
	if err != nil {
		http.Error(w, "failed to fetch declarations", http.StatusInternalServerError)
		slog.Error("failed to fetch declarations", "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(declarations)
}

func (t *DocumentsTransport) getDeclaration(w http.ResponseWriter, r *http.Request) {
	declarationID := chi.URLParam(r, "id")

	declaration, err := t.service.GetDeclarationByID(declarationID)
	if err != nil {
		http.Error(w, "failed to get declaration by id", http.StatusInternalServerError)
		slog.Error("failed to get declaration by id", "declarationID", declarationID, "error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(declaration)
}
