package members

import (
	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/members/repository"
	"github.com/lotarv/dozens_bot/internal/domains/members/service"
	"github.com/lotarv/dozens_bot/internal/domains/members/transport"
	"github.com/lotarv/dozens_bot/internal/storage"
)

type MembersController struct {
	repo      *repository.MembersRepository
	service   *service.MembersService
	transport *transport.MembersTransport
}

func NewMembersController(router *chi.Mux, storage *storage.Storage) *MembersController {
	repo := repository.New(storage)
	service := service.New(repo)
	transport := transport.New(router, service)

	return &MembersController{
		repo:      repo,
		service:   service,
		transport: transport,
	}
}

func (c *MembersController) Build() {
	c.transport.RegisterRoutes()
}

func (c *MembersController) Run() {

}
