package survivor

import (
	"context"
	"github.com/failuretoload/datamonster/store"
	"github.com/failuretoload/datamonster/web"
	"net/http"
	"strconv"

	"github.com/supertokens/supertokens-golang/recipe/session"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	repo *postgresRepo
}

func NewController(conn store.Connection) *Controller {
	repo := newRepo(conn)
	return &Controller{repo: repo}
}

func (c Controller) RegisterRoutes(r chi.Router) {
	r.Use(settlementIdExtractor)
	r.Get("/settlement/{id}/survivor", session.VerifySession(nil, c.getSurvivors))
}

func (c Controller) getSurvivors(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	settlementId, convErr := strconv.Atoi(param)
	if convErr != nil {
		web.MakeJsonResponse(w, http.StatusInternalServerError, "unable to convert query param")
		return
	}
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
	survivors := make([]SurvivorDTO, len(s))
	for i, v := range s {
		survivors[i] = dtoFromDomain(v)
	}
	return survivors
}

type ctxSettlementIdKey string

const SettlementIdKey ctxSettlementIdKey = "settlementId"

func settlementIdExtractor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		settlementIdString := chi.URLParam(r, "id")
		if settlementIdString != "" {
			settlementId, convErr := strconv.Atoi(settlementIdString)
			if convErr != nil {
				web.MakeJsonResponse(w, http.StatusBadRequest, "settlement id should be a number")
				return
			}
			ctx := context.WithValue(r.Context(), SettlementIdKey, settlementId)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
