package glossary

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/failuretoload/datamonster/logger"
	"github.com/failuretoload/datamonster/response"
	"github.com/go-chi/chi/v5"
)

type disorder struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Source     string `json:"source"`
	FlavorText string `json:"flavorText,omitempty"`
	Effect     string `json:"effect"`
}

func (d disorder) Key() string {
	return d.ID
}

type fightingArt struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Secret bool     `json:"secret"`
	Source string   `json:"source"`
	Text   []string `json:"text"`
	Notes  string   `json:"notes,omitempty"`
}

func (fa fightingArt) Key() string {
	return fa.ID
}

type innovation struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Source   string `json:"source"`
	Keywords string `json:"keywords"`
	Parent   string `json:"parent,omitempty"`
}

func (i innovation) Key() string {
	return i.ID
}

type knowledge struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Cost         int      `json:"cost"`
	Tenet        bool     `json:"tenet"`
	Type         string   `json:"type"`
	Description  []string `json:"description"`
	Condition    string   `json:"condition,omitempty"`
	Observations int      `json:"observations,omitempty"`
	Advance      string   `json:"advance,omitempty"`
	Activation   string   `json:"activation,omitempty"`
}

func (k knowledge) Key() string {
	return k.ID
}

type mappable interface {
	any
	Key() string
}

type glossary struct {
	Disorders    []disorder    `json:"disorders"`
	Fightingarts []fightingArt `json:"fightingArts"`
	Innovations  []innovation  `json:"innovations"`
	Knowledge    []knowledge   `json:"knowledge"`
}

type Controller struct {
	bulk         glossary
	disorders    map[string]disorder
	fightingarts map[string]fightingArt
	innovations  map[string]innovation
	knowledge    map[string]knowledge
}

func NewController(glossaryServerURL string) (*Controller, error) {
	glossary, err := fetchGlossary(glossaryServerURL + "/glossary")
	if err != nil {
		return nil, fmt.Errorf("unable to fetch glossary: %w", err)
	}

	return &Controller{
		bulk:         glossary,
		innovations:  toMap(glossary.Innovations),
		disorders:    toMap(glossary.Disorders),
		fightingarts: toMap(glossary.Fightingarts),
		knowledge:    toMap(glossary.Knowledge),
	}, nil
}

func (c Controller) RegisterRoutes(r chi.Router) {
	r.Get("/glossary", c.getGlossary)
	r.Get("/glossary/disorders", c.allDisorders)
	r.Get("/glossary/disorders/{id}", c.getDisorder)
	r.Get("/glossary/fightingarts", c.allFightingArts)
	r.Get("/glossary/fightingarts/{id}", c.getFightingArt)
	r.Get("/glossary/innovations", c.allInnovations)
	r.Get("/glossary/innovations/{id}", c.getInnovation)
	r.Get("/glossary/knowledge", c.allKnowledge)
	r.Get("/glossary/knowledge/{id}", c.getKnowledge)
}

func (c Controller) getGlossary(w http.ResponseWriter, r *http.Request) {
	response.OK(r.Context(), w, c.bulk)
}

func (c Controller) allDisorders(w http.ResponseWriter, r *http.Request) {
	response.OK(r.Context(), w, c.bulk.Disorders)
}

func (c Controller) getDisorder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := idParam(r)
	if id == "" {
		response.BadRequest(ctx, w, fmt.Errorf("invalid id"))
		return
	}

	response.OK(r.Context(), w, c.disorders[id])
}

func (c Controller) allFightingArts(w http.ResponseWriter, r *http.Request) {
	response.OK(r.Context(), w, c.bulk.Fightingarts)
}

func (c Controller) getFightingArt(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := idParam(r)
	if id == "" {
		response.BadRequest(ctx, w, fmt.Errorf("invalid id"))
		return
	}

	response.OK(r.Context(), w, c.fightingarts[id])
}

func (c Controller) allInnovations(w http.ResponseWriter, r *http.Request) {
	response.OK(r.Context(), w, c.bulk.Innovations)
}

func (c Controller) getInnovation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := idParam(r)
	if id == "" {
		response.BadRequest(ctx, w, fmt.Errorf("invalid id"))
		return
	}

	response.OK(r.Context(), w, c.innovations[id])
}

func (c Controller) allKnowledge(w http.ResponseWriter, r *http.Request) {
	response.OK(r.Context(), w, c.bulk.Knowledge)
}

func (c Controller) getKnowledge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := idParam(r)
	if id == "" {
		response.BadRequest(ctx, w, fmt.Errorf("invalid id"))
		return
	}

	response.OK(r.Context(), w, c.knowledge[id])
}

func fetchGlossary(uri string) (glossary, error) {
	_, err := url.Parse(uri)
	if err != nil {
		return glossary{}, fmt.Errorf("%s is not a valid glossary uri: %w", uri, err)
	}

	resp, err := http.Get(uri)
	if err != nil {
		return glossary{}, fmt.Errorf("unable to request glossary: %w", err)
	}
	defer tryClose(resp.Body)

	var result glossary
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return glossary{}, fmt.Errorf("unable to decode glossary from response: %w", err)
	}

	return result, nil
}

func toMap[T mappable](collection []T) map[string]T {
	target := make(map[string]T)

	if len(collection) == 0 {
		logger.Warn(context.Background(), fmt.Sprintf("glossary: %T is empty", collection))
	}
	for _, item := range collection {
		target[item.Key()] = item
	}

	return target
}

func tryClose(rc io.ReadCloser) {
	err := rc.Close()
	if err != nil {
		slog.Log(context.Background(), slog.LevelWarn, "error closing ReadCloser", slog.Any("error", err))
	}
}

func idParam(r *http.Request) string {
	return chi.URLParam(r, "id")
}
