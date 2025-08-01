package piggy_bank

import (
	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/piggy_bank/repository"
	"github.com/lotarv/dozens_bot/internal/domains/piggy_bank/service"
	"github.com/lotarv/dozens_bot/internal/domains/piggy_bank/transport"
	"github.com/lotarv/dozens_bot/internal/storage"
)

type PiggyBankController struct {
	repo      *repository.PiggyBankRepository
	service   *service.PiggyBankService
	transport *transport.PiggyBankTransport
}

func NewPiggyBankController(router *chi.Mux, storage *storage.Storage) *PiggyBankController {
	repo := repository.New(storage)
	service := service.New(repo)
	transport := transport.New(router, service)

	return &PiggyBankController{
		repo:      repo,
		service:   service,
		transport: transport,
	}
}

func (c *PiggyBankController) Build() {
	c.transport.RegisterRoutes()
}

func (c *PiggyBankController) Run() {

}

func (c *PiggyBankController) GetRepository() *repository.PiggyBankRepository {
	return c.repo
}
