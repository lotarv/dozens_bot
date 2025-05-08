package transport

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type service interface {
	SyncMembersWithNotion() error
}

type NotionSyncTransport struct {
	router  *chi.Mux
	service service
}

func New(router *chi.Mux, service service) *NotionSyncTransport {
	return &NotionSyncTransport{
		router:  router,
		service: service,
	}
}

func (t *NotionSyncTransport) RegisterRoutes() {
	t.router.Post("/api/sync-members", t.syncMembersWithNotion)
}

func (t *NotionSyncTransport) syncMembersWithNotion(w http.ResponseWriter, r *http.Request) {
	err := t.service.SyncMembersWithNotion()
	if err != nil {
		slog.Info("failed to synchronize members with notion", "error", err)
		http.Error(w, "failed to synchronize members with notion", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully synchronized members with notion"))

}
