package notionSync

import (
	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/notionSync/repository"
	"github.com/lotarv/dozens_bot/internal/domains/notionSync/service"
	"github.com/lotarv/dozens_bot/internal/domains/notionSync/transport"
	"github.com/lotarv/dozens_bot/internal/storage"
)

type NotionSyncController struct {
	repo      *repository.NotionSyncRepository
	service   *service.NotionSyncService
	transport *transport.NotionSyncTransport
}

func NewNotionSyncController(router *chi.Mux, storage *storage.Storage) *NotionSyncController {
	repo := repository.New(storage.DB())
	service := service.New(repo)
	transport := transport.New(router, service)

	return &NotionSyncController{
		repo:      repo,
		service:   service,
		transport: transport,
	}

}

func (c *NotionSyncController) Build() {
	c.transport.RegisterRoutes()
}

func (c *NotionSyncController) Run() {

}
