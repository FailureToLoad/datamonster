package settlement

import (
	postgres "datamonster/settlement/repo"
	"datamonster/web"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	repo *postgres.PostgresRepo
}

func NewController(repo *postgres.PostgresRepo) *Controller {
	return &Controller{repo: repo}
}

type SettlementDTO struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	SurvivalLimit       int    `json:"limit"`
	DepartingSurvival   int    `json:"departing"`
	CollectiveCognition int    `json:"cc"`
	Year                int    `json:"year"`
}

type SettlementsDTO struct {
	Settlements []SettlementDTO `json:"settlements"`
	Count       int             `json:"count"`
}

type CreateSettlementRequest struct {
	Name string `json:"name"`
}

func (c Controller) RegisterRoutes(r chi.Router) {
	r.Get("/settlement", session.VerifySession(nil, c.getSettlements))
	r.Post("/settlement", session.VerifySession(nil, c.createSettlement))
	r.Route("/settlement/{id}", func(r chi.Router) {
		r.Get("/", session.VerifySession(nil, c.getSettlement))
	})
}

func (c Controller) getSettlements(w http.ResponseWriter, r *http.Request) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	userID := sessionContainer.GetUserID()
	settlements, repoErr := c.repo.Select(r.Context(), userID)
	if repoErr != nil {
		web.MakeJsonResponse(w, http.StatusInternalServerError, "Error retrieving settlements")
		return
	}
	data := domainListToDto(settlements)
	web.MakeJsonResponse(w, http.StatusOK, data)
}

func (c Controller) createSettlement(w http.ResponseWriter, r *http.Request) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	userID := sessionContainer.GetUserID()
	var body CreateSettlementRequest
	err := web.DecodeJsonRequest(r.Body, &body)
	if err != nil {
		web.MakeJsonResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if body.Name == "" {
		web.MakeJsonResponse(w, http.StatusBadRequest, "Name is required")
		return
	}
	settlement := postgres.Settlement{
		Owner:               userID,
		Name:                body.Name,
		SurvivalLimit:       1,
		DepartingSurvival:   0,
		CollectiveCognition: 0,
		CurrentYear:         1,
	}
	newId, insertErr := c.repo.Insert(r.Context(), settlement)
	if insertErr != nil {
		web.MakeJsonResponse(w, http.StatusInternalServerError, "Unable to create settlement")
		return
	}

	settlement.Id = newId
	dto := domainToDto(settlement)
	web.MakeJsonResponse(w, http.StatusOK, dto)
}

func (c Controller) getSettlement(w http.ResponseWriter, r *http.Request) {
	settlementId := chi.URLParam(r, "id")
	settlement, repoErr := c.repo.Get(r.Context(), settlementId)
	if repoErr != nil {
		web.MakeJsonResponse(w, http.StatusInternalServerError, "Error retrieving settlement")
		return
	}
	dto := domainToDto(settlement)
	web.MakeJsonResponse(w, http.StatusOK, dto)
}

func domainListToDto(settlements []postgres.Settlement) SettlementsDTO {
	dtos := SettlementsDTO{}
	for _, s := range settlements {
		dto := SettlementDTO{
			Id:                  s.Id,
			Name:                s.Name,
			SurvivalLimit:       s.SurvivalLimit,
			DepartingSurvival:   s.DepartingSurvival,
			CollectiveCognition: s.CollectiveCognition,
			Year:                s.CurrentYear,
		}
		dtos.Settlements = append(dtos.Settlements, dto)
		dtos.Count++
	}
	return dtos
}

func domainToDto(s postgres.Settlement) SettlementDTO {
	return SettlementDTO{
		Id:                  s.Id,
		Name:                s.Name,
		SurvivalLimit:       s.SurvivalLimit,
		DepartingSurvival:   s.DepartingSurvival,
		CollectiveCognition: s.CollectiveCognition,
		Year:                s.CurrentYear,
	}
}
