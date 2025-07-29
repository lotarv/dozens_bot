package meetings

import (
	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/meetings/repository"
	"github.com/lotarv/dozens_bot/internal/domains/meetings/service"
	"github.com/lotarv/dozens_bot/internal/domains/meetings/transport"
	user_repo "github.com/lotarv/dozens_bot/internal/domains/user/repository"
	"github.com/lotarv/dozens_bot/internal/storage"
)

type MeetingsController struct {
	repo      *repository.MeetingsRepository
	service   *service.MeetingsService
	transport *transport.MeetingsTransport
}

func NewMeetingsController(storage *storage.Storage, router *chi.Mux, userRepo *user_repo.UsersRepository) *MeetingsController {
	repo := repository.New(storage, userRepo)
	service := service.New(repo)
	transport := transport.New(router, service)

	return &MeetingsController{
		repo:      repo,
		service:   service,
		transport: transport,
	}
}

func (c *MeetingsController) Build() {
	c.transport.RegisterRoutes()
}

func (c *MeetingsController) Run() {

}
