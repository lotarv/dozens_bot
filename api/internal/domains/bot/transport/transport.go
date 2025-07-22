package transport

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	global_types "github.com/lotarv/dozens_bot/internal/types"
)

type service interface {
	StartNewTransaction(userID int64)
}

type BotTransport struct {
	router  *chi.Mux
	service service
}

func New(router *chi.Mux, service service) *BotTransport {
	return &BotTransport{
		router:  router,
		service: service,
	}
}

func (t *BotTransport) RegisterRoutes() {
	t.router.Post("/api/bot/new-transaction", t.startNewTransaction)
}

func (t *BotTransport) startNewTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextKeyUserID).(int64)
	if !ok {
		http.Error(w, "userID not found in context", http.StatusBadRequest)
		slog.Error("userID not found in context")
		return
	}

	t.service.StartNewTransaction(userID)
	w.WriteHeader(http.StatusOK)
}
