package glossary

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

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

type collection[T mappable] struct {
	idMap map[string]T
	list  []T
}

func newCollection[T mappable](definitions []T) collection[T] {
	m := toMap(definitions)
	return collection[T]{
		list:  definitions,
		idMap: m,
	}
}

func (c collection[T]) All() []T {
	return c.list
}

func (c collection[T]) Get(id string) T {
	return c.idMap[id]
}

type Controller struct {
	disorders    collection[disorder]
	fightingarts collection[fightingArt]
	innovations  collection[innovation]
	knowledge    collection[knowledge]
}

func NewController(glossaryServerURL string) (*Controller, error) {
	disorderList, err := fetchAll[disorder](glossaryServerURL + "/disorders")
	if err != nil {
		return nil, fmt.Errorf("unable to fetch disorders: %w", err)
	}
	disorderCollection := newCollection(disorderList)

	fightingArtList, err := fetchAll[fightingArt](glossaryServerURL + "/fighting-arts")
	if err != nil {
		return nil, fmt.Errorf("unable to fetch fighting arts: %w", err)
	}
	fightingArtCollection := newCollection(fightingArtList)

	innovationList, err := fetchAll[innovation](glossaryServerURL + "/innovations")
	if err != nil {
		return nil, fmt.Errorf("unable to fetch innovations: %w", err)
	}
	innovationCollection := newCollection(innovationList)

	knowledgeList, err := fetchAll[knowledge](glossaryServerURL + "/knowledge")
	if err != nil {
		return nil, fmt.Errorf("unable to fetch knowledge: %w", err)
	}
	knowledgeCollection := newCollection(knowledgeList)

	return &Controller{
		innovations:  innovationCollection,
		disorders:    disorderCollection,
		fightingarts: fightingArtCollection,
		knowledge:    knowledgeCollection,
	}, nil
}

func (c Controller) RegisterRoutes(r chi.Router) {
	r.Get("/glossary/disorders", c.allDisorders)
	r.Get("/glossary/disorders/{id}", c.getDisorder)
	r.Get("/glossary/fightingarts", c.allFightingArts)
	r.Get("/glossary/fightingarts/{id}", c.getFightingArt)
	r.Get("/glossary/innovations", c.allInnovations)
	r.Get("/glossary/innovations/{id}", c.getInnovation)
	r.Get("/glossary/knowledge", c.allKnowledge)
	r.Get("/glossary/knowledge/{id}", c.getKnowledge)
}

func (c Controller) allDisorders(w http.ResponseWriter, r *http.Request) {
	response.OK(r.Context(), w, c.disorders.All())
}

func (c Controller) getDisorder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := idParam(r)
	if id == "" {
		response.BadRequest(ctx, w, fmt.Errorf("invalid id"))
		return
	}

	response.OK(r.Context(), w, c.disorders.Get(id))
}

func (c Controller) allFightingArts(w http.ResponseWriter, r *http.Request) {
	response.OK(r.Context(), w, c.fightingarts.All())
}

func (c Controller) getFightingArt(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := idParam(r)
	if id == "" {
		response.BadRequest(ctx, w, fmt.Errorf("invalid id"))
		return
	}

	response.OK(r.Context(), w, c.fightingarts.Get(id))
}

func (c Controller) allInnovations(w http.ResponseWriter, r *http.Request) {
	response.OK(r.Context(), w, c.innovations.All())
}

func (c Controller) getInnovation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := idParam(r)
	if id == "" {
		response.BadRequest(ctx, w, fmt.Errorf("invalid id"))
		return
	}

	response.OK(r.Context(), w, c.innovations.Get(id))
}

func (c Controller) allKnowledge(w http.ResponseWriter, r *http.Request) {
	response.OK(r.Context(), w, c.knowledge.All())
}

func (c Controller) getKnowledge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := idParam(r)
	if id == "" {
		response.BadRequest(ctx, w, fmt.Errorf("invalid id"))
		return
	}

	response.OK(r.Context(), w, c.knowledge.Get(id))
}

func fetchAll[T mappable](definitionURI string) ([]T, error) {
	_, err := url.Parse(definitionURI)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid definition uri: %w", definitionURI, err)
	}

	resp, err := http.Get(definitionURI)
	if err != nil {
		return nil, fmt.Errorf("unable to request definitions: %w", err)
	}
	defer tryClose(resp.Body)

	var result []T
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("unable to decode definitions from response")
	}

	return result, nil
}

func toMap[T mappable](collection []T) map[string]T {
	target := make(map[string]T)

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
