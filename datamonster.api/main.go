package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowOriginFunc:    func(r *http.Request, origin string) bool { return true },
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		MaxAge:             3599, // Maximum value not ignored by any of major browsers
	})
	r.Use(c.Handler)
	contentType := render.SetContentType(render.ContentTypeJSON)
	timeout := middleware.Timeout(20 * time.Second)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(contentType)
	r.Use(timeout)
	r.Get("/settlement", GetSettlement)
	http.ListenAndServe(":8888", r)
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

func GetSettlement(w http.ResponseWriter, r *http.Request) {
	resp := &SettlementResponse{
		SettlementName:    "Rad Dreamers",
		SettlementType:    "PotDK",
		SurvivalLimit:     1,
		DepartingSurvival: 1,
	}
	render.Render(w, r, resp)
}
