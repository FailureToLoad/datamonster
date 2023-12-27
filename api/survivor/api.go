package survivor

import (
	"datamonster/store"
	"datamonster/web"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	repo *postgresRepo
}

func NewController(conn store.Connection) *Controller {
	repo := newRepo(conn)
	return &Controller{repo: repo}
}

func (c Controller) RegisterRoutes(r chi.Router, authHandler func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		web.SetDefaultMiddleware(r)
		web.SetCorsHandler(r)
		r.Use(web.SettlementIdExtractor)
		r.Use(authHandler)
		r.Get("/settlement/{id}/survivor", c.getSurvivors)
	})
}

func (c Controller) getSurvivors(w http.ResponseWriter, r *http.Request) {
	settlementId := r.Context().Value(web.SettlementIdKey).(int)
	query := getAllSurvivorsForSettlement(settlementId)
	survivors, err := c.repo.Find(r.Context(), query)
	if err != nil {
		web.MakeJsonResponse(w, http.StatusInternalServerError, "Error retrieving survivors")
		return
	}
	data := dtoListFromDomain(survivors)
	web.MakeJsonResponse(w, http.StatusOK, data)
}

type SurvivorDTO struct {
	Id               int    `json:"id"`
	Settlement       int    `json:"settlement"`
	Name             string `json:"name"`
	Born             int    `json:"born"`
	Gender           string `json:"gender"`
	Status           string `json:"status"`
	HuntXp           int    `json:"huntXp"`
	Survival         int    `json:"survival"`
	Movement         int    `json:"movement"`
	Accuracy         int    `json:"accuracy"`
	Strength         int    `json:"strength"`
	Evasion          int    `json:"evasion"`
	Luck             int    `json:"luck"`
	Speed            int    `json:"speed"`
	Insanity         int    `json:"insanity"`
	SystemicPressure int    `json:"systemicPressure"`
	Torment          int    `json:"torment"`
	Lumi             int    `json:"lumi"`
	Courage          int    `json:"courage"`
	Understanding    int    `json:"understanding"`
}

func dtoFromDomain(s Survivor) SurvivorDTO {
	return SurvivorDTO(s)
}

func dtoListFromDomain(s []Survivor) []SurvivorDTO {
	dtos := make([]SurvivorDTO, len(s))
	for i, v := range s {
		dtos[i] = dtoFromDomain(v)
	}
	return dtos
}
