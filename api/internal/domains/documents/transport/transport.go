package transport

import (
	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/documents/types"
)

type service interface {
	GetAllReports(user_url string) ([]types.Report, error)
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
	// t.router.Get("/api/{user_id}/reports")
}
