package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/user/types"
)

type service interface {
	GetAll(ctx context.Context) ([]types.User, error)
}

type UserTransport struct {
	router  *chi.Mux
	service service
}

func New(router *chi.Mux, service service) *UserTransport {
	return &UserTransport{
		router:  router,
		service: service,
	}
}

func (t *UserTransport) RegisterRoutes() {
	t.router.Get("/api/users", t.GetAll)
}

func (t *UserTransport) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := t.service.GetAll(r.Context())
	if err != nil {
		fmt.Printf("error: %v", err)
		http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}
