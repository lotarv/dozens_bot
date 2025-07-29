package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/meetings/types"
	global_types "github.com/lotarv/dozens_bot/internal/types"
)

type service interface {
	GetDozenMeetings(ctx context.Context, userID int64) ([]types.Meeting, error)
}

type MeetingsTransport struct {
	router  *chi.Mux
	service service
}

func New(router *chi.Mux, service service) *MeetingsTransport {
	return &MeetingsTransport{
		router:  router,
		service: service,
	}
}

func (t *MeetingsTransport) RegisterRoutes() {
	t.router.Get("/api/meetings", t.handleGetMeetings)
}

func (t *MeetingsTransport) handleGetMeetings(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextKeyUserID).(int64)
	if !ok {
		http.Error(w, "userID not found in context", http.StatusBadRequest)
		slog.Error("userID not found in context")
		return
	}

	meetings, err := t.service.GetDozenMeetings(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not in dozen", http.StatusUnauthorized)
		slog.Error("user not in dozen", "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(meetings)
}
