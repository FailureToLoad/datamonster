package settlement

import (
	"datamonster/web"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

const BaseRoute = "/settlement"

type Controller struct {
	repo *Repo
}

func NewController(repo *Repo) *Controller {
	return &Controller{repo: repo}
}

type SettlementDTO struct {
	Id                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	SurvivalLimit       int       `json:"limit"`
	DepartingSurvival   int       `json:"departing"`
	CollectiveCognition int       `json:"cc"`
	Year                int       `json:"year"`
}

type SettlementsDTO struct {
	Settlements []SettlementDTO `json:"settlements"`
	Count       int             `json:"count"`
}

type SettlementsRequest struct {
	Owner string `json:"owner"`
}

func (c Controller) RegisterRoutes(r chi.Router) {
	r.Get("/", c.getSettlements)
}

func (c Controller) getSettlements(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(web.UserIdKey).(string)
	settlements, repoErr := c.repo.GetAllForUser(r.Context(), userId)
	if repoErr != nil {
		web.MakeJsonResponse(w, http.StatusInternalServerError, "Error retrieving settlements")
		return
	}
	data := domainListToDto(settlements)
	web.MakeJsonResponse(w, http.StatusOK, data)
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
