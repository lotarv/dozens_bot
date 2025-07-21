package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/piggy_bank/types"
)

type service interface {
	GetPiggyBank(ctx context.Context, bank_id int) (types.PiggyBank, error)
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
	t.router.Get("/api/piggy-bank", t.getPiggyBank)
}

func (t *PiggyBankTransport) getPiggyBank(w http.ResponseWriter, r *http.Request) {
	//TODO: сделать копилки для разных десяток
	piggy_bank, err := t.service.GetPiggyBank(r.Context(), 1)
	if err != nil {
		slog.Error("failed to get piggy bank: ", "error", err)
		http.Error(w, "failed to get piggy bank", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(piggy_bank)
}
