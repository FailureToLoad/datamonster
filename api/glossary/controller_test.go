package glossary_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/failuretoload/datamonster/glossary"
	"github.com/failuretoload/datamonster/server"
	"github.com/failuretoload/datamonster/testenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var requester *testenv.Requester

func TestMain(m *testing.M) {
	ctx := context.Background()

	glossaryContainer, err := testenv.NewGlossaryContainer(ctx)
	if err != nil {
		log.Fatalf("unable to set up glossary container: %v", err)
	}
	defer glossaryContainer.Cleanup(ctx)

	controller, err := glossary.NewController(glossaryContainer.URL)
	if err != nil {
		log.Fatalf("unable to create glossary controller: %v", err)
	}

	requester, err = testenv.NewRequester([]server.Controller{controller})
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

type disorder struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Source     string `json:"source"`
	FlavorText string `json:"flavorText,omitempty"`
	Effect     string `json:"effect"`
}

type fightingArt struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Secret bool     `json:"secret"`
	Source string   `json:"source"`
	Text   []string `json:"text"`
	Notes  string   `json:"notes,omitempty"`
}

type innovation struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Source   string `json:"source"`
	Keywords string `json:"keywords"`
	Parent   string `json:"parent,omitempty"`
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

type glossaryResponse struct {
	Disorders    []disorder    `json:"disorders"`
	FightingArts []fightingArt `json:"fightingArts"`
	Innovations  []innovation  `json:"innovations"`
	Knowledge    []knowledge   `json:"knowledge"`
}

func TestGetAllDisorders(t *testing.T) {
	body, status := requester.GetAllDisorders("test-user")
	require.Equal(t, http.StatusOK, status)

	var disorders []disorder
	require.NoError(t, json.NewDecoder(body).Decode(&disorders))
	require.Len(t, disorders, 2)

	validateDisorders(t, disorders)
}

func TestGetDisorder(t *testing.T) {
	body, status := requester.GetDisorder("test-user", "019412a0-0001-7000-8000-000000000001")
	require.Equal(t, http.StatusOK, status)

	var d disorder
	require.NoError(t, json.NewDecoder(body).Decode(&d))
	assert.Equal(t, "019412a0-0001-7000-8000-000000000001", d.ID)
	assert.Equal(t, "Test Disorder", d.Name)
	assert.Equal(t, "core", d.Source)
	assert.Equal(t, "Bad things happen", d.Effect)
}

func TestGetAllFightingArts(t *testing.T) {
	body, status := requester.GetAllFightingArts("test-user")
	require.Equal(t, http.StatusOK, status)

	var arts []fightingArt
	require.NoError(t, json.NewDecoder(body).Decode(&arts))
	require.Len(t, arts, 2)

	validateFightingArts(t, arts)
}

func TestGetFightingArt(t *testing.T) {
	body, status := requester.GetFightingArt("test-user", "019412a0-0003-7000-8000-000000000003")
	require.Equal(t, http.StatusOK, status)

	var art fightingArt
	require.NoError(t, json.NewDecoder(body).Decode(&art))
	assert.Equal(t, "019412a0-0003-7000-8000-000000000003", art.ID)
	assert.Equal(t, "Normal", art.Name)
	assert.Equal(t, false, art.Secret)
	assert.Equal(t, "core", art.Source)
	assert.Equal(t, []string{"You can swing a sword."}, art.Text)
}

func TestGetAllInnovations(t *testing.T) {
	body, status := requester.GetAllInnovations("test-user")
	require.Equal(t, http.StatusOK, status)

	var innovations []innovation
	require.NoError(t, json.NewDecoder(body).Decode(&innovations))
	require.Len(t, innovations, 2)

	validateInnovations(t, innovations)
}

func TestGetInnovation(t *testing.T) {
	body, status := requester.GetInnovation("test-user", "019412a0-0006-7000-8000-000000000006")
	require.Equal(t, http.StatusOK, status)

	var inn innovation
	require.NoError(t, json.NewDecoder(body).Decode(&inn))
	assert.Equal(t, "019412a0-0006-7000-8000-000000000006", inn.ID)
	assert.Equal(t, "Buffoonery", inn.Name)
	assert.Equal(t, "core", inn.Source)
	assert.Equal(t, "shame", inn.Keywords)
	assert.Equal(t, "019412a0-0005-7000-8000-000000000005", inn.Parent)
}

func TestGetAllKnowledge(t *testing.T) {
	body, status := requester.GetAllKnowledge("test-user")
	require.Equal(t, http.StatusOK, status)

	var items []knowledge
	require.NoError(t, json.NewDecoder(body).Decode(&items))
	require.Len(t, items, 2)
	validateKnowledge(t, items)
}

func TestGetKnowledge(t *testing.T) {
	body, status := requester.GetKnowledge("test-user", "019412a0-0008-7000-8000-000000000008")
	require.Equal(t, http.StatusOK, status)

	var k knowledge
	require.NoError(t, json.NewDecoder(body).Decode(&k))
	assert.Equal(t, "019412a0-0008-7000-8000-000000000008", k.ID)
	assert.Equal(t, "Advanced Knowledge", k.Name)
	assert.Equal(t, 3, k.Cost)
	assert.Equal(t, true, k.Tenet)
	assert.Equal(t, "nonsense", k.Type)
	assert.Equal(t, []string{"Deep thoughts."}, k.Description)
	assert.Equal(t, "requires 019412a0-0007-7000-8000-000000000007", k.Condition)
}

func TestGetDisorder_Unauthorized(t *testing.T) {
	t.Cleanup(requester.Unauthorized())
	_, status := requester.GetDisorder("unauthorized", "019412a0-0001-7000-8000-000000000001")
	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestGetGlossary(t *testing.T) {
	body, status := requester.GetGlossary("test-user")
	require.Equal(t, http.StatusOK, status)

	var g glossaryResponse
	require.NoError(t, json.NewDecoder(body).Decode(&g))
	require.Len(t, g.Disorders, 2)
	validateDisorders(t, g.Disorders)

	require.Len(t, g.FightingArts, 2)
	validateFightingArts(t, g.FightingArts)

	require.Len(t, g.Innovations, 2)
	validateInnovations(t, g.Innovations)

	require.Len(t, g.Knowledge, 2)
	validateKnowledge(t, g.Knowledge)
}

func validateDisorders(t *testing.T, items []disorder) {
	for _, d := range items {
		assert.NotEmpty(t, d.ID)
		assert.NotEmpty(t, d.Name)
		assert.NotEmpty(t, d.Source)
		assert.NotEmpty(t, d.Effect)
	}
}

func validateFightingArts(t *testing.T, items []fightingArt) {
	var hasSecret, hasNonSecret bool
	for _, art := range items {
		assert.NotEmpty(t, art.ID)
		assert.NotEmpty(t, art.Name)
		assert.NotEmpty(t, art.Source)
		assert.NotEmpty(t, art.Text)

		if art.Name == "Normal" {
			assert.False(t, art.Secret)
			hasNonSecret = true
		}
		if art.Name == "Secret" {
			assert.True(t, art.Secret)
			hasSecret = true
		}
	}
	assert.True(t, hasSecret, "expected a secret fighting art")
	assert.True(t, hasNonSecret, "expected a non-secret fighting art")
}

func validateInnovations(t *testing.T, items []innovation) {
	for _, inn := range items {
		assert.NotEmpty(t, inn.ID)
		assert.NotEmpty(t, inn.Name)
		assert.NotEmpty(t, inn.Source)
		assert.NotEmpty(t, inn.Keywords)
	}
}

func validateKnowledge(t *testing.T, items []knowledge) {
	for _, k := range items {
		assert.NotEmpty(t, k.ID)
		assert.NotEmpty(t, k.Name)
		assert.NotEmpty(t, k.Type)
		assert.NotEmpty(t, k.Description)
	}
}
