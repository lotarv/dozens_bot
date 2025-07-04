package documents

import (
	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/documents/repository"
	"github.com/lotarv/dozens_bot/internal/domains/documents/service"
	"github.com/lotarv/dozens_bot/internal/domains/documents/transport"
	"github.com/lotarv/dozens_bot/internal/storage"
)

type DocumentsController struct {
	repo      *repository.DocumentsRepository
	service   *service.DocumentsService
	transport *transport.DocumentsTransport
}

func NewDocumentsController(router *chi.Mux, db *storage.Storage) *DocumentsController {
	repo := repository.New(db.DB())
	service := service.New(repo)
	transport := transport.New(router, service)

	return &DocumentsController{
		repo:      repo,
		service:   service,
		transport: transport,
	}
}

func (c *DocumentsController) Build() {
	c.transport.RegisterRoutes()
}

func (c *DocumentsController) Run() {

}
