package transport

import (
	"github.com/go-chi/chi/v5"
)

type service interface {
}

type PiggyBankTransport struct {
	router  *chi.Mux
	service service
}

func New(router *chi.Mux, service service) *PiggyBankTransport {
	return &PiggyBankTransport{
		router:  router,
		service: service,
	}
}

func (t *PiggyBankTransport) RegisterRoutes() {
}
