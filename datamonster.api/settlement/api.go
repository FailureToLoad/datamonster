package settlement

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const BaseRoute = "/settlement"

func RegisterRoutes(r chi.Router) {
	r.Get("/", getSettlement)
}

type SettlementResponse struct {
	SettlementName    string `json:"name"`
	SettlementType    string `json:"type"`
	SurvivalLimit     int    `json:"limit"`
	DepartingSurvival int    `json:"departing"`
	Elapsed           int64  `json:"elapsed"`
}

func (rd *SettlementResponse) Render(w http.ResponseWriter, r *http.Request) error {
	rd.Elapsed = 10
	return nil
}

func getSettlement(w http.ResponseWriter, r *http.Request) {
	resp := &SettlementResponse{
		SettlementName:    "Rad Dreamers",
		SettlementType:    "PotDK",
		SurvivalLimit:     1,
		DepartingSurvival: 1,
	}
	render.Render(w, r, resp)
}
