package user

import (
	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/user/repository"
	"github.com/lotarv/dozens_bot/internal/domains/user/service"
	"github.com/lotarv/dozens_bot/internal/domains/user/transport"
	"github.com/lotarv/dozens_bot/internal/storage"
)

type UserController struct {
	repo      *repository.UsersRepository
	service   *service.UserService
	transport *transport.UserTransport
}

func NewUserController(router *chi.Mux, storage *storage.Storage) *UserController {
	repo := repository.New(storage)
	service := service.New(repo)
	transport := transport.New(router, service)

	return &UserController{
		repo:      repo,
		service:   service,
		transport: transport,
	}
}

func (c *UserController) Build() {
	c.transport.RegisterRoutes()
}

func (c *UserController) Run() {

}
