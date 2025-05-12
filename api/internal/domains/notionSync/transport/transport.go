package transport

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type service interface {
	SyncMembersWithNotion() error
	SyncDeclarationsWithNotion() error
	SyncReportsWithNotion() error
	SyncDocumentsWithNotion() error
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
	t.router.Post("/api/sync-declarations", t.syncDeclarationsWithNotion)
	t.router.Post("/api/sync-reports", t.syncReportsWithNotion)
	t.router.Post("/api/sync-documents", t.syncDocumentsWithNotion)
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

func (t *NotionSyncTransport) syncDeclarationsWithNotion(w http.ResponseWriter, r *http.Request) {
	err := t.service.SyncDeclarationsWithNotion()
	if err != nil {
		slog.Info("failed to synchronize declarations with notion", "error", err)
		http.Error(w, "failed to synchronize declarations with notion", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully synchronized declarations with notion"))
}

func (t *NotionSyncTransport) syncReportsWithNotion(w http.ResponseWriter, r *http.Request) {
	err := t.service.SyncReportsWithNotion()
	if err != nil {
		slog.Info("failed to synchronize reports with notion", "error", err)
		http.Error(w, "failed to synchronize reports with notion", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully synchronized reports with notion"))
}

func (t *NotionSyncTransport) syncDocumentsWithNotion(w http.ResponseWriter, r *http.Request) {
	err := t.service.SyncDocumentsWithNotion()
	if err != nil {
		slog.Info("failed to synchronize documents with notion", "error", err)
		http.Error(w, "failed to synchronize documents with notion", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully synchronized documents with notion"))
}
