package transport

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/auth"
	"github.com/lotarv/dozens_bot/internal/domains/user/types"
	global_types "github.com/lotarv/dozens_bot/internal/types"
)

type service interface {
	CreateUser(ctx context.Context, user *types.User) error
	UpdateUser(ctx context.Context, user *types.User) error
	GetUserByID(ctx context.Context, userID int64) (*types.User, error)
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
	t.router.Post("/api/users", t.CreateUser)
	t.router.Get("/api/users", t.GetAll)
}

func (t *UserTransport) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := t.service.GetAll(r.Context())
	if err != nil {
		fmt.Printf("error: %v", err)
		http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		return
	}
	if users == nil {
		users = []types.User{}
	}
	json.NewEncoder(w).Encode(users)
}

func (t *UserTransport) CreateUser(w http.ResponseWriter, r *http.Request) {
	creds, ok := r.Context().Value(global_types.ContextKeyCredentials).(auth.Credentials)
	if !ok {
		http.Error(w, "credentials not found in context", http.StatusInternalServerError)
		slog.Error("credentials not found in context")
		return
	}

	user, err := t.service.GetUserByID(r.Context(), creds.ID)
	if err != nil {
		// Если пользователь не существует, создаем нового
		user = &types.User{
			FullName: creds.FullName,
			AvatarURL: sql.NullString{
				String: creds.AvatarUrl,
				Valid:  creds.AvatarUrl != "",
			},
			ID:           creds.ID,
			Niche:        "не указана",
			AnnualIncome: 0,
		}
		if err := t.service.CreateUser(r.Context(), user); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			slog.Error("Failed to create user", "error", err)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
