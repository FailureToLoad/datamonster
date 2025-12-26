package settlement

import (
	"context"
	"fmt"
	"net/http"

	"github.com/failuretoload/datamonster/request"
	"github.com/failuretoload/datamonster/response"
	"github.com/failuretoload/datamonster/settlement/domain"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

type (
	Repo interface {
		All(ctx context.Context, userID string) ([]domain.Settlement, error)
		Insert(ctx context.Context, s domain.Settlement) (uuid.UUID, error)
		Get(ctx context.Context, userID string, settlementID uuid.UUID) (*domain.Settlement, error)
	}
	Controller struct {
		records Repo
	}
	CreateSettlementRequest struct {
		Name string `json:"name"`
	}
)

func NewController(r Repo) (*Controller, error) {
	if r == nil {
		return nil, fmt.Errorf("repo cannot be nil")
	}

	return &Controller{records: r}, nil
}

func (c Controller) RegisterRoutes(r chi.Router) {
	r.Get("/settlements", c.getSettlements)
	r.Post("/settlements", c.createSettlement)
	r.Get("/settlements/{id}", c.getSettlement)
}

func (c Controller) getSettlements(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := request.UserID(ctx)
	if userID == "" {
		response.BadRequest(ctx, w, fmt.Errorf("userID is required"))
		return
	}
	settlements, repoErr := c.records.All(ctx, userID)
	if repoErr != nil {
		response.InternalServerError(ctx, w, fmt.Errorf("unable to retrieve settlements: %w", repoErr))
		return
	}

	if settlements == nil {
		settlements = []domain.Settlement{}
	}

	response.OK(ctx, w, settlements)
}

func (c Controller) createSettlement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := request.UserID(ctx)
	if userID == "" {
		response.BadRequest(ctx, w, fmt.Errorf("userID is required"))
		return
	}
	var body CreateSettlementRequest
	err := request.DecodeJSON(r.Body, &body)
	if err != nil {
		response.BadRequest(ctx, w, fmt.Errorf("invalid request body: %w", err))
		return
	}
	if body.Name == "" {
		response.BadRequest(ctx, w, fmt.Errorf("name is required"))
		return
	}

	settlementID, err := c.records.Insert(ctx, domain.Settlement{Name: body.Name, Owner: userID})
	if err != nil {
		response.InternalServerError(ctx, w, fmt.Errorf("unable to persist settlement: %w", err))
	} else {
		response.OK(ctx, w, settlementID)
	}
}

func (c Controller) getSettlement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := request.UserID(ctx)
	if userID == "" {
		response.BadRequest(ctx, w, fmt.Errorf("userID is required"))
		return
	}

	settlementID, err := request.SettlementIDFromURL(r)
	if err != nil {
		response.BadRequest(ctx, w, err)
		return
	}

	settlement, repoErr := c.records.Get(ctx, userID, settlementID)
	if repoErr != nil {
		response.InternalServerError(ctx, w, fmt.Errorf("unable to retrieve settlement: %w", repoErr))
		return
	}
	if settlement == nil {
		response.NotFound(ctx, w, fmt.Errorf("settlement not found"))
		return
	}

	response.OK(ctx, w, settlement)
}
