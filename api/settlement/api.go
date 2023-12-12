package settlement

import (
	postgres "datamonster/settlement/repo"
	"datamonster/web"
	"fmt"
	"log"
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
	r.Group(func(r chi.Router) {
		web.SetDefaultMiddleware(r)
		web.SetCorsHandler(r)
		web.SetAuthHandler(r)
		r.Get("/settlement", c.getSettlements)
		r.Post("/settlement", c.createSettlement)
		r.Get("/settlement/{id}", c.getSettlement)
	})
}

func (c Controller) getSettlements(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(web.UserIdKey).(int)
	log.Default().Println("---")
	log.Default().Printf("Requesting settlements for user %d", userId)
	query := fmt.Sprintf("SELECT * FROM campaign.settlement WHERE owner = %d", userId)
	settlements, repoErr := c.repo.Select(r.Context(), query)
	if repoErr != nil {
		log.Default().Printf("Error retrieving settlements %s", repoErr.Error())
		web.MakeJsonResponse(w, http.StatusInternalServerError, "Error retrieving settlements")
		return
	}
	data := domainListToDto(settlements)
	log.Default().Printf("Returning settlements for user %d", userId)
	log.Default().Println("---")
	web.MakeJsonResponse(w, http.StatusOK, data)
}

func (c Controller) createSettlement(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(web.UserIdKey).(int)
	var body CreateSettlementRequest
	err := web.DecodeJson(r.Body, &body)
	if err != nil {
		web.MakeJsonResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if body.Name == "" {
		web.MakeJsonResponse(w, http.StatusBadRequest, "Name is required")
		return
	}
	log.Default().Printf("Creating settlement: %s", body)
	settlement := postgres.Settlement{
		Owner:               userId,
		Name:                body.Name,
		SurvivalLimit:       1,
		DepartingSurvival:   0,
		CollectiveCognition: 0,
		CurrentYear:         1,
	}
	newId, insertErr := c.repo.Insert(r.Context(), settlement)
	if insertErr != nil {
		log.Default().Printf("Error inserting settlement %s", insertErr.Error())
		web.MakeJsonResponse(w, http.StatusInternalServerError, "Unable to create settlement")
		return
	}

	settlement.Id = newId
	dto := domainToDto(settlement)
	web.MakeJsonResponse(w, http.StatusOK, dto)
}

func (c Controller) getSettlement(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(web.UserIdKey).(int)
	settlementId := chi.URLParam(r, "id")
	log.Default().Printf("Retrieving settlement %s for user %d", settlementId, userId)
	settlement, repoErr := c.repo.Get(r.Context(), settlementId)
	if repoErr != nil {
		log.Default().Printf("Error retrieving settlement %s", repoErr.Error())
		web.MakeJsonResponse(w, http.StatusInternalServerError, "Error retrieving settlement")
		return
	}
	if settlement.Owner != userId {
		log.Default().Printf("User %d does not own settlement %s", userId, settlementId)
		web.MakeJsonResponse(w, http.StatusForbidden, "You do not own this settlement")
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
