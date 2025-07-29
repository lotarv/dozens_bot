package transport

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/meetings/types"
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
	slog.Info("GET MEETINGS")
}
