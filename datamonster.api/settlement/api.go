package settlement

import (
	"datamonster/web"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const BaseRoute = "/settlement"

type Controller struct {
	repo *Repo
}

func NewController(repo *Repo) *Controller {
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
	r.Get("/", c.getSettlements)
	r.Post("/", c.createSettlement)
}

func (c Controller) getSettlements(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(web.UserIdKey).(string)
	log.Default().Printf("Retrieving settlements for user %s", userId)
	settlements, repoErr := c.repo.GetAllForUser(r.Context(), userId)
	if repoErr != nil {
		log.Default().Printf("Error retrieving settlements %s", repoErr.Error())
		web.MakeJsonResponse(w, http.StatusInternalServerError, "Error retrieving settlements")
		return
	}
	data := domainListToDto(settlements)
	web.MakeJsonResponse(w, http.StatusOK, data)
}

func (c Controller) createSettlement(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(web.UserIdKey).(string)
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
	settlement := Settlement{
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

func domainListToDto(settlements []Settlement) SettlementsDTO {
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

func domainToDto(s Settlement) SettlementDTO {
	return SettlementDTO{
		Id:                  s.Id,
		Name:                s.Name,
		SurvivalLimit:       s.SurvivalLimit,
		DepartingSurvival:   s.DepartingSurvival,
		CollectiveCognition: s.CollectiveCognition,
		Year:                s.CurrentYear,
	}
}
