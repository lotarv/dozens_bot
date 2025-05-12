package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/members/types"
)

type service interface {
	GetMembers() ([]types.Member, error)
}

type MembersTransport struct {
	router  *chi.Mux
	service service
}

func New(router *chi.Mux, service service) *MembersTransport {
	return &MembersTransport{
		router:  router,
		service: service,
	}
}

func (t *MembersTransport) RegisterRoutes() {
	t.router.Get("/api/members", t.handleGetMembers)
}

func (t *MembersTransport) handleGetMembers(w http.ResponseWriter, r *http.Request) {
	slog.Info("function called")
	members, err := t.service.GetMembers()
	if err != nil {
		http.Error(w, "failed to get members", http.StatusInternalServerError)
		slog.Error("failed to get members", "error", err)
		return
	}
	// w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(members)
	if err != nil {
		http.Error(w, "failed to send members", http.StatusInternalServerError)
		slog.Error("failed to encode members", "members", members, "error", err)
		return
	}
}
