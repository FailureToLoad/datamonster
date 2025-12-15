package settlement

import (
	"context"
	"errors"
	"net/http"

	"github.com/failuretoload/datamonster/request"
	"github.com/failuretoload/datamonster/response"
	"github.com/failuretoload/datamonster/settlement/domain"
	"github.com/go-chi/chi/v5"
)

type Repo interface {
	All(ctx context.Context, userID string) ([]domain.Settlement, error)
}

type Controller struct {
	records Repo
}

func NewController(r Repo) (*Controller, error) {
	if r == nil {
		return nil, errors.New("repo cannot be nil")
	}

	return &Controller{records: r}, nil
}

func (c Controller) RegisterRoutes(r chi.Router) {
	r.Get("/settlements", c.getSettlements)
}

func (c Controller) getSettlements(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := request.UserID(ctx)
	settlements, repoErr := c.records.All(ctx, userID)
	if repoErr != nil {
		response.InternalServerError(w, "unable to retrieve settlements", repoErr)
		return
	}

	if settlements == nil {
		settlements = []domain.Settlement{}
	}

	response.OK(w, settlements)
}
