package survivor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/failuretoload/datamonster/request"
	"github.com/failuretoload/datamonster/survivor/domain"
	"github.com/gofrs/uuid/v5"

	"github.com/failuretoload/datamonster/response"

	"github.com/go-chi/chi/v5"
)

type Repo interface {
	All(ctx context.Context, settlementID uuid.UUID) ([]domain.Survivor, error)
	Create(ctx context.Context, d domain.Survivor) (domain.Survivor, error)
	Update(ctx context.Context, settlementID, survivorID uuid.UUID, updates map[string]int) (domain.Survivor, error)
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
	r.Group(func(gr chi.Router) {
		gr.Use(settlementIDToContext)
		gr.Get("/settlements/{id}/survivors", c.getSurvivors)
		gr.Post("/settlements/{id}/survivors", c.createSurvivor)
		gr.Patch("/settlements/{id}/survivors/{survivorID}", c.updateSurvivor)
	})
}

func (c Controller) getSurvivors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	settlementID := request.SettlementID(ctx)
	survivors, err := c.db.All(ctx, settlementID)
	if err != nil {
		response.InternalServerError(ctx, w, fmt.Errorf("error retrieving survivors: %w", err))
		return
	}

	response.OK(ctx, w, survivors)
}

func (c Controller) createSurvivor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	survivorDTO := domain.Survivor{}
	decodeErr := request.DecodeJSON(r.Body, &survivorDTO)
	if decodeErr != nil {
		response.InternalServerError(ctx, w, fmt.Errorf("unable to decode request body: %w", decodeErr))
		return
	}

	if survivorDTO.Name == "" {
		response.BadRequest(ctx, w, fmt.Errorf("survivor name is required"))
		return
	}

	survivorDTO.SettlementID = request.SettlementID(ctx)
	survivor, err := c.db.Create(ctx, survivorDTO)
	if err != nil {
		response.InternalServerError(ctx, w, fmt.Errorf("error creating survivor: %w", err))
		return
	}

	response.OK(ctx, w, survivor)
}

func (c Controller) updateSurvivor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var updates map[string]int
	if err := request.DecodeJSON(r.Body, &updates); err != nil {
		response.BadRequest(ctx, w, fmt.Errorf("unable to decode request body: %w", err))
		return
	}

	settlementID := request.SettlementID(ctx)
	survivorID, err := uuid.FromString(chi.URLParam(r, "survivorID"))
	if err != nil {
		response.BadRequest(ctx, w, fmt.Errorf("invalid survivor id"))
		return
	}

	survivor, err := c.db.Update(ctx, settlementID, survivorID, updates)
	if err != nil {
		response.InternalServerError(ctx, w, fmt.Errorf("error updating survivor: %w", err))
		return
	}

	response.OK(ctx, w, survivor)
}

func settlementIDToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := request.SettlementIDFromURL(r)
		if err != nil {
			response.InternalServerError(ctx, w, err)
			return
		}

		ctx = request.SetSettlementID(ctx, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
