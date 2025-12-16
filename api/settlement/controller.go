package settlement

import (
	"context"
	"errors"
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
		Insert(ctx context.Context, s domain.Settlement) (int, error)
		Get(ctx context.Context, userID string, settlementID uuid.UUID) (domain.Settlement, error)
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
		return nil, errors.New("repo cannot be nil")
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
		msg := "userID is required"
		response.BadRequest(w, msg, errors.New(msg))
	}
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

func (c Controller) createSettlement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := request.UserID(ctx)
	if userID == "" {
		msg := "userID is required"
		response.BadRequest(w, msg, errors.New(msg))
	}
	var body CreateSettlementRequest
	err := request.DecodeJSONRequest(r.Body, &body)
	if err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}
	if body.Name == "" {
		response.BadRequest(w, "name is required", nil)
		return
	}

	settlementID, err := c.records.Insert(ctx, domain.Settlement{Name: body.Name})
	if err != nil {
		response.InternalServerError(w, "unable to persist settlement", err)
	} else {
		response.OK(w, settlementID)
	}

}

func (c Controller) getSettlement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := request.UserID(ctx)
	if userID == "" {
		msg := "userID is required"
		response.BadRequest(w, msg, errors.New(msg))
	}

	settlementID, err := request.IDParam(r)
	if err != nil {
		response.BadRequest(w, "invalid settlement id", err)
	}

	settlement, repoErr := c.records.Get(ctx, userID, settlementID)
	if repoErr != nil {
		response.InternalServerError(w, "unable to retrieve settlement", repoErr)
		return
	}

	response.OK(w, settlement)
}
