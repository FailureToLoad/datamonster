package survivor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/failuretoload/datamonster/survivor/domain"
	"github.com/gofrs/uuid/v5"

	"github.com/failuretoload/datamonster/response"

	"github.com/go-chi/chi/v5"
)

type Repo interface {
	All(ctx context.Context, settlementID uuid.UUID) ([]domain.Survivor, error)
	Create(ctx context.Context, d domain.Survivor) (domain.Survivor, error)
}

type Controller struct {
	db Repo
}

func NewController(r Repo) (*Controller, error) {
	if r == nil {
		return nil, fmt.Errorf("repo cannot be nil")
	}
	return &Controller{db: r}, nil
}

func (c Controller) RegisterRoutes(r chi.Router) {
	r.Get("/settlements/{id}/survivors", c.getSurvivors)
	r.Post("/settlements/{id}/survivors", c.createSurvivor)
}

func (c Controller) getSurvivors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	param := chi.URLParam(r, "id")
	settlementID, convErr := uuid.FromString(param)
	if convErr != nil {
		response.InternalServerError(ctx, w, fmt.Errorf("unable to convert query param: %w", convErr))
		return
	}
	survivors, err := c.db.All(ctx, settlementID)
	if err != nil {
		response.InternalServerError(ctx, w, fmt.Errorf("error retrieving survivors: %w", err))
		return
	}

	response.OK(ctx, w, survivors)
}

func (c Controller) createSurvivor(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	ctx := r.Context()
	settlementID, convErr := uuid.FromString(param)
	if convErr != nil {
		response.InternalServerError(ctx, w, fmt.Errorf("invalid settlement id: %w", convErr))
		return
	}
	survivorDTO := domain.Survivor{}
	decodeErr := json.NewDecoder(r.Body).Decode(&survivorDTO)
	if decodeErr != nil || survivorDTO.Name == "" {
		response.InternalServerError(ctx, w, fmt.Errorf("unable to decode request body: %w", decodeErr))
		return
	}
	survivorDTO.Settlement = settlementID
	survivor, err := c.db.Create(ctx, survivorDTO)
	if err != nil {
		response.InternalServerError(ctx, w, fmt.Errorf("error creating survivor: %w", err))
		return
	}

	response.OK(ctx, w, survivor)
}
